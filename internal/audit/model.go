package audit


// model里面是对audit进行建模
// 要存储：request_id,key_id,tenant_name,endpoint,耗时,create_at,statue
type AuditInfo struct{
	RequestID string `json:"request_id"`
	KeyID     string `json:"key_id"`
	TenantID  string `json:"tenant_id"`
	Status    int    `json:"status"`
	Endpoint  string `json:"endpoint"`
	CreatedAt string  `json:"created_at"`
	LatencyMS int64  `json:"latency_ms"`
}