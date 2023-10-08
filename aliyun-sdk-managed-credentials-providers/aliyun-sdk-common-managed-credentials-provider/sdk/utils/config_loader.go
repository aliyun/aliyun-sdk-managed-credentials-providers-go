package utils

var configName = ""

func SetConfigName(name string) {
	configName = name
}

func GetConfigName() string {
	return configName
}
