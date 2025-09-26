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
	r.GET("/polls/:id", h.GetPoll)
	r.PATCH("/polls/:id", h.PatchPoll)
	r.DELETE("/polls/:id", h.DeletePoll)
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

	// first create
	createBody := []byte(`{
		"title": "test poll 1",
		"options": [{"text": "a"}, {"text": "b"}]
	}`)
	createReq, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	createBody = []byte(`{
		"title": "test poll 2",
		"options": [{"text": "c"}, {"text": "d"}]
	}`)
	createReq, _ = http.NewRequest("POST", "/polls", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW = httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/polls", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "test poll 1")
	require.Contains(t, w.Body.String(), "\"text\":\"a\"")
	require.Contains(t, w.Body.String(), "test poll 2")
	require.Contains(t, w.Body.String(), "\"text\":\"d\"")
}

func TestGetPollHandler(t *testing.T) {
	router := setupRouter(t)

	// first create
	createBody := []byte(`{
		"title": "test poll",
		"options": [{"text": "a"}, {"text": "b"}]
	}`)
	createReq, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/polls/1", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "test poll")
	require.Contains(t, w.Body.String(), "\"text\":\"a\"")
	require.Contains(t, w.Body.String(), "\"text\":\"b\"")
}

func TestPatchPollHandler(t *testing.T) {
	router := setupRouter(t)

	// first create
	createBody := []byte(`{
		"title": "old title",
		"options": [{"text": "old a"}, {"text": "old b"}]
	}`)
	createReq, _ := http.NewRequest("POST", "/polls", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	// then patch
	patchBody := []byte(`{
		"title": "updated title",
		"options": [{"text": "updated a"}, {"text": "updated b"}, {"text": "new c"}]
	}`)
	patchReq, _ := http.NewRequest("PATCH", "/polls/1", bytes.NewBuffer(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchW := httptest.NewRecorder()
	router.ServeHTTP(patchW, patchReq)
	require.Equal(t, http.StatusOK, patchW.Code)

	// then check
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/polls/1", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "updated title")
	require.Contains(t, w.Body.String(), "updated a")
	require.Contains(t, w.Body.String(), "new c")
}

func TestDeletePollHandler(t *testing.T) {

}
