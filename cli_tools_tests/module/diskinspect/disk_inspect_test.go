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

package diskinspect

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/paramhelper"
	"github.com/GoogleCloudPlatform/compute-image-tools/common/runtime"
	"github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pb"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/disk"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/daisycommon"
	daisycompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
)

const (
	defaultNetwork = ""
	defaultSubnet  = ""
	workflowDir    = "../../../daisy_workflows"
)

var (
	zone = runtime.GetConfig("GOOGLE_CLOUD_ZONE", "compute/zone")
)

func TestInspectDisk(t *testing.T) {
	t.Parallel()

	project := runtime.GetConfig("GOOGLE_CLOUD_PROJECT", "project")

	client, err := daisycompute.NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		caseName string
		imageURI string
		expected *pb.InspectionResults
	}{
		{
			imageURI: "projects/opensuse-cloud/global/images/opensuse-leap-15-2-v20200702",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "opensuse-15",
					DistroId:     pb.Distro_OPENSUSE,
					Distro:       "opensuse",
					MajorVersion: "15",
					MinorVersion: "2",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
				BiosBootable: true,
			},
		},
		{
			imageURI: "projects/suse-sap-cloud/global/images/sles-15-sp1-sap-v20200803",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "sles-sap-15",
					DistroId:     pb.Distro_SLES_SAP,
					Distro:       "sles-sap",
					MajorVersion: "15",
					MinorVersion: "1",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
				BiosBootable: true,
			},
		},
		// UEFI
		{
			caseName: "UEFI inspection test for GPT UEFI",
			imageURI: "projects/gce-uefi-images/global/images/rhel-7-v20200403",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "rhel-7",
					DistroId:     pb.Distro_RHEL,
					Distro:       "rhel",
					MajorVersion: "7",
					MinorVersion: "8",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
			},
		},
		{
			caseName: "UEFI inspection test for MBR-only",
			imageURI: "projects/debian-cloud/global/images/debian-9-stretch-v20200714",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "debian-9",
					DistroId:     pb.Distro_DEBIAN,
					Distro:       "debian",
					MajorVersion: "9",
					MinorVersion: "12",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			caseName: "UEFI inspection test for GPT UEFI - windows",
			imageURI: "projects/gce-uefi-images/global/images/windows-server-2019-dc-core-v20200609",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "2019",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
			},
		},
		{
			caseName: "UEFI inspection test for GPT UEFI with BIOS boot",
			imageURI: "projects/gce-uefi-images/global/images/ubuntu-1804-bionic-v20200317",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "ubuntu-1804",
					DistroId:     pb.Distro_UBUNTU,
					Distro:       "ubuntu",
					MajorVersion: "18",
					MinorVersion: "04",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
				BiosBootable: true,
			},
		},
		{
			caseName: "UEFI inspection test for GPT UEFI with hybrid MBR",
			imageURI: "projects/compute-image-tools-test/global/images/image-ubuntu-2004-hybrid-mbr",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "ubuntu-2004",
					DistroId:     pb.Distro_UBUNTU,
					Distro:       "ubuntu",
					MajorVersion: "20",
					MinorVersion: "04",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
				BiosBootable: true,
			},
		},
		{
			caseName: "UEFI inspection test for MBR-only UEFI",
			imageURI: "projects/compute-image-tools-test/global/images/image-uefi-mbr-only",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					CliFormatted: "ubuntu-1604",
					DistroId:     pb.Distro_UBUNTU,
					Distro:       "ubuntu",
					MajorVersion: "16",
					MinorVersion: "04",
					Architecture: pb.Architecture_X64,
				},
				OsCount:      1,
				UefiBootable: true,
			},
		},

		// Windows Server
		{
			imageURI: "projects/windows-cloud/global/images/windows-server-2008-r2-dc-v20200114",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "2008",
					MinorVersion: "r2",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-2012-r2-vmware-import",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "2012",
					MinorVersion: "r2",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-2016-import",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "2016",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-2019",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "2019",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},

		// Windows Desktop
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-7-ent-x86-nodrivers",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "7",
					Architecture: pb.Architecture_X86,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-7-import",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "7",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-8-1-ent-x86-nodrivers",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "8",
					MinorVersion: "1",
					Architecture: pb.Architecture_X86,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-8-1-x64",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "8",
					MinorVersion: "1",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-10-1909-ent-x86-nodrivers",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "10",
					Architecture: pb.Architecture_X86,
				},
				OsCount: 1,
			},
		},
		{
			imageURI: "projects/compute-image-tools-test/global/images/windows-10-1709-import",
			expected: &pb.InspectionResults{
				OsRelease: &pb.OsRelease{
					DistroId:     pb.Distro_WINDOWS,
					Distro:       "windows",
					MajorVersion: "10",
					Architecture: pb.Architecture_X64,
				},
				OsCount: 1,
			},
		},
	} {
		// Without this, each parallel test will reference the last tt instance.
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		currentTest := tt
		name := currentTest.caseName
		if name == "" {
			name = currentTest.imageURI
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			inspector, err := disk.NewInspector(daisycommon.WorkflowAttributes{
				Project:           project,
				Zone:              zone,
				WorkflowDirectory: workflowDir,
			}, defaultNetwork, defaultSubnet)
			if err != nil {
				t.Fatal(err)
			}

			diskURI := createDisk(t, client, project, currentTest.imageURI)

			actual, err := inspector.Inspect(diskURI, true)
			assert.NoError(t, err)
			actual.ElapsedTimeMs = 0

			if diff := cmp.Diff(currentTest.expected, actual, protocmp.Transform()); diff != "" {
				t.Errorf("diff (-want +got):\n%s", diff)
			}
			deleteDisk(t, client, project, diskURI)
		})
	}
}

