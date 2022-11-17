package sdk

import (
	"errors"
	"fmt"
	"reflect"

	commonsdk "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	osssdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var credentialUpdaterSet map[string]struct{}

func init() {
	credentialUpdaterSet = make(map[string]struct{})
	credentialUpdaterSet[reflect.TypeOf(SecretsMangerOssPluginCredentialUpdater{}).Name()] = struct{}{}
}

type secretsManagerOssPlugin struct {
}

type ossClientBuilder struct {
	endpoint        string
	secretName      string
	provider        *OssPluginCredentialsProvider
	akExpireHandler service.AKExpireHandler
}

type ossPluginCredentials struct {
	accessKeyId     string
	accessKeySecret string
}

func NewSecretsManagerOssPlugin() *secretsManagerOssPlugin {
	return &secretsManagerOssPlugin{}
}

func NewOssClientBuilder(endpoint, secretName string, provider *OssPluginCredentialsProvider, akExpireHandler service.AKExpireHandler) *ossClientBuilder {
	return &ossClientBuilder{
		endpoint:        endpoint,
		secretName:      secretName,
		provider:        provider,
		akExpireHandler: akExpireHandler,
	}
}

func NewOssPluginCredentials(accessKeyId, accessKeySecret string) *ossPluginCredentials {
	return &ossPluginCredentials{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}

func (soc *secretsManagerOssPlugin) getOssClient(endpoint, secretName string, akExpireHandler service.AKExpireHandler) (*ProxyOssClient, error) {
	accessKey, err := commonsdk.GetAccessKey(secretName)
	if err != nil {
		return nil, err
	}
	provider := NewOssPluginCredentialsProvider(NewOssPluginCredentials(accessKey.GetAccessKeyId(), accessKey.GetAccessKeySecret()))
	if akExpireHandler == nil {
		akExpireHandler = NewOssAKExpireHandler()
	}
	proxyOssClient, err := NewOssClientBuilder(endpoint, secretName, provider, akExpireHandler).build()
	if err != nil {
		return nil, err
	}
	ossPluginCredentialUpdater := NewSecretsMangerOssPluginCredentialUpdater(proxyOssClient, provider)
	err = commonsdk.RegisterSecretsManagerUpdater(secretName, ossPluginCredentialUpdater)
	if err != nil {
		return nil, err
	}
	return proxyOssClient, nil
}

func (soc *secretsManagerOssPlugin) closeOssClient(client *ProxyOssClient, secretName string) error {
	commonsdk.CloseSecretsManagerPluginUpdaterAndClient(secretName, client)
	return nil
}

func (soc *secretsManagerOssPlugin) destroy() {
	commonsdk.CloseSecretsManagerPluginUpdaterAndClientByTypeName(credentialUpdaterSet)
}

func (ocb *ossClientBuilder) build() (*ProxyOssClient, error) {
	if ocb.provider == nil {
		return nil, errors.New(fmt.Sprintf("provider cannot be null."))
	}
	if ocb.endpoint == "" {
		return nil, errors.New(fmt.Sprintf("Missing parameter endpoint"))
	}
	ossClient, err := osssdk.New(ocb.endpoint, "", "", osssdk.SetCredentialsProvider(ocb.provider))
	if err != nil {
		return nil, err
	}
	return &ProxyOssClient{
		ossClient,
		ocb.secretName,
		ocb.akExpireHandler,
	}, nil
}

func (sc *ossPluginCredentials) GetAccessKeyID() string {
	return sc.accessKeyId
}

func (sc *ossPluginCredentials) GetAccessKeySecret() string {
	return sc.accessKeySecret
}

func (sc *ossPluginCredentials) GetSecurityToken() string {
	return ""
}
