package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceIface
}

func NewUserHandler(s ServiceIface) *Handler {
	return &Handler{service: s}
}

// PostUser godoc
// @Summary Create a new user
// @Description Create a new user with username
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User input"
// @Success 200 {object} map[string]string "user created"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed creating user"
// @Router /users [post]
func (h *Handler) PostUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.service.CreateUser(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed creating user": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} map[string]string "could not get users"
// @Router /users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a single user by its ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string "invalid request"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// PatchUser godoc
// @Summary Update a user
// @Description Update username by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body User true "User object"
// @Success 200 {object} map[string]string "user updated"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "could not update user"
// @Router /users/{id} [patch]
func (h *Handler) PatchUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.service.UpdateUser(id, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "user deleted"
// @Failure 400 {object} map[string]string "invalid request"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