func TestInspectionWorksWithNonDefaultNetwork(t *testing.T) {
	t.Parallel()

	project := "compute-image-test-custom-vpc"
	client, err := daisycompute.NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	network := "projects/compute-image-test-custom-vpc/global/networks/unrestricted-egress"
	region, err := paramhelper.GetRegion(zone)
	if err != nil {
		t.Fatal(err)
	}
	subnet := fmt.Sprintf("projects/compute-image-test-custom-vpc/regions/%s/subnetworks/unrestricted-egress", region)

	for _, tt := range []struct {
		caseName string
		network  string
		subnet   string
	}{
		{"network and subnet", network, subnet},
		{"network only", network, ""},
		{"subnet only", "", subnet},
	} {
		currentTest := tt
		t.Run(currentTest.caseName, func(t *testing.T) {
			t.Parallel()
			t.Logf("Network=%s, Subnet=%s, project=%s", network, subnet, project)
			inspector, err := disk.NewInspector(daisycommon.WorkflowAttributes{
				Project:           "compute-image-test-custom-vpc",
				Zone:              zone,
				WorkflowDirectory: workflowDir,
			}, currentTest.network, currentTest.subnet)
			if err != nil {
				t.Fatal(err)
			}

			diskURI := createDisk(t, client, project, "projects/opensuse-cloud/global/images/opensuse-leap-15-2-v20200702")
			defer deleteDisk(t, client, project, diskURI)
			actual, err := inspector.Inspect(diskURI, true)
			if err != nil {
				t.Fatalf("Inspect failed: %v", err)
			}
			if actual.OsRelease.Distro != "opensuse" {
				t.Errorf("expected=opensuse, actual=%s", actual.OsRelease.Distro)
			}
		})
	}
}

func createDisk(t *testing.T, client daisycompute.Client, project, srcImage string) string {
	name := "d" + uuid.New().String()
	err := client.CreateDisk(project, zone, &compute.Disk{
		Name:        name,
		SourceImage: srcImage,
	})
	if err != nil {
		t.Fatal(err)
	}
	diskURI := fmt.Sprintf("projects/%s/zones/%s/disks/%s", project, zone, name)
	t.Logf("created disk: %s", diskURI)
	return diskURI
}

func deleteDisk(t *testing.T, client daisycompute.Client, project, diskURI string) {
	name := diskURI[strings.LastIndex(diskURI, "/")+1:]
	err := client.DeleteDisk(project, zone, name)
	if err != nil {
		t.Logf("Failed to delete disk: %v", err)
	} else {
		t.Logf("Deleted disk")
	}
}
