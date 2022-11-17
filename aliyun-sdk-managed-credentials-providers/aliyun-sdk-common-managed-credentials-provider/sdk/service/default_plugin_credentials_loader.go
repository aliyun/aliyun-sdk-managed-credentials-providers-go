package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/service"
	cutils "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/utils"
)

type DefaultPluginCredentialsLoader struct {
}

func (loader *DefaultPluginCredentialsLoader) Load() (*SecretsManagerPluginCredentialsProvider, error) {
	credentialsProperties, err := cutils.LoadCredentialsProperties(constants.DefaultConfigName)
	if err != nil {
		return nil, err
	}
	if credentialsProperties != nil {
		var monitorPeriodMilliseconds int64
		var monitorCustomerMilliseconds int64
		monitorPeriod := credentialsProperties.SourceProperties[constants.PropertiesMonitorPeriodMillisecondsKey]
		monitorCustomer := credentialsProperties.SourceProperties[constants.PropertiesMonitorCustomerMillisecondsKey]
		if len(monitorPeriod) > 0 {
			monitorPeriodMilliseconds, err = strconv.ParseInt(monitorPeriod, 10, 64)
			if err != nil {
				logger.GetCommonLogger(constants.LoggerName).Warnf("action:ParseInt", err)
			}
		}
		if len(monitorCustomer) > 0 {
			monitorCustomerMilliseconds, err = strconv.ParseInt(monitorCustomer, 10, 64)
			if err != nil {
				logger.GetCommonLogger(constants.LoggerName).Warnf("action:ParseInt", err)
			}
		}
		blockingQueue := models.NewBlockingQueue(1000)
		return NewSecretsManagerPluginCredentialsProvider(
			credentialsProperties.Credential,
			credentialsProperties.RegionInfoSlice,
			credentialsProperties.SecretNameSlice,
			NewDefaultSecretExchange(),
			NewMonitorMemoryCacheSecretStoreStrategy(blockingQueue, monitorPeriodMilliseconds, monitorCustomerMilliseconds),
			NewDefaultSecretsManagerPluginCacheHook(blockingQueue, NewDefaultSecretRecoveryStrategy()),
			service.NewFullJitterBackoffStrategy(constants.RetryMaxAttempts, constants.RetryInitialIntervalMills, constants.Capacity),
			NewRotateAKSecretRefreshSecretStrategy(0, 0),
			credentialsProperties.DkmsConfigsMap,
		), nil
	} else {
		return nil, errors.New(fmt.Sprintf("missing default config [%s]", constants.DefaultConfigName))
	}

}