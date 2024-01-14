package middlewaresUsecase

import "github.com/DeepAung/deep-art/modules/middlewares/middlewaresRepository"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewaresRepository middlewaresRepository.IMiddlewaresRepository
}

func NewMiddlewaresUsecase(
	middlewaresRepository middlewaresRepository.IMiddlewaresRepository,
) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}
