package adapter

import (
	"context"
	"home24-technical-test/internal/user/model"
	"home24-technical-test/internal/user/service"
)

// GetUserAdapter encapsulate process for get user in adapter
type GetUserAdapter struct {
	service service.ServiceInterface
}

// NewGetUserAdapter build an adapter for get user
func NewGetUserAdapter(
	service service.ServiceInterface,
) GetUserAdapter {
	return GetUserAdapter{
		service: service,
	}
}

func (r GetUserAdapter) Execute(ctx context.Context, userID int) (*model.User, error) {
	result, err := r.service.GetUser(ctx, userID)

	return result, err
}
