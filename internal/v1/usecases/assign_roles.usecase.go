package usecases

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/pkg/arrays"
	"gorm.io/gorm"
)

type AssignRolesUsecase struct {
	UserID     uint
	NewRoleIDs []uint

	urRepo repositories.IUserRoleRepo
}

func InitAssignRolesUsecase(tx *gorm.DB) *AssignRolesUsecase {
	return &AssignRolesUsecase{
		urRepo: repositories.NewUserRoleRepo(tx),
	}
}

func (u *AssignRolesUsecase) Execute(c common.Context) error {
	oldRoles, err := u.urRepo.GetRoleIDsByUserID(u.UserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve old roles: %w", err)
	}

	idsToRemove, idsToAdd := arrays.DiffCheck(oldRoles, u.NewRoleIDs)

	if len(idsToRemove) > 0 {
		if err := u.urRepo.DeleteByUserIDAndRoleIDs(c, u.UserID, idsToRemove); err != nil {
			return fmt.Errorf("failed to delete old roles: %w", err)
		}
	}

	if len(idsToAdd) > 0 {
		userRoles := arrays.Map(idsToAdd, func(roleID uint) models.UserRoleInput {
			return models.UserRoleInput{
				UserID: u.UserID,
				RoleID: roleID,
			}
		})

		if err := u.urRepo.BulkCreate(c, userRoles); err != nil {
			return fmt.Errorf("failed to add new roles: %w", err)
		}
	}

	return nil
}
