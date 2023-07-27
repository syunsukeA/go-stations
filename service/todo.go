package service

import (
	"context"
	"database/sql"
	"time"
	_"strconv"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mattn/go-sqlite3"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	todo := model.TODO{
		Subject: 	 subject,
		Description: description,
	}

	if len(subject) == 0 {
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()
		return &todo, nil
	}

	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		panic(err)
	}

	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	row := s.db.QueryRowContext(ctx, confirm, id)

	var created_at, updated_at time.Time
	row.Scan(&subject, &description, &created_at, &updated_at)
	todo.ID = id
	todo.CreatedAt = created_at
	todo.UpdatedAt = updated_at

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	
	var rows *sql.Rows
	var err error
	todos := []*model.TODO{}

	if prevID > 0 && size > 0{
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else if size > 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		size = 5
		rows, err = s.db.QueryContext(ctx, read, size)
	}
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id int64
		var subject, description string
		var created_at, updated_at time.Time
		err := rows.Scan(&id, &subject, &description, &created_at, &updated_at)
		if err != nil {
			break
		}
		todo := model.TODO{
			ID: id,
			Subject: subject,
			Description: description,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	todo := model.TODO{
		ID: id,
		Subject: subject,
		Description: description,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if id <= 0 {
		return &todo, &model.ErrNotFound{}
	}

	if len(subject) == 0 {
		return &todo, sqlite3.Error{Code: 19}
	}

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		panic(err)
	}

	_, err = stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		panic(err)
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	var created_at, updated_at time.Time
	if row.Scan(&subject, &description, &created_at, &updated_at) != nil {
		return &todo, &model.ErrNotFound{}
	}

	todo.CreatedAt = created_at
	todo.UpdatedAt = updated_at

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	if len(ids) <= 0 {
		return nil
	}
	NewdeleteFmt := fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))
	stmt, err := s.db.PrepareContext(ctx, NewdeleteFmt)
	var new_ids []interface{}
    for _, id := range ids {
        new_ids = append(new_ids, id)
    }

	res, err := stmt.ExecContext(ctx, new_ids...)
	if err != nil {
		panic(err)
	}
	n_affected, err := res.RowsAffected()
	if n_affected == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
