package adapter

import (
	"context"
	"home24-technical-test/internal/user/public"
	"home24-technical-test/internal/user/service"
)

// LoginAdapter encapsulate process for login in adapter
type LoginAdapter struct {
	service service.ServiceInterface
}

// NewLoginAdapter build an adapter for login
func NewLoginAdapter(
	service service.ServiceInterface,
) LoginAdapter {
	return LoginAdapter{
		service: service,
	}
}

func (r LoginAdapter) Execute(ctx context.Context, params *public.LoginParams) (*public.LoginResponse, error) {
	result, err := r.service.Login(ctx, params)

	return result, err
}
