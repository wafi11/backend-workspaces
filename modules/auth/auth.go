package auth

import "context"

type UserRepository interface {
	Create(c context.Context, req RegisterUser) error
	Login(c context.Context, req LoginUser) (*LoginResponse, error)
}
