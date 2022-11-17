package sdk

import osssdk "github.com/aliyun/aliyun-oss-go-sdk/oss"

type OssPluginCredentialsProvider struct {
	credentials osssdk.Credentials
}

func NewOssPluginCredentialsProvider(credentials osssdk.Credentials) *OssPluginCredentialsProvider {
	return &OssPluginCredentialsProvider{
		credentials: credentials,
	}
}

func (sop *OssPluginCredentialsProvider) SetCredentials(credentials osssdk.Credentials) {
	sop.credentials = credentials
}

func (sop *OssPluginCredentialsProvider) GetCredentials() osssdk.Credentials {
	return sop.credentials
}
