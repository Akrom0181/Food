package models

type Order struct {
	Id         string  `json:"id"`
	UserId     string  `json:"user_id"`
	ProductId  string  `json:"product_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type CreateOrder struct {
	UserId     string  `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type UpdateOrder struct {
	UserId     string  `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type GetOrder struct {
	Id         string  `json:"id"`
	UserId     string  `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type GetAllOrdersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllOrdersResponse struct {
	Orders []Order `json:"orders"`
	Count  int64   `json:"count"`
}
