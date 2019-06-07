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

		nsMatch := strings.ToLower(client.GoSDKNamespace) == strings.ToLower(entry.Namespace.Name)
		if client.GoSDKNamespace == "documentdb" && entry.Namespace.Name == "cosmos-db" ||
			client.GoSDKNamespace == "appinsights" && entry.Namespace.Name == "applicationinsights" ||
			client.GoSDKNamespace == "insights" && entry.Namespace.Name == "monitor" ||
			client.GoSDKNamespace == "devices" && entry.Namespace.Name == "iothub" ||
			client.GoSDKNamespace == "dtl" && entry.Namespace.Name == "devtestlabs" ||
			client.GoSDKNamespace == "MsSql" && entry.Namespace.Name == "sql" ||
			client.GoSDKNamespace == "storeAccount" && entry.Namespace.Name == "datalake-store" ||
			client.GoSDKNamespace == "filesystem" && entry.Namespace.Name == "datalake-store" ||
			client.GoSDKNamespace == "analyticsAccount" && entry.Namespace.Name == "datalake-analytics" ||
			client.GoSDKNamespace == "media" && entry.Namespace.Name == "mediaservices" ||
			client.GoSDKNamespace == "backup" && entry.Namespace.Name == "recoveryservicesbackup" ||
			client.GoSDKNamespace == "locks" && entry.Namespace.Name == "resources" ||
			client.GoSDKNamespace == "resourcesprofile" && entry.Namespace.Name == "resources" ||
			client.GoSDKNamespace == "policy" && entry.Namespace.Name == "resources" ||
			client.GoSDKNamespace == "subscriptions" && entry.Namespace.Name == "subscription" {
			nsMatch = true
		}

		clientSDK := client.GoSDKClient[:len(client.GoSDKClient)-len("Client")]
		resMatch := strings.ToLower(clientSDK) == strings.ToLower(entry.ResourceName)
		if client.GoSDKNamespace == "automation" && clientSDK == "Account" && entry.ResourceName == "AutomationAccount" ||
			client.GoSDKNamespace == "redis" && clientSDK == "" && entry.ResourceName == "Redis" ||
			client.GoSDKNamespace == "apimanagement" && clientSDK == "Service" && entry.ResourceName == "ApiManagementService" ||
			client.GoSDKNamespace == "filesystem" && clientSDK == "" && entry.ResourceName == "FileSystem" ||
			client.GoSDKNamespace == "keyVault" && clientSDK == "Base" && entry.ResourceName == "Secrets" ||
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
			client.GoSDKNamespace == "policy" && clientSDK == "SetDefinitions" && entry.ResourceName == "PolicySetDefinitions" && entry.ProviderName == "Microsoft.Authorization" ||
			client.GoSDKNamespace == "batch" && clientSDK == "Account" && entry.ResourceName == "BatchAccount" {
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
