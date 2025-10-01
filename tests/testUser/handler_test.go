package testUser

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"voting/internal/user/crud"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupHandler(t *testing.T) crud.Handler {
	service := setupService(t)
	return *crud.NewUserHandler(service)
}

func setupRouter(t *testing.T) *gin.Engine {
	h := setupHandler(t)
	r := gin.Default()
	r.POST("/users", h.PostUser)
	r.GET("/users", h.GetUsers)
	r.GET("/users/:id", h.GetUser)
	r.PATCH("/users/:id", h.PatchUser)
	r.DELETE("/users/:id", h.DeleteUser)
	return r
}

func TestPostUserHandler(t *testing.T) {
	router := setupRouter(t)
	body := []byte(`{"username": "alice"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"user created"`)
}

func TestGetUsersHandler(t *testing.T) {
	router := setupRouter(t)

	// create first user
	createBody := []byte(`{"username": "alice"}`)
	createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// create second user
	createBody = []byte(`{"username": "bob"}`)
	createReq, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW = httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then get users
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "alice")
	require.Contains(t, w.Body.String(), "bob")
}

func TestGetUserHandler(t *testing.T) {
	router := setupRouter(t)

	// first create
	createBody := []byte(`{"username": "charlie"}`)
	createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "charlie")
}

func TestPatchUserHandler(t *testing.T) {
	router := setupRouter(t)

	// first create
	createBody := []byte(`{"username": "oldname"}`)
	createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then patch
	patchBody := []byte(`{"username": "newname"}`)
	patchReq, _ := http.NewRequest("PATCH", "/users/1", bytes.NewBuffer(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()
	router.ServeHTTP(patchW, patchReq)
	require.Equal(t, http.StatusOK, patchW.Code)

	// then check
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "newname")
}

func TestDeleteUserHandler(t *testing.T) {
	router := setupRouter(t)

	// first create
	createBody := []byte(`{"username": "todelete"}`)
	createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then delete
	deleteReq, _ := http.NewRequest("DELETE", "/users/1", nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	require.Equal(t, http.StatusOK, deleteW.Code)

	// then check
	getReq, _ := http.NewRequest("GET", "/users/1", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	require.Equal(t, http.StatusBadRequest, getW.Code)
}
