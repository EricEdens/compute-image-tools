This is a POC of using [jsonnet-go](https://github.com/google/go-jsonnet) to
render Daisy workflows.

## Demo

`$ go run workflows.go`

```json
{
   "Steps": {
      "import_disk": {
         "CreateInstances": [
            {
               "networkInterfaces": [
                  {
                     "network": "global/networks/default",
                     "subnetwork": "projects/edens/sub2"
                  }
               ]
            }
         ]
      }
   }
}
```
```json
{
   "Steps": {
      "import_disk": {
         "CreateInstances": [
            {
               "networkInterfaces": [ ]
            }
         ]
      }
   }
}
```
```json
{
   "Steps": {
      "create_image": {
         "CreateImages": [
            {
               "Licenses": [
                  "projects/compute-image-tools/global/licenses/virtual-disk-import",
                  "custom/license/from/user",
                  "projects/ubuntu-os-cloud/global/licenses/ubuntu-1604-xenial"
               ]
            }
         ]
      }
   }
}
```

## Template files

The examples use jsonnet's function style. Go  passes
arguments to the template engine, and the return value
from the function is emitted as JSON.

For more info, see [top-level arguments](https://jsonnet.org/learning/tutorial.html#parameterize-entire-config)
in their docs.

### external-ips.jsonnet

This example illustrates:
* Required and default arguments. Evaluation will fail if 
  `use_external_ip` is not initialized. The other args
  have default values and are optional.
* Conditional expressions, where the networkInterfaces
  block is rendered only if `use_external_ip` is true.

```javascript
function(
  use_external_ip,
  import_network='global/networks/default',
  import_subnet='',
) {
  Steps: {
    import_disk: {
      CreateInstances: [
        {
          networkInterfaces: if use_external_ip then [
            {
              network: import_network,
              subnetwork: import_subnet,
            },
          ] else [],
        },
      ],
    },
  },
}
```


### licenses.jsonnet

This example illustrates:
* String interpolation to create the license for Ubuntu 16.04.
* Combining two lists of licenses via std.setUnion (the default
  license plus whatever is passed at runtime.)

```javascript
function(
  ubuntu_version='1604-xenial',
  additional_licenses=[],
) {
  Steps: {
    create_image: {
      CreateImages: [
        {
          Licenses: std.setUnion(
            ['projects/ubuntu-os-cloud/global/licenses/ubuntu-%s' % ubuntu_version],
            additional_licenses,
          ),
        },
      ],
    },
  },
}

```

## Looking forward

The OS-detection will split the currently-monolithic import workflow
into three smaller workflows: inflation, detection, and translation.

To implement this with jsonnet, we'll have:

### inflation

When inflating disk files, we'd use a jsonnet template
based on `import_disk.wf.json`. The daisy variables
will become jsonnet function parameters.

When inflating from GCP images, we can make a direct API
call to make a PD from the image.

### detection

This will be a new workflow that takes a reference to a PD,
attaches it to an instance, and runs a startup-script containing
the logic of OS detection.

### translation

We'll need have one translation workflow per supported OS version,
but as seen in licenses.json, we can consolidate the existing
workflow files to *at least* one jsonnet file per distro, and perhaps
one jsonnet file for all of Linux.