# 阿里云托管凭据客户端配置文件设置 

`managed_credentials_providers.properties`(在程序运行目录下)初始化阿里云凭据管家动态RAM凭据客户端：

```properties
## 用户自定义的刷新频率, 默认为6小时，最小值为5分钟，单位为毫秒
## 下面的配置将凭据刷新频率设定为1小时
refresh_secret_ttl=3600000
## cache_client_dkms_config_info配置项说明:
## 1. cache_client_dkms_config_info配置项为json数组，支持配置多个region实例
## 2. ignoreSslCerts:是否忽略ssl证书 (true:忽略ssl证书,false:验证ssl证书)
## 3. caCert:CA证书路径，或证书pem内容。如果ignoreSslCerts是false，caCert为必填，如果ignoreSslCerts是true，忽略caCert
## 4. passwordFromFilePathName和passwordFromEnvVariable
## passwordFromFilePathName:client key密码配置从文件中获取，与passwordFromEnvVariable二选一
## 例:当配置passwordFromFilePathName:"client_key_password_from_file_path"时，
##    需在环境变量中添加client_key_password_from_file_path=<你的client key密码文件所在的绝对路径>，
##    以及对应写有password的文件。
##    passwordFromEnvVariable:client key密码配置从环境变量中获取，与passwordFromFilePathName二选一
## 例:当配置passwordFromEnvVariable:"client_key_password_from_env_variable"时，
##    需在环境变量中添加client_key_password_from_env_variable=<你的client key密码对应的环境变量名>
##    以及对应的环境变量(xxx_env_variable=<your password>)。
## 5. clientKeyFile:client key json文件的绝对路径
## 6. regionId:地域Id
## 7. endpoint:专属kms的域名地址
cache_client_dkms_config_info=[{"ignoreSslCerts":false,"caCert":"path/to/caCert","passwordFromFilePathName":"client_key_password_from_file_path","clientKeyFile":"<your client key file absolute path>","regionId":"<your dkms region>","endpoint":"<your dkms endpoint>"}]
```
