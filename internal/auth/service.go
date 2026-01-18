package auth

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"errors"

	"github.com/google/uuid"
	"github.com/zhenerxing/llm-gateway/internal/apperr"
)

// Service 必须存东西，不关心东西存在哪里，只关心是否实现了Keystore的方法
type Service struct{
	store KeyStore
}

// 将store KeyStore指针化为Service
func PointerService(store KeyStore) *Service{
	return &Service{store: store}
}

// 创建key的输入，包含租户id和quota
type CreateKeyInput struct{
	TenantID  string `json:"tenant_id"`
	Quota 	  Quota `json:"quota"`
}
/*
type Quota struct{
	DailyRequests int `json:"quota_daily_requests"`
	DailyTokens int `json:"quota_daily_tokens"`
}
*/

// 创建key的输出，一般来说只返回一次key
type CreateKeyOutput struct{
	KeyInfo
}
/*
type KeyInfo struct {
	Key       string `json:"key"`
	KeyID     string `json:"key_id"`
	TenantID  string `json:"tenant_id"`
	Active    bool `json:"active"`
	Quota 	  Quota `json:"quota"`
	ExpiresAt *time.Time `json:"expires_at"`
}
*/

// 在auth service 时，报错的规范，后续应该整合到error.go里面
type AppError struct {
	Code 	   string
	HTTPStatus int
	Message    string
	Fields     map[string]any
	Cause      error 
}

// 只关心报错代码和报错信息时的报错函数
func (e *AppError) Error() string {return e.Code + ": " + e.Message}

func (e *AppError) Unwrap() error { return e.Cause }

// 
func BadRequest(code,msg string, fields map[string]any) *AppError {
	return &AppError{
		Code:code,
		HTTPStatus:400,
		Message:msg,
		Fields:fields,
	}
}

// 创建key服务，主要目的是实现接口和数据的分离，通过service类，进行依赖注入
// 在数据传入的时候，在这个函数中统一封装为 KeyInfo（无Key），然后传入到keystore接口的create方法中
// 不关心在接口中存入了什么样的数据结构，接收create方法返回的KeyInfo（有key）
// 返回生成的KeyInfo
func (s *Service) CreateKey (input CreateKeyInput) (*CreateKeyOutput,error){
	tenant := strings.TrimSpace(input.TenantID)
	if tenant == ""{
		return nil, apperr.New(
            apperr.PLATFORM_REQUEST_INVALID,
            apperr.TypePlatform,
            "tenant_id is required",
        ).WithDetails(map[string]any{"tenant_id": "required"})
	}
	mrk,mrk_err := MustRandKey(32)
	if mrk_err != nil{
		return nil, apperr.Wrap(
            apperr.PLATFORM_INTERNAL_ERROR,
            apperr.TypePlatform,
            "random key not generated",
            mrk_err,
        )
	}
	apikey := "gw_" + mrk
	rec := KeyInfo{
		Key : apikey,
		KeyID: uuid.NewString(),
		TenantID: tenant,
		Active :true,
		Quota:Quota{
			DailyRequests: PickDefault(input.Quota.DailyRequests, 1000),
			DailyTokens:   PickDefault(input.Quota.DailyTokens, 200000),
		},
	}
	if err := s.store.Create(rec);err != nil{
		// 例：已存在 -> 409
		if errors.Is(err, ErrAlreadyExists) {
			return nil, apperr.Wrap(
                apperr.PLATFORM_CONFLICT,
                apperr.TypePlatform,
                "key already exists",
                err,
            )
		}
		// 其他错误 -> 500
		return nil, apperr.Wrap(
            apperr.PLATFORM_INTERNAL_ERROR,
            apperr.TypePlatform,
            "internal server error",
            err,
        )
	}
	return &CreateKeyOutput{KeyInfo: rec},nil

}

// 接收从接口中返回的list，将key清空，正面key存在，但不将其暴露

func (s *Service) ListKeys()([]KeyInfo,error){
	recs, err := s.store.List()
	if err != nil{
		return nil,err
	}
	for i := range recs{
		recs[i].Key = ""
	}
	return recs,nil
}

func PickDefault(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func MustRandKey(long int) (string,error){
	b := make([]byte,long)
	_,err := rand.Read(b)
	if  err != nil{
		return "",err
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	return s,err
}
