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
	fmt.Println("Azure Resource Coverage Analyzer [v0.1.0]")

	if valid, apiPath, tfPath, configPath := parseArguments(); valid {
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

		cov, err := coverage.NewCoverage(configPath)
		if err != nil {
			fmt.Fprint(os.Stderr, "%+v", err)
			os.Exit(-3)
		}
		cov.LoadFromSpec(spec)
		if err := cov.AnalyzeTerraformCoverage(tf); err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-3)
		}
		cov.OutputCsv()
	} else {
		usage()
	}
}

func parseArguments() (valid bool, apiPath, tfPath, configPath string) {
	flag.Usage = usage
	flag.StringVar(&apiPath, "api-spec-path", "", "Specify the local root folder path of azure-rest-api-specs Github repository")
	flag.StringVar(&tfPath, "terraform-path", "", "Specify the local root folder path of terraform-provider-azurerm Github repository")
	flag.StringVar(&configPath, "config", "resource-config.json", "Specify the resource configuration path")
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
	fmt.Printf("  %s [-config <resource configuration json file>]\n", strings.Repeat(" ", len(exe)))
	fmt.Println()
	fmt.Println("Arguments:")
	flag.PrintDefaults()
}
