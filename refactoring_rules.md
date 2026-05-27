# Instructions and Rules: Refactoring to Clean Layered Architecture (v1 API)

This document provides a set of architectural rules, code templates, and a step-by-step plan to refactor your Go backend project to follow a clean, multi-layered architecture based on **Fiber v3** and **GORM**. It mirrors the design pattern of the `kas-bon-v2` reference project.

> [!NOTE]
> This architecture does NOT require Casbin. Custom middleware or other lightweight authentication strategies can be placed in `internal/v1/middlewares/` as needed.

---

## 1. Directory Structure

All new API v1 logic must be isolated within the `internal/v1` directory, divided into distinct, decoupled responsibility layers:

```
internal/
├── common/                  # Shared utilities (JSON responses, token helpers)
└── v1/
    ├── models/              # GORM database models & request payload DTOs
    ├── repositories/        # Database CRUD queries (interface + struct)
    ├── services/            # Business logic orchestration (interface + struct)
    ├── usecases/            # Cross-repository workflows / transaction controllers
    ├── handlers/            # HTTP request binders & response routers (Fiber handlers)
    ├── routes/              # Routing groups configuration
    └── middlewares/         # Authn, CORS, logging, or other HTTP middlewares
```

---

## 2. Core Architecture Rules

### Rule 1: Strict Layer Separation and Dependency Inversion
- **Dependency Flow**: `routes -> handlers -> services -> repositories -> models/database`.
- Code at a higher level must access lower layers **only through Interfaces**.
  - A handler holds an interface of a service: `svc services.IUserService`.
  - A service holds an interface of a repository: `repo repositories.IUserRepository`.
- High-level orchestrations (multiple tables/repos, transactions) must live in **Usecases** or **Services**, not in Handlers.

### Rule 2: DB Connection and Context Flow
- A `*gorm.DB` instance (or a pointer to a transaction `*gorm.DB`) must be passed from the router/service down to the repositories or usecases during initialization.
- For queries requiring HTTP context (e.g. pagination query parsing), pass `fiber.Ctx` down.

### Rule 3: Uniform API Responses
- All JSON API responses (except for paginated listings) must be wrapped inside `common.JSONResponse` to ensure consistent success flags and client messaging:
  ```go
  type JSONResponse struct {
      Items     any    `json:"items"`
      IsSuccess bool   `json:"isSuccess"`
      Message   string `json:"message"`
  }
  ```

---

## 3. Code Patterns & Templates

### A. Shared Utilities (`internal/common`)
Create a shared response formatter.

**File:** `internal/common/json_response.go`
```go
package common

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type JSONResponse struct {
	Items     any    `json:"items"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// NewJSONResponse returns a wrapped API response.
// If the items argument is an error, it is treated as a failed action.
func NewJSONResponse(items any, message string) *JSONResponse {
	err, ok := items.(error)
	if ok {
		if message == "" {
			message = err.Error()
		} else {
			message = fmt.Sprintf("%s: %s", message, err.Error())
		}
		return &JSONResponse{
			Items:     nil,
			IsSuccess: false,
			Message:   message,
		}
	}

	return &JSONResponse{
		Items:     items,
		IsSuccess: true,
		Message:   message,
	}
}

// StatusCodeFromError resolves standard database error instances to HTTP status codes.
func StatusCodeFromError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.StatusNotFound
	}
	return fiber.StatusInternalServerError
}
```

---

### B. Models (`internal/v1/models`)
- Separate DB models from request payloads (Inputs).
- Embed a base `Input` helper in request payloads to suppress GORM metadata fields during JSON serialization.

**File:** `internal/v1/models/input.model.go`
```go
package models

import "time"

// Input serves as a base DTO that exposes ID, CreatedAt, and UpdatedAt 
// to GORM while ignoring them in JSON deserialization.
type Input struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
```

**File:** `internal/v1/models/app.model.go`
```go
package models

import "gorm.io/gorm"

const appTableName = "apps"

