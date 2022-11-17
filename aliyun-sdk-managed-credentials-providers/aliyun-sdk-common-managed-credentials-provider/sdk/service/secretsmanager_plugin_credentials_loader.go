package service

type SecretsManagerPluginCredentialsLoader interface {
	Load() (*SecretsManagerPluginCredentialsProvider, error)
}
