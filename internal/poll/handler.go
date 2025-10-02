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
// @Param poll body poll.CreateOrPatchPollRequest true "Poll user input"
// @Success 200 {object} map[string]string "poll created"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 500 {object} map[string]string "failed creating poll"
// @Router /polls [post]
// @Security ApiKeyAuth
func (h *Handler) PostPoll(c *gin.Context) {
	// fill poll schema
	var req CreateOrPatchPollRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	options := make([]Option, len(req.Options))
	for i, o := range req.Options {
		options[i] = Option{
			Text: o.Text,
		}
	}
	// get user id
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// send to service
	if err := h.service.CreatePoll(uid.(string), req.Title, options); err != nil {
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
// @Failure 404 {object} map[string]string "poll not found"
// @Router /polls/{id} [get]
func (h *Handler) GetPoll(c *gin.Context) {
	id := c.Param("id")
	poll, err := h.service.GetPollByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
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
// @Param poll body poll.CreateOrPatchPollRequest true "CreateOrPatchPollRequest user input"
// @Success 200 {object} map[string]string "poll updated"
// @Failure 404 {object} map[string]string "poll not found"
// @Failure 500 {object} map[string]string "could not update poll"
// @Router /polls/{id} [patch]
// @Security ApiKeyAuth
func (h *Handler) PatchPoll(c *gin.Context) {
	id := c.Param("id")
	var req CreateOrPatchPollRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	options := make([]Option, len(req.Options))
	for i, o := range req.Options {
		options[i] = Option{
			Text: o.Text,
		}
	}

	pollToUpdate := Poll{
		Title:   req.Title,
		Options: options,
	}
	// get user id
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := uid.(string)

	// get poll form db
	pollID := c.Param("id")
	poll, err := h.service.GetPollByID(pollID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	// validate that user which is trying to update is the creator of the poll
	if poll.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this poll"})
		return
	}

	err = h.service.UpdatePoll(id, pollToUpdate)
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
// @Failure 404 {object} map[string]string "poll not found"
// @Router /polls/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) DeletePoll(c *gin.Context) {
	id := c.Param("id")
	// get user id
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := uid.(string)

	// get poll form db
	pollID := c.Param("id")
	poll, err := h.service.GetPollByID(pollID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	// validate that user which is trying to delete is the creator of the poll
	if poll.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this poll"})
		return
	}
	if err := h.service.DeletePoll(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "poll deleted"})
}
