package sdk

import (
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	commonsdk "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
)

var sdkCorePlugin *secretsManagerSdkCorePlugin
var once sync.Once

func GetClient(client interface{}, regionId string, secretName string) (interface{}, error) {
	return GetClientWithOptions(client, regionId, nil, secretName)
}
func GetClientWithOptions(client interface{}, regionId string, config *sdk.Config, secretName string) (interface{}, error) {
	err := initSecretsManagerPlugin()
	if err != nil {
		return nil, err
	}
	client, err = sdkCorePlugin.getClient(client, regionId, config, secretName)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CloseSDKCoreClient(client interface{}, secretName string) error {
	err := initSecretsManagerPlugin()
	if err != nil {
		return err
	}
	return sdkCorePlugin.closeSDKCoreClient(client, secretName)
}

func Destroy() {
	if sdkCorePlugin != nil {
		sdkCorePlugin.destroy()
	}
}

func initSecretsManagerPlugin() (err error) {
	err = commonsdk.Init()
	if err != nil {
		return
	}
	once.Do(func() {
		sdkCorePlugin = newSecretsManagerSdkCorePlugin()
	})
	return
}