// App is the database schema model
type App struct {
	gorm.Model
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

// AppInput is the schema model for creating/binding new resources
type AppInput struct {
	Input
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

// TableName matches GORM conventions
func (AppInput) TableName() string {
	return appTableName
}
```

---

### C. Repositories (`internal/v1/repositories`)
- Expose an interface `IResourceRepo` and a struct implementation `ResourceRepo`.
- If pagination is needed, use `github.com/morkid/paginate`.

**File:** `internal/v1/repositories/app.repo.go`
```go
package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IAppRepo interface {
	Create(app models.AppInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.App)
	Update(id any, app models.App) error
	Delete(id any) error
	GetByID(id any) (models.App, error)
}

type AppRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewAppRepo(db *gorm.DB) *AppRepo {
	return &AppRepo{
		db:        db,
		pageQuery: `select * from apps where deleted_at is null`,
	}
}

func (r *AppRepo) Create(app models.AppInput) (uint, error) {
	err := r.db.Create(&app).Error
	return app.ID, err
}

func (r *AppRepo) Read(c fiber.Ctx) (paginate.Page, []models.App) {
	var apps []models.App
	pg := paginate.New()
	stmt := r.db.Raw(r.pageQuery)
	page := pg.With(stmt).Request(c.Request()).Response(&apps)
	return page, apps
}

func (r *AppRepo) Update(id any, app models.App) error {
	return r.db.Model(&models.App{}).Where("id = ?", id).Updates(&app).Error
}

func (r *AppRepo) Delete(id any) error {
	return r.db.Delete(&models.App{}, id).Error
}

func (r *AppRepo) GetByID(id any) (models.App, error) {
	var app models.App
	err := r.db.Where("id = ?", id).First(&app).Error
	return app, err
}
```

---

### D. Services (`internal/v1/services`)
- Handle business logic validation, transactions initialization, and repo execution.
- Accept interfaces in the constructor `NewAppService(repo repositories.IAppRepo)`.

**File:** `internal/v1/services/app.service.go`
```go
package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IAppService interface {
	Create(app models.AppInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.App)
	Update(id any, app models.App) error
	Delete(id any) error
	GetByID(id any) (models.App, error)
}

type AppService struct {
	repo repositories.IAppRepo
}

func NewAppService(repo repositories.IAppRepo) *AppService {
	return &AppService{repo: repo}
}

func (s *AppService) Create(app models.AppInput) (uint, error) {
	return s.repo.Create(app)
}

func (s *AppService) Read(c fiber.Ctx) (paginate.Page, []models.App) {
	return s.repo.Read(c)
}

func (s *AppService) Update(id any, app models.App) error {
	return s.repo.Update(id, app)
}

func (s *AppService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *AppService) GetByID(id any) (models.App, error) {
	return s.repo.GetByID(id)
}
```

---

### E. Usecases (`internal/v1/usecases`)
- Usecases represent transactional boundaries or coordinate operations between multiple domains.
- They are initialized with a transaction DB reference: `InitRegisterUserUsecase(tx *gorm.DB)`.

**File:** `internal/v1/usecases/register_user.usecase.go`
```go
package usecases

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"gorm.io/gorm"
)

type RegisterUserUsecase struct {
	User models.UserInput

	uRepo repositories.IUserRepo
	dRepo repositories.IDomainRepo
}

func InitRegisterUserUsecase(tx *gorm.DB) *RegisterUserUsecase {
	return &RegisterUserUsecase{
		uRepo: repositories.NewUserRepo(tx),
		dRepo: repositories.NewDomainRepo(tx),
	}
}

func (u *RegisterUserUsecase) Execute() (uint, error) {
	// 1. Business Logic / Multi-Repository validations
	// 2. Write operations
	id, err := u.uRepo.Create(u.User)
	if err != nil {
		return 0, fmt.Errorf("failed creating user: %w", err)
	}

	return id, nil
}
```

**Service Invocation pattern:**
```go
func (s *UserService) Create(user models.UserInput) (uint, error) {
	tx := s.db.Begin() // Starts transaction
	defer tx.Rollback()

	registerUser := usecases.InitRegisterUserUsecase(tx)
	registerUser.User = user
	id, err := registerUser.Execute()
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return id, nil
}
```

---

### F. Handlers (`internal/v1/handlers`)
- Rely on Service interfaces.
- Handle HTTP binding/deserialization via `c.Bind().JSON()`.
- Return HTTP responses utilizing `common.JSONResponse`.

**File:** `internal/v1/handlers/app.handler.go`
```go
package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type AppHandler struct {
	svc services.IAppService
}

func NewAppHandler(svc services.IAppService) *AppHandler {
	return &AppHandler{svc: svc}
}

func (h *AppHandler) Create(c fiber.Ctx) error {
	var app models.AppInput
	if err := c.Bind().JSON(&app); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.svc.Create(app)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating app"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "App created"))
}

func (h *AppHandler) Read(c fiber.Ctx) error {
	page, _ := h.svc.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

func (h *AppHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	app, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting app"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(app, "App found"))
}
```

---

### G. Routes (`internal/v1/routes`)
- Separate route groupings per resource.
- Wire dependencies and middlewares inside route files.
- Wire master route setup in `routes.go`.

**File:** `internal/v1/routes/app.route.go`
```go
package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupAppRoutes(v1 fiber.Router) {
	apps := v1.Group("/apps")

	h := handlers.NewAppHandler(services.NewAppService(repositories.NewAppRepo(configs.DB)))

	// Example router setup:
	apps.Get("/", middlewares.Authn(), h.Read)
	apps.Get("/:id", middlewares.Authn(), h.GetByID)
	apps.Post("/", middlewares.Authn(), h.Create)
}
```

**File:** `internal/v1/routes/routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	// Setup individual routes
	SetupAppRoutes(v1)
	// SetupUserRoutes(v1)
	// etc.
}
```

---

## 4. Step-by-Step Refactoring Plan

Follow these steps when migrating existing endpoints to the new structure:

1. **Setup Core Directories**:
   - Create directories under `internal/v1/`: `models/`, `repositories/`, `services/`, `usecases/`, `handlers/`, `routes/`, `middlewares/`.
   - Setup `internal/common/json_response.go`.

2. **Migrate Domain Models**:
   - Move or redefine DB entities in `models/`.
   - Create the `<Resource>Input` struct (embedding `Input` base) for each corresponding request payload.

3. **Migrate DB Operations (Repositories)**:
   - Extract raw SQL / GORM query strings from handlers or controllers.
   - Design the `IRepo` interface with DB methods.
   - Implement the repo struct wrapping `*gorm.DB`.

4. **Define Service Layer**:
   - Create business logic processes in `services/`.
   - Implement constructor injection to accept repositories.
   - Identify database-heavy processes that span multiple actions or tables. Migrate them into GORM Transaction wrapper blocks calling `usecases.Init...`.

5. **Construct Handlers**:
   - Write handlers matching the resource. Bind incoming JSON into `models.<Input>` structs.
   - Invoke Service methods, wrap outcomes with `common.NewJSONResponse(...)` or directly output listings.

6. **Establish Routes and Bind to main.go**:
   - Map routes in `internal/v1/routes/`.
   - Setup routing dependencies, endpoints, and authentication middleware.
   - Mount `routes.SetupRoutes(api)` to your main server router in `main.go`.
