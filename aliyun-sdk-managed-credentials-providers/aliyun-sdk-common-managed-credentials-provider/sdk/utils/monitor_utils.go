package utils

import (
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
)

func LogAndAddMonitorMessage(blockingQueue *models.BlockingQueue, monitorMessageInfo *models.MonitorMessageInfo) {
	logger.GetCommonLogger(constants.LoggerName).Errorf(monitorMessageInfo.Action, monitorMessageInfo.SecretName, monitorMessageInfo.ErrorMessage)
	blockingQueue.Offer(monitorMessageInfo)
}
