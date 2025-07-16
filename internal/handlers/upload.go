package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	// TODO: Add file service dependency
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	// TODO: Implement file upload
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded"})
}