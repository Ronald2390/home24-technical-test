package adapter

import (
	"context"
	"home24-technical-test/internal/user/service"
)

// LogoutAdapter encapsulate process for logout in adapter
type LogoutAdapter struct {
	service service.ServiceInterface
}

// NewLogoutAdapter build an adapter for logout
func NewLogoutAdapter(
	service service.ServiceInterface,
) LogoutAdapter {
	return LogoutAdapter{
		service: service,
	}
}

func (r LogoutAdapter) Execute(ctx context.Context, token string) error {
	err := r.service.Logout(ctx, token)

	return err
}
