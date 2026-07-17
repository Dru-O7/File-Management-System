package admin

import (
	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetStats() (*SystemStats, error)
	GetAllUsers() ([]UserResponse, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
	GetAllDocumentTypes() ([]DocumentTypeResponse, error)
	CreateDocumentType(dt *models.DocumentType) error
	GetDocumentTypeByID(id uuid.UUID) (*models.DocumentType, error)
	UpdateDocumentType(dt *models.DocumentType) error
	DeleteDocumentType(id uuid.UUID) error
	GetAllSchools() ([]SchoolResponse, error)
	GetSchoolByID(id uuid.UUID) (*models.School, error)
	UpdateSchool(school *models.School) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetStats() (*SystemStats, error) {
	var stats SystemStats

	r.db.Model(&models.User{}).Count(&stats.TotalUsers)
	r.db.Model(&models.Document{}).Count(&stats.TotalDocuments)
	r.db.Model(&models.School{}).Count(&stats.TotalSchools)
	r.db.Model(&models.DocumentType{}).Count(&stats.TotalDocumentTypes)
	r.db.Model(&models.Document{}).Where("status = ?", models.StatusPendingApproval).Count(&stats.PendingDocuments)
	r.db.Model(&models.Document{}).Where("status = ?", models.StatusApproved).Count(&stats.ApprovedDocuments)
	r.db.Model(&models.Document{}).
		Where("status NOT IN ?", []string{string(models.StatusClosed), string(models.StatusArchived), string(models.StatusRejected)}).
		Count(&stats.ActiveDocuments)
	r.db.Model(&models.WorkflowHistory{}).Where("event_type = ?", "sla_breach").Count(&stats.SLABreaches)

	return &stats, nil
}

func (r *repository) GetAllUsers() ([]UserResponse, error) {
	var users []models.User
	if err := r.db.Preload("School").Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}

	var resp []UserResponse
	for _, u := range users {
		schoolName := ""
		if u.School != nil {
			schoolName = u.School.Name
		}
		resp = append(resp, UserResponse{
			ID:           u.ID,
			Name:         u.Name,
			Email:        u.Email,
			Role:         u.Role,
			SchoolID:     u.SchoolID,
			SchoolName:   schoolName,
			ClassSection: u.ClassSection,
			Subject:      u.Subject,
			Phone:        u.Phone,
			CreatedAt:    u.CreatedAt,
		})
	}
	return resp, nil
}

func (r *repository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *repository) GetAllDocumentTypes() ([]DocumentTypeResponse, error) {
	var dts []models.DocumentType
	if err := r.db.Preload("School").Order("created_at desc").Find(&dts).Error; err != nil {
		return nil, err
	}

	var resp []DocumentTypeResponse
	for _, dt := range dts {
		resp = append(resp, DocumentTypeResponse{
			ID:                dt.ID,
			SchoolID:          dt.SchoolID,
			SchoolName:        dt.School.Name,
			Name:              dt.Name,
			Slug:              dt.Slug,
			WorkflowStages:    dt.WorkflowStages,
			RequiredFields:    dt.RequiredFields,
			SlaHours:          dt.SlaHours,
			Active:            dt.Active,
		})
	}
	return resp, nil
}

func (r *repository) CreateDocumentType(dt *models.DocumentType) error {
	return r.db.Create(dt).Error
}

func (r *repository) GetDocumentTypeByID(id uuid.UUID) (*models.DocumentType, error) {
	var dt models.DocumentType
	if err := r.db.First(&dt, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}

func (r *repository) UpdateDocumentType(dt *models.DocumentType) error {
	return r.db.Save(dt).Error
}

func (r *repository) DeleteDocumentType(id uuid.UUID) error {
	return r.db.Delete(&models.DocumentType{}, "id = ?", id).Error
}

func (r *repository) GetAllSchools() ([]SchoolResponse, error) {
	var schools []models.School
	if err := r.db.Order("created_at desc").Find(&schools).Error; err != nil {
		return nil, err
	}

	var resp []SchoolResponse
	for _, s := range schools {
		resp = append(resp, SchoolResponse{
			ID:        s.ID,
			Name:      s.Name,
			Slug:      s.Slug,
			Settings:  s.Settings,
			CreatedAt: s.CreatedAt,
		})
	}
	return resp, nil
}

func (r *repository) GetSchoolByID(id uuid.UUID) (*models.School, error) {
	var s models.School
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) UpdateSchool(school *models.School) error {
	return r.db.Save(school).Error
}
