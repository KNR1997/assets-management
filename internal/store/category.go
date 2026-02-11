package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"size:100;uniqueIndex;not null"`
	Description string `gorm:"size:255"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryStore struct {
	db *gorm.DB
}

func (s *CategoryStore) GetByID(ctx context.Context, id int64) (*Category, error) {
	var category Category

	err := s.db.WithContext(ctx).
		First(&category, id).
		Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s CategoryStore) Create(ctx context.Context, category *Category) error {
	return s.db.WithContext(ctx).Create(category).Error
}

func (s *CategoryStore) Update(ctx context.Context, category *Category) error {
	result := s.db.WithContext(ctx).
		Model(&Category{}).
		Where("id = ?", category.ID).
		Updates(map[string]interface{}{
			"name":        category.Name,
			"description": category.Description,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *CategoryStore) Delete(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).
		Delete(&Category{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
