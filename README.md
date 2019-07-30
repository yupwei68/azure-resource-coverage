# Azure Resource Coverage

A console tool to calculate the coverage rate of Azure resources in OSS tools.

## Usage

This tool will calculate the Azure resource coverage rate based on the source code of [Azure Rest API Spec](https://github.com/Azure/azure-rest-api-specs) and [Terraform provider for Azure](https://github.com/terraform-providers/terraform-provider-azurerm). So please clone these two repositories to your local filesystem:

* `git clone https://github.com/Azure/azure-rest-api-specs.git`
* `git clone https://github.com/terraform-providers/terraform-provider-azurerm.git`

Next please download and unzip the latest binary of this tool in [`Releases`](https://github.com/JunyiYi/azure-resource-coverage/releases) tab for your system. Finally, run the following command line to generate the result:

```sh
./azure-resource-coverage -api-spec-path "<local full path to azure-rest-api-specs>" -terraform-path "<local full path to terraform-provider-azurerm>" > coverage.csv
```

And starting from the second line of `coverage.csv`, we have the coverage result.
