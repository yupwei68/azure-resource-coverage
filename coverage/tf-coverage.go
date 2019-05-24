package coverage

import (
	"fmt"
	"strings"

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
		nsMatch := strings.ToLower(client.GoSDKNamespace) == strings.ToLower(entry.NamespaceName)
		if client.GoSDKNamespace == "documentdb" && entry.NamespaceName == "cosmos-db" ||
			client.GoSDKNamespace == "appinsights" && entry.NamespaceName == "applicationinsights" ||
			client.GoSDKNamespace == "insights" && entry.NamespaceName == "monitor" ||
			client.GoSDKNamespace == "devices" && entry.NamespaceName == "iothub" ||
			client.GoSDKNamespace == "dtl" && entry.NamespaceName == "devtestlabs" ||
			client.GoSDKNamespace == "MsSql" && entry.NamespaceName == "sql" ||
			client.GoSDKNamespace == "storeAccount" && entry.NamespaceName == "datalake-store" ||
			client.GoSDKNamespace == "filesystem" && entry.NamespaceName == "datalake-store" ||
			client.GoSDKNamespace == "analyticsAccount" && entry.NamespaceName == "datalake-analytics" ||
			client.GoSDKNamespace == "media" && entry.NamespaceName == "mediaservices" ||
			client.GoSDKNamespace == "backup" && entry.NamespaceName == "recoveryservicesbackup" ||
			client.GoSDKNamespace == "locks" && entry.NamespaceName == "resources" ||
			client.GoSDKNamespace == "resourcesprofile" && entry.NamespaceName == "resources" ||
			client.GoSDKNamespace == "policy" && entry.NamespaceName == "resources" ||
			client.GoSDKNamespace == "subscriptions" && entry.NamespaceName == "subscription" {
			nsMatch = true
		}

		clientSDK := client.GoSDKClient[:len(client.GoSDKClient)-len("Client")]
		resMatch := strings.ToLower(clientSDK) == strings.ToLower(entry.ResourceName)
		if client.GoSDKNamespace == "automation" && clientSDK == "Account" && entry.ResourceName == "AutomationAccount" ||
			client.GoSDKNamespace == "redis" && clientSDK == "" && entry.ResourceName == "Redis" ||
			client.GoSDKNamespace == "apimanagement" && clientSDK == "Service" && entry.ResourceName == "ApiManagementService" ||
			client.GoSDKNamespace == "filesystem" && clientSDK == "" && entry.ResourceName == "FileSystem" ||
			client.GoSDKNamespace == "keyVault" && clientSDK == "Base" && entry.ResourceName == "<unknown>" ||
			client.GoSDKNamespace == "managementgroups" && clientSDK == "" && entry.ResourceName == "ManagementGroups" ||
			client.GoSDKNamespace == "managementgroups" && clientSDK == "Subscriptions" && entry.ResourceName == "ManagementGroupSubscriptions" ||
			client.GoSDKNamespace == "network" && clientSDK == "Interfaces" && entry.ResourceName == "NetworkInterfaces" ||
			client.GoSDKNamespace == "network" && clientSDK == "Profiles" && entry.ResourceName == "NetworkProfiles" ||
			client.GoSDKNamespace == "network" && clientSDK == "SecurityGroups" && entry.ResourceName == "NetworkSecurityGroups" ||
			client.GoSDKNamespace == "network" && clientSDK == "Watchers" && entry.ResourceName == "NetworkWatchers" ||
			client.GoSDKNamespace == "notificationhubs" && clientSDK == "" && entry.ResourceName == "NotificationHubs" ||
			client.GoSDKNamespace == "backup" && clientSDK == "ProtectedItemsGroup" && entry.ResourceName == "ProtectableItems" ||
			client.GoSDKNamespace == "resources" && clientSDK == "DeploymentsGroup" && entry.ResourceName == "Deployments" ||
			client.GoSDKNamespace == "resources" && clientSDK == "GroupsGroup" && entry.ResourceName == "ResourceGroups" ||
			client.GoSDKNamespace == "resources" && clientSDK == "Group" && entry.ResourceName == "Resources" ||
			client.GoSDKNamespace == "subscriptions" && clientSDK == "Group" && entry.ResourceName == "Subscriptions" ||
			client.GoSDKNamespace == "security" && clientSDK == "Contacts" && entry.ResourceName == "SecurityContacts" ||
			client.GoSDKNamespace == "signalr" && clientSDK == "" && entry.ResourceName == "SignalR" ||
			client.GoSDKNamespace == "storage" && clientSDK == "Accounts" && entry.ResourceName == "StorageAccounts" ||
			client.GoSDKNamespace == "web" && clientSDK == "Apps" && entry.ResourceName == "Sites" ||
			client.GoSDKNamespace == "policy" && clientSDK == "Assignments" && entry.ResourceName == "PolicyAssignments" && entry.ProviderName == "Microsoft.Authorization" ||
			client.GoSDKNamespace == "policy" && clientSDK == "Definitions" && entry.ResourceName == "PolicyDefinitions" && entry.ProviderName == "Microsoft.Authorization" ||
			client.GoSDKNamespace == "policy" && clientSDK == "SetDefinitions" && entry.ResourceName == "PolicySetDefinitions" && entry.ProviderName == "Microsoft.Authorization" {
			resMatch = true
		}
		if client.GoSDKNamespace == "network" && clientSDK == "PublicIPAddresses" && entry.ResourceName == "PublicIpAddresses" {
			resMatch = false
		}

		if nsMatch && resMatch {
			if found != nil {
				return nil, fmt.Errorf("Found more than one client %v in coverage", client)
			}
			found = entry
		}
	}
	if found == nil {
		return nil, fmt.Errorf("Cannot find client %v in coverage", client)
	}
	return found, nil
}
