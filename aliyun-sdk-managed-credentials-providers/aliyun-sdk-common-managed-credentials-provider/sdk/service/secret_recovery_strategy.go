package service

import "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"

type SecretRecoveryStrategy interface {
	RecoverGetSecret(secretName string) (*models.SecretInfo, error)
}

type DefaultSecretRecoveryStrategy struct {
}

func NewDefaultSecretRecoveryStrategy() *DefaultSecretRecoveryStrategy {
	return &DefaultSecretRecoveryStrategy{}
}

func (d *DefaultSecretRecoveryStrategy) RecoverGetSecret(secretName string) (*models.SecretInfo, error) {
	return nil, nil
}
