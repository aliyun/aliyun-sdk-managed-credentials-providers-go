 # Configuration File Setting For The Aliyun Security Client 

You can use the configuration file named `managed_credentials_providers.properties` (it exists in the program running directory) to initialize the Aliyun SDK Managed Credentials Providers:

```
## the custom refresh time interval of the secret, by default 6 hour, the minimum value is 5 minutes，the time unit is milliseconds
## the config item to set 1 hour with the custom refresh time interval of the secret 
refresh_secret_ttl=3600000
## The details of the configuration item named cache_client_dkms_config_info:
## 1. The configuration item named cache_client_dkms_config_info must be configured as a json array, you can configure multiple region instances
## 2. ignoreSslCerts:If ignore ssl certs (true: Ignores the ssl certificate, false: Validates the ssl certificate)
## 3. caCert:CA certificate file path, or certificate pem content. If ignoreSslCerts is false, caCert is required. if ignoreSslCerts is true, ignore caCert.
## 4. passwordFromFilePathName and passwordFromEnvVariable
##   passwordFromFilePathName:The client key password configuration is obtained from the file,choose one of the two with passwordFromEnvVariable.
##   e.g. while configuring passwordFromFilePathName: "client_key_password_from_file_path",
##                You need to add client_key_password_from_file_path=< your password file absolute path > in env.
##                and correspond to a file with a password written on it.
##   passwordFromEnvVariable:The client key password configuration is obtained from the environment variable,choose one of the two with passwordFromFilePathName.
##   e.g. while configuring passwordFromEnvVariable: "client_key_password_from_env_variable",
##                You need to add client_key_password_from_env_variable=< your client key private key password from environment variable > in env
##                and the corresponding env variable (xxx_env_variable=<your password>).
## 5. clientKeyFile:The absolute path to the client key json file
## 6. regionId:Region id
## 7. endpoint:Domain address of dkms
cache_client_dkms_config_info=[{"ignoreSslCerts":false,"caCert":"path/to/caCert","passwordFromFilePathName":"client_key_password_from_file_path","clientKeyFile":"<your client key file absolute path>","regionId":"<your dkms region>","endpoint":"<your dkms endpoint>"}]
```
