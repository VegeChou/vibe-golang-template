package memory

import (
	"fmt"
	"sync"

	"vibe-golang-template/internal/model"
)

type UserRepository struct {
	mu    sync.RWMutex
	next  int
	users []model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{next: 1}
}

func (r *UserRepository) List() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.User, len(r.users))
	copy(result, r.users)
	return result
}

func (r *UserRepository) Create(user model.User) model.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("u-%03d", r.next)
	r.next++
	r.users = append(r.users, user)
	return user
}
