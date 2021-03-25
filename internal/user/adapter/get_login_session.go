package adapter

import (
	"context"
	"home24-technical-test/internal/user/model"
	"home24-technical-test/internal/user/service"
)

// GetLoginSessionAdapter encapsulate process for get login session in adapter
type GetLoginSessionAdapter struct {
	service service.ServiceInterface
}

// NewGetLoginSessionAdapter build an adapter for get login session
func NewGetLoginSessionAdapter(
	service service.ServiceInterface,
) GetLoginSessionAdapter {
	return GetLoginSessionAdapter{
		service: service,
	}
}

func (r GetLoginSessionAdapter) Execute(ctx context.Context, token string) (*model.Session, error) {
	result, err := r.service.GetLoginSession(ctx, token)

	return result, err
}
