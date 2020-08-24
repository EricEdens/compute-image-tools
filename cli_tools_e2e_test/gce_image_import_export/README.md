## Build and Execution

The suite is built by Cloud Build using the `/prowjobs_cloudbuild.yaml`
configuration, and is run by Prow. The Prow job's configuration
is in `/test-infra/prow/config.yaml`.

For more details, including UI links, see [test-infra/README.md](../../test-infra/README.md).


## Ad-hoc test running

To run the tests against your locally-checked out codebase, customize
**test_case_filter** and **test_suite_filter** and execute *from the root of 
the compute-image-tools repo*:

```yaml
yaml=$(mktemp)
cat <<EOF > $yaml
timeout: 7200s
steps:
  - name: gcr.io/cloud-builders/docker
    args: [
      'build',
      '--file=gce_image_import_export_tests.Dockerfile',
      '--tag=gce_image_import_export_tests',
      '.'
    ]
  - name: gce_image_import_export_tests
    dir: /
    args: [
      '-test_case_filter=ubuntu',
      '-test_suite_filter=ImageImport',
      '-test_project_id=$(gcloud config get-value project)',
      '-test_zone=$(gcloud config get-value compute/zone)'
    ]
EOF
gcloud builds submit --config $yaml
```

There are two Cloud Build steps. The first builds the local code base and
creates an image from the results. The creates a container from the image and
runs all tests with a name containing "ubuntu" in the ImageImport suite.
Project and zone are retrieved from the gcloud's default config values.

Execute uses your project's [cloud build service account](https://cloud.google.com/cloud-build/docs/securing-builds/configure-access-for-cloud-build-service-account).
