package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/tsmdev/software-development/backend/go-project/constants"
	"gitlab.com/tsmdev/software-development/backend/go-project/domains"
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/validation"
	"gorm.io/gorm"
)

// JWTAuthController struct
type JWTAuthController struct {
	env         lib.Env
	signature   lib.Signature
	logger      lib.Logger
	service     domains.AuthService
	userService domains.UserService
}

// NewJWTAuthController creates new controller
func NewJWTAuthController(
	env lib.Env,
	signature lib.Signature,
	logger lib.Logger,
	service domains.AuthService,
	userService domains.UserService,
) JWTAuthController {
	return JWTAuthController{
		env:         env,
		logger:      logger,
		service:     service,
		userService: userService,
		signature:   signature,
	}
}

// SignIn signs in user
func (jwt JWTAuthController) SignIn(c *gin.Context) {
	jwt.logger.Info("Signing in user")
	request := dto.AuthRequestDto{}

	if err := c.ShouldBind(&request); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthRequestDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	result, err := jwt.service.ValidateAuth(&request, constants.TypeAuthToken)
	if err != nil {
		jwt.logger.Error("validation auth error: ", err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		result,
		http.StatusText(http.StatusOK),
	)
}

// SignIn V2 signs in user
func (jwt JWTAuthController) SignInV2(c *gin.Context) {
	jwt.logger.Info("Signing in user")
	request := dto.AuthRequestDto{}

	if err := c.ShouldBind(&request); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthRequestDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	encryptedString := jwt.signature.Encrypt(request.Username + request.Password + request.Client)

	jwt.logger.Debug("request signature: ", request.Signature)
	jwt.logger.Debug("encrypted value: ", encryptedString)
	if request.Signature != encryptedString {
		globalResponse(
			c,
			http.StatusUnauthorized,
			nil,
			nil,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	result, err := jwt.service.ValidateAuth(&request, constants.TypeAuthToken)
	if err != nil {
		jwt.logger.Error("validation auth error: ", err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		result,
		http.StatusText(http.StatusOK),
	)
}

// IssueAccessToken issues refresh token
func (jwt JWTAuthController) IssueAccessToken(c *gin.Context) {
	jwt.logger.Info("Issuing refresh token")
	request := dto.AuthIssueRefreshTokenDto{}
	if err := c.ShouldBind(&request); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthIssueRefreshTokenDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	result, err := jwt.service.ValidateAuth(nil, request.RefreshToken)
	if err != nil {
		jwt.logger.Error("validation auth error: ", err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		result,
		http.StatusText(http.StatusOK),
	)
}

// Register registers user
func (jwt JWTAuthController) Register(c *gin.Context) {

	user := dto.CreateUserRequest{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	if err := jwt.userService.WithTrx(trxHandle).CreateUser(user); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		"User created successfully",
	)
}

// ChangePassword changes password
func (jwt JWTAuthController) ChangePassword(c *gin.Context) {
	var request dto.AuthChangePasswordDto

	userId, err := strconv.Atoi(c.Request.Header.Get("user_id"))
	if err != nil {
		c.Error(err)
		return
	}

	if err := c.ShouldBind(&request); err != nil {
		jwt.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthChangePasswordDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	if err = jwt.service.ChangePassword(userId, &request); err != nil {
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		http.StatusText(http.StatusOK),
	)
}

// ForgotPassword handles the forgot password functionality.
// It receives a POST request with the user's email and sends a password reset email.
// If the request is invalid, it returns a validation error response.
// If an error occurs during the password reset process, it logs the error and returns an error response.
func (d JWTAuthController) ForgotPassword(c *gin.Context) {
	var request dto.AuthForgotPasswordDto

	d.logger.Info("Validate forgot password request")
	if err := c.ShouldBind(&request); err != nil {
		d.logger.Error(err.Error())
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthForgotPasswordDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	d.logger.Info("Calling forgot password method")
	if err := d.service.ForgotPassword(request); err != nil {
		d.logger.Error(err.Error())
		c.Error(err)
		return
	}

	d.logger.Info("Forgot password handled successfully")
	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		http.StatusText(http.StatusOK),
	)
}

// ValidateOTP validates the OTP provided in the request.
// It binds the request body to the ValidateOTPRequestDto struct,
// handles validation errors, and calls the ValidateOTP method of the service.
// If the validation or service call fails, it logs the error and returns an error response.
// Otherwise, it logs the success and returns a success response.
func (d JWTAuthController) ValidateOTP(c *gin.Context) {
	var request dto.ValidateOTPRequestDto

	d.logger.Info("Validate OTP request")
	if err := c.ShouldBind(&request); err != nil {
		d.logger.Error(err.Error())
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.ValidateOTPRequestDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	d.logger.Info("Calling validate OTP method")
	response, err := d.service.ValidateOTP(request)
	if err != nil {
		d.logger.Error(err.Error())
		c.Error(err)
		return
	}

	d.logger.Info("Validate OTP handled successfully")
	globalResponse(
		c,
		http.StatusOK,
		nil,
		response,
		http.StatusText(http.StatusOK),
	)
}

// ResetPassword handles the reset password functionality.
// It validates the reset password request, calls the reset password method,
// and returns a global response indicating the success or failure of the operation.
func (d JWTAuthController) ResetPassword(c *gin.Context) {
	var request dto.AuthResetPasswordDto
	var userIdString = c.Request.Header.Get("user_id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		d.logger.Error(err.Error())
		c.Error(err)
		return
	}

	d.logger.Info("Validate reset password request")
	if err := c.ShouldBind(&request); err != nil {
		d.logger.Error(err.Error())
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.AuthResetPasswordDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	d.logger.Debug("User password: ", request.Password)

	request.UserId = userId
	request.ResetToken = c.Request.Header.Get("reset_token")

	d.logger.Info("Calling reset password method")
	if err := d.service.ResetPassword(request); err != nil {
		d.logger.Error(err.Error())
		c.Error(err)
		return
	}

	d.logger.Info("Reset password handled successfully")
	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		http.StatusText(http.StatusOK),
	)
}

func (d JWTAuthController) Logout(c *gin.Context) {
	d.logger.Info("Getting auth bearer token")

	var request dto.LogoutRequestDto

	authorization := c.Request.Header.Get("Authorization")
	token := strings.Replace(authorization, "Bearer ", "", -1)

	if err := c.ShouldBind(&request); err != nil {
		d.logger.Error(err.Error())
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.LogoutRequestDto{}),
			nil,
			http.StatusText(http.StatusBadRequest),
		)
	}

	d.logger.Info("Removing token from DB")
	err := d.service.Logout(token, request.RefreshToken)
	if err != nil {
		d.logger.Error("Failed to remove token from DB")
		c.Error(err)
		return
	}

	d.logger.Info("Token removed successfully")

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		http.StatusText(http.StatusOK),
	)
}
