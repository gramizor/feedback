package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"feedback/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Migrate() error {
	// Вызываем AutoMigrate для создания таблиц
	if err := r.db.AutoMigrate(&ds.Group{}, &ds.Request{}, &ds.User{}); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetActiveGroups() (*[]ds.Group, error) {
	var groups []ds.Group

	if err := r.db.Where("status = ?", "active").Find(&groups).Error; err != nil {
		return nil, err
	}
	return &groups, nil
}

func (r *Repository) GetActiveGroupById(id int) (*ds.Group, error) {
	group := &ds.Group{}
	if err := r.db.First(group, id).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func (r *Repository) DeactivateGroupByID(id int) error {
	if err := r.db.Exec("UPDATE groups SET status='deleted' WHERE id= ?", id).Error; err != nil {
		return err
	}
	return nil
}
