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

func (o *OrderRepo) Create(ctx context.Context, order *models.Order) (*models.Order, error) {

	id := uuid.New()
	query := `INSERT INTO "orders" (
		id,
		user_id,
		total_price,
		status,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := o.db.Exec(context.Background(), query,
		id.String(),
		order.UserId,
		order.TotalPrice,
		order.Status,
	)

	if err != nil {
		o.log.Error("error while creating order in strg" + err.Error())
		return &models.Order{}, err
	}
	return &models.Order{
		Id:         id.String(),
		UserId:     order.UserId,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}, nil
}

func (o *OrderRepo) Update(ctx context.Context, order *models.Order) (*models.Order, error) {
	query := `UPDATE "orders" SET 
		user_id=$1,
		total_price=$2,
		status=$3,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := o.db.Exec(context.Background(), query,
		order.UserId,
		order.TotalPrice,
		order.Status,
		order.Id,
	)
	if err != nil {
		o.log.Error("error while updating in strg" + err.Error())
		return &models.Order{}, err
	}
	return &models.Order{
		Id:         order.Id,
		UserId:     order.UserId,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}, nil
}

func (o *OrderRepo) GetAll(ctx context.Context, req *models.GetAllOrdersRequest) (*models.GetAllOrdersResponse, error) {
	var (
		resp   = &models.GetAllOrdersResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (user_id ILIKE '%%%v%%' OR status ILIKE '%%%v%%') `, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := o.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        user_id,
        total_price,
        status,
        created_at,
        updated_at FROM "orders" WHERE id=1`+filter)
	if err != nil {
		o.log.Error("error while getting all order in strg")
		return resp, err
	}

	for rows.Next() {
		var (
			order      = models.Order{}
			userID     sql.NullString
			totalPrice sql.NullFloat64
			status     sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&order.Id,
			&userID,
			&totalPrice,
			&status,
			&createdAt,
			&updatedAt); err != nil {
			return resp, err
		}

		resp.Orders = append(resp.Orders, models.Order{
			Id:         order.Id,
			UserId:     userID.String,
			TotalPrice: totalPrice.Float64,
			Status:     status.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
		})
	}
	return resp, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var (
		order      = models.Order{}
		userID     sql.NullString
		totalPrice sql.NullFloat64
		status     sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
	)
	if err := o.db.QueryRow(context.Background(), `SELECT id, user_id, total_price, status, created_at, updated_at FROM "orders" WHERE id = $1`, id).Scan(
		&order.Id,
		&userID,
		&totalPrice,
		&status,
		&createdAt,
		&updatedAt,
	); err != nil {
		o.log.Error("error while getbyid order in strg")
		return &models.Order{}, err
	}
	return &models.Order{
		Id:         order.Id,
		UserId:     userID.String,
		TotalPrice: totalPrice.Float64,
		Status:     status.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
	}, nil
}

func (o *OrderRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "orders" WHERE id = $1`
	_, err := o.db.Exec(context.Background(), query, id)
	if err != nil {
		o.log.Error("error while deleting order in strg")
		return err
	}
	return nil
}
