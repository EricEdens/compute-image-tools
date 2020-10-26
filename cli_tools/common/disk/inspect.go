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

package disk

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/distro"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/daisycommon"
	"github.com/GoogleCloudPlatform/compute-image-tools/daisy"
	"github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pb"
)

const (
	workflowFile = "image_import/inspection/boot-inspect.wf.json"
)

// Inspector finds partition and boot-related properties for a disk.
type Inspector interface {
	// Inspect finds partition and boot-related properties for a disk and
	// returns an InspectionResult. The reference is implementation specific.
	Inspect(reference string, inspectOS bool) (InspectionResult, error)
	Cancel(reason string) bool
	TraceLogs() []string
}

// InspectionResult contains the partition and boot-related properties of a disk.
type InspectionResult struct {
	// UEFIBootable indicates whether the disk is bootable with UEFI.
	UEFIBootable bool

	// BIOSBootableWithHybridMBROrProtectiveMBR indicates whether the disk is BIOS-bootable
	// or "hybrid MBR" mode or "protective MBR" mode.
	BIOSBootableWithHybridMBROrProtectiveMBR bool

	// RootFS indicates the file system type of the partition containing
	// the root directory ("/").
	RootFS string

	Architecture, Distro, Major, Minor string
}

// NewInspector creates an Inspector that can inspect GCP disks.
// A GCE instance runs the inspection; network and subnet are used
// for its network interface.
func NewInspector(wfAttributes daisycommon.WorkflowAttributes, network string, subnet string) (Inspector, error) {
	wf, err := daisy.NewFromFile(path.Join(wfAttributes.WorkflowDirectory, workflowFile))
	if err != nil {
		return nil, err
	}
	daisycommon.SetWorkflowAttributes(wf, wfAttributes)
	wf.Vars["network"] = daisy.Var{Value: network}
	wf.Vars["subnet"] = daisy.Var{Value: subnet}
	return &bootInspector{[]string{}, &defaultDaisyWorker{wf}}, nil
}

// bootInspect implements disk.Inspector using the Python boot-inspect package,
// executed on a worker VM using Daisy.
type bootInspector struct {
	traceLogs []string
	worker    daisyWorker
}

// daisyWorker is a facade over daisy.Workflow to facilitate mocking.
type daisyWorker interface {
	runAndReadSerialValue(key string, vars map[string]string) (string, error)
	cancel(reason string) bool
	traceLogs() []string
}

func (i *bootInspector) Cancel(reason string) bool {
	i.tracef("Canceling with reason: %q", reason)
	return i.worker.cancel(reason)
}

func (i *bootInspector) TraceLogs() []string {
	var combined []string
	combined = append(combined, i.traceLogs...)
	combined = append(combined, i.worker.traceLogs()...)
	return combined
}

// Inspect finds partition and boot-related properties for a GCP persistent disk, and
// returns an InspectionResult. `reference` is a fully-qualified PD URI, such as
// "projects/project-name/zones/us-central1-a/disks/disk-name". `inspectOS` is a flag
// to determine whether to inspect OS on the disk.
func (i *bootInspector) Inspect(reference string, inspectOS bool) (InspectionResult, error) {
	results := &pb.InspectionResults{}

	// Run the inspection worker.
	vars := map[string]string{
		"pd_uri":        reference,
		"is_inspect_os": strconv.FormatBool(inspectOS),
	}
	encodedProto, err := i.worker.runAndReadSerialValue("inspect_pb", vars)
	if err != nil {
		return i.assembleErrors(reference, results, pb.InspectionResults_RUNNING_WORKER, err)
	}

	// Decode the base64-encoded proto.
	bytes, err := base64.StdEncoding.DecodeString(encodedProto)
	if err == nil {
		err = proto.Unmarshal(bytes, results)
	}
	if err != nil {
		return i.assembleErrors(reference, results, pb.InspectionResults_DECODING_WORKER_RESPONSE, err)
	}
	i.tracef("Detection results: %s", results.String())

	// Validate the results.
	if err = i.validate(results); err != nil {
		return i.assembleErrors(reference, results, pb.InspectionResults_INTERPRETING_INSPECTION_RESULTS, err)
	}

	if err = i.populate(results); err != nil {
		return i.assembleErrors(reference, results, pb.InspectionResults_INTERPRETING_INSPECTION_RESULTS, err)
	}

	return createLegacyResults(results), nil
}

// tracef formats according to a format specifier and appends the results to the trace logs.
func (i *bootInspector) tracef(format string, a ...interface{}) {
	i.traceLogs = append(i.traceLogs, fmt.Sprintf(format, a...))
}

