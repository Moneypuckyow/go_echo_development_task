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

func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := c.Validate(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := h.service.Create(u); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Buat struct patch (pakai pointer biar bisa deteksi nil)
	var u UserUpdate
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Kalau PUT → semua field harus ada (validate required)
	if c.Request().Method == http.MethodPut {
		if err := c.Validate(&u); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	u.ID = id

	// Service layer yang handle logic beda PUT vs PATCH
	var rows int64
	var err error
	if c.Request().Method == http.MethodPut {
		rows, err = h.service.UpdateFull(u) // replace full
	} else {
		rows, err = h.service.UpdatePartial(u) // update sebagian
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, u)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, err := h.service.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}

func RegisterRoutes(g *echo.Group, db interface {
	Query(string, ...any) (*sql.Rows, error)
}) {
	// Now the type assertion should work
	repo := NewRepository(db.(*sql.DB))
	service := NewService(repo)
	h := NewHandler(service)

	g.GET("users", h.GetUsers)
	g.GET("users/:id", h.GetUserByID)
	g.POST("users", h.CreateUser)
	g.PATCH("users/:id", h.UpdateUser)
	g.PUT("users/:id", h.UpdateUser)
	g.DELETE("users/:id", h.DeleteUser)
}
