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

func (h *Handler) PostUser(c *gin.Context) {
	var req = User{}
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

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, user)
}

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

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