// assembleErrors sets the errorWhen field, and generates an error object.
func (i *bootInspector) assembleErrors(pdURI string, results *pb.InspectionResults,
	errorWhen pb.InspectionResults_ErrorWhen, err error) (InspectionResult, error) {
	results.ErrorWhen = errorWhen
	if err != nil {
		err = fmt.Errorf("failed to inspect %v: %w", pdURI, err)
	} else {
		err = fmt.Errorf("failed to inspect %v", pdURI)
	}
	return createLegacyResults(results), err
}

// createLegacyResults converts pb.InspectionResults to InspectionResult.
func createLegacyResults(pbResults *pb.InspectionResults) (results InspectionResult) {
	if pbResults.OsCount == 1 && pbResults.OsRelease != nil {
		results = InspectionResult{

			Distro: pbResults.OsRelease.GetDistro(),
			Major:  pbResults.OsRelease.GetMajorVersion(),
			Minor:  pbResults.OsRelease.GetMinorVersion(),
		}
		if pbResults.OsRelease.Architecture != pb.Architecture_ARCHITECTURE_UNKNOWN {
			results.Architecture = strings.ToLower(pbResults.OsRelease.Architecture.String())
		}
	}
	results.UEFIBootable = pbResults.GetUefiBootable()
	results.BIOSBootableWithHybridMBROrProtectiveMBR = pbResults.GetBiosBootable()
	results.RootFS = pbResults.RootFs
	return results
}

// validate checks the fields from a pb.InspectionResults object for logical consistency, returning
// an error if an issue is found.
func (i *bootInspector) validate(results *pb.InspectionResults) error {
	// Only populate OsRelease when one OS is found.
	if results.OsCount != 1 {
		if results.OsRelease != nil {
			return fmt.Errorf(
				"Worker should not return OsRelease when NumOsFound != 1: NumOsFound=%d", results.OsCount)
		}
		return nil
	}

	if results.OsRelease.CliFormatted != "" {
		return errors.New("Worker should not return CliFormatted")
	}

	if results.OsRelease.Distro != "" {
		return errors.New("Worker should not return Distro name, only DistroId")
	}

	if results.OsRelease.MajorVersion == "" {
		return errors.New("Missing MajorVersion")
	}

	if results.OsRelease.Architecture == pb.Architecture_ARCHITECTURE_UNKNOWN {
		return errors.New("Missing Architecture")
	}

	if results.OsRelease.DistroId == pb.Distro_DISTRO_UNKNOWN {
		return errors.New("Missing DistroId")
	}

	return nil
}

// populate fills the fields in the pb.InspectionResults that are not returned by the worker.
// This is required since the worker is unaware of import-specific idioms, such as the formatting
// used by gcloud's --os argument.
func (i *bootInspector) populate(results *pb.InspectionResults) error {
	if results.ErrorWhen == pb.InspectionResults_NO_ERROR && results.OsCount == 1 {
		distroEnum, major, minor := results.OsRelease.DistroId,
			results.OsRelease.MajorVersion, results.OsRelease.MinorVersion

		distroName := strings.ReplaceAll(strings.ToLower(results.OsRelease.GetDistroId().String()), "_", "-")

		results.OsRelease.Distro = distroName
		version, err := distro.FromComponents(distroName, major, minor)
		if err != nil {
			i.tracef("Failed to interpret version distro=%q, major=%q, minor=%q: %v",
				distroEnum, major, minor, err)
		} else {
			results.OsRelease.CliFormatted = version.AsGcloudArg()
		}
	}
	return nil
}

type defaultDaisyWorker struct {
	wf *daisy.Workflow
}

// runAndReadSerialValue runs the daisy workflow with the supplied vars, and returns the serial
// output value associated with the supplied key.
func (w *defaultDaisyWorker) runAndReadSerialValue(key string, vars map[string]string) (string, error) {
	for k, v := range vars {
		w.wf.AddVar(k, v)
	}
	err := w.wf.Run(context.Background())
	if err != nil {
		return "", err
	}
	return w.wf.GetSerialConsoleOutputValue(key), nil
}

func (w *defaultDaisyWorker) cancel(reason string) bool {
	if w.wf != nil {
		w.wf.CancelWithReason(reason)
		return true
	}

	//indicate cancel was not performed
	return false
}

func (w *defaultDaisyWorker) traceLogs() []string {
	if w.wf != nil && w.wf.Logger != nil {
		return w.wf.Logger.ReadSerialPortLogs()
	}
	return []string{}
}
