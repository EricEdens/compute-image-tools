//  Copyright 2020 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ovfimporter

import (
	"errors"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/disk"
	ovfutils "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_import/ovf_utils"
	"sync"
)

type MultiDiskImporter interface {
	Import(osID string, bootDisk *ovfutils.DiskInfo, dataDisks []*ovfutils.DiskInfo) (bootDiskURI string, dataDiskURIs []string, err error)
	Cancel(reason string) bool
}

func newDefaultMultiDiskImporter(
	singleDiskImporterProvider disk.ImporterProvider,
	environment disk.Environment,
	overrides disk.TranslationOverrides) MultiDiskImporter {

	panic("TODO")
}

type state int32

const (
	runnable  = state(0)
	running   = state(1)
	cancelled = state(2)
	done      = state(3)
)

type importResult struct {
	bootURI, dataURI string
	err              error
}
type defaultMultiDiskImporter struct {
	singleDiskImporterProvider disk.ImporterProvider
	environment                disk.Environment
	overrides                  disk.TranslationOverrides

	importState      state
	importStateMutex sync.Mutex

	// cancelChan is monitored by the main thread. When an error is found, the main thread
	// cancels all pending imports.
	cancelChan chan error

	// inFlightDisks in non-zero when there are disks pending to be returned to the caller.
	// It is decremented when:
	//   - A disk is cleaned up (following cancellation or other error).
	//   - A disk is returned to the caller (it's now the caller's responsibility to clean it up).
	inFlightDisks sync.WaitGroup
}

func (m *defaultMultiDiskImporter) Import(osID string, bootDisk *ovfutils.DiskInfo, dataDisks []*ovfutils.DiskInfo) (bootDiskURI string, dataDiskURIs []string, err error) {
	// Quit early if:
	//  - Another thread called `Cancel`, or
	//  - `Import` has already been called
	// The first is plausible, the second is a programming bug.
	m.importStateMutex.Lock()
	if m.importState == cancelled {
		return "", nil, errors.New("Import cancelled")
	} else if m.importState != runnable {
		panic("Import can only be run a single time")
	}
	m.importState = running
	m.importStateMutex.Unlock()

	bootDiskURI, dataDiskURIs, err = m.runImports(bootDisk, dataDisks)

	m.importStateMutex.Lock()
	defer m.importStateMutex.Unlock()
	if err != nil {
		m.importState = cancelled
	} else {
		m.importState = done
	}
	return bootDiskURI, dataDiskURIs, err
}

func (m *defaultMultiDiskImporter) runImports(bootDisk *ovfutils.DiskInfo, dataDisks []*ovfutils.DiskInfo) (bootDiskURI string, dataDiskURIs []string, err error) {

	m.inFlightDisks.Add(1 + len(dataDisks))

	importers, resultsChan := m.startImportsInGoRoutines(bootDisk, dataDisks)
	results, cancelReason := m.waitForCompletion(importers, resultsChan)

	if cancelReason != nil {
		m.cleanup(results)
		bootDiskURI = ""
		dataDiskURIs = nil
	} else {
		for _, result := range results {
			m.inFlightDisks.Done()
			if result.bootURI != "" {
				bootDiskURI = result.bootURI
			} else if result.dataURI != "" {
				dataDiskURIs = append(dataDiskURIs, result.dataURI)
			}
		}
	}
	return bootDiskURI, dataDiskURIs, err
}

func (m *defaultMultiDiskImporter) waitForCompletion(importers []disk.Importer, resultsChan chan importResult) ([]importResult, error) {
	var results []importResult
	var cancelReason error
	for len(results) < len(importers) {
		var currError error
		select {
		case currError = <-m.cancelChan:
		case res := <-resultsChan:
			results = append(results, res)
			currError = res.err
		}
		if currError != nil && cancelReason == nil {
			for _, importer := range importers {
				importer.Cancel(currError.Error())
			}
			cancelReason = currError
		}
	}
	return results, cancelReason
}

func (m *defaultMultiDiskImporter) startImportsInGoRoutines(bootDisk *ovfutils.DiskInfo, dataDisks []*ovfutils.DiskInfo) ([]disk.Importer, chan importResult) {
	var importers []disk.Importer
	resultsChan := make(chan importResult)

	importer := m.singleDiskImporterProvider.NewTranslatingDiskImporter(m.environment, m.overrides)
	importers = append(importers, importer)
	go func() {
		uri, err := importer.Import(disk.FileOrImageReference(bootDisk.FilePath), nil)
		resultsChan <- importResult{bootURI: string(uri), err: err}
	}()
	for _, dataDisk := range dataDisks {
		importer := m.singleDiskImporterProvider.NewDataDiskImporter(m.environment)
		importers = append(importers, importer)
		toImport := disk.FileOrImageReference(dataDisk.FilePath)
		go func() {
			uri, err := importer.Import(toImport, nil)
			resultsChan <- importResult{dataURI: string(uri), err: err}
		}()
	}
	return importers, resultsChan
}

func (m *defaultMultiDiskImporter) Cancel(reason string) bool {
	m.importStateMutex.Lock()
	switch m.importState {
	case done:
		return true
	case runnable:
		m.importState = cancelled
		return false
	case cancelled:
		return false
	}
	m.importStateMutex.Unlock()

	// state is running
	m.cancelChan <- errors.New(reason)
	m.inFlightDisks.Wait()
	m.importStateMutex.Lock()
	defer m.importStateMutex.Unlock()
	return m.importState == done
}

func (m *defaultMultiDiskImporter) cleanup(results []importResult) {
	for _, result := range results {
		var toDelete string
		if result.bootURI != "" {
			toDelete = result.bootURI
		} else if result.dataURI != "" {
			toDelete = result.dataURI
		}
		if toDelete != "" {
			panic("TODO: Implement delete disk")
		}
		m.inFlightDisks.Done()
	}
}
