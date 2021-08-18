package todos

import "context"

type ListOption int

const (
	All         ListOption = 0
	Completed   ListOption = 1
	Incompleted ListOption = 2
)

type Todo struct {
	ID        string
	Label     string
	Completed bool
}

type Service interface {
	Add(ctx context.Context, todo *Todo) (*Todo, error)
	ChangeStatus(ctx context.Context, id string) (*Todo, error)
	List(ctx context.Context, option ListOption) ([]*Todo, error)
	Remove(ctx context.Context, id string) error
}
