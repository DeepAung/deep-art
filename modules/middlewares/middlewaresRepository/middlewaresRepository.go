package middlewaresRepository

type IMiddlewaresRepository interface {
}

type middlewaresRepository struct {
}

func NewMiddlewaresRepository() IMiddlewaresRepository {
	return &middlewaresRepository{}
}
