# Aliyun SDK Managed Credentials Providers for Go

The Aliyun SDK Managed Credentials Providers for Go enables golang developers to easily access to other Aliyun Services
using managed RAM credentials stored in Aliyun Secrets Manager.

Read this in other languages: [English](README.md), [简体中文](README.zh-cn.md)

- [Aliyun Secrets Manager Managed RAM Credentials Summary](https://www.alibabacloud.com/help/doc-detail/152001.htm)
- [Issues](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/issues)
- [Release](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/releases)

## License

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)

## Features

* Provide an easy method to other AliCould Services using managed RAM credentials
* Provides multiple access ways such as ECS instance RAM Role or Client Key to obtain a managed RAM credentials
* Provides the Aliyun Service client to refresh the RAM credentials automatically

## Requirements

- Your secret must be a managed RAM credentials
- Golang 1.13 or later

## Background

When applications use Aliyun SDK to call Alibaba Cloud Open APIs, access keys are traditionally used to authenticate to
the cloud service. While access keys are easy to use they present security risks that could be leveraged by adversarial
developers or external threats.

Alibaba Cloud SecretsManager is a solution that helps mitigate the risks by allowing organizations centrally manage
access keys for all applications, allowing automatically or mannually rotating them without interrupting the
applications. The managed access keys in SecretsManager is
called [Managed RAM Credentials](https://www.alibabacloud.com/help/doc-detail/212421.htm).

For more advantages of using SecretsManager, refer
to [SecretsManager Overview](https://www.alibabacloud.com/help/doc-detail/152001.htm).

## Client Mechanism

Applications use the access key that is managed by SecretsManager via the 'Secret Name' representing the access key.

The Managed Credentials Provider periodically obtains the Access Key represented by the secret name and supply it to
Aliyun SDK when accessing Alibaba Cloud Open APIs. The provider normally refreshes the locally cached access key at a
specified interval, which is configurable.

However, there are circumstances that the cached access key is no longer valid, which typically happens when emergent
access key rotation is performed by adminstrators in SecretsManager to respond to an leakage incident. Using invalid
access key to call Open APIs usually results in an exception that corresponds to an API error code. The Managed
Credentials Provider will immediately refresh the cached access key and retry the failed Open API if the corresponding
error code is `InvalidAccessKeyId.NotFound` or `InvalidAccessKeyId`.

Application developers can override or extend this behavior for specific cloud services if the APIs return other error
codes for using expired access keys. Refer
to [Modifying the default expire handler](#modifying-the-default-expire-handler).

## Install

If you use `go mod` to manage your dependence, You can declare the dependence in the go.mod file. Take an example to 
import Aliyun OSS SDK Managed Credentials Provider:

```
require (
	github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider v0.0.6
)
```

Or, Run the following command to get the remote code package:

```
$ go get -u github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider
```

## Support Aliyun Services

The Aliyun SDK Managed Credentials Providers for Go supports the following Aliyun Services:

|                           Aliyun SDK Name                           |                                                                                      Plugin Name                                                                                       |
|:-------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
| [Alibaba Cloud SDK](https://github.com/aliyun/alibaba-cloud-sdk-go) | [Managed Credentials Provider for Aliyun Go SDK](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/tree/master/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider)  |
|       [OSS SDK](https://github.com/aliyun/aliyun-oss-go-sdk)        | [Managed Credentials Provider for Aliyun Go OSS SDK](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/tree/master/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider) |

## Aliyun Managed Credentials Providers Sample

### Step 1: Configure the credentials provider

Use configuration file(`managed_credentials_providers.properties`)to access
KMS([Configuration file setting for details](README_config.md))，You could use the recommended way to access KMS with
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

You cloud use the following code to access Aliyun services with managed RAM credentials(taking an example to access OSS)
.

```go
package sample

import (
	"fmt"
	ossprovider "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider/sdk"
	//"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
)

func main() {
	secretName := "********"
	endpoint := "https://oss-cn-hangzhou.aliyuncs.com"
	//custom configuration
	//utils.SetConfigName("custom-config")
	client, err := ossprovider.New(endpoint, secretName)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.ListBuckets()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range result.Buckets {
		// do something with bucket
	}
	
	client.Shutdown()
}

```

## Modifying the default expire handler

With Aliyun SDK Managed Credentials Provider that supports customed error retry, you can customize the error retry 
judgment of the client due to manual rotation of credentials in extreme scenarios, you only implement the following 
interface.

```go

type AKExpireHandler interface {
	// judge whether the exception is caused by AccessKey expiration
	JudgeAKExpire(err error) bool
}

```

The sample codes below show customed judgment exception interface and use it to call aliyun services.

```go
package sample

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	ossprovider "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider/sdk"
	//"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
)

type MyOssAKExpireHandler struct {
	Code string
}

func (handler *MyOssAKExpireHandler) JudgeAKExpire(err error) bool {
	if e, ok := err.(oss.ServiceError); ok {
		if e.Code == handler.Code {
			return true
		}
	}
	return false
}

const AkExpireErrorCode = "InvalidAccessKeyId"

func main() {
	secretName := "********"
	endpoint := "https://oss-cn-hangzhou.aliyuncs.com"
	
	akExpireHandler := &MyOssAKExpireHandler{Code: AkExpireErrorCode}
	//custom configuration
	//utils.SetConfigName("custom-config")
	client, err := ossprovider.NewProxyOssClientWithHandler(endpoint, secretName, akExpireHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.ListBuckets()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range result.Buckets {
		// do something with bucket
	}
	
	client.Shutdown()
}

```
