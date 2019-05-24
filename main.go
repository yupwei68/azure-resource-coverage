package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

func main() {
	fmt.Println("Azure Resource Coverage v0.0.1")

	if valid, apiPath := parseArguments(); valid {
		spec, err := apispec.LoadFrom(apiPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-1)
		}

		totalCnt := 0
		mngCnt := 0
		ctlCnt := 0
		resCnt := 0
		for _, ns := range spec.Namespaces {
			if ns.Type == apispec.Management {
				mngCnt++
			}
			if ns.Type == apispec.ControlPlane {
				ctlCnt++
			}
			totalCnt++
			fmt.Printf("%s\n", ns.RelativePath)
			for pname, p := range ns.Providers {
				fmt.Printf("\t%q (%s)\n", pname, p.RelativePath)
				for name, r := range p.Resources {
					ops := ""
					if r.SupportCreate() {
						ops += "C"
					}
					if r.SupportRead() {
						ops += "R"
					}
					if r.SupportUpdate() {
						ops += "U"
					}
					if r.SupportDelete() {
						ops += "D"
					}
					if r.SupportList() {
						ops += "L"
					}
					fmt.Printf("\t\t%q [%s] {%s}\n", name, ops, strings.Join(r.AdditionalOperations(), ","))
					resCnt++
				}
			}
		}
		fmt.Printf("Total %d namespaces, %d management namespaces, %d control namespaces, %d resources\n", totalCnt, mngCnt, ctlCnt, resCnt)
	} else {
		usage()
	}
}

func parseArguments() (valid bool, apiPath string) {
	flag.Usage = usage
	flag.StringVar(&apiPath, "api-spec-path", "", "Specify the root folder of azure-rest-api-specs Github repository")
	flag.Parse()

	valid = true
	if apiPath == "" {
		fmt.Println("missing required flag: -api-spec-path")
		valid = false
	}

	return
}

func usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	children, _ := ioutil.ReadDir("D:\\Source\\Azure\\azure-rest-api-specs\\specification\\frontdoor\\resource-manager\\Microsoft.Network\\preview")
	for _, f := range children {
		fmt.Println(f.Name())
	}
}
