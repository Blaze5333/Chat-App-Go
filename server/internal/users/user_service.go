package user

import "time"

type service struct {
	repo    Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	return &service{
		repo:    repo,
		timeout: 2 * time.Second,
	}
}
