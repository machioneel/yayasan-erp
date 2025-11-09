package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *models.Role) error
	GetByID(id uuid.UUID) (*models.Role, error)
	GetByCode(code string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	GetAllWithPermissions() ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id uuid.UUID) error
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissions(roleID uuid.UUID) error
	GetPermissions(roleID uuid.UUID) ([]models.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("RolePermissions.Permission").First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByCode(code string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("RolePermissions.Permission").First(&role, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) GetAllWithPermissions() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("RolePermissions.Permission").Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	// Check if role is system role
	var role models.Role
	if err := r.db.First(&role, "id = ?", id).Error; err != nil {
		return err
	}
	
	if role.IsSystemRole {
		return errors.New("cannot delete system role")
	}
	
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}

func (r *roleRepository) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	// Remove existing permissions
	if err := r.RemovePermissions(roleID); err != nil {
		return err
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := r.db.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *roleRepository) RemovePermissions(roleID uuid.UUID) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

func (r *roleRepository) GetPermissions(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	
	return permissions, err
}
