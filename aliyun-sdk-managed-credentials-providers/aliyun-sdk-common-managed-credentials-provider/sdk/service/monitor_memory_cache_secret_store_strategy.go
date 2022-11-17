package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/models"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/utils"
	secretsmanagerclient "github.com/aliyun/aliyun-secretsmanager-client-go/sdk"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
	cmodels "github.com/aliyun/aliyun-secretsmanager-client-go/sdk/models"
	cmap "github.com/orcaman/concurrent-map"
)

const (
	DefaultMonitorPeriod = 30 * 60 * 1000
	DefaultSendPeriod    = 120 * 60 * 1000
)

type MonitorMemoryCacheSecretStoreStrategy struct {
	BlockingQueue               *models.BlockingQueue
	MonitorPeriodMilliseconds   int64
	MonitorCustomerMilliseconds int64
	cacheSecretInfoMap          cmap.ConcurrentMap
	secretsManagerPluginMonitor *secretsManagerPluginMonitor
	secretCacheClient           *secretsmanagerclient.SecretManagerCacheClient
}

type CustomerTask struct {
	blockingQueue *models.BlockingQueue
}

type MonitorTask struct {
	*MonitorMemoryCacheSecretStoreStrategy
	monitorPeriod int64
	blockingQueue *models.BlockingQueue
}

type secretsManagerPluginMonitor struct {
	monitorPeriod         int64
	monitorCustomerPeriod int64
	blockingQueue         *models.BlockingQueue
	scheduledTimer        *time.Timer
	monitorTicker         *time.Ticker
	sendTicker            *time.Ticker
}

func newSecretsManagerPluginMonitor(monitorPeriod int64, monitorCustomerPeriod int64, blockingQueue *models.BlockingQueue) *secretsManagerPluginMonitor {
	return &secretsManagerPluginMonitor{monitorPeriod: monitorPeriod, monitorCustomerPeriod: monitorCustomerPeriod, blockingQueue: blockingQueue}
}

func NewMonitorMemoryCacheSecretStoreStrategy(blockingQueue *models.BlockingQueue, monitorPeriodMilliseconds, monitorCustomerMilliseconds int64) *MonitorMemoryCacheSecretStoreStrategy {
	return &MonitorMemoryCacheSecretStoreStrategy{cacheSecretInfoMap: cmap.New(), BlockingQueue: blockingQueue, MonitorPeriodMilliseconds: monitorPeriodMilliseconds, MonitorCustomerMilliseconds: monitorCustomerMilliseconds, secretsManagerPluginMonitor: newSecretsManagerPluginMonitor(monitorPeriodMilliseconds, monitorCustomerMilliseconds, blockingQueue)}
}

func (m *MonitorMemoryCacheSecretStoreStrategy) Init() error {
	m.secretsManagerPluginMonitor.Init(m)
	return nil
}

func (m *MonitorMemoryCacheSecretStoreStrategy) StoreSecret(cacheSecretInfo *cmodels.CacheSecretInfo) error {
	m.cacheSecretInfoMap.Set(cacheSecretInfo.SecretInfo.SecretName, cacheSecretInfo)
	return nil
}

func (m *MonitorMemoryCacheSecretStoreStrategy) GetCacheSecretInfo(secretName string) (*cmodels.CacheSecretInfo, error) {
	if cacheSecretInfoI, ok := m.cacheSecretInfoMap.Get(secretName); ok {
		if cacheSecretInfo, ok := cacheSecretInfoI.(*cmodels.CacheSecretInfo); ok {
			return cacheSecretInfo, nil
		} else {
			return nil, errors.New(fmt.Sprintf("invalid type [CacheSecretInfo]"))
		}
	}
	return nil, errors.New(fmt.Sprintf("invalid cacheSecretInfoMap key [%s]", secretName))
}

func (m *MonitorMemoryCacheSecretStoreStrategy) Close() error {
	if m.cacheSecretInfoMap != nil {
		m.cacheSecretInfoMap.Clear()
	}
	if m.secretsManagerPluginMonitor != nil {
		m.secretsManagerPluginMonitor.close()
	}
	return nil
}

func (m *MonitorMemoryCacheSecretStoreStrategy) AddRefreshHook(secretCacheClient *secretsmanagerclient.SecretManagerCacheClient) {
	m.secretCacheClient = secretCacheClient
}

func (smpm *secretsManagerPluginMonitor) Init(monitorMemoryCacheSecretStoreStrategy *MonitorMemoryCacheSecretStoreStrategy) {
	if smpm.monitorPeriod < DefaultMonitorPeriod {
		smpm.monitorPeriod = DefaultMonitorPeriod
	}
	if smpm.monitorCustomerPeriod < DefaultSendPeriod {
		smpm.monitorCustomerPeriod = DefaultSendPeriod
	}
	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(DefaultMonitorPeriod)
	logger.GetCommonLogger(constants.LoggerName).Debugf("secretsManagerPluginMonitor create monitor timer")
	smpm.scheduledTimer = time.AfterFunc(time.Duration(start), func() {
		go func() {
			monitor := NewMonitorTask(monitorMemoryCacheSecretStoreStrategy, smpm.monitorPeriod, smpm.blockingQueue)
			smpm.monitorTicker = time.NewTicker(time.Duration(smpm.monitorPeriod) * time.Millisecond)
			for range smpm.monitorTicker.C {
				monitor.Run()
			}
		}()
		go func() {
			send := NewCustomerTask(smpm.blockingQueue)
			smpm.sendTicker = time.NewTicker(time.Duration(smpm.monitorCustomerPeriod) * time.Millisecond)
			for range smpm.sendTicker.C {
				send.Run()
			}
		}()
	})
}
func (smpm *secretsManagerPluginMonitor) close() {
	if smpm.monitorTicker != nil {
		smpm.monitorTicker.Stop()
	}
	if smpm.sendTicker != nil {
		smpm.sendTicker.Stop()
	}
}

