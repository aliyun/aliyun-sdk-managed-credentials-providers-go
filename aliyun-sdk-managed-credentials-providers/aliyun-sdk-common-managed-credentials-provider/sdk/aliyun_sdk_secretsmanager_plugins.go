package sdk

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/auth"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	alogger "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/logger"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	secretsmanager_client "github.com/aliyun/aliyun-secretsmanager-client-go/sdk"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
	cservice "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/service"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type aliyunSdkSecretsManagerPlugin struct {
	loader              service.SecretsManagerPluginCredentialsLoader
	provider            *service.SecretsManagerPluginCredentialsProvider
	secretCacheClient   *secretsmanager_client.SecretManagerCacheClient
	refreshTimestampMap cmap.ConcurrentMap
}

type TokenBucket struct {
	// 当前的Token数量
	currentTokens int64
	// 最大容量
	maxTokens int64
	//每隔多长时间,增加一个Token
	rate int64
	//上次补充token时间,ms
	lastUpdateTime int64

	lock sync.Mutex
}

func newAliyunSdkSecretsManagerPlugin(loader service.SecretsManagerPluginCredentialsLoader) *aliyunSdkSecretsManagerPlugin {
	return &aliyunSdkSecretsManagerPlugin{
		loader:              loader,
		refreshTimestampMap: cmap.New(),
	}
}

func newTokenBucket(maxTokens, rate int64) *TokenBucket {
	return &TokenBucket{
		maxTokens:      maxTokens,
		currentTokens:  maxTokens,
		rate:           rate,
		lastUpdateTime: time.Now().UnixNano() / 1e6,
	}
}

func (assmp *aliyunSdkSecretsManagerPlugin) init() error {
	if assmp.loader != nil {
		provider, err := assmp.loader.Load()
		if err != nil {
			return err
		}
		assmp.provider = provider
	} else {
		loader := &service.DefaultPluginCredentialsLoader{}
		provider, err := loader.Load()
		if err != nil {
			return err
		}
		assmp.provider = provider
	}
	err := assmp.initLogger()
	if err != nil {
		return err
	}
	err = assmp.initSecretManagerClient()
	if err != nil {
		return err
	}
	// 异步刷新
	go func() {
		for _, secretName := range assmp.provider.SecretNames {
			e := assmp.refreshSecretInfo(secretName)
			if e != nil {
				logger.GetCommonLogger(constants.LoggerName).Errorf("action:refreshSecret", e)
			}
		}
	}()
	logger.GetCommonLogger(constants.LoggerName).Infof("aliyunSdkSecretsManagerPlugin init success")
	return nil
}

func (assmp *aliyunSdkSecretsManagerPlugin) initLogger() error {
	if !logger.IsRegistered(constants.LoggerName) {
		u, err := user.Current()
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
			return err
		}
		logPath := filepath.Join(u.HomeDir, "secretsmanager", "secretsmanager_plugin.log")
		fileStat, err := os.Stat(filepath.Dir(logPath))
		if err != nil || !fileStat.IsDir() {
			err = os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
			if err != nil {
				logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
				return err
			}
		}
		logrusLog := logrus.New()
		log := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     28,
			Compress:   true,
		}
		logrusLog.SetFormatter(&logrus.JSONFormatter{})
		logrusLog.SetOutput(log)
		err = logger.RegisterLogger(constants.LoggerName, alogger.NewLogrusLogger(logrusLog))
		if err != nil {
			logger.GetCommonLogger("").Errorf(err.Error())
			return err
		}
	}
	return nil
}

