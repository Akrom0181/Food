package models

type GetAllOrdersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllOrdersResponse struct {
	Orders []Order `json:"orders"`
	Count  int64   `json:"count"`
}

type ChangeStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type GetOrderStatus struct {
	Status string `json:"status"`
}

type Order struct {
	Id         string      `json:"id,omitempty"`
	UserId     string      `json:"user_id,omitempty"`
	TotalPrice float64     `json:"total_price,omitempty"`
	Status     string      `json:"status,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
	UpdatedAt  string      `json:"updated_at,omitempty"`
	OrderItems []OrderItem `json:"order_items,omitempty"`
}

type OrderCreate struct {
	UserId string `json:"user_id,omitempty"`
	Status string `json:"status"`
}

type SwaggerOrderCreate struct {
	UserId string `json:"user_id,omitempty"`
}

type OrderUpdate struct {
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items,omitempty"`
}

type OrderUpdateS struct {
	TotalPrice float64           `json:"total_price"`
	Status     string            `json:"status"`
	OrderItems []UpdateOrderItem `json:"order_items,omitempty"`
}

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type OrderGetListRequest struct {
	UserId string `json:"user_id,omitempty"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type OrderGetListResponse struct {
	Count int      `json:"count"`
	Order []*Order `json:"order"`
}

type OrderCreateRequest struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}

type SwaggerOrderCreateRequest struct {
	Order SwaggerOrderCreate  `json:"order"`
	Items []SwaggerOrderItems `json:"items"`
}

