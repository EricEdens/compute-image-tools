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
	mock_disk "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/disk/mocks"
	ovfutils "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_import/ovf_utils"
	"github.com/golang/mock/gomock"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/disk"
)

func Test_MultiDiskImporter_CancelPropagatesToEachImporter_SingleBootDisk(t *testing.T) {

	osID := "ubuntu-1804"
	diskURI := "disk/uri"
	env := disk.Environment{}
	overrides := disk.TranslationOverrides{}
	cancelReason := "cancel"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	bootDiskImporter := mock_disk.NewMockImporter(mockCtrl)
	bootDiskImporter.EXPECT().Import(diskURI, nil).Do(func() {
		time.Sleep(time.Hour)
	})
	mockProvider := mock_disk.NewMockImporterProvider(mockCtrl)
	mockProvider.EXPECT().NewTranslatingDiskImporter(env, overrides).Return(bootDiskImporter)

	underTest := newDefaultMultiDiskImporter(mockProvider, env, overrides)
	go func() {
		underTest.Import(osID, &ovfutils.DiskInfo{FilePath: diskURI}, nil)
	}()
	underTest.Cancel(cancelReason)
}
