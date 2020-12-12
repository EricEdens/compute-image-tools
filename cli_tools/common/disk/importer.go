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
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/logging"
	"google.golang.org/api/compute/v1"
)

type FileOrImageReference string
type DiskURI string
type ImageURI string

// Importer creates a PD from a disk file or a GCE disk image.
//
// To rebuild the mock, run `go generate ./...`
//go:generate go run github.com/golang/mock/mockgen -source $GOFILE -destination mocks/mock_importer.go
type Importer interface {
	// Import creates a PD from the disk file or GCE disk image
	// specified by source. Fields from destinationPrototype will
	// be included in the final PD. This is useful for specifying the
	//type of disk, licenses, and guest OS features.
	Import(source FileOrImageReference, destinationPrototype *compute.Disk) (DiskURI, error)
	Cancel(reason string) bool
}

type ImporterProvider interface {
	NewDataDiskImporter(environment Environment) Importer
	NewTranslatingDiskImporter(environment Environment, overrides TranslationOverrides) Importer
}

type TranslationOverrides struct {
	// BYOL determines which license is included in the final disk.
	BYOL bool

	// GcloudOSFlag overrides detection results. When empty,
	// OS detection results are used.
	GcloudOSFlag string

	// UEFI overrides detection results. When specified, UEFI_COMPATIBLE is
	// included in the GuestOSAttributes.
	UEFI bool

	// NoGuestEnvironment disables the installation of the guest
	// environment.
	NoGuestEnvironment bool
}

type Environment struct {
	Logger                                                                     logging.Logger
	Project, Zone, GCSPath, OAuth, Timeout, ComputeEndpoint, WorkflowDirectory string
	DisableGCSLogs, DisableCloudLogs, DisableStdoutLogs, NoExternalIP          bool
}
