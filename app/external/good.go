package external

import (
	"database/sql"
	"encoding/json"
	"hezzl_t/app/internal/entity/global"
	"hezzl_t/app/internal/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type GoodHandlers struct {
	GoodUsecase usecase.GoodUsecase
	DB          *sqlx.DB
}

func NewGoodHandlers(goodUsecase usecase.GoodUsecase, db *sqlx.DB) *GoodHandlers {
	return &GoodHandlers{
		GoodUsecase: goodUsecase,
		DB:          db,
	}

}

func (e *GoodHandlers) CreateGood(w http.ResponseWriter, r *http.Request) {
	pID := r.URL.Query().Get("projectID")
	projectID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var param struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.CreateCood(tx, projectID, param.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (e *GoodHandlers) UpdateGood(w http.ResponseWriter, r *http.Request) {
	var (
		projectID int
		id        int
		pID       string
		rID       string
		err       error
		param     struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		desc sql.NullString
	)

	rID = r.URL.Query().Get("id")
	id, err = strconv.Atoi(rID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pID = r.URL.Query().Get("projectID")
	projectID, err = strconv.Atoi(pID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(param.Description) != "" {
		desc = sql.NullString{Valid: true, String: param.Description}
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.UpdateGood(tx, id, projectID, param.Name, desc)
	if err != nil {
		if err == global.ErrGoodNotFount {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (e *GoodHandlers) FindGood(w http.ResponseWriter, r *http.Request) {
	var (
		projectID int
		id        int
		pID       string
		idStr     string
		err       error
	)

	idStr = r.URL.Query().Get("id")
	id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pID = r.URL.Query().Get("projectID")
	projectID, err = strconv.Atoi(pID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.FindGood(tx, id, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (e *GoodHandlers) RemoveGood(w http.ResponseWriter, r *http.Request) {
	var (
		projectID int
		id        int
		pID       string
		rID       string
		err       error
	)

	rID = r.URL.Query().Get("id")
	id, err = strconv.Atoi(rID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pID = r.URL.Query().Get("projectID")
	projectID, err = strconv.Atoi(pID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.RemoveGood(tx, id, projectID)
	if err != nil {
		if err == global.ErrGoodNotFount {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (e *GoodHandlers) LoadGoodInfo(w http.ResponseWriter, r *http.Request) {
	var (
		limit     int
		offset    int
		limitStr  string
		offsetStr string
		err       error
	)

	limitStr = r.URL.Query().Get("limit")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offsetStr = r.URL.Query().Get("offset")
	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.LoadGoodInfo(tx, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (e *GoodHandlers) UpdatePriority(w http.ResponseWriter, r *http.Request) {
	var (
		projectID int
		id        int
		pID       string
		idStr     string
		err       error
		param     struct {
			NewPriority int `json:"newPriority"`
		}
	)

	idStr = r.URL.Query().Get("id")
	id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pID = r.URL.Query().Get("projectID")
	projectID, err = strconv.Atoi(pID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := e.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	data, err := e.GoodUsecase.UpdatePriority(tx, id, projectID, param.NewPriority)
	if err != nil {
		if err == global.ErrGoodNotFount {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
