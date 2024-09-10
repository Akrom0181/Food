package handler

// import (
// 	"context"
// 	"food/api/models"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// // @ID 			create_order
// // @Router 		/food/api/v1/createorder [POST]
// // @Summary 	Create Order
// // @Description Create a new order
// // @Tags 		order
// // @Accept 		json
// // @Produce 	json
// // @Param 		Order body models.CreateOrder true "Order"
// // @Success 	200 {object} models.Order
// // @Response 	400 {object} Response{data=string} "Bad Request"
// // @Failure 	500 {object} Response{data=string} "Server error"
// func (h *Handler) CreateOrder(c *gin.Context) {
// 	var order models.Order

// 	if err := c.ShouldBindJSON(&order); err != nil {
// 		h.log.Error(err.Error() + " : " + "error Order Should Bind Json!")
// 		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
// 		return
// 	}

// 	resp, err := h.storage.Order().Create(c.Request.Context(), &order)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "Error Order Create")
// 		c.JSON(http.StatusInternalServerError, "Server error!")
// 		return
// 	}

// 	h.log.Info("Order created successfully!")
// 	c.JSON(http.StatusCreated, resp)
// }

// // @ID 			update_order
// // @Router 		/food/api/v1/updateorder/{id} [PUT]
// // @Summary 	Update Order
// // @Description Update an existing order
// // @Tags 		order
// // @Accept 		json
// // @Produce 	json
// // @Param 		id path string true "Order ID"
// // @Param 		Order body models.UpdateOrder true "UpdateOrderRequest"
// // @Success 	200 {object} models.Order
// // @Response 	400 {object} Response{data=string} "Bad Request"
// // @Failure 	500 {object} Response{data=string} "Server error"
// func (h *Handler) UpdateOrder(c *gin.Context) {
// 	var updateOrder models.UpdateOrder

// 	if err := c.ShouldBindJSON(&updateOrder); err != nil {
// 		h.log.Error(err.Error() + " : " + "error Order Should Bind Json!")
// 		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
// 		return
// 	}

// 	id := c.Param("id")
// 	order, err := h.storage.Order().GetByID(c.Request.Context(), id)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "Error Order Not Found")
// 		c.JSON(http.StatusBadRequest, "Order not found!")
// 		return
// 	}

// 	order.UserId = updateOrder.UserId
// 	order.TotalPrice = updateOrder.TotalPrice
// 	order.Status = updateOrder.Status

// 	resp, err := h.storage.Order().Update(c.Request.Context(), order)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "Error Order Update")
// 		c.JSON(http.StatusInternalServerError, "Server error!")
// 		return
// 	}

// 	h.log.Info("Order updated successfully!")
// 	c.JSON(http.StatusOK, resp)
// }

// // @ID 			get_order
// // @Router 		/food/api/v1/getorder/{id} [GET]
// // @Summary 	Get Order by ID
// // @Description Retrieve an order by its ID
// // @Tags 		order
// // @Accept 		json
// // @Produce 	json
// // @Param 		id path string true "Order ID"
// // @Success 	200 {object} models.Order
// // @Response 	400 {object} Response{data=string} "Bad Request"
// // @Failure 	500 {object} Response{data=string} "Server error"
// func (h *Handler) GetOrderByID(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		h.log.Error("missing order id")
// 		c.JSON(http.StatusBadRequest, "you must fill the ID")
// 		return
// 	}

// 	order, err := h.storage.Order().GetByID(context.Background(), id)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "Error while getting order by ID")
// 		c.JSON(http.StatusInternalServerError, "Server Error")
// 		return
// 	}

// 	h.log.Info("Order retrieved successfully by ID")
// 	c.JSON(http.StatusOK, order)
// }

// // @ID 			get_all_orders
// // @Router 		/food/api/v1/getallorders [GET]
// // @Summary 	Get All Orders
// // @Description Retrieve all orders
// // @Tags 		order
// // @Accept 		json
// // @Produce 	json
// // @Param 		search query string false "Search orders by status"
// // @Param 		page   query uint64 false "Page number"
// // @Param 		limit  query uint64 false "Limit number of results per page"
// // @Success 	200 {object} models.GetAllOrdersResponse
// // @Response 	400 {object} Response{data=string} "Bad Request"
// // @Failure 	500 {object} Response{data=string} "Server error"
// func (h *Handler) GetAllOrders(c *gin.Context) {
// 	var req = &models.GetAllOrdersRequest{}

// 	req.Search = c.Query("search")

// 	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "error while parsing page")
// 		c.JSON(http.StatusBadRequest, "BadRequest at paging")
// 		return
// 	}

// 	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "error while parsing limit")
// 		c.JSON(http.StatusInternalServerError, "Internal server error while parsing limit")
// 		return
// 	}

// 	req.Page = page
// 	req.Limit = limit

// 	orders, err := h.storage.Order().GetAll(context.Background(), req)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "Error while getting all orders")
// 		c.JSON(http.StatusInternalServerError, "Error while getting all orders")
// 		return
// 	}

// 	h.log.Info("Orders retrieved successfully")
// 	c.JSON(http.StatusOK, orders)
// }

// // @ID 			delete_order
// // @Router 		/food/api/v1/deleteorder/{id} [DELETE]
// // @Summary 	Delete Order by ID
// // @Description Delete an order by its ID
// // @Tags 		order
// // @Accept 		json
// // @Produce 	json
// // @Param 		id path string true "Order ID"
// // @Success 	200 {object} Response{data=string} "Success Request"
// // @Response 	400 {object} Response{data=string} "Bad Request"
// // @Failure 	500 {object} Response{data=string} "Server error"
// func (h *Handler) DeleteOrder(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		h.log.Error("missing order id")
// 		c.JSON(http.StatusBadRequest, "fill the gap with id")
// 		return
// 	}

// 	err := uuid.Validate(id)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "error while validating id")
// 		c.JSON(http.StatusBadRequest, "please enter a valid id")
// 		return
// 	}

// 	err = h.storage.Order().Delete(context.Background(), id)
// 	if err != nil {
// 		h.log.Error(err.Error() + ":" + "error while deleting order")
// 		c.JSON(http.StatusBadRequest, "please input valid data")
// 		return
// 	}

// 	h.log.Info("Order deleted successfully!")
// 	c.JSON(http.StatusOK, id)
// }
