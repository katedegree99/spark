package repository

import "context"

type ListUsersParams struct {
	Q      *string
	TagIDs []uint
	Offset int
	Limit  int
}

type UserListRecord struct {
	UserID      uint
	Name        string
	Bio         *string
	IconImageID *uint
	ThingIDs    []uint
}

type UserRepository interface {
	ListUsers(ctx context.Context, excludeUserID uint, params ListUsersParams) ([]*UserListRecord, error)
}
