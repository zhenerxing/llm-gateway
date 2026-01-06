package server

import(
	"log"
	"os"
	"go.uber.org/zap"
	"fmt"

	"github.com/zhenerxing/llm-gateway/internal/auth"
)

func Run()  error {
	// 接收http
	
	//创建zap.Logger实例
	logger,err_zap := zap.NewProduction() // 更好解析
	if err_zap != nil{
		return fmt.Errorf("init zap failed: %w", err_zap)
	}

	//defer增加收尾功能，
	defer func() { _ = logger.Sync() }()

	store := auth.NewInMemoryKeyStore(map[string]auth.KeyInfo{
		"dev-key-123": {
			Key:      "dev-key-123",
			TenantID: "dev",
			Active:   true,
		},
	})



	// 将http的具体信息交给router去解析分发
	// 同时使用logger将日志结构传入router中
	r:= Router(logger,store)

	// 解析地址，如果没有给出地址，则默认8080
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != ""{
		addr = v
	}
	
	// 监听端口，如果返回错误则报错
	log.Printf("listening on %s",addr)
	if err := r.Run(addr); err != nil{
		return fmt.Errorf("gin run(%s): %w", addr, err)
	}

	return nil

}