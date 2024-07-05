package services

import (
	"errors"

	initi "easyauthapi/configs"
	pl "easyauthapi/models/payload"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//===================================================================

type AdmissionService struct {
	DB *gorm.DB
}

func (s *AdmissionService) NewService() *AdmissionService {
	return &AdmissionService{
		DB: initi.DB,
	}
}

//===================================================================

// Create creates a record
func (s *AdmissionService) Create(InputData *pl.Admission) (*pl.Admission, error) {
	passwordGenerate, err := bcrypt.GenerateFromPassword([]byte(InputData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	InputData.Password = string(passwordGenerate)
	if err := s.DB.Create(InputData).Error; err != nil {
		return nil, errors.New("cannot create new Admission")
	}

	return InputData, nil
}

//===================================================================

// FindByUUID finds by UUID
func (s *AdmissionService) FindByUUID(UuId uuid.UUID) (*pl.Admission, error) {
	var admission pl.Admission

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuid": UuId,
	}
	if err := s.DB.Where(filterSrc).First(&admission).Error; err != nil {
		return nil, errors.New("cannot find admission")
	}
	return &admission, nil
}

//=================

// FindByValue finds by Value
func (s *AdmissionService) FindByValue(filterSrc map[string]interface{}) (*pl.Admission, error) {
	var admission pl.Admission
	if err := s.DB.Where(filterSrc).First(&admission).Error; err != nil {
		return nil, errors.New("cannot find Admission")
	}
	return &admission, nil
}

//=================

// CheckByValue finds by Value
func (s *AdmissionService) CheckByValue(field, value string) error {
	var admission pl.Admission
	filterSrc := map[string]interface{}{
		field: value,
	}
	if err := s.DB.Where(filterSrc).First(&admission).Error; err == nil {
		return errors.New("type " + field + ", user " + value + " is already in use")
	}
	return nil
}

//=================

// GetByQuery finds by Query
func (s *AdmissionService) GetByQuery(filterSrc map[string]interface{}, page int, limit int) ([]pl.Admission, error) {
	var admissions []pl.Admission
	if err := s.DB.Where(filterSrc).Offset(page * limit).Limit(limit + 1).Find(&admissions).Error; err != nil {
		return nil, errors.New("cannot find Admissions")
	}
	return admissions, nil
}

//==================================

// Update updates
func (s *AdmissionService) Update(filterSrc map[string]interface{}, UpdateData *pl.Admission) (*pl.Admission, error) {
	var admission *pl.Admission
	if err := s.DB.Where(filterSrc).First(&admission).Error; err != nil {
		return nil, errors.New("cannot find Admission")
	}

	if err := s.DB.Model(&admission).Updates(UpdateData).Error; err != nil {
		return nil, errors.New("cannot update Admission")
	}

	return admission, nil
}

//==================================

// Delete deletes a One
func (s *AdmissionService) DeleteOne(filterSrc map[string]interface{}) error {
	//menghapus secara fisik Unscoped()
	if err := s.DB.Unscoped().Where(filterSrc).Delete(&pl.Admission{}).Error; err != nil {
		return errors.New("cannot delete Admission")
	}
	return nil
}

//===================================================================
