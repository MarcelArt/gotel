package handlers

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type HousekeepingTaskHandler struct {
	service services.IHousekeepingTaskService
}

var _ = entities.HousekeepingTask{}

func NewHousekeepingTaskHandler(service services.IHousekeepingTaskService) *HousekeepingTaskHandler {
	return &HousekeepingTaskHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new housekeeping task
// @Description  Create a new housekeeping task with the provided details
// @Tags         housekeeping-tasks
// @Accept       json
// @Produce      json
// @Param        housekeepingTask  body      models.HousekeepingTaskInput  true  "Housekeeping task details"
// @Success      201               {object}  common.JSONResponse{items=uint}
// @Failure      400               {object}  common.JSONResponse
// @Failure      500               {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/housekeeping-tasks [post]
func (h *HousekeepingTaskHandler) Create(c fiber.Ctx) error {
	var task models.HousekeepingTaskInput
	if err := c.Bind().JSON(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, task)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating housekeeping task"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Housekeeping task created"))
}

// Read godoc
// @Summary      List housekeeping tasks
// @Description  Get a paginated list of housekeeping tasks
// @Tags         housekeeping-tasks
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.HousekeepingTaskPage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/housekeeping-tasks [get]
func (h *HousekeepingTaskHandler) Read(c fiber.Ctx) error {
	tasks, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(tasks)
}

// Update godoc
// @Summary      Update housekeeping task
// @Description  Update an existing housekeeping task's details
// @Tags         housekeeping-tasks
// @Accept       json
// @Produce      json
// @Param        id                path      string                        true  "Housekeeping Task ID"
// @Param        housekeepingTask  body      models.HousekeepingTaskInput  true  "Updated housekeeping task details"
// @Success      200               {object}  common.JSONResponse
// @Failure      400               {object}  common.JSONResponse
// @Failure      401               {object}  common.JSONResponse
// @Failure      500               {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/housekeeping-tasks/{id} [put]
func (h *HousekeepingTaskHandler) Update(c fiber.Ctx) error {
	var task models.HousekeepingTaskInput
	if err := c.Bind().JSON(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), task); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating housekeeping task"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Housekeeping task updated"))
}

// Delete godoc
// @Summary      Delete housekeeping task
// @Description  Delete a housekeeping task by ID
// @Tags         housekeeping-tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Housekeeping Task ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/housekeeping-tasks/{id} [delete]
func (h *HousekeepingTaskHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting housekeeping task"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Housekeeping task deleted"))
}

// GetByID godoc
// @Summary      Get housekeeping task by ID
// @Description  Get detailed information about a housekeeping task by its ID
// @Tags         housekeeping-tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Housekeeping Task ID"
// @Success      200  {object}  common.JSONResponse{items=entities.HousekeepingTask}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/housekeeping-tasks/{id} [get]
func (h *HousekeepingTaskHandler) GetByID(c fiber.Ctx) error {
	task, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting housekeeping task"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(task, "Housekeeping task found"))
}

func (h *HousekeepingTaskHandler) SetupRoutes(v1 fiber.Router) {
	tasks := v1.Group("/housekeeping-tasks")

	tasks.Use(middlewares.Authn())

	tasks.Post("/", middlewares.Authz("housekeepingTasks#create"), h.Create)
	tasks.Get("/", middlewares.Authz("housekeepingTasks#read"), h.Read)
	tasks.Get("/:id", middlewares.Authz("housekeepingTasks#read"), h.GetByID)
	tasks.Put("/:id", middlewares.Authz("housekeepingTasks#update"), h.Update)
	tasks.Delete("/:id", middlewares.Authz("housekeepingTasks#delete"), h.Delete)
}
