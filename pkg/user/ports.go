package user

import "context"

type CreateUserPort interface {
	CreateUser(ctx context.Context, info PublicInfo) error
}
