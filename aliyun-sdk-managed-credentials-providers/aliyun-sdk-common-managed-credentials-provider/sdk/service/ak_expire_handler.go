package service

type AKExpireHandler interface {
	// 判断异常是否由AK过期引起
	JudgeAKExpire(err error) bool
}
