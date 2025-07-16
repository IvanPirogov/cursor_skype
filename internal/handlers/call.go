package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"github.com/gin-gonic/gin"
)

type CallHandler struct {
	// TODO: Add call service dependency
}

func NewCallHandler() *CallHandler {
	return &CallHandler{}
}

func (h *CallHandler) GetCalls(c *gin.Context) {
	// TODO: Implement get calls
	c.JSON(http.StatusOK, gin.H{"calls": []models.Call{}})
}

func (h *CallHandler) InitiateCall(c *gin.Context) {
	// TODO: Implement initiate call
	c.JSON(http.StatusCreated, gin.H{"message": "Call initiated"})
}

func (h *CallHandler) AnswerCall(c *gin.Context) {
	// TODO: Implement answer call
	c.JSON(http.StatusOK, gin.H{"message": "Call answered"})
}

func (h *CallHandler) RejectCall(c *gin.Context) {
	// TODO: Implement reject call
	c.JSON(http.StatusOK, gin.H{"message": "Call rejected"})
}

func (h *CallHandler) EndCall(c *gin.Context) {
	// TODO: Implement end call
	c.JSON(http.StatusOK, gin.H{"message": "Call ended"})
}