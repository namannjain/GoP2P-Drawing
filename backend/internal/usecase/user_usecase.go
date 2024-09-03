package usecase

import (
	"goP2Pbackend/internal/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: ur,
	}
}

func (u *userUsecase) Create(user *domain.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) GetByID(id string) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userUsecase) GetByEmail(email string) (*domain.User, error) {
	return u.userRepo.GetByEmail(email)
}

func (u *userUsecase) Update(user *domain.User) error {
	return u.userRepo.Update(user)
}
