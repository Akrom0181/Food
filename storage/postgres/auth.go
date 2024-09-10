package postgres

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"food/api/models"
// 	"food/config"
// 	"food/pkg"
// 	"food/pkg/jwt"
// 	"food/pkg/logger"
// 	"food/pkg/smtp"

// 	// "food/pkg/password"
// 	"food/storage"
// 	"time"

// 	"github.com/jackc/pgx/v4/pgxpool"
// )

// type AuthRepo struct {
// 	db    *pgxpool.Pool
// 	user  storage.IUserStorage
// 	log   logger.LoggerI
// 	redis storage.IRedisStorage
// }

// func NewAuthRepo(db *pgxpool.Pool, log logger.LoggerI, redis storage.IRedisStorage, user storage.IUserStorage) AuthRepo {
// 	return AuthRepo{
// 		db:    db,
// 		log:   log,
// 		redis: redis,
// 		user: user,
// 	}
// }

// // func (a authService) CustomerLogin(ctx context.Context, loginRequest models.CustomerLoginRequest) (models.CustomerLoginResponse, error) {
// // 	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
// // 	customer, err := a.storage.Customer().GetByLogin(ctx, loginRequest.Login)
// // 	if err != nil {
// // 		a.log.Error("error while getting customer credentials by login", logger.Error(err))
// // 		return models.CustomerLoginResponse{}, err
// // 	}

// // 	if err = password.CompareHashAndPassword(customer.Password, loginRequest.Password); err != nil {
// // 		a.log.Error("error while comparing password", logger.Error(err))
// // 		return models.CustomerLoginResponse{}, err
// // 	}

// // 	m := make(map[interface{}]interface{})

// // 	m["user_id"] = customer.ID
// // 	m["user_role"] = config.USER_ROLE

// // 	accessToken, refreshToken, err := jwt.GenJWT(m)
// // 	if err != nil {
// // 		a.log.Error("error while generating tokens for customer login", logger.Error(err))
// // 		return models.CustomerLoginResponse{}, err
// // 	}

// // 	return models.CustomerLoginResponse{
// // 		AccessToken:  accessToken,
// // 		RefreshToken: refreshToken,
// // 	}, nil
// // }

// func (a *AuthRepo) UserRegister(ctx context.Context, loginRequest models.UserRegisterRequest) error {
// 	fmt.Println(" loginRequest.Login: ", loginRequest.Mail)

// 	otpCode := pkg.GenerateOTP()

// 	msg := fmt.Sprintf("Your otp code is: %v, for registering Khorezm_Shashlik. Don't give it to anyone", otpCode)

// 	err := a.redis.SetX(ctx, loginRequest.Mail, otpCode, time.Minute*2)
// 	if err != nil {
// 		a.log.Error("error while setting otpCode to redis customer register", logger.Error(err))
// 		return err
// 	}

// 	err = smtp.SendMail(loginRequest.Mail, msg)
// 	if err != nil {
// 		a.log.Error("error while sending otp code to customer register", logger.Error(err))
// 		return err
// 	}
// 	return nil
// }

// func (a *AuthRepo) UserRegisterConfirm(ctx context.Context, req models.UserRegisterConfRequest) (models.UserLoginResponse, error) {
// 	resp := models.UserLoginResponse{}

// 	otp, err := a.redis.Get(ctx, req.Mail)
// 	if err != nil {
// 		a.log.Error("error while getting otp code for customer register confirm", logger.Error(err))
// 		return resp, err
// 	}
// 	if req.Otp != otp {
// 		a.log.Error("incorrect otp code for customer register confirm", logger.Error(err))
// 		return resp, errors.New("incorrect otp code")
// 	}
// 	req.User.Email = req.Mail

// 	id, err := a.user.Create(ctx, req.User)
// 	if err != nil {
// 		a.log.Error("error while creating customer", logger.Error(err))
// 		return resp, err
// 	}
// 	var m = make(map[interface{}]interface{})

// 	m["user_id"] = id
// 	m["user_role"] = config.USER_ROLE

// 	accessToken, refreshToken, err := jwt.GenJWT(m)
// 	if err != nil {
// 		a.log.Error("error while generating tokens for customer register confirm", logger.Error(err))
// 		return resp, err
// 	}
// 	resp.AccessToken = accessToken
// 	resp.RefreshToken = refreshToken

// 	return resp, nil
// }
