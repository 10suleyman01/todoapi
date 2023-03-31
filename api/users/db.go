package users

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"todoproject/db"
)

var QueryGetAll = `SELECT * FROM users;`
var QueryGetById = `SELECT * FROM users WHERE id = $1`
var QueryGetByNameAndPassword = `SELECT * FROM users WHERE name = $1 AND password_hash=$2`
var QueryCreate = `INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id`
var QueryUpdate = `UPDATE users SET name = $1 WHERE id = $2`
var QueryDelete = `DELETE FROM users WHERE id = $1`

type Storage struct {
	db  db.Client
	log *logrus.Logger
}

func NewStorage(db db.Client, log *logrus.Logger) *Storage {
	return &Storage{db: db, log: log}
}

func (s *Storage) GetAll(ctx context.Context) ([]User, error) {
	rows, err := s.db.Query(ctx, QueryGetAll)
	if err != nil {
		s.log.Errorf("failed to query all users. due to error: %v", err)
		return nil, err
	}
	// make users array
	users := make([]User, 0)
	for rows.Next() {
		var user User
		errScan := rows.Scan(&user.Id, &user.Name, &user.Password)
		if errScan != nil {
			s.log.Errorf("failed to scan users. due to error: %v", errScan)
			return nil, errScan
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		s.log.Errorf("failed to next row. due to error: %v", err)
		return nil, err
	}
	return users, nil
}

func (s *Storage) GetById(ctx context.Context, id string) (user User, errQuery error) {
	if errQuery := s.db.QueryRow(ctx, QueryGetById, id).Scan(&user.Id, &user.Name, &user.Password); errQuery != nil {
		s.TraceQueryError(errQuery)
		s.log.Errorf("failed to query get by id=(%s), due to error: %v", id, errQuery)
		return User{}, errQuery
	}
	return user, nil
}

func (s *Storage) Create(ctx context.Context, user User) error {
	if err := s.db.QueryRow(ctx, QueryCreate, user.Name, user.Password).Scan(&user.Id); err != nil {
		s.TraceQueryError(err)
		return err
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, user User) error {
	_, errUpdate := s.db.Exec(ctx, QueryUpdate, user.Name, user.Id)
	if errUpdate != nil {
		s.TraceQueryError(errUpdate)
		s.log.Errorf("failed to update users=(%v). due to error: %v", user, errUpdate)
		return errUpdate
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	_, errDelete := s.db.Exec(ctx, QueryDelete, id)
	if errDelete != nil {
		s.TraceQueryError(errDelete)
		s.log.Errorf("failed to delete user_id=(%s). due to error: %v", id, errDelete)
		return errDelete
	}
	return nil
}

func (s *Storage) GetByNameAndPassword(ctx context.Context, name, password string) (User, error) {
	var user User
	err := s.db.QueryRow(ctx, QueryGetByNameAndPassword, name, password).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		s.TraceQueryError(err)
		return User{}, err
	}
	return user, nil
}

func (s *Storage) TraceQueryError(err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		s.log.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
	} else {
		s.log.Error(err)
	}
}
