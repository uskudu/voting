package poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceIface
}

func NewPollHandler(s ServiceIface) *Handler {
	return &Handler{service: s}
}

// PostPoll godoc
// @Summary Create a new poll
// @Description Create a new poll with title and options
// @Tags polls
// @Accept json
// @Produce json
// @Param poll body poll.CreatePollRequest true "Poll user input"
// @Success 200 {object} map[string]string "poll created"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed creating poll"
// @Router /polls [post]
func (h *Handler) PostPoll(c *gin.Context) {
	var req CreatePollRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.service.CreatePoll(req.Title, req.Options); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed creating poll": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "poll created"})
}

// GetPolls godoc
// @Summary Get all polls
// @Description Retrieve all polls
// @Tags polls
// @Produce json
// @Success 200 {array} Poll
// @Failure 500 {object} map[string]string "could not get polls"
// @Router /polls [get]
func (h *Handler) GetPolls(c *gin.Context) {
	polls, err := h.service.GetPolls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get polls"})
		return
	}
	c.JSON(http.StatusOK, polls)
}

// GetPoll godoc
// @Summary Get poll by ID
// @Description Retrieve a single poll by its ID
// @Tags polls
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} Poll
// @Failure 400 {object} map[string]string "invalid request"
// @Router /polls/{id} [get]
func (h *Handler) GetPoll(c *gin.Context) {
	id := c.Param("id")
	poll, err := h.service.GetPollByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, poll)
}

// PatchPoll godoc
// @Summary Update a poll
// @Description Update poll title or options by ID
// @Tags polls
// @Accept json
// @Produce json
// @Param id path string true "Poll ID"
// @Param poll body Poll true "Poll object"
// @Success 200 {object} map[string]string "poll updated"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "could not update poll"
// @Router /polls/{id} [patch]
func (h *Handler) PatchPoll(c *gin.Context) {
	id := c.Param("id")
	var poll Poll
	if err := c.Bind(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.service.UpdatePoll(id, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update poll"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "poll updated"})

}

// DeletePoll godoc
// @Summary Delete a poll
// @Description Delete poll by ID
// @Tags polls
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} map[string]string "poll deleted"
// @Failure 400 {object} map[string]string "invalid request"
// @Router /polls/{id} [delete]
func (h *Handler) DeletePoll(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePoll(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "poll deleted"})
}
