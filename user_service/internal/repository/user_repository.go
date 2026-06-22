package repository

import (
	"gorm.io/gorm"

	"chat_app/internal/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
) *UserRepository {

	return &UserRepository{
		DB: db,
	}
}
func (r *UserRepository) Create(
	user *model.User,
) error {

	return r.DB.Create(user).Error
}
func (r *UserRepository) GetByID(
	id string,
) (*model.User, error) {

	user := &model.User{}

	err := r.DB.
		Where("id = ?", id).
		First(user).
		Error

	return user, err
}
func (r *UserRepository) GetByEmail(
	email string,
) (*model.User, error) {

	user := &model.User{}

	err := r.DB.
		Where("email = ?", email).
		First(user).
		Error

	return user, err
}
func (r *UserRepository) Update(
	user *model.User,
) error {

	return r.DB.Save(user).Error
}
func (r *UserRepository) Delete(
	id string,
) error {

	return r.DB.
		Delete(
			&model.User{},
			"id = ?",
			id,
		).
		Error
}
func (r *UserRepository) Search(
	query string,
) ([]model.User, error) {

	var users []model.User

	err := r.DB.
		Where(
			"username ILIKE ?",
			"%"+query+"%",
		).
		Find(&users).
		Error

	return users, err
}
