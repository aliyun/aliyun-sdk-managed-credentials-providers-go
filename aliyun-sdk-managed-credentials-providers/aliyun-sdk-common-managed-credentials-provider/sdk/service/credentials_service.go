package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/auth"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
	cutils "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/utils"
)

func GenerateCredentialsBySecret(secretData string) (*auth.SecretsManagerPluginCredentials, error) {
	if secretData == "" {
		return nil, errors.New(fmt.Sprintf("Missing param secretData"))
	}
	var accessKeyInfo models.AccessKeyInfo
	err := json.Unmarshal([]byte(secretData), &accessKeyInfo)
	if err != nil {
		return nil, err
	}
	var accessKeyId string
	var accessKeySecret string
	if len(accessKeyInfo.AccessKeyId) > 0 && len(accessKeyInfo.AccessKeySecret) > 0 {
		accessKeyId = accessKeyInfo.AccessKeyId
		accessKeySecret = accessKeyInfo.AccessKeySecret
	} else {
		return nil, errors.New(fmt.Sprintf("illegal secret data[%s]", secretData))
	}
	expireTimestamp := int64(constants.NotSupportTampAkTimestamp)
	if len(accessKeyInfo.ExpireTimestamp) > 0 {
		expireTimestamp, err = utils.ParseDate(accessKeyInfo.ExpireTimestamp, utils.TimezoneDatePattern)
		if err != nil {
			return nil, err
		}
	}
	generateTimestamp := int64(math.MaxInt64)
	if len(accessKeyInfo.GenerateTimestamp) > 0 {
		generateTimestamp, err = utils.ParseDate(accessKeyInfo.GenerateTimestamp, utils.TimezoneDatePattern)
		if err != nil {
			return nil, err
		}
	}
	return &auth.SecretsManagerPluginCredentials{
		AccessKeyId:       accessKeyId,
		AccessKeySecret:   accessKeySecret,
		ExpireTimestamp:   expireTimestamp,
		GenerateTimestamp: generateTimestamp,
	}, nil
}

func GenerateSecretInfoByCredentials(securityCredentials *auth.SecretsManagerPluginCredentials, secretName string) (*cmodels.SecretInfo, error) {
	if securityCredentials == nil {
		return nil, errors.New(fmt.Sprintf("Missing param securityCredentials"))
	}
	accessKeyInfo := models.AccessKeyInfo{
		AccessKeyId:     securityCredentials.AccessKeyId,
		AccessKeySecret: securityCredentials.AccessKeySecret,
	}
	secretValue, err := json.Marshal(&accessKeyInfo)
	if err != nil {
		return nil, err
	}
	versionId, err := utils.GenStandardUuid()
	if err != nil {
		return nil, err
	}
	createTime := utils.FormatDate(time.Now(), utils.TimezoneDatePattern)
	nextRotationDate := utils.FormatDate(time.Unix(0, time.Now().UnixNano()+24*60*60*1000*int64(time.Millisecond)), utils.TimezoneDatePattern)
	return &cmodels.SecretInfo{
		SecretName:        secretName,
		VersionId:         versionId,
		SecretValue:       string(secretValue),
		SecretDataType:    cutils.TextDataType,
		CreateTime:        createTime,
		SecretType:        constants.AccessKeySecretType,
		AutomaticRotation: constants.DefaultAutomaticRotation,
		ExtendedConfig:    "",
		RotationInterval:  constants.DefaultRotationInterval,
		NextRotationDate:  nextRotationDate,
	}, nil
}

func ParseTTL(secretInfo *cmodels.SecretInfo, defaultRotationInterval int64) int64 {
	if !strings.EqualFold(constants.AccessKeySecretType, secretInfo.SecretType) {
		var secretValue models.SecretValue
		err := json.Unmarshal([]byte(secretInfo.SecretValue), &secretValue)
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
			return -1
		}
		if len(secretValue.RefreshInterval) == 0 {
			return defaultRotationInterval
		}
		refreshInterval, err := strconv.ParseInt(strings.Replace(secretValue.RefreshInterval, "s", "", 1), 10, 64)
		if err != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
			return -1
		}
		return refreshInterval * 1000
	}
	rotationInterval := secretInfo.RotationInterval
	if len(rotationInterval) == 0 {
		return defaultRotationInterval
	}
	ttl, err := strconv.ParseInt(strings.ReplaceAll(rotationInterval, "s", ""), 10, 64)
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
		return -1
	}
	return ttl * 1000
}
