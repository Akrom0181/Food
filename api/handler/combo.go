package handler

import (
	"encoding/json"
	"food/api/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)
// Create Combo godoc
// @ID          create_combo
// @Router      /food/api/v1/combo [POST]
// @Summary     Create Combo
// @Description Create a new combo with a set of items
// @Tags        combo
// @Accept      json
// @Produce     json
// @Param       Combo body models.SwaggerComboCreateRequest true "CreateComboRequest"
// @Success     201 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *Handler) CreateCombo(c *gin.Context) {
	var (
		request models.ComboCreateRequest
	)

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("error reading body: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}
	h.log.Info("Incoming JSON: " + string(body))

	// Unmarshal the request body into the ComboCreateRequest struct
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.log.Error("error unmarshalling JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid JSON!"})
		return
	}

	// Validate the Combo data
	if request.Combo.Name == "" {
		h.log.Error("Combo name is empty!")
		c.JSON(http.StatusBadRequest, Response{Data: "Combo name is required!"})
		return
	}
	if request.Combo.Price <= 0 {
		h.log.Error("Invalid combo price!")
		c.JSON(http.StatusBadRequest, Response{Data: "Valid combo price is required!"})
		return
	}

	// Validate each item in the combo
	for _, item := range request.Combo.ComboItems {
		if item.ProductId == "" {
			h.log.Error("Product ID is empty for one of the items!")
			c.JSON(http.StatusBadRequest, Response{Data: "Product ID is required for each item!"})
			return
		}
		if item.Quantity <= 0 {
			h.log.Error("Invalid quantity for product: " + item.ProductId)
			c.JSON(http.StatusBadRequest, Response{Data: "Valid quantity is required for each item!"})
			return
		}
	}

	// Call the Create method in the repository to insert the combo into the database
	combo, err := h.storage.Combo().Create(c.Request.Context(), &request)
	if err != nil {
		h.log.Error("error in Combo.Create: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	// Respond with the created combo's ID
	h.log.Info("Combo Created Successfully!")
	c.JSON(http.StatusCreated, Response{Data: combo})
}
