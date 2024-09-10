package handler

import (
	"fmt"
	"food/api/models"
	check "food/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogin godoc
// @Router       /food/api/v1/user/login [POST]
// @Summary      User login
// @Description  Login to Khorezm_Shashlik
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.UserLoginRequest true "login"
// @Success      201  {object}  models.UserLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) UserLogin(c *gin.Context) {
	loginReq := models.UserLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.log, "error while binding body", http.StatusInternalServerError, err)
		return
	}

	fmt.Println("loginReq: ", loginReq)

	//TODO: need validate login & password

	loginResp, err := h.service.Auth().UserLogin(c.Request.Context(), loginReq)
	if err != nil {
		handleResponseLog(c, h.log, "unauthorized", http.StatusUnauthorized, err)
		return
	}

	handleResponseLog(c, h.log, "Success", http.StatusOK, loginResp)

}

// UserRegister godoc
// @Router       /food/api/v1/sendcode [POST]
// @Summary      User register
// @Description  Registering to Khorezm_Shashlik
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterRequest true "register"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) UserRegister(c *gin.Context) {
	loginReq := models.UserRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq: ", loginReq)

	if err := check.ValidateEmailAddress(loginReq.Mail); err != nil {
		handleResponseLog(c, h.log, "error while validating email" + loginReq.Mail, http.StatusBadRequest, err.Error())
		return
	}

	// if err := check.CheckEmailExists(loginReq.Mail); err != nil {
	// 	handleResponseLog(c, h.log, "error email does not exist" + loginReq.Mail, http.StatusBadRequest, err.Error())
	// }

	err := h.service.Auth().UserRegister(c.Request.Context(), loginReq)
	if err != nil {
		handleResponseLog(c, h.log, "", http.StatusInternalServerError, err)
		return
	}

	handleResponseLog(c, h.log, "Otp sent successfull", http.StatusOK, "")
}

// UserRegisterConfirm godoc
// @Router       /food/api/v1/user/verifycode [POST]
// @Summary      User register
// @Description  Registering to Khorezm_Shashlik
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterConfRequest true "register"
// @Success      201  {object}  models.UserLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) UserRegisterConfirm(c *gin.Context) {
	req := models.UserRegisterConfRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponseLog(c, h.log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("req: ", req)

	//TODO: need validate login & password

	confResp, err := h.service.Auth().UserRegisterConfirm(c.Request.Context(), req)
	if err != nil {
		handleResponseLog(c, h.log, "error while confirming", http.StatusUnauthorized, err.Error())
		return
	}

	handleResponseLog(c, h.log, "Success", http.StatusOK, confResp)

}
