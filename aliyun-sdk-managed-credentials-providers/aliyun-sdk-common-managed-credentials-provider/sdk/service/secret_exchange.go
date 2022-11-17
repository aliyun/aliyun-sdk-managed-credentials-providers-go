package service

type SecretExchange interface {
	ExchangeSecretName(userSecretName string) (string, error)
}
type DefaultSecretExchange struct {
}

func NewDefaultSecretExchange() *DefaultSecretExchange {
	return &DefaultSecretExchange{}
}

func (dse *DefaultSecretExchange) ExchangeSecretName(userSecretName string) (string, error) {
	return userSecretName, nil
}
