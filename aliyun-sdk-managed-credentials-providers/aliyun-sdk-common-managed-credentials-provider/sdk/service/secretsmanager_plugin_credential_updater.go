package service

import cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"

type SecretsManagerPluginCredentialUpdater interface {
	// 获取云产品Client
	GetClient() interface{}

	// 更新TmpAK信息
	UpdateCredential(secretInfo *cmodels.SecretInfo) error

	// 获取类型名称
	GetTypeName() string

	// 关闭，释放资源
	Close() error
}
