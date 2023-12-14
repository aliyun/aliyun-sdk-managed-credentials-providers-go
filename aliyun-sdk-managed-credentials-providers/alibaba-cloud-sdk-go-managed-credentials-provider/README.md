# Managed Credentials Provider for Aliyun Go SDK

The Managed Credentials Provider for Aliyun Go SDK enables Golang developers to easily access to other Aliyun Services using managed RAM credentials stored in Aliyun Secrets Manager.

Read this in other languages: [English](README.md), [简体中文](README.zh-cn.md)

## Background
When applications use Aliyun SDK to call Alibaba Cloud Open APIs, access keys are traditionally used to authenticate to the cloud service. While access keys are easy to use they present security risks that could be leveraged by adversarial developers or external threats.

Alibaba Cloud SecretsManager is a solution that helps mitigate the risks by allowing organizations centrally manage access keys for all applications, allowing automatically or mannually rotating them without interrupting the applications. The managed access keys in SecretsManager is called [Managed RAM Credentials](https://www.alibabacloud.com/help/doc-detail/212421.htm).

For more advantages of using SecretsManager, refer to [SecretsManager Overview](https://www.alibabacloud.com/help/doc-detail/152001.htm).

## Client Mechanism
Applications use the access key that is managed by SecretsManager via the `Secret Name` representing the access key.

The Managed Credentials Provider periodically obtains the Access Key represented by the secret name and supply it to Aliyun SDK when accessing Alibaba Cloud Open APIs. The provider normally refreshes the locally cached access key at a specified interval, which is configurable.

## Install

If you use `go mod` to manage your dependence, You can declare the dependence in the go.mod file:

```
require (
	github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider v0.0.6
)
```

Or, Run the following command to get the remote code package:

```
$ go get -u github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider
```

## Aliyun SDK Managed Credentials Provider Sample

### Step 1: Configure the credentials provider

Use configuration file(`managed_credentials_providers.properties`)to access
KMS([Configuration file setting for details](../../README_config.md))，You could use the recommended way to access KMS with
Client Key.

```properties
 cache_client_dkms_config_info=[{"regionId":"<your dkms region>","endpoint":"<your dkms endpoint>","passwordFromFilePath":"< your password file path >","clientKeyFile":"<your client key file path>","ignoreSslCerts":false,"caFilePath":"<your CA certificate file path>"}]
```
```
    The details of the configuration item named cache_client_dkms_config_info:
    1. The configuration item named cache_client_dkms_config_info must be configured as a json array, you can configure multiple region instances
    2. regionId:Region id 
    3. endpoint:Domain address of dkms
    4. passwordFromFilePath and passwordFromEnvVariable
      passwordFromFilePath:The client key password configuration is obtained from the file,choose one of the two with passwordFromEnvVariable.
      e.g. while configuring passwordFromFilePath: < your password file path >, you need to configure a file with password written under the configured path
      passwordFromEnvVariable:The client key password configuration is obtained from the environment variable,choose one of the two with passwordFromFilePath.
      e.g. while configuring passwordFromEnvVariable: "your_password_env_variable",
           You need to add your_password_env_variable=< your client key private key password > in env.
    5. clientKeyFile:The path to the client key json file
    6. ignoreSslCerts:If ignore ssl certs (true: Ignores the ssl certificate, false: Validates the ssl certificate)
    7. caFilePath:The path of the CA certificate of the dkms
```

Note: By default, the Aliyun Managed Credentials Provider loads the configuration file from the current directory of the program.

### Step 2: Use the credentials provider in Aliyun SDK

You could use the following code to access Aliyun services with managed RAM credentials。

```go

package sample

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	sdkcoreprovider "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider/sdk"
	//"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
)

func main() {
	secretName := "********"
	regionId := "cn-hangzhou"

	//custom configuration
	//utils.SetConfigName("custom-config")
	client, err := sdkcoreprovider.GetClient(&ecs.Client{}, regionId, secretName)
	if err != nil {
		fmt.Println(err)
		return
	}
	ecsClient := client.(*ecs.Client)

	request := ecs.CreateDescribeInstancesRequest()
	instancesResponse, err := ecsClient.DescribeInstances(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, instance := range instancesResponse.Instances.Instance {
		// do something with instance
	}
}

```
