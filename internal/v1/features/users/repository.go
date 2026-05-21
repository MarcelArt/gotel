package users

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserRepo interface {
	common.IBaseCrudRepo[entities.User, UserInput, UserPage]
	GetByUsernameOrEmail(c common.Context, usernameOrEmail string) (entities.User, error)
}

type UserRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IUserRepo = &UserRepo{}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db:        db,
		pageQuery: "SELECT id, username, email FROM users",
	}
}

func (r *UserRepo) Create(c common.Context, input UserInput) (uint, error) {
	ctx := c.Context()

	user, err := common.Cast[entities.User](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}
	user.Password = input.Password

	err = gorm.G[entities.User](r.db).Create(ctx, &user)

	return user.ID, err
}

func (r *UserRepo) Read(c fiber.Ctx) (paginate.Page, []UserPage) {
	users := make([]UserPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&users)

	return page, users
}

func (r *UserRepo) Update(c common.Context, id any, input UserInput) error {
	ctx := c.Context()
	user, err := common.Cast[entities.User](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.User](r.db).Where("id = ?", id).Updates(ctx, user)

	return err
}

func (r *UserRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.User](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *UserRepo) GetByID(c common.Context, id any) (entities.User, error) {
	var user entities.User
	ctx := c.Context()

	user, err := gorm.G[entities.User](r.db).Where("id = ?", id).First(ctx)

	return user, err
}

func (r *UserRepo) GetByUsernameOrEmail(c common.Context, usernameOrEmail string) (entities.User, error) {
	ctx := c.Context()

	return gorm.G[entities.User](r.db).Where("username = $1 or email = $1", usernameOrEmail).First(ctx)
}
