package service

import (
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/cache"
)

type SecretsManagerPluginCacheHook interface {
	cache.SecretCacheHook
	RegisterSecretsManagerUpdater(secretName string, securityUpdater SecretsManagerPluginCredentialUpdater) error
	CloseSecurityUpdaterAndClientByClient(secretName string, client interface{}) error
	CloseSecurityUpdaterAndClientByTypeName(updaterClasses map[string]struct{}) error
}
