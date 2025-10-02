package user

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
// @Param user body user.CreateOrLoginUserRequest true "User input"
// @Success 200 {object} map[string]string "user created"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed creating user"
// @Router /users [post]
func (h *Handler) PostUser(c *gin.Context) {
	var req CreateOrLoginUserRequest
	if c.Bind(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.service.CreateUser(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed creating user": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

// Login godoc
// @Summary Authenticate user and create JWT
// @Description Authenticates a user by username and returns a JWT token in a cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.CreateOrLoginUserRequest true "Login request payload"
// @Success 200 {object} map[string]string "login successful"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed to create jwt token"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req CreateOrLoginUserRequest
	if c.Bind(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	got, err := h.service.Authenticate(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// set jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": got.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create jwt token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*42*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

// Validate godoc
// @Summary      Validate JWT and return user info
// @Description  Returns the authenticated user info extracted from the JWT token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "User info"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Router       /validate [get]
// @Security     ApiKeyAuth
func (h *Handler) Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": user})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Tags users
// @Produce json
// @Success 200 {array} user.User
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
// @Success 200 {object} user.User
// @Failure 404 {object} map[string]string "user not found"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	got, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, got)
}

// PatchUser godoc
// @Summary Update a user
// @Description Update username by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.PatchUserRequest true "PatchUserRequest user input"
// @Success 200 {object} map[string]string "user updated"
// @Failure 404 {object} map[string]string "user not found"
// @Failure 500 {object} map[string]string "could not update user"
// @Router /users/{id} [patch]
// @Security ApiKeyAuth
func (h *Handler) PatchUser(c *gin.Context) {
	id := c.Param("id")
	var req PatchUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// get user id
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := uid.(string)

	// get user form db
	usr, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// validate that user which is trying to update is the owner of the account
	if usr.ID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this account"})
		return
	}

	err = h.service.UpdateUser(id, req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
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
// @Failure 404 {object} map[string]string "user not found"
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	// get user id
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := uid.(string)

	// get user form db
	usr, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// validate that user which is trying to update is the owner of the account
	if usr.ID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this account"})
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
