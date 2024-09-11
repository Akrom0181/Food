package storage

import (
	"context"
	"food/api/models"
	"time"
)

type IStorage interface {
	CloseDB()
	// Auth() IAuthStorage
	User() IUserStorage
	Branch() IBranchStorage
	Banner() IBannerStorage
	Category() ICategoryStorage
	Product() IProductStorage
	Order() IOrderStorage
	OrderItem() IOrderItemStorage
	CourierAssignment() ICourierAssignmentStorage
	Notification() INotificationStorage
	DeliveryHistory() IDeliveryHistoryStorage
	Redis() IRedisStorage
}

type IUserStorage interface {
	Create(context.Context, *models.User) (*models.User, error)
	GetAll(ctx context.Context, request *models.GetAllUsersRequest) (*models.GetAllUsersResponse, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(context.Context, *models.User) (*models.User, error)
	Delete(context.Context, string) error
	GetByLogin(ctx context.Context, login string) (models.User, error)
}

type IBannerStorage interface {
	Create(context.Context, *models.Banner) (*models.Banner, error)
	GetAll(ctx context.Context, request *models.GetAllBannerRequest) (*models.GetAllBannerResponse, error)
	Delete(ctx context.Context, id string) error
}

type IBranchStorage interface {
	Create(ctx context.Context, branch *models.Branch) (*models.Branch, error)
	GetAll(ctx context.Context, request *models.GetAllBranchesRequest) (*models.GetAllBranchesResponse, error)
	GetByID(ctx context.Context, id string) (*models.Branch, error)
	Update(ctx context.Context, branch *models.Branch) (*models.Branch, error)
	Delete(ctx context.Context, id string) error
}

type ICategoryStorage interface {
	Create(context.Context, *models.Category) (*models.Category, error)
	GetAll(ctx context.Context, request *models.GetAllCategoriesRequest) (*models.GetAllCategoriesResponse, error)
	GetByID(ctx context.Context, id string) (*models.Category, error)
	Update(context.Context, *models.Category) (*models.Category, error)
	Delete(context.Context, string) error
}

type IProductStorage interface {
	Create(context.Context, *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, request *models.GetAllProductsRequest) (*models.GetAllProductsResponse, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Update(context.Context, *models.Product) (*models.Product, error)
	Delete(context.Context, string) error
}

type IOrderStorage interface {
	Create(context.Context, *models.Order) (*models.Order, error)
	GetAll(ctx context.Context, request *models.GetAllOrdersRequest) (*models.GetAllOrdersResponse, error)
	GetByID(ctx context.Context, id string) (*models.Order, error)
	Update(context.Context, *models.Order) (*models.Order, error)
	Delete(context.Context, string) error
}

type IOrderItemStorage interface {
	Create(context.Context, *models.OrderItem) (*models.OrderItem, error)
	GetAll(ctx context.Context, request *models.GetAllOrderItemsRequest) (*models.GetAllOrderItemsResponse, error)
	GetByID(ctx context.Context, id string) (*models.OrderItem, error)
	Update(context.Context, *models.OrderItem) (*models.OrderItem, error)
	Delete(context.Context, string) error
}

type ICourierAssignmentStorage interface {
	Create(context.Context, *models.CourierAssignment) (*models.CourierAssignment, error)
	GetAll(ctx context.Context, request *models.GetAllCourierAssignmentsRequest) (*models.GetAllCourierAssignmentsResponse, error)
	GetByID(ctx context.Context, id string) (*models.CourierAssignment, error)
	Update(context.Context, *models.CourierAssignment) (*models.CourierAssignment, error)
	Delete(context.Context, string) error
}

type INotificationStorage interface {
	Create(context.Context, *models.Notification) (*models.Notification, error)
	GetAll(ctx context.Context, request *models.GetAllNotificationsRequest) (*models.GetAllNotificationsResponse, error)
	GetByID(ctx context.Context, id string) (*models.Notification, error)
	Update(context.Context, *models.Notification) (*models.Notification, error)
	Delete(context.Context, string) error
}

type IDeliveryHistoryStorage interface {
	Create(context.Context, *models.DeliveryHistory) (*models.DeliveryHistory, error)
	GetAll(ctx context.Context, request *models.GetAllDeliveryHistoriesRequest) (*models.GetAllDeliveryHistoriesResponse, error)
	GetByID(ctx context.Context, id string) (*models.DeliveryHistory, error)
	Update(context.Context, *models.DeliveryHistory) (*models.DeliveryHistory, error)
	Delete(context.Context, string) error
}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
}

// type IAuthStorage interface {
// 	UserRegister(ctx context.Context, loginRequest models.UserRegisterRequest) error
// 	UserRegisterConfirm(ctx context.Context, req models.UserRegisterConfRequest) (models.UserLoginResponse, error)
// }
