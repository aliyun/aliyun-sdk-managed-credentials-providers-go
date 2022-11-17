package sdk

import (
	"reflect"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
)

type SecretsMangerOssPluginCredentialUpdater struct {
	ossClient *ProxyOssClient
	provider  *OssPluginCredentialsProvider
}

func NewSecretsMangerOssPluginCredentialUpdater(ossClient *ProxyOssClient, provider *OssPluginCredentialsProvider) *SecretsMangerOssPluginCredentialUpdater {
	return &SecretsMangerOssPluginCredentialUpdater{
		ossClient: ossClient,
		provider:  provider,
	}
}

func (ocu *SecretsMangerOssPluginCredentialUpdater) GetClient() interface{} {
	return ocu.ossClient
}

func (ocu *SecretsMangerOssPluginCredentialUpdater) UpdateCredential(secretInfo *cmodels.SecretInfo) error {
	credential, err := service.GenerateCredentialsBySecret(secretInfo.SecretValue)
	if err != nil {
		return err
	}
	ocu.provider.SetCredentials(NewOssPluginCredentials(credential.GetAccessKeyId(), credential.GetAccessKeySecret()))
	return nil
}

func (ocu *SecretsMangerOssPluginCredentialUpdater) GetTypeName() string {
	return reflect.TypeOf(*ocu).Name()
}

func (ocu *SecretsMangerOssPluginCredentialUpdater) Close() error {
	return nil
}
