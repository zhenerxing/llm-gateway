package audit

import(
	"context"
	"database/sql"
	"time"
)

type AuditStore interface{
	InitSchema(ctx context.Context) error
	Insert(ctx context.Context, r AuditInfo) error
	Query(ctx context.Context, TenantID, from, to string, limit int) ([]AuditInfo, error)
}

type SQLiteStore struct{
	DB *sql.DB
}



func (s *SQLiteStore) InitSchema(ctx context.Context) error{
	schema := `
	CREATE TABLE IF NOT EXISTS audit_log(
	id INTEGER PRIMARY KEY  AUTOINCREMENT,
	request_id TEXT NOT NULL,
	key_id TEXT NOT NULL,
	tenant_id TEXT NOT NULL,
	status INTEGER NOT NULL,
	endpoint TEXT NOT NULL,
	created_at TEXT NOT NULL,
	latency_ms INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_audit_key_time ON audit_log(tenant_id, created_at);
	CREATE INDEX IF NOT EXISTS idx_audit_time ON audit_log(created_at);
	`
	_, err := s.DB.ExecContext(ctx, schema)
	return err
}

func (s *SQLiteStore) Insert(ctx context.Context, r AuditInfo) error{
	q := `
	INSERT INTO audit_log(request_id,key_id,tenant_id,status,endpoint,created_at,latency_ms)
	VALUES(?,?,?,?,?,?,?);
	`
	_, err := s.DB.ExecContext(ctx, q,
		r.RequestID,
		r.KeyID,
		r.TenantID,
		r.Status,
		r.Endpoint,
		r.CreatedAt,
		r.LatencyMS,
	)
	return err
}

func (s *SQLiteStore) Query(ctx context.Context, TenantID, from, to string, limit int) ([]AuditInfo, error){
	q := `SELECT request_id , key_id ,tenant_id , status , endpoint , created_at , latency_ms
	FROM audit_log WHERE tenant_id = ?
	`
	args := []any{TenantID}
	if from != ""{
		args = append(args,from)
		q +=  ` AND created_at >= ? `
	}

	if to != "" {
		args = append(args,to)
		q +=  ` AND created_at <= ? `
	}

	args = append(args,limit)
	q += ` ORDER BY created_at DESC LIMIT ? `

	rows,err := s.DB.QueryContext(ctx,q,args...)
	if err != nil{
		return nil,err
	}
	defer rows.Close()

	out := make([]AuditInfo,0,32)
	for rows.Next(){
		var r AuditInfo
		if err := rows.Scan(&r.RequestID,&r.KeyID,&r.TenantID,&r.Status,&r.Endpoint,&r.CreatedAt,&r.LatencyMS);err != nil{
			return nil,err
		}
		out = append(out,r)
	}
	return out , rows.Err()
}
func NowUTC() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}