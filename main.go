package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/coverage"
	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func main() {
	fmt.Println("Azure Resource Coverage v0.0.1")

	if valid, apiPath, tfPath := parseArguments(); valid {
		spec, err := apispec.LoadFrom(apiPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-1)
		}

		tf, err := tfprovider.LoadConfig(tfPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(-1)
		}
		fmt.Println(tf.FullPath)

		cov := coverage.ToCoverage(spec)
		fmt.Printf("Total coverage entries: %d\n", len(cov))
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
	fmt.Println("Usage:")
	flag.PrintDefaults()
}
