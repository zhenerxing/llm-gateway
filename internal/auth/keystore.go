package auth

import (
	"time"
)

type Quota struct{
	DailyRequests int `json:"quota_daily_requests"`
	DailyTokens int `json:"quota_daily_tokens"`
}

// 构建数据类型，key信息，里面包含key用于匹配，active用来表示是否在线，ExpiresAt表示是否key过期
type KeyInfo struct {
	Key       string `json:"key"`
	KeyID     string `json:"key_id"`
	TenantID  string `json:"tenant_id"`
	Active    bool `json:"active"`
	Quota 	  Quota `json:"quota"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// 构建了一个接口，这个接口的作用是从KeyInfo类型的map里面去匹配出对应的KeyInfo，bool值的作用是判断是否存在
type KeyStore interface {
	Get(apiKey string) (KeyInfo, bool)
	Create(info KeyInfo) error
	List()([]KeyInfo,error)
}