func (assmp *aliyunSdkSecretsManagerPlugin) initSecretManagerClient() error {
	clientBuilder := cservice.NewDefaultSecretManagerClientBuilder().
		WithCredentials(assmp.provider.Credentials).
		WithBackoffStrategy(assmp.provider.BackOffStrategy)
	for _, regionInfo := range assmp.provider.RegionInfos {
		clientBuilder.AddRegionInfo(regionInfo)
	}
	cacheClient, err := secretsmanager_client.NewSecretCacheClientBuilder(clientBuilder.Build()).
		WithCacheSecretStrategy(assmp.provider.CacheSecretStoreStrategy).
		WithRefreshSecretStrategy(assmp.provider.RefreshSecretStrategy).
		WithSecretCacheHook(assmp.provider.CacheHook).Build()
	if err != nil {
		return err
	}
	assmp.secretCacheClient = cacheClient
	assmp.provider.CacheSecretStoreStrategy.AddRefreshHook(assmp.secretCacheClient)
	return nil
}

func (assmp *aliyunSdkSecretsManagerPlugin) getSecretInfo(secretName string) (*cmodels.SecretInfo, error) {
	return assmp.secretCacheClient.GetSecretInfo(secretName)
}

func (assmp *aliyunSdkSecretsManagerPlugin) getSecretName(userSecretName string) (string, error) {
	return assmp.findSecretName(userSecretName)
}

func (assmp *aliyunSdkSecretsManagerPlugin) getAccessKey(secretName string) (*auth.SecretsManagerPluginCredentials, error) {
	secretInfo, err := assmp.getSecretInfo(secretName)
	if err != nil {
		return nil, err
	}
	return service.GenerateCredentialsBySecret(secretInfo.SecretValue)
}

func (assmp *aliyunSdkSecretsManagerPlugin) registerSecretsManagerUpdater(secretName string, securityUpdater service.SecretsManagerPluginCredentialUpdater) {
	assmp.provider.CacheHook.RegisterSecretsManagerUpdater(secretName, securityUpdater)
}

func (assmp *aliyunSdkSecretsManagerPlugin) closeSecurityUpdaterAndClientByClient(secretName string, client interface{}) {
	assmp.provider.CacheHook.CloseSecurityUpdaterAndClientByClient(secretName, client)
}

func (assmp *aliyunSdkSecretsManagerPlugin) closeSecurityUpdaterAndClientByTypeName(updaterClasses map[string]struct{}) {
	assmp.provider.CacheHook.CloseSecurityUpdaterAndClientByTypeName(updaterClasses)
}

func (assmp *aliyunSdkSecretsManagerPlugin) shutdown() {
	err := assmp.secretCacheClient.Close()
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf("action:shutdown", err)
		return
	}
}

func (assmp *aliyunSdkSecretsManagerPlugin) refreshSecretInfo(secretName string) error {
	if assmp.judgeRefreshSecretInfo(secretName) {
		_, err := assmp.secretCacheClient.RefreshNow(secretName)
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf("action:RefreshNow", err)
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (assmp *aliyunSdkSecretsManagerPlugin) judgeRefreshSecretInfo(secretName string) bool {
	newTokenBucket := newTokenBucket(constants.DefaultMaxTokenNumber, constants.DefaultRateLimitPeriod)
	obj, exists := assmp.refreshTimestampMap.Get(secretName)
	if exists {
		if tokenBucket, ok := obj.(*TokenBucket); ok {
			return tokenBucket.hasQuota()
		}
	}
	assmp.refreshTimestampMap.Set(secretName, newTokenBucket)
	return newTokenBucket.hasQuota()
}

func (assmp *aliyunSdkSecretsManagerPlugin) findSecretName(userSecretName string) (string, error) {
	if userSecretName == "" {
		return "", errors.New(fmt.Sprintf("userSecretName cannot be null"))
	}
	return assmp.provider.SecretExchange.ExchangeSecretName(userSecretName)
}

func (tb *TokenBucket) hasQuota() bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()
	now := time.Now().UnixNano() / 1e6
	// 增加的Token数量
	newTokens := (now - tb.lastUpdateTime) / tb.rate
	if newTokens > 0 {
		tb.lastUpdateTime = now
	}
	tb.currentTokens += newTokens
	if tb.currentTokens > tb.maxTokens {
		tb.currentTokens = tb.maxTokens
	}
	remaining := tb.currentTokens - 1
	if remaining >= 0 {
		tb.currentTokens = remaining
		return true
	}
	return false
}
