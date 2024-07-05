package services

import (
	"errors"

	"easyauthapi/configs"
	"easyauthapi/models/datastore"
	pl "easyauthapi/models/payload"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

//===================================================================

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) NewService() *UserService {
	return &UserService{
		DB: configs.DB,
	}
}

//===================================================================

// Create creates a record
func (s *UserService) Create(InputData *pl.User) (*pl.User, error) {
	datastore := datastore.UserDatastore(*InputData)
	user := datastore.Post()

	if err := s.DB.Preload("Roles").Create(user).Error; err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

//==================================

// FindByID finds by ID
func (s *UserService) FindByID(idUser string) (*pl.User, error) {
	var user pl.User
	if err := s.DB.Preload("Roles").First(&user, idUser).Error; err != nil {
		return nil, errors.New("cannot find user")
	}
	return &user, nil
}

//=================

// FindByUUID finds by UUID
func (s *UserService) FindByUUID(UuId uuid.UUID) (*pl.User, error) {
	var user pl.User

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuid": UuId,
	}
	if err := s.DB.Where(filterSrc).Preload("Roles").First(&user).Error; err != nil {
		return nil, errors.New("cannot find user")
	}
	return &user, nil
}

//=================

// GetByQuery gets a paginated list with filters
func (s *UserService) GetByQuery(filterSrc map[string]interface{}, page int, limit int) ([]pl.User, error) {
	var users []pl.User
	query := s.DB.Where(filterSrc).Preload("Roles")

	if err := query.Offset(page * limit).Limit(limit + 1).Find(&users).Error; err != nil {
		return nil, errors.New("cannot find users")
	}

	return users, nil
}

//==================================

// Update updates
func (s *UserService) Update(filterSrc map[string]interface{}, UpdateData *pl.User) (*pl.User, error) {
	var user pl.User
	if err := s.DB.Where(filterSrc).Preload("Roles").First(&user).Error; err != nil {
		return nil, errors.New("cannot find user")
	}

	if err := s.DB.Preload("Roles").Model(&user).Updates(UpdateData).Error; err != nil {
		return nil, errors.New("cannot update user")
	}

	return &user, nil
}

//==================================

// Delete deletes a
func (s *UserService) Delete(filterSrc map[string]interface{}) error {
	//menghapus secara fisik Unscoped()
	if err := s.DB.Unscoped().Where(filterSrc).Preload("Roles").Delete(&pl.User{}).Error; err != nil {
		return errors.New("cannot delete user")
	}
	return nil
}

//===================================================================
