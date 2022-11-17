package service

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
)

const (
	DefaultDelayInterval          = 5 * 60 * 1000
	DefaultRandomDisturbanceRange = 10 * 60 * 1000
)

type RotateAKSecretRefreshSecretStrategy struct {
	rotationInterval  int64
	delayInterval     int64
	randomDisturbance int64
}

func NewRotateAKSecretRefreshSecretStrategy(rotationInterval, delayInterval int64) *RotateAKSecretRefreshSecretStrategy {
	rand.Seed(time.Now().UnixNano())
	refreshSecretStrategy := &RotateAKSecretRefreshSecretStrategy{
		rotationInterval:  constants.DefaultRotationIntervalInMs,
		delayInterval:     DefaultDelayInterval,
		randomDisturbance: int64(rand.Intn(DefaultRandomDisturbanceRange)),
	}
	if rotationInterval > 0 {
		refreshSecretStrategy.rotationInterval = rotationInterval
	}
	if delayInterval > 0 {
		refreshSecretStrategy.delayInterval = delayInterval
	}
	return refreshSecretStrategy
}

func (rrs *RotateAKSecretRefreshSecretStrategy) Init() error {
	return nil
}

func (rrs *RotateAKSecretRefreshSecretStrategy) GetNextExecuteTime(secretName string, ttl, offsetTimestamp int64) int64 {
	now := time.Now().UnixNano() / 1e6
	if ttl+offsetTimestamp > now {
		return ttl + offsetTimestamp + rrs.randomDisturbance
	} else {
		return now + ttl + rrs.randomDisturbance
	}
}

func (rrs *RotateAKSecretRefreshSecretStrategy) ParseNextExecuteTime(cacheSecretInfo *cmodels.CacheSecretInfo) int64 {
	secretInfo := cacheSecretInfo.SecretInfo
	nextRotationDate := rrs.parseNextRotationDate(secretInfo)
	now := time.Now().UnixNano() / 1e6
	if nextRotationDate >= now+rrs.rotationInterval+rrs.randomDisturbance || nextRotationDate <= now {
		return now + rrs.rotationInterval + rrs.randomDisturbance
	} else {
		return nextRotationDate + rrs.delayInterval + rrs.randomDisturbance
	}
}

func (rrs *RotateAKSecretRefreshSecretStrategy) ParseTTL(secretInfo *cmodels.SecretInfo) int64 {
	if secretInfo != nil && secretInfo.SecretType != "" && strings.EqualFold(constants.AccessKeySecretType, secretInfo.SecretType) {
		extendedConfigStr := secretInfo.ExtendedConfig
		var extendedConfig models.ExtendedConfig
		err := json.Unmarshal([]byte(extendedConfigStr), &extendedConfig)
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
			return -1
		}
		if strings.EqualFold(constants.RamUserAccessKeySecretSubType, extendedConfig.SecretSubType) {
			rotationInterval := secretInfo.RotationInterval
			if rotationInterval == "" {
				return rrs.rotationInterval + rrs.randomDisturbance
			}
			rotationIntervalInt, err := strconv.ParseInt(strings.Replace(rotationInterval, "s", "", 1), 10, 64)
			if err != nil {
				logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
				return -1
			}
			return rotationIntervalInt*1000 + rrs.randomDisturbance
		}
	}
	var secretValue models.SecretValue
	err := json.Unmarshal([]byte(secretInfo.SecretValue), &secretValue)
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
		return -1
	}
	if len(secretValue.RefreshInterval) == 0 {
		return -1
	}
	refreshIntervalInt, err := strconv.ParseInt(strings.Replace(secretValue.RefreshInterval, "s", "", 1), 10, 64)
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
		return -1
	}
	return refreshIntervalInt*1000 + rrs.randomDisturbance
}

func (rrs *RotateAKSecretRefreshSecretStrategy) Close() error {
	return nil
}

func (rrs *RotateAKSecretRefreshSecretStrategy) parseNextRotationDate(secretInfo *cmodels.SecretInfo) int64 {
	nextRotationDateStr := secretInfo.NextRotationDate
	if nextRotationDateStr == "" {
		secretValue := models.SecretValue{ScheduleRotateTimestamp: -1}
		err := json.Unmarshal([]byte(secretInfo.SecretValue), &secretValue)
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
			return -1
		}
		return secretValue.ScheduleRotateTimestamp * 1000
	}
	nextRotationDate, err := utils.ParseDate(nextRotationDateStr, utils.TimezoneDatePattern)
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
		return -1
	}
	return nextRotationDate
}
