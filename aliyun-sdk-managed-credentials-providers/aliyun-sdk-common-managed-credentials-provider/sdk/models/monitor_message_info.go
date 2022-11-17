package models

import (
	"time"
)

type MonitorMessageInfo struct {
	Action       string
	SecretName   string
	AccessKeyId  string
	ErrorMessage string
	Alarm        bool
	Timestamp    time.Time
}

func NewMonitorMessageInfo(action, secretName, accessKeyId, errorMessage string, alarm bool) *MonitorMessageInfo {
	return &MonitorMessageInfo{
		Action:       action,
		SecretName:   secretName,
		AccessKeyId:  accessKeyId,
		ErrorMessage: errorMessage,
		Alarm:        alarm,
		Timestamp:    time.Now(),
	}
}
