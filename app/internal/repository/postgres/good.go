package postgres

import (
	"database/sql"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/entity/global"
	"hezzl_t/app/internal/repository"

	"github.com/jmoiron/sqlx"
)

type goodRepo struct{}

func NewGoodRepo() repository.Good {
	return &goodRepo{}
}

func (*goodRepo) CreateGood(tx *sqlx.Tx, projectID int, name string) (data entity.Good, err error) {
	err = tx.QueryRow(`INSERT INTO goods (name, project_id) VALUES ($1, $2) 
        RETURNING id, project_id, name, description, priority, removed, created_at`, name, projectID).Scan(
		&data.ID, &data.ProjectID, &data.Name, &data.Description, &data.Priority, &data.Removed, &data.CreatedAt,
	)
	return
}

func (*goodRepo) UpdateGood(tx *sqlx.Tx, id, projectID int, name string, description sql.NullString) (data entity.Good, err error) {
	err = tx.QueryRow(`
	UPDATE goods
	SET name = $1,
		description = $2
	WHERE id = $3
	RETURNING id, project_id, name, description, priority, removed, created_at
	`, name, description, id).Scan(&data.ID, &data.ProjectID, &data.Name, &data.Description, &data.Priority, &data.Removed, &data.CreatedAt)
	return
}

func (*goodRepo) FindGood(tx *sqlx.Tx, id, projectID int) (data entity.Good, err error) {
	err = tx.Get(&data, `
	select id, project_id, name, description, priority, removed, created_at
	from goods
	where id = $1 and project_id = $2
	`, id, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = global.ErrNoData
		}
	}
	return
}

func (*goodRepo) RemoveGood(tx *sqlx.Tx, id, projectID int) (data entity.GoodRemoveInfo, err error) {
	err = tx.QueryRow(`
	UPDATE goods
	SET removed = true
	WHERE id = $1 and project_id = $2
	RETURNING id, project_id, removed
	`, id, projectID).Scan(&data.ID, &data.ProjectID, &data.Removed)
	return
}

func (*goodRepo) LoadGoodInfo(tx *sqlx.Tx, offset, limit int) (data []entity.Good, err error) {
	err = tx.Select(&data, `
	SELECT id, project_id, name, description, priority, removed, created_at
	FROM goods
	OFFSET $1
	LIMIT $2
	order by id
	`, offset, limit)
	if err == nil && len(data) == 0 {
		err = global.ErrNoData
	}

	return
}

func (*goodRepo) FindPriorityList(tx *sqlx.Tx, id int) (data []entity.PriorityInfo, err error) {
	err = tx.Select(&data, `
	select id, priority from goods where id >= $1 order by id
	`, id)
	if err == nil && len(data) == 0 {
		err = global.ErrNoData
		return
	}

	return
}

func (*goodRepo) EditPriority(tx *sqlx.Tx, id, newPriority int) (err error) {
	_, err = tx.Exec(`
	update goods
	set priority = $1
	where id = $2
	`, newPriority, id)

	return
}
