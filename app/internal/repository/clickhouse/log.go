package clickhouse_repo

import (
	"context"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/repository"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type logRepo struct {
	conn driver.Conn
	ctx  context.Context
}

func NewLogRepo(conn driver.Conn, ctx context.Context) repository.Log {
	return &logRepo{
		conn: conn,
		ctx:  ctx,
	}
}

func (r *logRepo) LogGood(g entity.Good) (err error) {
	err = r.conn.Exec(r.ctx, `
	insert into table logs(ID, ProjectID, Name, Description, Priority, Removed, EventTime)
	values($1, $2, $3, $4, $5, $6, $7)
	`, g.ID, g.ProjectID, g.Name, g.Description, g.Priority, g.Removed, time.Now())
	return
}
