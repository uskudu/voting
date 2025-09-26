package testPoll

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"voting/internal/poll"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupHandler(t *testing.T) poll.Handler {
	service := setupService(t)
	return *poll.NewPollHandler(service)
}

func setupRouter(t *testing.T) *gin.Engine {
	h := setupHandler(t)
	r := gin.Default()
	r.POST("/polls", h.PostPoll)
	r.GET("/polls", h.GetPolls)
	r.GET("/poll/:id", h.GetPoll)
	r.PATCH("/poll/:id", h.PatchPoll)
	r.DELETE("/poll/:id", h.DeletePoll)
	return r
}

func TestPostPollHandler(t *testing.T) {
	router := setupRouter(t)
	body := []byte(`{
		"title": "test poll",
		"options": [
			{"text": "a"},
			{"text": "b"}
		]
	}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"poll created"`)
}

func TestGetPollsHandler(t *testing.T) {
	router := setupRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/polls", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "[]")
}

func TestGetPollHandler(t *testing.T) {

}
func TestPatchPollHandler(t *testing.T) {

}
func TestDeletePollHandler(t *testing.T) {

}
