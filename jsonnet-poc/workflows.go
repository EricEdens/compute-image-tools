package main

import (
	"encoding/json"
	"github.com/google/go-jsonnet"
	"io/ioutil"
	"log"
)

func main() {
	yesExternalIps()
	noExternalIps()

	moreLicenses()
}

func moreLicenses() {
	licenses := []string{
		"projects/compute-image-tools/global/licenses/virtual-disk-import",
		"custom/license/from/user",
	}
	vm := jsonnet.MakeVM()
	variable, _ := json.Marshal(licenses)
	vm.TLACode("additional_licenses", string(variable))
	renderAndPrint("templates/licenses.jsonnet", vm)
}

func yesExternalIps() {
	vm := jsonnet.MakeVM()
	vm.TLACode("use_external_ip", "true")
	vm.TLAVar("import_subnet", "projects/edens/sub2")
	renderAndPrint("templates/external-ips.jsonnet", vm)
}

func noExternalIps() {
	vm := jsonnet.MakeVM()
	vm.TLACode("use_external_ip", "false")
	renderAndPrint("templates/external-ips.jsonnet", vm)
}

func renderAndPrint(fname string, vm *jsonnet.VM) {
	template, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	rendered, err := vm.EvaluateSnippet("", string(template))
	if err != nil {
		log.Fatal(err)
	}
	println(rendered)
}
