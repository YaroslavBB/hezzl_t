package entity

import (
	"database/sql"
	"time"
)

type Good struct {
	ID          int            `db:"id" json:"id"`
	ProjectID   int            `db:"project_id" json:"project_id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	Priority    int            `db:"priority" json:"priority"`
	Removed     bool           `db:"removed" json:"removed"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

type GoodRemoveInfo struct {
	ID        int  `db:"id" json:"id"`
	ProjectID int  `db:"project_id" json:"project_id"`
	Removed   bool `db:"removed" json:"removed"`
}

type GoodInfo struct {
	Meta struct {
		Total   int
		Removed int
		Limit   int
		Offset  int
	} `json:"meta"`
	Goods []Good `json:"goods"`
}

func (g *GoodInfo) CountMeta() {
	g.Meta.Total = len(g.Goods)

	for _, item := range g.Goods {
		if item.Removed {
			g.Meta.Removed++
		}
	}
}

type GoodInfoCashKey struct {
	Limit  int
	Offset int
}

type PriorityInfo struct {
	ID       int `db:"id" json:"id"`
	Priority int `db:"priority" json:"priority"`
}
