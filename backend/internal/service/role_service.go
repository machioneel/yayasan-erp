package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type RoleService interface {
	GetAll() ([]models.Role, error)
	GetAllWithPermissions() ([]models.Role, error)
	GetByID(id uuid.UUID) (*models.Role, error)
	GetByCode(code string) (*models.Role, error)
	Create(req *models.CreateRoleRequest) (*models.Role, error)
	Update(id uuid.UUID, req *models.UpdateRoleRequest) (*models.Role, error)
	Delete(id uuid.UUID) error
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetPermissions(roleID uuid.UUID) ([]models.Permission, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (s *roleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *roleService) GetAllWithPermissions() ([]models.Role, error) {
	return s.roleRepo.GetAllWithPermissions()
}

func (s *roleService) GetByID(id uuid.UUID) (*models.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *roleService) GetByCode(code string) (*models.Role, error) {
	return s.roleRepo.GetByCode(code)
}

func (s *roleService) Create(req *models.CreateRoleRequest) (*models.Role, error) {
	// Check if code exists
	existing, _ := s.roleRepo.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("role code already exists")
	}

	role := &models.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsSystemRole: false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepo.AssignPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *roleService) Update(id uuid.UUID, req *models.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Cannot update system roles
	if role.IsSystemRole {
		return nil, errors.New("cannot update system role")
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	// Update permissions if provided
	if req.PermissionIDs != nil {
		if err := s.roleRepo.AssignPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *roleService) Delete(id uuid.UUID) error {
	return s.roleRepo.Delete(id)
}

func (s *roleService) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

func (s *roleService) GetPermissions(roleID uuid.UUID) ([]models.Permission, error) {
	return s.roleRepo.GetPermissions(roleID)
}
