package web

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type HousekeepingTaskWebViewModel struct {
	ID            uint
	Priority      uint
	PriorityText  string
	PriorityStyle string
	StartedAt     string
	CompletedAt   string
	Note          string
	Assignee      string
	Assigner      string
	RoomFloor     string
	RoomNumber    string
	Status        string
	StatusStyle   string
}

type HousekeepingTasksViewModel struct {
	BaseViewModel
	Tasks      []HousekeepingTaskWebViewModel
	Pagination PaginationInfo
	RoomNumber string
	Assignee   string
	Assigner   string
	Note       string
	Priority   string
	Status     string
	Sort       string
	Filters    string
	Error      string
	Success    string
}

func (h *WebHandler) getHousekeepingTasksViewModel(c fiber.Ctx, userID any) (HousekeepingTasksViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return HousekeepingTasksViewModel{}, err
	}

	roomNumber := c.Query("roomNumber")
	assignee := c.Query("assignee")
	assigner := c.Query("assigner")
	note := c.Query("note")
	priority := c.Query("priority")
	status := c.Query("status")
	sort := c.Query("sort")

	if sort == "" {
		sort = "-id"
	}

	filters := c.Query("filters")
	if filters == "" {
		filters = buildHousekeepingFilters(roomNumber, assignee, assigner, note, priority, status)
	}

	originalFilters := string(c.Request().URI().QueryArgs().Peek("filters"))
	originalSort := string(c.Request().URI().QueryArgs().Peek("sort"))

	if filters != "" {
		c.Request().URI().QueryArgs().Set("filters", filters)
	} else {
		c.Request().URI().QueryArgs().Del("filters")
	}
	c.Request().URI().QueryArgs().Set("sort", sort)

	pageInfo, tasksList := h.housekeepingTaskService.Read(c)

	// Restore original query args
	if originalFilters != "" {
		c.Request().URI().QueryArgs().Set("filters", originalFilters)
	} else {
		c.Request().URI().QueryArgs().Del("filters")
	}
	if originalSort != "" {
		c.Request().URI().QueryArgs().Set("sort", originalSort)
	} else {
		c.Request().URI().QueryArgs().Del("sort")
	}

	webTasks := make([]HousekeepingTaskWebViewModel, len(tasksList))
	for i, t := range tasksList {
		var priorityText, priorityStyle string
		switch t.Priority {
		case 1:
			priorityText = "High"
			priorityStyle = "background-color: rgba(239, 68, 68, 0.08); border-color: rgba(239, 68, 68, 0.15); color: #f87171;"
		case 2:
			priorityText = "Medium"
			priorityStyle = "background-color: rgba(245, 158, 11, 0.08); border-color: rgba(245, 158, 11, 0.15); color: #fbbf24;"
		case 3:
			priorityText = "Low"
			priorityStyle = "background-color: rgba(185, 200, 222, 0.08); border-color: rgba(185, 200, 222, 0.15); color: var(--color-secondary);"
		default:
			priorityText = fmt.Sprintf("Priority %d", t.Priority)
			priorityStyle = "background-color: rgba(185, 200, 222, 0.08); border-color: rgba(185, 200, 222, 0.15); color: var(--color-secondary);"
		}

		var statusText, statusStyle string
		if t.CompletedAt != nil {
			statusText = "Completed"
			statusStyle = "background-color: rgba(16, 185, 129, 0.08); border-color: rgba(16, 185, 129, 0.15); color: #34d399;"
		} else if t.StartedAt != nil {
			statusText = "In Progress"
			statusStyle = "background-color: rgba(59, 130, 246, 0.08); border-color: rgba(59, 130, 246, 0.15); color: #60a5fa;"
		} else {
			statusText = "Pending"
			statusStyle = "background-color: rgba(245, 158, 11, 0.08); border-color: rgba(245, 158, 11, 0.15); color: #fbbf24;"
		}

		startedAtStr := "-"
		if t.StartedAt != nil {
			startedAtStr = t.StartedAt.Format("2006-01-02 15:04:05")
		}

		completedAtStr := "-"
		if t.CompletedAt != nil {
			completedAtStr = t.CompletedAt.Format("2006-01-02 15:04:05")
		}

		webTasks[i] = HousekeepingTaskWebViewModel{
			ID:            t.ID,
			Priority:      t.Priority,
			PriorityText:  priorityText,
			PriorityStyle: priorityStyle,
			StartedAt:     startedAtStr,
			CompletedAt:   completedAtStr,
			Note:          t.Note,
			Assignee:      t.Assignee,
			Assigner:      t.Assigner,
			RoomFloor:     t.RoomFloor,
			RoomNumber:    t.RoomNumber,
			Status:        statusText,
			StatusStyle:   statusStyle,
		}
	}

	prevPage := pageInfo.Page - 1
	if prevPage < 0 {
		prevPage = 0
	}

	pagination := PaginationInfo{
		Page:        pageInfo.Page,
		CurrentPage: pageInfo.Page + 1,
		Size:        pageInfo.Size,
		TotalPages:  pageInfo.TotalPages,
		Total:       pageInfo.Total,
		Last:        pageInfo.Last,
		First:       pageInfo.First,
		NextPage:    pageInfo.Page + 1,
		PrevPage:    prevPage,
	}

	return HousekeepingTasksViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Housekeeping Tasks Directory - Gotel",
			ActiveTab:   "housekeeping_tasks",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Tasks:      webTasks,
		Pagination: pagination,
		RoomNumber: roomNumber,
		Assignee:   assignee,
		Assigner:   assigner,
		Note:       note,
		Priority:   priority,
		Status:     status,
		Sort:       sort,
		Filters:    filters,
	}, nil
}

func (h *WebHandler) HousekeepingTasksGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getHousekeepingTasksViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "housekeeping_tasks", vm)
}

func buildHousekeepingFilters(roomNumber, assignee, assigner, note, priority, status string) string {
	var parts []string

	if roomNumber != "" {
		parts = append(parts, fmt.Sprintf(`["room_number", "like", "%s"]`, roomNumber))
	}
	if assignee != "" {
		parts = append(parts, fmt.Sprintf(`["assignee", "like", "%s"]`, assignee))
	}
	if assigner != "" {
		parts = append(parts, fmt.Sprintf(`["assigner", "like", "%s"]`, assigner))
	}
	if note != "" {
		parts = append(parts, fmt.Sprintf(`["note", "like", "%s"]`, note))
	}
	if priority != "" {
		parts = append(parts, fmt.Sprintf(`["priority", "=", %s]`, priority))
	}
	if status != "" {
		if status == "pending" {
			parts = append(parts, `["started_at", "is", null]`)
		} else if status == "in_progress" {
			parts = append(parts, `["started_at", "is not", null]`, `["completed_at", "is", null]`)
		} else if status == "completed" {
			parts = append(parts, `["completed_at", "is not", null]`)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	var result bytes.Buffer
	result.WriteString("[")
	for i, part := range parts {
		if i > 0 {
			result.WriteString(`, ["and"], `)
		}
		result.WriteString(part)
	}
	result.WriteString("]")
	return result.String()
}
