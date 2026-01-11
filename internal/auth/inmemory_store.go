package auth

import (
	"errors"
	"sync"
)

var ErrAlreadyExists = errors.New("api key already exists")

// 构建内部的KeyInfo信息map，防止传入者还能修改key，同时增加一个锁
type InMemoryStore struct{
	mu sync.RWMutex
	data map[string]KeyInfo
}

// 将传入的seed map[string]KeyInfo保存到内存中，防止seed的传入者还能修改KeyInfo map
func NewInMemoryKeyStore(seed map[string]KeyInfo) *InMemoryStore{
	cp := make(map[string]KeyInfo,len(seed))
	for i ,v := range seed{
		cp[i] = v
	}
	return &InMemoryStore{data : cp}
}

// 对InMemoryStore Struct实现了TypeStore接口

// 作用是通过原子操作，将传入的apikey对应的KeyInfo返回，如果不存在则通过bool体现
func (s *InMemoryStore) Get(apiKey string) (KeyInfo,bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ki, ok := s.data[apiKey]
	return ki,ok
}

/*
type KeyInfo struct {
	Key       string `json:"Key"`
	TenantID  string
	Active    bool
	QUota 	  Quota
	ExpiresAt *time.Time
}
*/
func (s *InMemoryStore) Create(info KeyInfo) error{
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _,exist := s.data[info.Key]; exist{
		return ErrAlreadyExists
	}

	s.data[info.Key] = info

	return nil

}

func (s *InMemoryStore) List()([]KeyInfo,error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	cu := make([]KeyInfo,0,len(s.data))
	for _,v := range s.data{
		cu = append(cu,v)
	}
	return cu,nil
}