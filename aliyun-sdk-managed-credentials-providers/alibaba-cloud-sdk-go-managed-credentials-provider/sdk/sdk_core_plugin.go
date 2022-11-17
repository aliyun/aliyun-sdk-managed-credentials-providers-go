package sdk

import (
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	commonmanager "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
)

var credentialUpdaterSet map[string]struct{}

const (
	InitWithOptionsFiledName = "InitWithOptions"
	EndpointMapFiledName     = "EndpointMap"
)

func init() {
	credentialUpdaterSet = make(map[string]struct{})
	credentialUpdaterSet[reflect.TypeOf(SecretsMangerSdkCorePluginCredentialUpdater{}).Name()] = struct{}{}
}

type secretsManagerSdkCorePlugin struct {
}

func newSecretsManagerSdkCorePlugin() *secretsManagerSdkCorePlugin {
	return &secretsManagerSdkCorePlugin{}
}

func (sac *secretsManagerSdkCorePlugin) getClient(client interface{}, regionId string, config *sdk.Config, secretName string) (interface{}, error) {
	accessKey, err := commonmanager.GetAccessKey(secretName)
	if err != nil {
		return nil, err
	}
	accessKeyCredential := credentials.NewAccessKeyCredential(accessKey.AccessKeyId, accessKey.AccessKeySecret)
	if config == nil {
		config = sdk.NewConfig()
	}
	_, err = utils.Call(client, InitWithOptionsFiledName, regionId, config, accessKeyCredential)
	if err != nil {
		return nil, err
	}
	err = utils.SetUnExportedField(client, EndpointMapFiledName, GetEndpointMap())
	if err != nil {
		return nil, err
	}
	err = utils.SetUnExportedField(client, "EndpointType", GetEndpointType())
	if err != nil {
		return nil, err
	}
	securityCredentialUpdater := NewSecretsMangerSdkCorePluginCredentialUpdater(client, accessKeyCredential)
	err = commonmanager.RegisterSecretsManagerUpdater(secretName, securityCredentialUpdater)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (sac *secretsManagerSdkCorePlugin) closeSDKCoreClient(client interface{}, secretName string) error {
	commonmanager.CloseSecretsManagerPluginUpdaterAndClient(secretName, client)
	return nil
}

func (sac *secretsManagerSdkCorePlugin) destroy() {
	commonmanager.CloseSecretsManagerPluginUpdaterAndClientByTypeName(credentialUpdaterSet)
}
