package handler

import (
	"fmt"
	"food/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload Multiple Files godoc
// @ID           upload_multiple_files
// @Router       /food/api/v1/uploadfiles [POST]
// @Summary      Upload Multiple Files
// @Description  Upload Multiple Files
// @Tags         Upload File
// @Accept       multipart/form-data
// @Procedure    json
// @Param        file formData []file true "File to upload"
// @Success      200 {object} Response{data=string} "Success Request"
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.log.Error("Multipart form error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	resp, err := helper.UploadFiles(form)
	if err != nil {
		h.log.Error("Upload error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload files"})
		return
	}

	// Log the generated URLs
	for _, url := range resp.Url {
		fmt.Printf("Uploaded file ID: %s, URL: %s\n", url.Id, url.Url)
	}

	c.JSON(http.StatusOK, resp)
}


// delete file godoc
// @ID           delete_file
// @Router       /food/api/v1/deletefile [DELETE]
// @Summary      Delete File
// @Description  Delete File
// @Tags         Upload File
// @Accept       multipart/form-data
// @Procedure    json
// @Param        id query string true "id"
// @Success      200 {object} Response{data=string} "Success Request"
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteFile(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file ID"})
		return
	}

	err := helper.DeleteFile(id)
	if err != nil {
		h.log.Error("Delete error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "File deleted successfully"})
}
