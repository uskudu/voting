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

func (h *Handler) PostPoll(c *gin.Context) {
	var req = Poll{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.service.CreatePoll(req.Title, req.Options); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed creating testPoll": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "testPoll created"})
}

func (h *Handler) GetPolls(c *gin.Context) {
	polls, err := h.service.GetPolls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get polls"})
		return
	}
	c.JSON(http.StatusOK, polls)
}

func (h *Handler) GetPoll(c *gin.Context) {
	id := c.Param("id")
	poll, err := h.service.GetPollByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, poll)
}

func (h *Handler) PatchPoll(c *gin.Context) {
	id := c.Param("id")
	var poll Poll
	if err := c.Bind(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.service.UpdatePoll(id, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update testPoll"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "testPoll updated"})

}

func (h *Handler) DeletePoll(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePoll(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "testPoll deleted"})
}
