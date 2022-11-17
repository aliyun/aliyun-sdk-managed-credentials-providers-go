package sdk

import (
	"errors"
	"fmt"
	"sync"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/auth"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
)

var (
	aliyunSdkSecretsManagerPluginInstance *aliyunSdkSecretsManagerPlugin
	once                                  sync.Once
)

func Init() error {
	return InitAliyunSdkSecretsManagerPlugin(nil)
}

func InitAliyunSdkSecretsManagerPlugin(loader service.SecretsManagerPluginCredentialsLoader) (err error) {
	once.Do(func() {
		aliyunSdkSecretsManagerPluginInstance = newAliyunSdkSecretsManagerPlugin(loader)
		err = aliyunSdkSecretsManagerPluginInstance.init()
		if err != nil {
			return
		}
	})
	return
}

func GetAccessKey(secretName string) (*auth.SecretsManagerPluginCredentials, error) {
	client, err := GetAliyunSdkSecretsManagerPlugin()
	if err != nil {
		return nil, err
	}
	return client.getAccessKey(secretName)
}

func GetAccessKeyId(secretName string) (string, error) {
	securityCredentials, err := GetAccessKey(secretName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can not find access key id by the secretName[%s]", secretName))
	}
	return securityCredentials.AccessKeyId, nil
}

func GetAccessKeySecret(secretName string) (string, error) {
	securityCredentials, err := GetAccessKey(secretName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can not find access key secret by the secretName[%s]", secretName))
	}
	return securityCredentials.AccessKeySecret, nil
}

func RefreshSecretInfo(secretName string) error {
	client, err := GetAliyunSdkSecretsManagerPlugin()
	if err != nil {
		return err
	}
	err = client.refreshSecretInfo(secretName)
	if err != nil {
		return err
	}
	return nil
}

func RegisterSecretsManagerUpdater(secretName string, securityUpdater service.SecretsManagerPluginCredentialUpdater) error {
	securityCloudClient, err := GetAliyunSdkSecretsManagerPlugin()
	if err != nil {
		return err
	}
	securityCloudClient.registerSecretsManagerUpdater(secretName, securityUpdater)
	return nil
}

func GetSecretName(userSecretName string) (string, error) {
	securityCloudClient, err := GetAliyunSdkSecretsManagerPlugin()
	if err != nil {
		return "", err
	}
	return securityCloudClient.getSecretName(userSecretName)
}

func CloseSecretsManagerPluginUpdaterAndClient(secretName string, client interface{}) {
	securityCloudClient, _ := GetAliyunSdkSecretsManagerPlugin()
	if securityCloudClient != nil {
		securityCloudClient.closeSecurityUpdaterAndClientByClient(secretName, client)
	}
}

func CloseSecretsManagerPluginUpdaterAndClientByTypeName(updaterClasses map[string]struct{}) {
	securityCloudClient, _ := GetAliyunSdkSecretsManagerPlugin()
	if securityCloudClient != nil {
		securityCloudClient.closeSecurityUpdaterAndClientByTypeName(updaterClasses)
	}
}

func Shutdown() error {
	if aliyunSdkSecretsManagerPluginInstance == nil {
		return errors.New(fmt.Sprintf("Not initialize secrets manager plugin"))
	}
	aliyunSdkSecretsManagerPluginInstance.shutdown()
	return nil
}

func GetAliyunSdkSecretsManagerPlugin() (*aliyunSdkSecretsManagerPlugin, error) {
	if aliyunSdkSecretsManagerPluginInstance == nil {
		return nil, errors.New(fmt.Sprintf("Not initialize secrets manager plugin"))
	}
	return aliyunSdkSecretsManagerPluginInstance, nil
}
