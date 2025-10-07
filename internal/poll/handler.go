package poll

import (
	"net/http"
	"voting/notifications/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// validate user logged in
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// fill poll schema
	var req CreateOrPatchPollRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	pollID := uuid.NewString()
	options := make([]Option, len(req.Options))
	for i, o := range req.Options {
		options[i] = Option{
			PollID: pollID,
			Text:   o.Text,
			Votes:  0,
		}
	}
	poll := Poll{
		ID:      pollID,
		Title:   req.Title,
		UserID:  uid.(string),
		Options: options,
	}
	// send to service
	if err := h.service.CreatePoll(&poll); err != nil {
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
	// get user id from cookie
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
	// update poll
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
	pollToUpdate := Poll{
		Title:   req.Title,
		Options: options,
	}
	err = h.service.UpdatePoll(pollID, pollToUpdate)
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
	// get user id from cookie
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
	if err = h.service.DeletePoll(pollID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "poll deleted"})
}

// Vote godoc
// @Summary Vote for poll option
// @Description Allows an authenticated user to vote for a specific option in a poll
// @Tags polls
// @Accept json
// @Produce json
// @Param pollID path string true "Poll ID"
// @Param optionID path string true "option ID"
// @Success 200 {object} map[string]string "vote added"
// @Failure 400 {object} map[string]string "option not found"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "you cant vote at your own polls"
// @Failure 404 {object} map[string]string "poll not found"
// @Router /polls/{pollID}/options/{optionID}/vote [post]
// @Security ApiKeyAuth
func (h *Handler) Vote(c *gin.Context) {
	// get user id from cookie
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
	// validate user cant vote at his own votes
	if poll.UserID == userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cant vote at your own polls"})
		return
	}
	// add vote
	optionID := c.Param("optionID")
	err = h.service.AddVote(pollID, optionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "option not found"})
		return
	}
	// websoket notification
	message := "User " + userid + " voted on your poll '" + poll.Title + "'"
	ws.HubInstance.Notify(poll.UserID, message)

	c.JSON(http.StatusOK, gin.H{"message": "vote added"})
}
