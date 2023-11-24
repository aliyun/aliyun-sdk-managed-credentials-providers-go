# 阿里云SDK托管凭据Go插件

阿里云SDK托管凭据Go插件可以使Golang开发者通过托管RAM凭据快速使用阿里云服务。

*其他语言版本: [English](README.md), [简体中文](README.zh-cn.md)*

- [阿里云托管RAM凭据主页](https://help.aliyun.com/document_detail/212421.html)
- [Issues](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/issues)
- [Release](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/releases)

## 许可证

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)


## 优势
* 支持用户快速通过托管RAM凭据快速使用阿里云服务
* 支持多种认证鉴权方式如ECS实例RAM Role和Client Key
* 支持阿里云服务客户端自动刷新AK信息

## 软件要求

- 您的凭据必须是托管RAM凭据
- Golang 1.13 或以上版本

## 背景
当您通过阿里云SDK访问服务时，访问凭证(Access Keys)是被用于认证用户身份. 然而访问凭据容易在使用中存在安全风险，可能会被敌对的开发人员或外部威胁所利用.

阿里云凭据管家具备一种帮助降低风险的解决方案，它允许组织集中管理所有应用程序的访问凭据，允许在不中断应用程序的情况下自动或手动旋转这些凭据。 凭据管家中的托管访问凭据称为[托管RAM凭据](https://help.aliyun.com/document_detail/212421.html) 。

有关使用凭据管家的更多优势信息，请参阅 [凭据管家概述](https://help.aliyun.com/document_detail/152001.html) 。

## 客户端机制
应用程序使用由凭据管家通过代表访问凭据的`凭据名称`管理的访问凭据。

托管凭证插件在访问阿里云开放API时，定期获取由`凭据名称`表示的访问凭据，并提供给阿里云SDK。提供程序通常以指定的间隔（可配置）刷新本地缓存的访问凭据。

在某些情况下，缓存的访问凭据不再有效，这通常发生在管理员在凭据管家中执行紧急访问凭据轮转以响应泄漏事件时。使用无效访问凭据调用OpenAPI通常会导致与API错误代码对应的异常。如果相应的错误代码为`InvalidAccessKeyId.NotFound`或`InvalidAccessKeyId`，则托管凭据插件提供程序将立即刷新缓存的访问评剧馆并重试失败的OpenAPI。

如果API返回使用过期访问凭据的其他错误代码，应用程序开发人员可以覆盖或扩展特定云服务的此行为。请参阅[修改默认过期处理程序](#修改默认过期处理程序示例)。

## 安装

您可以使用`go mod`管理您的依赖，下面以OSS插件导入为示例:

```
require (
	github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider v0.0.5
)
```

或者，通过`go get`命令安装远程代码包:

```
$ go get -u github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider
```


## 支持阿里云云产品

阿里云SDK托管凭据Go插件支持以下云产品:

|                         阿里云SDK名称                          |                                                                          插件名称                                                                          |
|:---------------------------------------------------------:|:------------------------------------------------------------------------------------------------------------------------------------------------------:|
| [阿里云SDK](https://github.com/aliyun/alibaba-cloud-sdk-go)  | [阿里云Go SDK托管凭据插件](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/tree/master/aliyun-sdk-managed-credentials-providers/alibaba-cloud-sdk-go-managed-credentials-provider) |  
| [OSS Go SDK](https://github.com/aliyun/aliyun-oss-go-sdk) |  [OSS Go SDK托管凭据插件](https://github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/tree/master/aliyun-sdk-managed-credentials-providers/aliyun-oss-go-sdk-managed-credentials-provider)  | 


## 使用凭据管家托管RAM凭据方式访问云产品

### 步骤1 配置托管凭据插件

通过配置文件(`managed_credentials_providers.properties`)指定访问凭据管家([配置文件设置详情](README_config.zh-cn.md))，推荐采用Client Key方式访问凭据管家

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

### 步骤 2 使用托管凭据插件访问云服务

下面以托管RAM凭据访问OSS服务为例。

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
	//自定义配置文件
	//utils.SetConfigName("custom-config")
	// 获取Proxy Oss Client
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

	// 通过下面方法关闭客户端来释放插件关联的资源
	client.Shutdown()
}

```

## 修改默认过期处理程序示例

在支持用户自定义错误重试的托管凭据go插件中，用户可以自定义客户端因凭据手动轮转极端场景下的错误重试判断逻辑，只实现以下接口即可。

```go

type AKExpireHandler interface {
	// 判断异常是否由AK过期引起
	JudgeAKExpire(err error) bool
}

```

下面代码示例是用户自定义判断异常接口和使用自定义判断异常实现访问云服务。

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
	//自定义配置文件
	//utils.SetConfigName("custom-config")
	// 获取Proxy Oss Client
	akExpireHandler := &MyOssAKExpireHandler{Code: AkExpireErrorCode}
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

	// 通过下面方法关闭客户端来释放插件关联的资源
	client.Shutdown()
}

```
