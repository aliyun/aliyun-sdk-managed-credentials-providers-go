package sdk

import (
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
)

type SecretsMangerSdkCorePluginCredentialUpdater struct {
	client              interface{}
	accessKeyCredential auth.Credential
}

func NewSecretsMangerSdkCorePluginCredentialUpdater(client interface{}, accessKeyCredential auth.Credential) *SecretsMangerSdkCorePluginCredentialUpdater {
	return &SecretsMangerSdkCorePluginCredentialUpdater{
		client:              client,
		accessKeyCredential: accessKeyCredential,
	}
}

func (acu *SecretsMangerSdkCorePluginCredentialUpdater) GetClient() interface{} {
	return acu.client
}

func (acu *SecretsMangerSdkCorePluginCredentialUpdater) UpdateCredential(secretInfo *cmodels.SecretInfo) error {
	securityCredentials, err := service.GenerateCredentialsBySecret(secretInfo.SecretValue)
	if err != nil {
		return err
	}
	credential, _ := acu.accessKeyCredential.(*credentials.AccessKeyCredential)
	credential.AccessKeyId = securityCredentials.GetAccessKeyId()
	credential.AccessKeySecret = securityCredentials.GetAccessKeySecret()
	return nil
}

func (acu *SecretsMangerSdkCorePluginCredentialUpdater) GetTypeName() string {
	return reflect.TypeOf(*acu).Name()
}

func (acu *SecretsMangerSdkCorePluginCredentialUpdater) Close() error {
	return nil
}
