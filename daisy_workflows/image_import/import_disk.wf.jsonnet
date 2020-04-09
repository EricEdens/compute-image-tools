function(prefix, brunch=false)
{
  "Name": "import-disk",
  "DefaultTimeout": "90m",
  "Vars": {
    "source_disk_file": {
      "Required": true,
      "Description": "The GCS path to the virtual disk to import."
    },
    "importer_instance_disk_size": {
      "Value": "10",
      "Description": "size of the importer instance disk, additional disk space is unused for the import but a larger size increases PD write speed"
    },
    "import_instance_disk_image": {
      "Value": "projects/compute-image-tools/global/images/family/debian-9-worker",
      "Description": "image to use for the importer instance"
    },
    "disk_name": "imported-disk-${ID}",
    "import_network": {
      "Value": "global/networks/default",
      "Description": "Network to use for the import instance"
    },
    "import_subnet": {
      "Value": "",
      "Description": "SubNetwork to use for the import instance"
    },
    "is_windows": {
      "Value": "false",
      "Description": "If enabled, WINDOWS will be added to GuestOsFeatures for the disk."
    },
    "import_license": {
      "Value": "projects/compute-image-tools/global/licenses/virtual-disk-import",
      "Description": "Import License used for tracking migration workflow use."
    }
  },
  "Sources": {
    "import_image.sh": "./import_image.sh",
    "source_disk_file": "${source_disk_file}"
  },
  "Steps": {
    "setup-disks": {
      "CreateDisks": [
        {
          "Name": "disk-importer",
          "SourceImage": "${import_instance_disk_image}",
          "SizeGb": "${importer_instance_disk_size}",
          "Type": "pd-ssd",
          "FallbackToPdStandard": true
        },
        {
          "Name": "${disk_name}",
          "SizeGb": "10",
          "Type": "pd-ssd",
          "ExactName": true,
          "NoCleanup": true,
          "isWindows": "${is_windows}",
          "FallbackToPdStandard": true,
          "Licenses": ["${import_license}"]
        },
        {
          "Name": "disk-${NAME}-scratch-${ID}",
          "SizeGb": "10",
          "Type": "pd-ssd",
          "ExactName": true,
          "FallbackToPdStandard": true
        }
      ]
    },
    "import-virtual-disk": {
      "CreateInstances": [
        {
          "Name": "inst-importer",
          "Disks": [
            {"Source": "disk-importer"},
            {"Source": "disk-${NAME}-scratch-${ID}"},
            {"Source": "${disk_name}"}
          ],
          "MachineType": "n1-standard-4",
          "Metadata": {
            "block-project-ssh-keys": "true",
            "disk_name": "${disk_name}",
            "scratch_disk_name": "disk-${NAME}-scratch-${ID}",
            "source_disk_file": "${source_disk_file}",
            "shutdown-script": "echo 'Worker instance terminated'",
            "startup-script": "${SOURCE:import_image.sh}"
          },
          "networkInterfaces": [
            {
              "network": "${import_network}",
              "subnetwork": "${import_subnet}"
            }
          ],
          "Scopes": [
            "https://www.googleapis.com/auth/devstorage.read_write",
            "https://www.googleapis.com/auth/compute"
          ]
        }
      ]
    },
    "wait-for-signal": {
      "WaitForInstancesSignal": [
        {
          "Name": "inst-importer",
          "SerialOutput": {
            "Port": 1,
            "SuccessMatch": "ImportSuccess:",
            "FailureMatch": [
              "ImportFailed:",
              "WARNING Failed to download metadata script",
              "Worker instance terminated"
            ],
            "StatusMatch": "Import:"
          }
        }
      ]
    },
    "cleanup": {
      "DeleteResources": {
        "Instances":["inst-importer"]
      }
    }
  },
  "Dependencies": {
    "import-virtual-disk": ["setup-disks"],
    "wait-for-signal": ["import-virtual-disk"],
    "cleanup": ["wait-for-signal"]
  }
}
