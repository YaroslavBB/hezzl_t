package usecase

import (
	"database/sql"
	"fmt"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/entity/global"
	repository_import "hezzl_t/app/internal/repository/import"
	"log"

	"github.com/jmoiron/sqlx"
)

type GoodUsecase struct {
	ri *repository_import.RepositoryImports
}

func NewGoodUsecase(rp *repository_import.RepositoryImports) *GoodUsecase {
	return &GoodUsecase{
		rp,
	}
}

func (u *GoodUsecase) CreateCood(tx *sqlx.Tx, projectID int, name string) (data entity.Good, err error) {
	data, err = u.ri.Good.CreateGood(tx, projectID, name)
	if err != nil {
		log.Println("CreateCood", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Cash.DelGoodInfo(); err != nil {
		log.Println("Cash DelGoodInfo", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Broker.SendLog(data); err != nil {
		log.Println("LogGood", err.Error())
		err = global.ErrInternalError
		return
	}

	return
}

func (u *GoodUsecase) UpdateGood(tx *sqlx.Tx, id, projectID int, name string, description sql.NullString) (data entity.Good, err error) {
	_, err = u.FindGood(tx, id, projectID)
	if err != nil {
		return
	}

	data, err = u.ri.Good.UpdateGood(tx, id, projectID, name, description)
	if err != nil {
		log.Println("UpdateGood", err.Error())
		err = global.ErrInternalError
		return
	}

	err = u.ri.Cash.DelGood(fmt.Sprintf("%d", id))
	if err != nil {
		log.Println("Cash DelGood", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Cash.DelGoodInfo(); err != nil {
		log.Println("Cash DelGoodInfo", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Broker.SendLog(data); err != nil {
		log.Println("LogGood", err.Error())
		err = global.ErrInternalError
		return
	}

	return
}

func (u *GoodUsecase) FindGood(tx *sqlx.Tx, id, projectID int) (data entity.Good, err error) {
	data, err = u.ri.Cash.FindGood(fmt.Sprintf("%d", id))
	switch err {
	case nil:
		return

	case global.ErrNoData:
		//
	default:
		log.Println("Cash FindGood", err.Error())
	}

	data, err = u.ri.Good.FindGood(tx, id, projectID)
	switch err {
	case nil:
		//

	case global.ErrNoData:
		err = global.ErrGoodNotFount
		return

	default:
		log.Println("FindGood", err.Error())
		err = global.ErrInternalError
		return
	}

	pErr := u.ri.Cash.CashGood(fmt.Sprintf("%d", data.ID), data)
	if pErr != nil {
		log.Println("CashGood", pErr.Error())
	}
	return
}

func (u *GoodUsecase) RemoveGood(tx *sqlx.Tx, id, projectID int) (data entity.GoodRemoveInfo, err error) {
	g, err := u.FindGood(tx, id, projectID)
	if err != nil {
		return
	}

	data, err = u.ri.Good.RemoveGood(tx, id, projectID)
	if err != nil {
		log.Println("RemoveGood", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Cash.DelGood(fmt.Sprintf("%d", data.ID)); err != nil {
		log.Println("Cash DelGood", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Cash.DelGoodInfo(); err != nil {
		log.Println("Cash DelGoodInfo", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Broker.SendLog(entity.Good{
		ID:          data.ID,
		ProjectID:   data.ProjectID,
		Name:        g.Name,
		Description: g.Description,
		Priority:    g.Priority,
		Removed:     data.Removed,
		CreatedAt:   g.CreatedAt,
	}); err != nil {
		log.Println("LogGood", err.Error())
		err = global.ErrInternalError
		return
	}

	return
}

func (u *GoodUsecase) LoadGoodInfo(tx *sqlx.Tx, offset, limit int) (data entity.GoodInfo, err error) {
	if offset == 0 {
		offset = global.DefaultOffset
	}

	if limit == 0 {
		limit = global.DefaultLimit
	}

	var (
		giKey = entity.GoodInfoCashKey{
			Limit:  limit,
			Offset: offset,
		}
	)

	data, err = u.ri.Cash.FindGoodInfo(giKey)
	switch err {
	case nil:
		return

	case global.ErrNoData:
		//
	default:
		log.Println("Cash FindGoodInfo", err.Error())
	}

	data.Goods, err = u.ri.Good.LoadGoodInfo(tx, offset, limit)
	switch err {
	case nil:
		data.CountMeta()
	case global.ErrNoData:
		//
	default:
		log.Println("LoadGoodInfo", err.Error())
		err = global.ErrInternalError
		return
	}

	err = u.ri.Cash.CashGoodInfo(giKey, data)
	if err != nil {
		log.Println("Cash CashGoodInfo", err.Error())
	}

	return
}

func (u *GoodUsecase) UpdatePriority(tx *sqlx.Tx, id, projectID, newPriority int) (data []entity.PriorityInfo, err error) {
	_, err = u.FindGood(tx, id, projectID)
	if err != nil {
		return
	}

	data, err = u.ri.Good.FindPriorityList(tx, id)
	if err != nil {
		log.Println("FindPriorityList", err.Error())
		err = global.ErrInternalError
		return
	}

	for _, item := range data {
		err = u.ri.Good.EditPriority(tx, item.ID, newPriority)
		if err != nil {
			log.Println("EditPriority", err.Error())
			err = global.ErrInternalError
			return
		}

		err = u.ri.Cash.DelGood(fmt.Sprintf("%d", item.ID))
		if err != nil {
			log.Println("DelGood", err.Error())
			err = global.ErrInternalError
			return
		}

		newPriority++

	}

	data, err = u.ri.Good.FindPriorityList(tx, id)
	if err != nil {
		log.Println("FindPriorityList", err.Error())
		err = global.ErrInternalError
		return
	}

	if err = u.ri.Cash.DelGoodInfo(); err != nil {
		log.Println("Cash DelGoodInfo", err.Error())
		err = global.ErrInternalError
		return
	}

	return
}
