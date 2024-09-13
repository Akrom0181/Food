package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewOrder(db *pgxpool.Pool, log logger.LoggerI) OrderRepo {
	return OrderRepo{
		db:  db,
		log: log,
	}
}

func (o *OrderRepo) Create(ctx context.Context, order *models.OrderCreateRequest) (*models.OrderCreateRequest, error) {
	tx, err := o.db.Begin(context.Background())
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Generate a new UUID for the order
	orderId := uuid.New().String()

	var totalSum float64
	for i, item := range order.Items {
		if item.Quantity <= 0 {
			return &models.OrderCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = o.db.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		order.Items[i].Price = productPrice
		order.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		totalSum += order.Items[i].TotalPrice
		order.Items[i].Id = item.ProductId
		order.Items[i].OrderId = orderId
		order.Items[i].CreatedAt = item.CreatedAt
	}

	// Insert the order
	orderQuery := `INSERT INTO "order" (id, user_id, total_price, created_at, updated_at) 
					  VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), orderQuery, orderId, order.Order.UserId, totalSum)
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}

	// Insert the order items
	itemQuery := `INSERT INTO "orderiteam" (id, quantity, order_id, product_id, price, total, created_at, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	for _, item := range order.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, orderId, item.ProductId, item.Price, item.TotalPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, err
		}
	}

	order.Order.Id = orderId
	order.Order.TotalPrice = totalSum

	return order, tx.Commit(context.Background())
}

func (o *OrderRepo) GetAll(ctx context.Context, request *models.GetAllOrdersRequest) (*[]models.OrderCreateRequest, error) {
	var (
		orders     []models.OrderCreateRequest
		created_at sql.NullString
		updated_at sql.NullString
	)

	// Query to retrieve all orders
	orderQuery := `
		SELECT id, user_id, total_price, status, created_at, updated_at
		FROM "order"
	`
	rows, err := o.db.Query(ctx, orderQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	defer rows.Close()

	// Iterate over the retrieved orders
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &created_at, &updated_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Query to retrieve order items for the current order
		orderItemQuery := `
			SELECT id, product_id, order_id, quantity, price, total, created_at, updated_at
			FROM "orderiteam"
			WHERE order_id = $1
		`
		itemRows, err := o.db.Query(ctx, orderItemQuery, order.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve items for order %s: %w", order.Id, err)
		}
		defer itemRows.Close()

		var orderItems []models.OrderItem
		for itemRows.Next() {
			var item models.OrderItem
			err = itemRows.Scan(&item.Id, &item.ProductId, &item.OrderId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at, &updated_at)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order item: %w", err)
			}
			orderItems = append(orderItems, models.OrderItem{
				Id:         item.Id,
				ProductId:  item.ProductId,
				OrderId:    item.OrderId,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
				CreatedAt:  created_at.String,
				UpdatedAt:  updated_at.String,
			})
		}

		// Add order items to the order struct
		// order.OrderItems = orderItems

		// Append the order to the result set
		orders = append(orders, models.OrderCreateRequest{
			Order: models.Order{
				Id:         order.Id,
                UserId:     order.UserId,
                TotalPrice: order.TotalPrice,
                Status:     order.Status,
                CreatedAt:  created_at.String,
                UpdatedAt:  updated_at.String,
			},
			Items: orderItems,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}
