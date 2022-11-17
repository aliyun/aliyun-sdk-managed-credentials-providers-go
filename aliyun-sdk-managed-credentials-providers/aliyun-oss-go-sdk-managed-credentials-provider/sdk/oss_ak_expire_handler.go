package sdk

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

const (
	AkExpireErrorCode = "InvalidAccessKeyId"
)

type OssAKExpireHandler struct {
	akExpireErrorCode string
}

func NewOssAKExpireHandler() *OssAKExpireHandler {
	return &OssAKExpireHandler{
		akExpireErrorCode: AkExpireErrorCode,
	}
}

func (handler *OssAKExpireHandler) JudgeAKExpire(err error) bool {
	if e, ok := err.(oss.ServiceError); ok {
		if e.Code == handler.akExpireErrorCode {
			return true
		}
	}
	return false
}
