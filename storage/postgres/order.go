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

// Create method to insert a new order into the database
func (o *OrderRepo) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	id := uuid.New()
	query := `INSERT INTO "order" (
		id, user_id, orderitem_id, total_price, status, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING created_at, updated_at`

	// Execute query and get the created_at and updated_at from the DB
	var createdAt, updatedAt string
	err := o.db.QueryRow(context.Background(), query,
		id.String(),
		order.UserId,
		order.OrderItemId,
		order.TotalPrice,
		order.Status,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		o.log.Error("error while creating order in storage: " + err.Error())
		return &models.Order{}, err
	}

	return &models.Order{
		Id:          id.String(),
		UserId:      order.UserId,
		OrderItemId: order.OrderItemId,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (o *OrderRepo) Update(ctx context.Context, order *models.Order) (*models.Order, error) {
	query := `UPDATE "order" SET 
		user_id = $1, orderitem_id = $2, total_price = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at`

	var updatedAt string
	err := o.db.QueryRow(context.Background(), query,
		order.UserId,
		order.OrderItemId,
		order.TotalPrice,
		order.Status,
		order.Id,
	).Scan(&updatedAt)

	if err != nil {
		o.log.Error("error while updating order in storage: " + err.Error())
		return &models.Order{}, err
	}

	return &models.Order{
		Id:          order.Id,
		UserId:      order.UserId,
		OrderItemId: order.OrderItemId,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt, 
		UpdatedAt:   updatedAt,
	}, nil
}

func (o *OrderRepo) GetAll(ctx context.Context, req *models.GetAllOrdersRequest) (*models.GetAllOrdersResponse, error) {
	resp := &models.GetAllOrdersResponse{}
	offset := (req.Page - 1) * req.Limit

	// Build query with optional search
	filter := ""
	if req.Search != "" {
		filter += fmt.Sprintf(` WHERE (user_id ILIKE '%%%v%%' OR status ILIKE '%%%v%%') `, req.Search, req.Search)
	}
	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := o.db.Query(context.Background(), `SELECT count(id) OVER(), 
		id, user_id, orderitem_id, total_price, status, created_at, updated_at
		FROM "order"`+filter)

	if err != nil {
		o.log.Error("error while getting all orders in storage: " + err.Error())
		return resp, err
	}

	for rows.Next() {
		var (
			order       = models.Order{}
			userID      sql.NullString
			orderItemId sql.NullString
			totalPrice  sql.NullFloat64
			status      sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&order.Id,
			&userID,
			&orderItemId,
			&totalPrice,
			&status,
			&createdAt,
			&updatedAt); err != nil {
			return resp, err
		}

		resp.Orders = append(resp.Orders, models.Order{
			Id:          order.Id,
			UserId:      userID.String,
			OrderItemId: orderItemId.String,
			TotalPrice:  totalPrice.Float64,
			Status:      status.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}
	return resp, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var (
		order       = models.Order{}
		userID      sql.NullString
		orderItemId sql.NullString
		totalPrice  sql.NullFloat64
		status      sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)
	if err := o.db.QueryRow(context.Background(), `SELECT id, user_id, orderitem_id, total_price, status, created_at, updated_at FROM "order" WHERE id = $1`, id).Scan(
		&order.Id,
		&userID,
		&orderItemId,
		&totalPrice,
		&status,
		&createdAt,
		&updatedAt,
	); err != nil {
		o.log.Error("error while getting order by ID in storage: " + err.Error())
		return &models.Order{}, err
	}
	return &models.Order{
		Id:          order.Id,
		UserId:      userID.String,
		OrderItemId: orderItemId.String,
		TotalPrice:  totalPrice.Float64,
		Status:      status.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}, nil
}

func (o *OrderRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "order" WHERE id = $1`
	result, err := o.db.Exec(context.Background(), query, id)
	if err != nil {
		o.log.Error("error while deleting order in storage: " + err.Error())
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no order found with id %v", id)
	}
	return nil
}
