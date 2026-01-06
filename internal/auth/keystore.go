package auth

import (
	"time"
)

// 构建数据类型，key信息，里面包含key用于匹配，active用来表示是否在线，ExpiresAt表示是否key过期
type KeyInfo struct {
	Key       string
	TenantID  string
	Active    bool
	ExpiresAt *time.Time
}

// 构建了一个接口，这个接口的作用是从KeyInfo类型的map里面去匹配出对应的KeyInfo，bool值的作用是判断是否存在
type KeyStore interface {
	Get(apiKey string) (KeyInfo, bool)
}