module github.com/GoogleCloudPlatform/compute-image-tools/cli_tools_e2e_test

go 1.13

require (
	cloud.google.com/go/storage v1.10.0
	github.com/GoogleCloudPlatform/compute-image-tools/cli_tools v0.0.0-20200813223603-3672ca27e050
	github.com/GoogleCloudPlatform/compute-image-tools/cli_tools_e2e_test/common v0.0.0
	github.com/GoogleCloudPlatform/compute-image-tools/daisy v0.0.0-20200813213118-5b8ac20eab97
	github.com/GoogleCloudPlatform/compute-image-tools/go/e2e_test_utils v0.0.0-20200813223603-3672ca27e050
	github.com/aws/aws-sdk-go v1.34.4
	google.golang.org/api v0.30.0
)

replace github.com/GoogleCloudPlatform/compute-image-tools/cli_tools_e2e_test/common => ./common
