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

type OrderItemRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewOrderItem(db *pgxpool.Pool, log logger.LoggerI) OrderItemRepo {
	return OrderItemRepo{
		db: db,
		log: log,
	}
}

func (o *OrderItemRepo) Create(ctx context.Context, orderItem *models.OrderItem) (*models.OrderItem, error) {

	id := uuid.New()
	query := `INSERT INTO "order_items" (
		id,
		order_id,
		product_id,
		quantity,
		price,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,$5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := o.db.Exec(context.Background(), query,
		id.String(),
		orderItem.OrderId,
		orderItem.ProductId,
		orderItem.Quantity,
		orderItem.Price,
	)

	if err != nil {
		return &models.OrderItem{}, err
	}
	return &models.OrderItem{
		Id:        id.String(),
		OrderId:   orderItem.OrderId,
		ProductId: orderItem.ProductId,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}, nil
}

func (o *OrderItemRepo) Update(ctx context.Context, orderItem *models.OrderItem) (*models.OrderItem, error) {
	query := `UPDATE "order_items" SET 
		order_id=$1,
		product_id=$2,
		quantity=$3,
		price=$4,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := o.db.Exec(context.Background(), query,
		orderItem.OrderId,
		orderItem.ProductId,
		orderItem.Quantity,
		orderItem.Price,
		orderItem.Id,
	)
	if err != nil {
		return &models.OrderItem{}, err
	}
	return &models.OrderItem{
		Id:        orderItem.Id,
		OrderId:   orderItem.OrderId,
		ProductId: orderItem.ProductId,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}, nil
}

func (o *OrderItemRepo) GetAll(ctx context.Context, req *models.GetAllOrderItemsRequest) (*models.GetAllOrderItemsResponse, error) {
	var (
		resp   = &models.GetAllOrderItemsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (order_id ILIKE '%%%v%%' OR product_id ILIKE '%%%v%%') `, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := o.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        order_id,
        product_id,
        quantity,
        price,
        created_at,
        updated_at FROM "order_items"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			orderItem  = models.OrderItem{}
			order_id   sql.NullString
			product_id sql.NullString
			quantity   sql.NullInt64
			price      sql.NullFloat64
			created_at sql.NullString
			updated_at sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&orderItem.Id,
			&order_id,
			&product_id,
			&quantity,
			&price,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}

		resp.OrderItems = append(resp.OrderItems, models.OrderItem{
			Id:        orderItem.Id,
			OrderId:   order_id.String,
			ProductId: product_id.String,
			Quantity:  int(quantity.Int64),
			Price:     price.Float64,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}
	return resp, nil
}

func (o *OrderItemRepo) GetByID(ctx context.Context, id string) (*models.OrderItem, error) {
	var (
		orderItem  = models.OrderItem{}
		order_id   sql.NullString
		product_id sql.NullString
		quantity   sql.NullInt64
		price      sql.NullFloat64
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := o.db.QueryRow(context.Background(), `
	    SELECT 
		 id, 
		 order_id, 
		 product_id, 
		 quantity, 
		 price, 
		 created_at, 
		 updated_at 
		FROM "order_items" WHERE id = $1`, id).Scan(
		&orderItem.Id,
		&order_id,
		&product_id,
		&quantity,
		&price,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.OrderItem{}, err
	}
	return &models.OrderItem{
		Id:        orderItem.Id,
		OrderId:   order_id.String,
		ProductId: product_id.String,
		Quantity:  int(quantity.Int64),
		Price:     price.Float64,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (o *OrderItemRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM order_items WHERE id = $1`
	_, err := o.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
