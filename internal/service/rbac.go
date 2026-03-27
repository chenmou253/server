package service

import (
	"context"
	"fmt"
	"slices"

	"go-admin/server/internal/model"
	"go-admin/server/internal/repository"
)

type RBACService struct {
	repo *repository.Repository
}

func NewRBACService(repo *repository.Repository) *RBACService {
	return &RBACService{repo: repo}
}

func (s *RBACService) GetProfile(ctx context.Context, userID uint) (*model.User, []string, []model.Menu, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	if user.Status != 1 {
		return nil, nil, nil, fmt.Errorf("user disabled")
	}

	permissionSet := make([]string, 0)
	menuMap := make(map[uint]model.Menu)
	for _, role := range user.Roles {
		if role.Status != 1 {
			continue
		}
		for _, permission := range role.Permissions {
			if !slices.Contains(permissionSet, permission.Code) {
				permissionSet = append(permissionSet, permission.Code)
			}
		}
		for _, menu := range role.Menus {
			menuMap[menu.ID] = menu
		}
	}

	menus := make([]model.Menu, 0, len(menuMap))
	for _, menu := range menuMap {
		menus = append(menus, menu)
	}
	slices.SortFunc(menus, func(a, b model.Menu) int { return a.Sort - b.Sort })

	return user, permissionSet, menus, nil
}

func (s *RBACService) HasPermission(permissions []string, required string) bool {
	if required == "" {
		return true
	}
	return slices.Contains(permissions, required)
}
