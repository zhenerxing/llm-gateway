package auth

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"errors"

	"github.com/google/uuid"
)

type Service struct{
	store KeyStore
}

func PointerService(store KeyStore) *Service{
	return &Service{store: store}
}

type CreateKeyInput struct{
	TenantID string `json:"tenant_id"`
	Quota 	  Quota `json:"quota"`
}
/*
type Quota struct{
	DailyRequests int `json:"quota_daily_requests"`
	DailyTokens int `json:"quota_daily_tokens"`
}
*/

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

type AppError struct {
	Code 	   string
	HTTPStatus int
	Message    string
	Fields     map[string]any
	Cause      error 
}

func (e *AppError) Error() string {return e.Code + ": " + e.Message}

func BadRequest(code,msg string, fields map[string]any) *AppError {
	return &AppError{
		Code:code,
		HTTPStatus:400,
		Message:msg,
		Fields:fields,
	}
}

func (s *Service) CreateKey (input CreateKeyInput) (*CreateKeyOutput,error){
	tenant := strings.TrimSpace(input.TenantID)
	if tenant == ""{
		return nil,BadRequest(
			"TENANT_ID_REQUIRED",
			"tenant_id is required",
			map[string]any{"tenant_id": "required"},
		)
	}
	mrk,mrk_err := MustRandKey(32)
	if mrk_err != nil{
		return nil,&AppError{
			Code:       "panic",
			HTTPStatus: 400,
			Message:    "random key not generated",
			Cause:      mrk_err,
		}
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
			return nil, &AppError{
				Code:       "ALREADY_EXISTS",
				HTTPStatus: 409,
				Message:    "key already exists",
				Cause:      err,
			}
		}
		// 其他错误 -> 500
		return nil, &AppError{
			Code:       "INTERNAL",
			HTTPStatus: 500,
			Message:    "internal server error",
			Cause:      err,
		}
	}
	return &CreateKeyOutput{KeyInfo: rec},nil

}


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