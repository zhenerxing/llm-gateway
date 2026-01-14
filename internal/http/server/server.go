package server

import(
	"log"
	"os"
	"go.uber.org/zap"
	"fmt"
	"database/sql"
	"context"
	_ "github.com/mattn/go-sqlite3"

	"github.com/zhenerxing/llm-gateway/internal/auth"
	"github.com/zhenerxing/llm-gateway/internal/audit"
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
	svc := auth.PointerService(store)

	// 1) SQLite 路径（建议可配置）
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./data/audit.db"
	}
	_ = os.MkdirAll("./data", 0o755)

	// 2) Open SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3) 可选：简单设置连接参数（SQLite 不需要太复杂）
	db.SetMaxOpenConns(1) // SQLite 常见推荐：避免并发写导致锁冲突
	db.SetConnMaxLifetime(0)

	// 4) Init audit store + schema
	auditStore := &audit.SQLiteStore{DB: db}
	if err := auditStore.InitSchema(context.Background()); err != nil {
		log.Fatal(err)
	}


	// 将http的具体信息交给router去解析分发
	// 同时使用logger将日志结构传入router中
	r:= Router(logger,store,svc, auditStore)

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