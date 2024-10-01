package postgres

import (
	"context"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ComboRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCombo(db *pgxpool.Pool, log logger.LoggerI) ComboRepo {
	return ComboRepo{
		db:  db,
		log: log,
	}
}
func (c *ComboRepo) Create(ctx context.Context, combo *models.ComboCreateRequest) (*models.ComboCreateRequest, error) {
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		return &models.ComboCreateRequest{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Generate a new UUID for the combo
	comboId := uuid.New().String()

	var totalSum float64
	for i, item := range combo.Items {
		if item.Quantity <= 0 {
			return &models.ComboCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = c.db.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.ComboCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		combo.Items[i].Price = productPrice
		combo.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		totalSum += combo.Items[i].TotalPrice
		combo.Items[i].ComboId = comboId
		combo.Items[i].CreatedAt = item.CreatedAt
	}

	// Insert the combo
	comboQuery := `INSERT INTO "combo" (id, name, price, description, created_at) 
					  VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), comboQuery, comboId, combo.Combo.Name, combo.Combo.Price, combo.Combo.Description)
	if err != nil {
		return &models.ComboCreateRequest{}, err
	}

	// Insert the combo items
	itemQuery := `INSERT INTO "combo_items" (id, quantity, combo_id, product_id, price, total_price, created_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)`

	for _, item := range combo.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, comboId, item.ProductId, item.Price, item.TotalPrice)
		for i := range combo.Items {
			combo.Items[i].Id = itemId
			combo.Items[i].CreatedAt = item.CreatedAt
		}
		if err != nil {
			return &models.ComboCreateRequest{}, err
		}
	}


	combo.Combo.Id = comboId
	combo.Combo.TotalPrice = totalSum

	return combo, tx.Commit(context.Background())
}
