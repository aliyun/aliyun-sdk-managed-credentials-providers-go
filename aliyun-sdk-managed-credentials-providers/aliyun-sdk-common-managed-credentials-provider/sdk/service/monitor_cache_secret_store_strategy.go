package service

import (
	secretsmanagerclient "github.com/aliyun/aliyun-secretsmanager-client-go/sdk"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/cache"
)

type MonitorCacheSecretStoreStrategy interface {
	cache.SecretCacheStoreStrategy
	AddRefreshHook(client *secretsmanagerclient.SecretManagerCacheClient)
}
