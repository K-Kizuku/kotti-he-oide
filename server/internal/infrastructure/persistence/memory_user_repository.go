package persistence

import (
	"context"
	"sync"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

type MemoryUserRepository struct {
	users  map[int]*model.User
	nextID int
	mutex  sync.RWMutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	repo := &MemoryUserRepository{
		users:  make(map[int]*model.User),
		nextID: 1,
		mutex:  sync.RWMutex{},
	}

	email, _ := valueobject.NewEmail("john@example.com")
	userID, _ := valueobject.NewUserID(1)
	sampleUser := model.NewUser(userID, "John Doe", email)
	repo.users[1] = sampleUser
	repo.nextID = 2

	return repo
}

func (r *MemoryUserRepository) Save(ctx context.Context, user *model.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID().Value()] = user
	return nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email().Equals(email) {
			return user, nil
		}
	}
	return nil, nil
}

func (r *MemoryUserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *MemoryUserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id.Value()]; !exists {
		return errors.ErrUserNotFound
	}

	delete(r.users, id.Value())
	return nil
}

func (r *MemoryUserRepository) NextIdentity(ctx context.Context) (valueobject.UserID, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id, err := valueobject.NewUserID(r.nextID)
	if err != nil {
		return valueobject.UserID{}, err
	}

	r.nextID++
	return id, nil
}
