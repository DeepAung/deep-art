package middlewaresUsecase

import "github.com/DeepAung/deep-art/modules/middlewares/middlewaresRepository"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId int, accessToken string) bool
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

func (u *middlewaresUsecase) FindAccessToken(userId int, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}