func NewCustomerTask(blockingQueue *models.BlockingQueue) *CustomerTask {
	return &CustomerTask{
		blockingQueue: blockingQueue,
	}
}

func (st *CustomerTask) Run() {
	for st.blockingQueue.Size() > 0 {
		obj := st.blockingQueue.Pop()
		monitorMessageInfo, ok := obj.(*models.MonitorMessageInfo)
		if !ok {
			logger.GetCommonLogger(constants.LoggerName).Errorf(fmt.Sprintf("blockingQueue unknown type, expect *models.MonitorMessageInfo"))
			continue
		}
		logger.GetCommonLogger(constants.LoggerName).Warnf("SecretsManagerPluginMonitor occur some problems secretName:{},action:{},errorMessage:{},timestamp:{} ",
			monitorMessageInfo.SecretName, monitorMessageInfo.Action, monitorMessageInfo.ErrorMessage, utils.FormatDate(monitorMessageInfo.Timestamp, utils.TimezoneDatePattern))
	}
}

func NewMonitorTask(monitorMemoryCacheSecretStoreStrategy *MonitorMemoryCacheSecretStoreStrategy, monitorPeriod int64, blockingQueue *models.BlockingQueue) *MonitorTask {
	return &MonitorTask{
		MonitorMemoryCacheSecretStoreStrategy: monitorMemoryCacheSecretStoreStrategy,
		monitorPeriod:                         monitorPeriod,
		blockingQueue:                         blockingQueue,
	}
}

func (mt *MonitorTask) Run() {
	for secretName, _ := range mt.cacheSecretInfoMap.Items() {
		mt.monitorTempAKStatus(secretName)
	}
}

func (mt *MonitorTask) monitorTempAKStatus(secretName string) error {
	secretInfo, err := mt.secretCacheClient.GetSecretInfo(secretName)
	if err != nil {
		return err
	}
	if secretInfo.SecretType != "" && strings.EqualFold(constants.RamCredentialsSecretType, secretInfo.SecretType) {
		extendedConfigStr := secretInfo.ExtendedConfig
		if extendedConfigStr != "" {
			var extendedConfig models.ExtendedConfig
			err = json.Unmarshal([]byte(extendedConfigStr), &extendedConfig)
			if err != nil {
				logger.GetCommonLogger(constants.LoggerName).Errorf(err.Error())
				return err
			}
			if strings.EqualFold(constants.RamUserAccessKeySecretSubType, extendedConfig.SecretSubType) {
				if obj, ok := mt.cacheSecretInfoMap.Get(secretName); ok {
					if cacheSecretInfo, ok := obj.(*cmodels.CacheSecretInfo); ok {
						expired, err := mt.judgeSecretExpired(secretInfo.RotationInterval, cacheSecretInfo.RefreshTimestamp)
						if err != nil {
							return err
						}
						if expired {
							finished, err := mt.secretCacheClient.RefreshNow(secretName)
							if err != nil {
								utils.LogAndAddMonitorMessage(mt.blockingQueue, models.NewMonitorMessageInfo(constants.MonitorAkStatusAction, secretName, "", fmt.Sprintf("secret[%s] ak expire and refresh with err[%s]", secretName, err.Error()), true))
							}
							if !finished {
								utils.LogAndAddMonitorMessage(mt.blockingQueue, models.NewMonitorMessageInfo(constants.MonitorAkStatusAction, secretName, "", fmt.Sprintf("secret[%s] ak expire and fail to refresh", secretName), true))
							} else {
								utils.LogAndAddMonitorMessage(mt.blockingQueue, models.NewMonitorMessageInfo(constants.MonitorAkStatusAction, secretName, "", fmt.Sprintf("secret[%s] ak expire,but success to refresh", secretName), false))
							}
						} else {
							utils.LogAndAddMonitorMessage(mt.blockingQueue, models.NewMonitorMessageInfo(constants.MonitorAkStatusAction, secretName, "", fmt.Sprintf("the status of secret[%s] is normal", secretName), false))
						}
					}
				}
			}
		}
	}
	return nil
}

func (mt *MonitorTask) judgeSecretExpired(interval string, refreshTimestamp int64) (bool, error) {
	if interval == "" {
		return false, errors.New("RotationInterval is nil")
	}
	rotationInterval, err := strconv.ParseInt(strings.Replace(interval, "s", "", 1), 10, 64)
	if err != nil {
		return false, err
	}
	if rotationInterval < 0 {
		return false, errors.New("RotationInterval is invalid")
	}
	return time.Now().UnixNano()/1e6 > refreshTimestamp+rotationInterval*1000, nil
}
