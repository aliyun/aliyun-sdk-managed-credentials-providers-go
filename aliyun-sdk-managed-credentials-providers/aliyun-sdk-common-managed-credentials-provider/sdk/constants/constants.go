package constants

const (
	PropertyNameKeyOldAccessKeyId          = "accessKey"
	PropertyNameKeyOldAccessKeySecret      = "accessSecret"
	PropertyNameKeyAccessKeyId             = "AccessKeyId"
	PropertyNameKeyAccessKeySecret         = "AccessKeySecret"
	PropertyNameKeyExpireTimestamp         = "ExpireTimestamp"
	PropertyNameKeyGenerateTimestamp       = "GenerateTimestamp"
	KmsSecretCurrentStageVersion           = "ACSCurrent"
	PropertyNameKeyRefreshInterval         = "refreshInterval"
	PropertyNameKeyScheduleRotateTimestamp = "scheduleRotateTimestamp"

	AccessKeySecretType = "RAMCredentials"

	// 重试时间间隔，单位ms
	RetryInitialIntervalMills = 2000

	// 最大等待时间，单位ms
	Capacity = 10000

	// 重试最大尝试次数
	RetryMaxAttempts = 6 * 3600 * 1000 / Capacity

	LoggerName = "CacheClient"

	NotSupportTampAkTimestamp = -1

	MonitorAkStatusAction = "monitorAkStatus"

	SendBatchMessageAction = "sendBatchMessage"

	UpdateCredentialAction = "updateCredential"

	RecoveryGetSecretAction = "recoveryGetSecret"

	NoKnownMonitorAccessKeyId = "NotKnown"

	// 默认令牌数
	DefaultMaxTokenNumber = 5

	DefaultRateLimitPeriod = 10 * 60 * 1000

	DefaultRotationInterval = "86400s"

	DefaultAutomaticRotation = "Enabled"

	DefaultRotationIntervalInMs = 6 * 60 * 60 * 1000

	DefaultConfigName = "managed_credentials_providers.properties"

	PropertiesMonitorPeriodMillisecondsKey = "monitor_period_milliseconds"

	PropertiesMonitorCustomerMillisecondsKey = "monitor_customer_milliseconds"

	RamCredentialsSecretType = "RamCredentials"

	ExtendedConfigPropertySecretSubType = "SecretSubType"

	RamUserAccessKeySecretSubType = "RamUserAccessKey"

	CustomMessageAction = "customMessage"
)
