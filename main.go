package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/coverage"
	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func main() {
	fmt.Println("Azure Resource Coverage Analyzer [v0.0.2]")

	if valid, apiPath, tfPath := parseArguments(); valid {
		spec, err := apispec.LoadFrom(apiPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-1)
		}

		tf, err := tfprovider.LoadConfig(tfPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-2)
		}

		cov := coverage.ToCoverage(spec)
		if err := cov.AnalyzeTerraformCoverage(tf); err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-3)
		}
		fmt.Println("Namespace,Provider,Resource,Operations,Terraform Support")
		for _, entry := range cov {
			tfStatus := ""
			if entry.InTerraform {
				tfStatus = "yes"
			}
			ops := ""
			if entry.Resource.SupportCreate() {
				ops += "C"
			}
			if entry.Resource.SupportRead() {
				ops += "R"
			}
			if entry.Resource.SupportUpdate() {
				ops += "U"
			}
			if entry.Resource.SupportDelete() {
				ops += "D"
			}
			if entry.Resource.SupportList() {
				ops += "L"
			}
			fmt.Printf("%s,%s,%s,%s,%s\n", entry.NamespaceName, entry.ProviderName, entry.ResourceName, ops, tfStatus)
		}
	} else {
		usage()
	}
}

func parseArguments() (valid bool, apiPath string, tfPath string) {
	flag.Usage = usage
	flag.StringVar(&apiPath, "api-spec-path", "", "Specify the local root folder path of azure-rest-api-specs Github repository")
	flag.StringVar(&tfPath, "terraform-path", "", "Specify the local root folder path of terraform-provider-azurerm Github repository")
	flag.Parse()

	valid = true
	if apiPath == "" {
		fmt.Println("missing required flag: -api-spec-path")
		valid = false
	}
	if tfPath == "" {
		fmt.Println("missing required flag: -terraform-path")
		valid = false
	}

	return
}

func usage() {
	flag.CommandLine.SetOutput(os.Stdout)
	fmt.Println("Usage:")
	exe := filepath.Base(os.Args[0])
	fmt.Printf("  %s -api-spec-path <local path to azure-rest-api-specs>\n", exe)
	fmt.Printf("  %s -terraform-path <local path to terraform-provider-azurerm>\n", strings.Repeat(" ", len(exe)))
	fmt.Println()
	fmt.Println("Arguments:")
	flag.PrintDefaults()
}
