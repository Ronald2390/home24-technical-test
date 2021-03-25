package adapter

import (
	"context"
	"home24-technical-test/internal/user/service"
)

// ChangePasswordAdapter encapsulate process for change password in adapter
type ChangePasswordAdapter struct {
	service service.ServiceInterface
}

// NewChangePasswordAdapter build an adapter for change password
func NewChangePasswordAdapter(
	service service.ServiceInterface,
) ChangePasswordAdapter {
	return ChangePasswordAdapter{
		service: service,
	}
}

func (r ChangePasswordAdapter) Execute(ctx context.Context, userID int, oldPassword, newPassword string) error {
	err := r.service.ChangePassword(ctx, userID, oldPassword, newPassword)

	return err
}
