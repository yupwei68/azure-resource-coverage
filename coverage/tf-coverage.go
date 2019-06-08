package coverage

import (
	"fmt"
	"strings"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func (cov ResourceCoverage) AnalyzeTerraformCoverage(tf *tfprovider.TerraformConfig) error {
	for _, client := range tf.Clients {
		entry, err := cov.findExactOneEntry(client)
		if err != nil {
			return err
		}
		entry.InTerraform = true
	}
	return nil
}

func (cov ResourceCoverage) findExactOneEntry(client *tfprovider.ReferencedClient) (*CoverageEntry, error) {
	var found *CoverageEntry = nil
	for _, entry := range cov {
		if entry.Namespace.Type == apispec.DataPlane && entry.Namespace.Name != "graphrbac" && entry.Namespace.Name != "datalake-store" {
			continue
		}

		if entry.Namespace.Name == "apimanagement" && entry.Namespace.Type == apispec.ControlPlane {
			continue
		}

		nsMatch := strings.ToLower(client.Package.BaseName()) == strings.ToLower(entry.Namespace.Name)
		if client.Package.BaseName() == "documentdb" && entry.Namespace.Name == "cosmos-db" ||
			client.Package.BaseName() == "insights" && entry.Namespace.Name == "applicationinsights" ||
			client.Package.BaseName() == "insights" && entry.Namespace.Name == "monitor" ||
			client.Package.BaseName() == "devices" && entry.Namespace.Name == "iothub" ||
			client.Package.BaseName() == "dtl" && entry.Namespace.Name == "devtestlabs" ||
			client.Package.BaseName() == "MsSql" && entry.Namespace.Name == "sql" ||
			client.Package.BaseName() == "account" && entry.Namespace.Name == "datalake-store" ||
			client.Package.BaseName() == "filesystem" && entry.Namespace.Name == "datalake-store" ||
			client.Package.BaseName() == "analyticsAccount" && entry.Namespace.Name == "datalake-analytics" ||
			client.Package.BaseName() == "media" && entry.Namespace.Name == "mediaservices" ||
			client.Package.BaseName() == "backup" && entry.Namespace.Name == "recoveryservicesbackup" ||
			client.Package.BaseName() == "locks" && entry.Namespace.Name == "resources" ||
			client.Package.BaseName() == "resourcesprofile" && entry.Namespace.Name == "resources" ||
			client.Package.BaseName() == "policy" && entry.Namespace.Name == "resources" ||
			client.Package.BaseName() == "subscriptions" && entry.Namespace.Name == "subscription" {
			nsMatch = true
		}

		clientSDK := client.GoSDKClient[:len(client.GoSDKClient)-len("Client")]
		resMatch := strings.ToLower(clientSDK) == strings.ToLower(entry.ResourceName)
		if client.Package.BaseName() == "automation" && clientSDK == "Account" && entry.ResourceName == "AutomationAccount" ||
			client.Package.BaseName() == "redis" && clientSDK == "" && entry.ResourceName == "Redis" ||
			client.Package.BaseName() == "apimanagement" && clientSDK == "Service" && entry.ResourceName == "ApiManagementService" ||
			client.Package.BaseName() == "filesystem" && clientSDK == "" && entry.ResourceName == "FileSystem" ||
			client.Package.BaseName() == "keyVault" && clientSDK == "Base" && entry.ResourceName == "Secrets" ||
			client.Package.BaseName() == "managementgroups" && clientSDK == "" && entry.ResourceName == "ManagementGroups" ||
			client.Package.BaseName() == "managementgroups" && clientSDK == "Subscriptions" && entry.ResourceName == "ManagementGroupSubscriptions" ||
			client.Package.BaseName() == "network" && clientSDK == "Interfaces" && entry.ResourceName == "NetworkInterfaces" ||
			client.Package.BaseName() == "network" && clientSDK == "Profiles" && entry.ResourceName == "NetworkProfiles" ||
			client.Package.BaseName() == "network" && clientSDK == "SecurityGroups" && entry.ResourceName == "NetworkSecurityGroups" ||
			client.Package.BaseName() == "network" && clientSDK == "Watchers" && entry.ResourceName == "NetworkWatchers" ||
			client.Package.BaseName() == "notificationhubs" && clientSDK == "" && entry.ResourceName == "NotificationHubs" ||
			client.Package.BaseName() == "backup" && clientSDK == "ProtectedItemsGroup" && entry.ResourceName == "ProtectableItems" ||
			client.Package.BaseName() == "resources" && clientSDK == "DeploymentsGroup" && entry.ResourceName == "Deployments" ||
			client.Package.BaseName() == "resources" && clientSDK == "GroupsGroup" && entry.ResourceName == "ResourceGroups" ||
			client.Package.BaseName() == "resources" && clientSDK == "Group" && entry.ResourceName == "Resources" ||
			client.Package.BaseName() == "subscriptions" && clientSDK == "Group" && entry.ResourceName == "Subscriptions" ||
			client.Package.BaseName() == "security" && clientSDK == "Contacts" && entry.ResourceName == "SecurityContacts" ||
			client.Package.BaseName() == "signalr" && clientSDK == "" && entry.ResourceName == "SignalR" ||
			client.Package.BaseName() == "storage" && clientSDK == "Accounts" && entry.ResourceName == "StorageAccounts" ||
			client.Package.BaseName() == "web" && clientSDK == "Apps" && entry.ResourceName == "Sites" ||
			client.Package.BaseName() == "policy" && clientSDK == "Assignments" && entry.ResourceName == "PolicyAssignments" && entry.ProviderName == "Microsoft.Authorization" ||
			client.Package.BaseName() == "policy" && clientSDK == "Definitions" && entry.ResourceName == "PolicyDefinitions" && entry.ProviderName == "Microsoft.Authorization" ||
			client.Package.BaseName() == "policy" && clientSDK == "SetDefinitions" && entry.ResourceName == "PolicySetDefinitions" && entry.ProviderName == "Microsoft.Authorization" ||
			client.Package.BaseName() == "batch" && clientSDK == "Account" && entry.ResourceName == "BatchAccount" {
			resMatch = true
		}
		if client.Package.BaseName() == "network" && clientSDK == "PublicIPAddresses" && entry.ResourceName == "PublicIpAddresses" {
			resMatch = false
		}

		if nsMatch && resMatch {
			if found != nil {
				return nil, fmt.Errorf("Found more than one client (%v).%s in coverage", client.Package, client.GoSDKClient)
			}
			found = entry
		}
	}
	if found == nil {
		return nil, fmt.Errorf("Cannot find client (%v).%s in coverage", client.Package, client.GoSDKClient)
	}
	return found, nil
}
