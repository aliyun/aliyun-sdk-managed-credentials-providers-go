package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
)

type DefaultSecretsManagerPluginCacheHook struct {
	PluginUpdaterMap       map[string][]SecretsManagerPluginCredentialUpdater
	updaterMtx             sync.Mutex
	BlockingQueue          *models.BlockingQueue
	SecretRecoveryStrategy SecretRecoveryStrategy
}

func NewDefaultSecretsManagerPluginCacheHook(blockingQueue *models.BlockingQueue, secretRecoveryStrategy SecretRecoveryStrategy) *DefaultSecretsManagerPluginCacheHook {
	return &DefaultSecretsManagerPluginCacheHook{
		BlockingQueue:          blockingQueue,
		SecretRecoveryStrategy: secretRecoveryStrategy,
		PluginUpdaterMap:       make(map[string][]SecretsManagerPluginCredentialUpdater),
	}
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) Init() error {
	return nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) Put(secretInfo *cmodels.SecretInfo) (*cmodels.CacheSecretInfo, error) {
	secretName := secretInfo.SecretName
	updaterList := dsmpch.PluginUpdaterMap[secretName]
	if updaterList != nil {
		for _, updater := range updaterList {
			err := updater.UpdateCredential(secretInfo)
			if err != nil {
				utils.LogAndAddMonitorMessage(dsmpch.BlockingQueue, models.NewMonitorMessageInfo(constants.UpdateCredentialAction, secretName, "", err.Error(), true))
			}
		}
	}
	return &cmodels.CacheSecretInfo{
		SecretInfo:       secretInfo,
		RefreshTimestamp: time.Now().UnixNano() / 1e6,
		Stage:            constants.KmsSecretCurrentStageVersion,
	}, nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) Get(cacheSecretInfo *cmodels.CacheSecretInfo) (*cmodels.SecretInfo, error) {
	return cacheSecretInfo.SecretInfo, nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) RecoveryGetSecret(secretName string) (*cmodels.SecretInfo, error) {
	secretInfo, err := dsmpch.SecretRecoveryStrategy.RecoverGetSecret(secretName)
	if err != nil {
		return nil, err
	}
	if secretInfo != nil {
		utils.LogAndAddMonitorMessage(dsmpch.BlockingQueue, models.NewMonitorMessageInfo(constants.RecoveryGetSecretAction, secretName, "", fmt.Sprintf("The secret named [%s] recovery success", secretName), true))
		return secretInfo, nil
	}
	utils.LogAndAddMonitorMessage(dsmpch.BlockingQueue, models.NewMonitorMessageInfo(constants.RecoveryGetSecretAction, secretName, "", fmt.Sprintf("The secret named [%s] recovery fail", secretName), true))
	return nil, nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) Close() error {
	if len(dsmpch.PluginUpdaterMap) > 0 {
		for _, updaterList := range dsmpch.PluginUpdaterMap {
			for i, updater := range updaterList {
				err := updater.Close()
				if err != nil {
					logger.GetCommonLogger(constants.LoggerName).Errorf("action:Close", err)
				}
				updaterList = append(updaterList[:i], updaterList[i+1:]...)
			}
		}
	}
	return nil
}
func (dsmpch *DefaultSecretsManagerPluginCacheHook) RegisterSecretsManagerUpdater(secretName string, securityUpdater SecretsManagerPluginCredentialUpdater) error {
	dsmpch.updaterMtx.Lock()
	defer dsmpch.updaterMtx.Unlock()
	var updaterList []SecretsManagerPluginCredentialUpdater
	updaterList, _ = dsmpch.PluginUpdaterMap[secretName]
	updaterList = append(updaterList, securityUpdater)
	dsmpch.PluginUpdaterMap[secretName] = updaterList
	return nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) CloseSecurityUpdaterAndClientByClient(secretName string, client interface{}) error {
	dsmpch.updaterMtx.Lock()
	defer dsmpch.updaterMtx.Unlock()
	updaterList, ok := dsmpch.PluginUpdaterMap[secretName]
	if ok {
		k := 0
		for _, updater := range updaterList {
			if updater.GetClient() != client {
				updaterList[k] = updater
				k++
			}
		}
		updaterList = updaterList[:k]
		dsmpch.PluginUpdaterMap[secretName] = updaterList
	}
	return nil
}

func (dsmpch *DefaultSecretsManagerPluginCacheHook) CloseSecurityUpdaterAndClientByTypeName(updaterClasses map[string]struct{}) error {
	dsmpch.updaterMtx.Lock()
	defer dsmpch.updaterMtx.Unlock()
	for secretName, credentialUpdaters := range dsmpch.PluginUpdaterMap {
		k := 0
		for _, credentialUpdater := range credentialUpdaters {
			if _, ok := updaterClasses[credentialUpdater.GetTypeName()]; !ok {
				credentialUpdaters[k] = credentialUpdater
				k++
			}
		}
		credentialUpdaters = credentialUpdaters[:k]
		dsmpch.PluginUpdaterMap[secretName] = credentialUpdaters
	}
	return nil
}
