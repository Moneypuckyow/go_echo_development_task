package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieve list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {array} User
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users [get]
func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve a single user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} User
// @Failure      404 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *Handler) GetUserByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create new user
// @Description  Create a new user with provided information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body User true "User creation data"
// @Success      201 {object} User
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
	}
	if err := c.Validate(&u); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	if err := h.service.Create(u); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, u)
}

// UpdateUserFull godoc
// @Summary      Update user (PUT)
// @Description  Replace all user fields with new values
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        user body User true "Complete user data"
// @Success      200 {object} User
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *Handler) UpdateUserFull(c echo.Context) error {
	return h.UpdateUser(c)
}

// UpdateUserPartial godoc
// @Summary      Update user (PATCH)
// @Description  Partially update user fields
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        user body UserUpdate true "Partial user data"
// @Success      200 {object} UserUpdate
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [patch]
func (h *Handler) UpdateUserPartial(c echo.Context) error {
	return h.UpdateUser(c)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var u UserUpdate
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
	}

	// Kalau PUT → semua field harus ada (validate required)
	if c.Request().Method == http.MethodPut {
		if err := c.Validate(&u); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
	}

	u.ID = id

	var rows int64
	var err error
	if c.Request().Method == http.MethodPut {
		rows, err = h.service.UpdateFull(u)
	} else {
		rows, err = h.service.UpdatePartial(u)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}

	return c.JSON(http.StatusOK, u)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} SuccessResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, err := h.service.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "User deleted successfully"})
}

func RegisterRoutes(g *echo.Group, db interface {
	Query(string, ...any) (*sql.Rows, error)
}) {
	repo := NewRepository(db.(*sql.DB))
	service := NewService(repo)
	h := NewHandler(service)

	g.GET("users", h.GetUsers)
	g.GET("users/:id", h.GetUserByID)
	g.POST("users", h.CreateUser)
	g.PATCH("users/:id", h.UpdateUserPartial)
	g.PUT("users/:id", h.UpdateUserFull)
	g.DELETE("users/:id", h.DeleteUser)
}
