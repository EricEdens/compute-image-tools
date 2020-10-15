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

package import_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/imagefile"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/logging"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/storage"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_vm_image_import/importer"
	daisycompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
)

const (
	workflowDir = "../../../daisy_workflows"
)

func TestDaisyFileInflation(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fileURI string
	}{
		{
			name:    "no resizing required",
			fileURI: "gs://compute-image-tools-test-resources/file-inflation-test/virt-8G.vmdk",
		},
		{
			name:    "resize dest",
			fileURI: "gs://compute-image-tools-test-resources/file-inflation-test/virt-12G.vmdk",
		},
		{
			name:    "resize scratch",
			fileURI: "gs://compute-image-tools-test-resources/file-inflation-test/raw-10G-virt-10G.img",
		},
		{
			name:    "resize scratch and dest",
			fileURI: "gs://compute-image-tools-test-resources/file-inflation-test/raw-12G-virt-12G.img",
		},
	} {
		currentTest := tt
		t.Run(currentTest.name, func(t *testing.T) {
			t.Parallel()
			diskID := runDaisyInflate(t, currentTest.fileURI)
			assertDiskExists(t, diskID)
			deleteDisk(t, diskID)
		})
	}

}

func assertDiskExists(t *testing.T, diskID string) {
	client, err := daisycompute.NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	disk, err := client.GetDisk(project, zone, diskID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found disk: %v", disk)
}

func deleteDisk(t *testing.T, diskID string) {
	client, err := daisycompute.NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteDisk(project, zone, diskID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Deleted disk: %v", diskID)
}

func runDaisyInflate(t *testing.T, fileURI string) string {

	namespace := uuid.New().String()
	expectedDiskName := "disk-" + namespace

	ctx := context.Background()
	storageClient, err := storage.NewStorageClient(ctx, logging.NewDefaultLogger())
	if err != nil {
		t.Fatal(err)
	}

	sourceObj, err := importer.NewSourceFactory(storageClient).Init(fileURI, "")
	if err != nil {
		t.Fatal(err)
	}

	args := importer.ImportArguments{
		ExecutionID: namespace,
		WorkflowDir: workflowDir,
		Project:     project,
		Source:      sourceObj,
		Timeout:     time.Hour,
		Zone:        zone,
	}
	inflater, err := importer.CreateDaisyInflater(args, imagefile.NewGCSInspector())
	if err != nil {
		t.Fatal(err)
	}

	pd, inflateInfo, err := inflater.Inflate()
	t.Logf("Finished inflation: pd=%v, inflateInfo=%v, err=%v", pd, inflateInfo, err)
	if err != nil {
		t.Fatal(err)
	}

	return expectedDiskName
}
