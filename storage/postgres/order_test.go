package postgres

import (
	"context"
	"food/api/models"
	"food/pkg/logger"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestOrderRepo_Create(t *testing.T) {
	type fields struct {
		db  *pgxpool.Pool
		log logger.LoggerI
	}
	type args struct {
		ctx   context.Context
		order *models.OrderCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.OrderCreateRequest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderRepo{
				db:  tt.fields.db,
				log: tt.fields.log,
			}
			got, err := o.Create(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepo_Update(t *testing.T) {
	type fields struct {
		db  *pgxpool.Pool
		log logger.LoggerI
	}
	type args struct {
		ctx          context.Context
		id           string
		updatedOrder *models.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.OrderCreateRequest
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepo{
				db:  tt.fields.db,
				log: tt.fields.log,
			}
			got, err := r.Update(tt.args.ctx, tt.args.id, tt.args.updatedOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepo.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepo_GetAll(t *testing.T) {
	type fields struct {
		db  *pgxpool.Pool
		log logger.LoggerI
	}
	type args struct {
		ctx     context.Context
		request *models.GetAllOrdersRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]models.OrderCreateRequest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderRepo{
				db:  tt.fields.db,
				log: tt.fields.log,
			}
			got, err := o.GetAll(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepo.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepo_Delete(t *testing.T) {
	type fields struct {
		db  *pgxpool.Pool
		log logger.LoggerI
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepo{
				db:  tt.fields.db,
				log: tt.fields.log,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
