package auth

import (
	"time"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
)

type CloudCredentials interface {
	GetAccessKeyId() string
	GetAccessKeySecret() string
}

type SecretsManagerPluginCredentials struct {
	AccessKeyId       string
	AccessKeySecret   string
	ExpireTimestamp   int64
	GenerateTimestamp int64
}

func (sc *SecretsManagerPluginCredentials) GetAccessKeyId() string {
	return sc.AccessKeyId
}

func (sc *SecretsManagerPluginCredentials) GetAccessKeySecret() string {
	return sc.AccessKeySecret
}

func (sc *SecretsManagerPluginCredentials) IsExpire() bool {
	if sc.ExpireTimestamp == constants.NotSupportTampAkTimestamp {
		return false
	}
	return time.Now().Unix()/1e6 >= sc.ExpireTimestamp
}

func (sc *SecretsManagerPluginCredentials) GetValidDuration() int64 {
	return sc.ExpireTimestamp - sc.GenerateTimestamp
}
