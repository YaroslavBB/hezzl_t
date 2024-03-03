package repository

import (
	"database/sql"
	"hezzl_t/app/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Good interface {
	CreateGood(tx *sqlx.Tx, projectID int, name string) (data entity.Good, err error)
	UpdateGood(tx *sqlx.Tx, id, projectID int, name string, description sql.NullString) (data entity.Good, err error)
	FindGood(tx *sqlx.Tx, id, projectID int) (data entity.Good, err error)
	RemoveGood(tx *sqlx.Tx, id, projectID int) (data entity.GoodRemoveInfo, err error)
	LoadGoodInfo(tx *sqlx.Tx, offset, limit int) (data []entity.Good, err error)
	FindPriorityList(tx *sqlx.Tx, id int) (data []entity.PriorityInfo, err error)
	EditPriority(tx *sqlx.Tx, id, newPriority int) (err error)
}

type Cash interface {
	CashGood(key string, value entity.Good) (err error)
	FindGood(key string) (g entity.Good, err error)
	DelGood(key string) (err error)
	CashGoodInfo(key entity.GoodInfoCashKey, value entity.GoodInfo) error
	DelGoodInfo() (err error)
	FindGoodInfo(key entity.GoodInfoCashKey) (data entity.GoodInfo, err error)
}

type Log interface {
	LogGood(g entity.Good) (err error)
}

type Broker interface {
	SendLog(g entity.Good) (err error)
}
