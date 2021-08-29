package storage

import (
	"context"
	"todos"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Add(ctx context.Context, todo *todos.Todo) (*todos.Todo, error) {
	addSQL := `INSERT INTO todos(id, label) VALUES ($1, $2)`
	todo.ID = uuid.NewString()
	_, err := p.db.ExecContext(ctx, addSQL, todo.ID, todo.Label)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (p *Postgres) ChangeStatus(ctx context.Context, id string) (*todos.Todo, error) {
	listSQL := `SELECT todos.id, todos.label, todos.completed FROM todos WHERE todos.id = $1`
	row := p.db.QueryRowContext(ctx, listSQL, id)
	var todo todos.Todo
	err := row.Scan(&todo.ID, &todo.Label, &todo.Completed)
	if err != nil {
		return nil, err
	}

	todo.Completed = !todo.Completed

	changeStatusSQL := `UPDATE todos SET todos.completed = $2 WHERE todos.id = $1`
	_, err = p.db.ExecContext(ctx, changeStatusSQL, todo.ID, todo.Completed)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}
func (p *Postgres) List(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
	listSQL := `SELECT todos.id, todos.label, todos.completed FROM todos`
	rows, err := p.db.QueryContext(ctx, listSQL)
	if err != nil {
		return nil, err
	}
	var result []*todos.Todo
	for rows.Next() {
		var todo todos.Todo
		err := rows.Scan(&todo.ID, &todo.Label, &todo.Completed)
		if err != nil {
			return nil, err
		}
		result = append(result, &todo)
	}
	return result, nil
}
func (p *Postgres) Remove(ctx context.Context, id string) error {
	removeSQL := `DELETE FROM todos WHERE todos.id = $`
	_, err := p.db.ExecContext(ctx, removeSQL, id)
	if err != nil {
		return err
	}

	return nil
}
