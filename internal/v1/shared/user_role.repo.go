package shared

import (
	"github.com/MarcelArt/gotel/internal/common"
	"gorm.io/gorm"
)

type IUserRoleRepoTx interface {
	BeginTx(tx *gorm.DB) IUserRoleRepoTx
	GetRoleIDsByUserID(userID any) ([]uint, error)
	DeleteByUserIDAndRoleIDs(c common.Context, userID any, roleIDs []uint) error
	BulkCreate(c common.Context, input []UserRoleInput) error
}
