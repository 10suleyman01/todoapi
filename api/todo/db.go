package todo

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"todoproject/db"
)

var QueryGetAll = `SELECT * FROM todo WHERE user_id = $1;`
var QueryGetTodoById = `SELECT * FROM todo WHERE id = $1 and user_id = $2`
var QueryCreate = `INSERT INTO todo (title, user_id) VALUES ($1, $2) RETURNING id`
var QueryUpdate = `UPDATE todo SET title = $1 WHERE id = $2 AND user_id = $3`
var QueryDelete = `DELETE FROM todo WHERE id = $1 AND user_id = $2`

type Storage struct {
	db  db.Client
	log *logrus.Logger
}

func NewStorage(db db.Client, log *logrus.Logger) *Storage {
	return &Storage{db: db, log: log}
}

func (s *Storage) GetAllTodoByUserId(ctx context.Context, userId string) (t []Todo, err error) {
	rows, err := s.db.Query(ctx, QueryGetAll, userId)
	if err != nil {
		s.TraceQueryError(err)
		s.log.Errorf("failed to query all todos. due to error: %v", err)
		return t, err
	}
	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		errScan := rows.Scan(&todo.Id, &todo.Title, &todo.UserId)
		if errScan != nil {
			s.log.Errorf("failed to scan todo. due to error: %v", errScan)
			return nil, errScan
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("failed to next row. due to error: %v", err)
		return nil, err
	}

	return todos, nil
}

func (s *Storage) GetTodoById(ctx context.Context, id, userId string) (todo Todo, err error) {
	errQuery := s.db.QueryRow(ctx, QueryGetTodoById, id, &todo.UserId).Scan(&todo.Id, &todo.Title, &todo.UserId)
	if errQuery != nil {
		s.TraceQueryError(errQuery)
		s.log.Errorf("failed to get todo by id=(%s), due to error: %v", id, errQuery)
		return Todo{}, errQuery
	}
	return todo, nil
}

func (s *Storage) Create(ctx context.Context, todo *Todo) error {
	if err := s.db.QueryRow(ctx, QueryCreate, todo.Title, todo.UserId).Scan(&todo.Id); err != nil {
		s.TraceQueryError(err)
		return err
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, todo Todo) error {
	_, errUpdate := s.db.Exec(ctx, QueryUpdate, todo.Title, todo.Id, todo.UserId)
	if errUpdate != nil {
		s.TraceQueryError(errUpdate)
		s.log.Errorf("failed to update todo id=(%s). due to error: %v", todo.Id, errUpdate)
		return errUpdate
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string, userId string) error {
	_, errDelete := s.db.Exec(ctx, QueryDelete, id, userId)
	if errDelete != nil {
		s.TraceQueryError(errDelete)
		s.log.Errorf("failed to delete todo id=(%s). due to error: %v", id, errDelete)
		return errDelete
	}
	return nil
}

func (s *Storage) TraceQueryError(err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		s.log.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
	} else {
		s.log.Error(err)
	}
}
