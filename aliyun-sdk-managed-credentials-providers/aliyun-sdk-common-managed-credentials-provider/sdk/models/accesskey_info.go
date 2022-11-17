package models

type AccessKeyInfo struct {
	AccessKeyId       string `json:"AccessKeyId,omitempty"`
	AccessKeySecret   string `json:"AccessKeySecret,omitempty"`
	ExpireTimestamp   string `json:"ExpireTimestamp,omitempty"`
	GenerateTimestamp string `json:"GenerateTimestamp,omitempty"`
}
