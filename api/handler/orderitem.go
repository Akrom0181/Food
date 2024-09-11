package handler

import (
	"food/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			create_order_item
// @Router 		/food/api/v1/createorderitem [POST]
// @Summary 	Create Order Item
// @Description Create a new order item
// @Tags 		order_item
// @Accept 		json
// @Produce 	json
// @Param 		OrderItem body models.CreateOrderItem true "OrderItem"
// @Success 	200 {object} models.OrderItem
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) CreateOrderItem(c *gin.Context) {
	var orderItem models.OrderItem

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		h.log.Error("error while binding OrderItem json: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid input data")
		return
	}

	resp, err := h.storage.OrderItem().Create(c.Request.Context(), &orderItem)
	if err != nil {
		h.log.Error("error while creating order item: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	h.log.Info("Order item created successfully")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			update_order_item
// @Router 		/food/api/v1/updateorderitem/{id} [PUT]
// @Summary 	Update Order Item
// @Description Update an existing order item
// @Tags 		order_item
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "OrderItem ID"
// @Param 		OrderItem body models.UpdateOrderItem true "UpdateOrderItemRequest"
// @Success 	200 {object} models.OrderItem
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateOrderItem(c *gin.Context) {
	id := c.Param("id")

	if err := uuid.Validate(id); err != nil {
		h.log.Error("invalid order item id: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}

	var updateOrderItem models.UpdateOrderItem
	if err := c.ShouldBindJSON(&updateOrderItem); err != nil {
		h.log.Error("error while binding OrderItem json: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid input data")
		return
	}

	orderItem, err := h.storage.OrderItem().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error while fetching order item: " + err.Error())
		c.JSON(http.StatusNotFound, "Order item not found")
		return
	}

	// Update order item fields
	orderItem.ProductId = updateOrderItem.ProductId
	orderItem.Quantity = updateOrderItem.Quantity
	orderItem.Price = updateOrderItem.Price

	resp, err := h.storage.OrderItem().Update(c.Request.Context(), orderItem)
	if err != nil {
		h.log.Error("error while updating order item: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	h.log.Info("Order item updated successfully")
	c.JSON(http.StatusOK, resp)
}

// @ID 			get_order_item
// @Router 		/food/api/v1/getbyorderitem/{id} [GET]
// @Summary 	Get Order Item by ID
// @Description Retrieve an order item by its ID
// @Tags 		order_item
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "OrderItem ID"
// @Success 	200 {object} models.OrderItem
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) GetOrderItemByID(c *gin.Context) {
	id := c.Param("id")

	if err := uuid.Validate(id); err != nil {
		h.log.Error("invalid order item id: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}

	orderItem, err := h.storage.OrderItem().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error while fetching order item: " + err.Error())
		c.JSON(http.StatusNotFound, "Order item not found")
		return
	}

	h.log.Info("Order item retrieved successfully")
	c.JSON(http.StatusOK, orderItem)
}

// @ID 			get_all_order_items
// @Router 		/food/api/v1/getallorderitems [GET]
// @Summary 	Get All Order Items
// @Description Retrieve all order items
// @Tags 		order_item
// @Accept 		json
// @Produce 	json
// @Param 		search query string false "Search order items by product ID"
// @Param 		page   query uint64 false "Page number"
// @Param 		limit  query uint64 false "Limit number of results per page"
// @Success 	200 {object} models.GetAllOrderItemsResponse
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllOrderItems(c *gin.Context) {
	var req = &models.GetAllOrderItemsRequest{}

	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		h.log.Error("error while parsing page: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid page value")
		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		h.log.Error("error while parsing limit: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid limit value")
		return
	}

	req.Page = page
	req.Limit = limit

	orderItems, err := h.storage.OrderItem().GetAll(c.Request.Context(), req)
	if err != nil {
		h.log.Error("error while fetching order items: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	h.log.Info("Order items retrieved successfully")
	c.JSON(http.StatusOK, orderItems)
}

// @ID 			delete_order_item
// @Router 		/food/api/v1/deleteorderitem/{id} [DELETE]
// @Summary 	Delete Order Item by ID
// @Description Delete an order item by its ID
// @Tags 		order_item
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "OrderItem ID"
// @Success 	200 {object} Response{data=string} "Success Request"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteOrderItem(c *gin.Context) {
	id := c.Param("id")

	if err := uuid.Validate(id); err != nil {
		h.log.Error("invalid order item id: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}

	err := h.storage.OrderItem().Delete(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error while deleting order item: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	h.log.Info("Order item deleted successfully")
	c.JSON(http.StatusOK, id)
}
