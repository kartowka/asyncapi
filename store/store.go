package store

import "gorm.io/gorm"

type Store struct {
	Users *UserStore
}

func New(db *gorm.DB) *Store {
	return &Store{
		Users: NewUserStore(db),
	}
}
