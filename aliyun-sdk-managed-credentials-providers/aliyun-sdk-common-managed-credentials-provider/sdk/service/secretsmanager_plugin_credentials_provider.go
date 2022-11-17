package service

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
	cservice "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/service"
)

type SecretsManagerPluginCredentialsProvider struct {
	Credentials              auth.Credential
	RegionInfos              []*cmodels.RegionInfo
	SecretNames              []string
	SecretExchange           SecretExchange
	CacheSecretStoreStrategy MonitorCacheSecretStoreStrategy
	CacheHook                SecretsManagerPluginCacheHook
	BackOffStrategy          cservice.BackoffStrategy
	RefreshSecretStrategy    cservice.RefreshSecretStrategy
	DkmsConfigsMap		 map[*cmodels.RegionInfo]*cmodels.DkmsConfig
}

func NewSecretsManagerPluginCredentialsProvider(credentials auth.Credential, regionInfos []*cmodels.RegionInfo, secretNames []string, secretExchange SecretExchange, cacheSecretStoreStrategy MonitorCacheSecretStoreStrategy, cacheHook SecretsManagerPluginCacheHook, backOffStrategy cservice.BackoffStrategy, refreshSecretStrategy cservice.RefreshSecretStrategy, dkmsConfigsMap map[*cmodels.RegionInfo]*cmodels.DkmsConfig) *SecretsManagerPluginCredentialsProvider {
	return &SecretsManagerPluginCredentialsProvider{
		Credentials:              credentials,
		RegionInfos:              regionInfos,
		SecretNames:              secretNames,
		SecretExchange:           secretExchange,
		CacheSecretStoreStrategy: cacheSecretStoreStrategy,
		CacheHook:                cacheHook,
		BackOffStrategy:          backOffStrategy,
		RefreshSecretStrategy:    refreshSecretStrategy,
		DkmsConfigsMap:	          dkmsConfigsMap,
	}
}
