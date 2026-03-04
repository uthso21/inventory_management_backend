package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

type UserHandler struct {
	userService usecases.UserService
}

func NewUserHandler(userService usecases.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
users, err := h.userService.ListUsers(c.Request().Context())
if err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
var user entities.User
if err := c.Bind(&user); err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
}

if err := h.userService.CreateUser(c.Request().Context(), &user); err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
}

user, err := h.userService.GetUser(c.Request().Context(), id)
if err != nil {
return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
}
return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
var user entities.User
if err := c.Bind(&user); err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
}

if err := h.userService.UpdateUser(c.Request().Context(), &user); err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
}

if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
return c.NoContent(http.StatusNoContent)
}
