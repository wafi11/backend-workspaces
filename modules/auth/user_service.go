package auth

import "context"

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return Service{repo: repo}
}

func (s *Service) Register(c context.Context, req RegisterUser) error {
	var err error
	err = validateEmail(req.Email)
	if err != nil {
		return err
	}
	err = validateUsername(req.Username)
	if err != nil {
		return err
	}

	phoneNumber, err := validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return err
	}

	return s.repo.Create(c, RegisterUser{
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: phoneNumber,
		Password:    req.Password,
	})
}
