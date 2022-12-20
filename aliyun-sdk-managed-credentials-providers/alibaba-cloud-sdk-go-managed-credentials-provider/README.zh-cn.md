# 阿里云Go SDK托管凭据插件

阿里云Go SDK的托管凭据插件可以帮助Golang开发者更方便的利用托管在凭据管家的RAM凭据，来访问阿里云服务的开放API。

*其他语言版本: [English](README.md), [简体中文](README.zh-cn.md)*

## 背景
当您的应用程序通过阿里云SDK访问云服务时，访问凭证(Access Keys)被用于认证应用的身份。访问凭据在使用中存在一定的安全风险，可能会被恶意的开发人员或外部威胁所利用。

阿里云凭据管家提供了帮助降低风险的解决方案，允许企业和组织集中管理所有应用程序的访问凭据，允许在不中断应用程序的情况下自动或手动轮转或者更新这些凭据。托管在凭据管家的Access Key被称为[托管RAM凭据](https://help.aliyun.com/document_detail/212421.html) 。

使用凭据管家的更多优势，请参阅 [凭据管家概述](https://help.aliyun.com/document_detail/152001.html) 。

## 客户端机制
应用程序引用托管RAM凭据（Access Key）的`凭据名称` 。

托管凭据插件定期从凭据管家获取由`凭据名称`代表的Access Key，并提供给阿里云SDK，应用则使用SDK访问阿里云开放API。插件以指定的间隔（可配置）刷新缓存在内存中的Access Key。

## 安装

您可以使用`go mod`管理您的依赖:

```
require (
	github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider v0.0.2
)
```

或者，通过`go get`命令安装远程代码包:

```
$ go get -u github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider
```

## 使用示例

### 步骤1：配置托管凭据插件

通过配置文件(`managed_credentials_providers.properties`)指定访问凭据管家([配置文件设置详情](../../README_config.zh-cn.md))，推荐采用Client Key方式访问凭据管家
- 访问DKMS, 您必须设置以下配置:

```properties
 cache_client_dkms_config_info=[{"regionId":"<your dkms region>","endpoint":"<your dkms endpoint>","passwordFromFilePath":"< your password file path >","clientKeyFile":"<your client key file path>","ignoreSslCerts":false,"caFilePath":"<your CA certificate file path>"}]
```
```
    cache_client_dkms_config_info配置项说明:
    1. cache_client_dkms_config_info配置项为json数组，支持配置多个region实例
    2. regionId:地域Id
    3. endpoint:专属kms的域名地址
    4. passwordFromFilePath和passwordFromEnvVariable
       passwordFromFilePath:client key密码配置从文件中获取，与passwordFromEnvVariable二选一
       例:当配置passwordFromFilePath:<你的client key密码文件所在的路径>,需在配置的路径下配置写有password的文件
       passwordFromEnvVariable:client key密码配置从环境变量中获取，与passwordFromFilePath二选一
       例:当配置"passwordFromEnvVariable":"your_password_env_variable"时，
         需在环境变量中添加your_password_env_variable=<你的client key对应的密码>
    5. clientKeyFile:client key json文件的路径
    6. ignoreSslCerts:是否忽略ssl证书 (true:忽略ssl证书,false:验证ssl证书)
    7. caFilePath:专属kms的CA证书路径
```

说明：插件默认从程序当前目录加载配置文件。

### 步骤 2：使用托管凭据插件访问云服务

您可以通过以下代码通过凭据管家托管RAM凭据使用阿里云SDK客户端。

```go

package sample

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	sdkcoreprovider "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider/sdk"
)

func main() {
	secretName := "********"
	regionId := "cn-hangzhou"

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
