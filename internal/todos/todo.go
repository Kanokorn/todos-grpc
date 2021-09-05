package todos

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
