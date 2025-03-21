package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/tsmdev/software-development/backend/go-project/constants"
	"gitlab.com/tsmdev/software-development/backend/go-project/domains"
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	httperror "gitlab.com/tsmdev/software-development/backend/go-project/error"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gitlab.com/tsmdev/software-development/backend/go-project/repository"
	"gitlab.com/tsmdev/software-development/backend/go-project/utils"
	"gorm.io/gorm"
)

// JWTAuthService service relating to authorization
type JWTAuthService struct {
	env               lib.Env
	logger            lib.Logger
	userService       domains.UserService
	userRepository    repository.UserRepository
	tokenStoreService domains.TokenStoreService
	otpService        domains.OTPService
	smtp              lib.SMTP
}

// NewJWTAuthService creates a new auth service
func NewJWTAuthService(
	env lib.Env,
	logger lib.Logger,
	userService domains.UserService,
	UserRepository repository.UserRepository,
	otpService domains.OTPService,
	smtp lib.SMTP,
	tokenStoreService domains.TokenStoreService) domains.AuthService {
	return JWTAuthService{
		env:               env,
		logger:            logger,
		userService:       userService,
		userRepository:    UserRepository,
		tokenStoreService: tokenStoreService,
		otpService:        otpService,
		smtp:              smtp,
	}
}

// WithTrx delegates transaction to repository database
func (s JWTAuthService) WithTrx(trxHandle *gorm.DB) domains.AuthService {
	s.tokenStoreService = s.tokenStoreService.WithTrx(trxHandle)
	return s
}

// Authorize authorizes the generated token
func (s JWTAuthService) Authorize(tokenString string, claimType string) (*dto.AuthIdentityDto, error) {
	s.logger.Info("Checking token")
	var ve *jwt.ValidationError
	var jwtSecretKey = s.env.JWTSecret
	if claimType == constants.TypeRefreshToken {
		jwtSecretKey = s.env.JwtRefreshSecret
	} else if claimType == constants.TypeResetToken {
		jwtSecretKey = s.env.JwtResetSecret
	}

	err := s.tokenStoreService.ValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httperror.NewHttpError("invalid token", "", http.StatusUnauthorized)
		}
		return nil, err
	}

	s.logger.Info("Parsing token")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("Validating token")
	if token.Valid {
		// get id from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("couldn't retrieve sub as an integer")
		}
		userId := fmt.Sprintf("%v", claims["sub"])
		s.logger.Debug("Claimed subject: ", userId)
		authIdentity := dto.AuthIdentityDto{
			UserID: userId,
			// RoleID:             fmt.Sprintf("%v", claims["role_id"]),
		}
		return &authIdentity, nil
	} else if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("token expired")
		}
	}
	return nil, errors.New("couldn't handle token")
}

// CreateToken creates jwt auth token
func (s JWTAuthService) CreateToken(user *models.User, claimType string) (tokenString string, err error) {
	var (
		token                 *jwt.Token
		tokenLifeTimeDuration = s.env.TokenLifetime
		jwtSecretKey          = s.env.JWTSecret
		jwtClaims             = jwt.MapClaims{
			"sub": user.ID,
		}
	)

	s.logger.Info("Validating claimType: ", claimType)
	if claimType == constants.TypeRefreshToken {
		tokenLifeTimeDuration = s.env.RefreshTokenLifetime
		jwtSecretKey = s.env.JwtRefreshSecret
	} else if claimType == constants.TypeResetToken {
		jwtSecretKey = s.env.JwtResetSecret
	} else {
		s.logger.Info("Generation maps claim for ", claimType)
		jwtClaims["sub"] = user.ID
		jwtClaims["name"] = user.Name
		jwtClaims["email"] = user.Email
	}

	s.logger.Info("Claims: ", jwtClaims)
	tokenLifetime := time.Now().Add(time.Second * time.Duration(tokenLifeTimeDuration))
	jwtClaims["exp"] = tokenLifetime.Unix()

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err = token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		s.logger.Error("JWT validation failed: ", err)
		return tokenString, err
	}

	return tokenString, nil
}

// CreateRefreshToken jwt refresh token
func (s JWTAuthService) CreateRefreshToken(user models.User) (refreshTokenString string, err error) {
	tokenLifetime := time.Now().Add(time.Second * time.Duration(s.env.RefreshTokenLifetime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": tokenLifetime.Unix(),
	})

	refreshTokenString, err = token.SignedString([]byte(s.env.JwtRefreshSecret))

	if err != nil {
		s.logger.Error("JWT validation failed: ", err)
		return refreshTokenString, err
	}

	return refreshTokenString, nil
}

// ValidateAuth validates user by email & password
func (s JWTAuthService) ValidateAuth(request *dto.AuthRequestDto, refreshTokenString string) (*dto.AuthResponseDto, error) {
	var user *models.User
	var err error

	if request != nil {
		s.logger.Info("Getting user data")
		lowercaseUsername := strings.ToLower(request.Username)
		user, err = s.userRepository.GetByUsername(lowercaseUsername)
		if err != nil {
			return nil, httperror.NewHttpError("invalid username or password", "", http.StatusUnauthorized)
		}

		s.logger.Info("Validating user password")
		s.logger.Debug("Password from request: %s", request.Password)
		s.logger.Debug("Hashed password from DB: %s", user.Password)

		s.logger.Info("Validating user password")
		if !utils.CheckPasswordHash(user.Password, request.Password) {
			return nil, httperror.NewHttpError("invalid username or password", "", http.StatusUnauthorized)
		}
	} else {

		s.logger.Info("Validating refresh token")
		authIdentity, err := s.Authorize(refreshTokenString, constants.TypeRefreshToken)
		if err != nil {
			return nil, err
		}

		err = s.tokenStoreService.ValidateToken(refreshTokenString)
		if err != nil {
			return nil, httperror.NewHttpError("invalid refresh token", "", http.StatusUnauthorized)
		}

		// Delete old token
		err = s.tokenStoreService.DeleteToken(refreshTokenString)
		if err != nil {
			return nil, err
		}

		s.logger.Info("Getting user data")
		userId, _ := strconv.Atoi(authIdentity.UserID)
		user, err = s.userRepository.GetByID(userId)
		if err != nil {
			return nil, err
		}
	}

	s.logger.Debug("User data: ", user)

	s.logger.Info("Generating auth token")
	token, err := s.CreateToken(user, constants.TypeAuthToken)
	if err != nil {
		return nil, err
	}

	err = s.SaveToken(token)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Generating refresh token")
	refreshToken, err := s.CreateToken(user, constants.TypeRefreshToken)
	if err != nil {
		return nil, err
	}

	err = s.SaveToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// SaveToken CreateToken creates jwt auth token
func (s JWTAuthService) SaveToken(token string) error {
	var tokenStore = models.TokenStore{
		Token: token,
	}
	return s.tokenStoreService.CreateToken(tokenStore)
}

// ChangePassword change user password
func (s JWTAuthService) ChangePassword(id int, request *dto.AuthChangePasswordDto) error {
	s.logger.Info("Getting user by ID: ", id)
	user, err := s.userService.GetOneUserById(id)
	if err != nil {
		return httperror.NewHttpError("User not found", "", http.StatusNotFound)
	}

	s.logger.Info("Validating user password")
	if !utils.CheckPasswordHash(user.Password, request.OldPassword) {
		return httperror.NewHttpError("invalid old password", "", http.StatusBadRequest)
	}

	return s.userService.UpdateUser(user.ID, dto.CreateUserRequest{
		Email:    user.Email,
		Password: request.NewPassword,
	})
}

// ForgotPassword validate user by email and send OTP code to email
func (s JWTAuthService) ForgotPassword(request dto.AuthForgotPasswordDto) error {
	user, err := s.userRepository.GetByEmail(request.Email)
	if err != nil {
		return httperror.NewHttpError(constants.ERROR_USER_NOT_FOUND, "", http.StatusNotFound)
	}

	_, err = s.otpService.GetByUserIdAndIsUsed(int(user.ID), false)
	if err == nil {
		return httperror.NewHttpError("OTP code already sent", "", http.StatusBadRequest)
	}

	s.logger.Info("Generating OTP code with random string number")
	otpCode := utils.GenerateRandomNumberString()

	// saving to db
	s.logger.Info("Saving OTP data to DB")
	err = s.otpService.Create(int(user.ID), otpCode)
	if err != nil {
		return httperror.NewHttpError(constants.ERROR_CREATING_OTP, "", http.StatusInternalServerError)
	}

	s.logger.Info("Sending OTP code to email")
	go func() {
		err := s.smtp.SendEmail(dto.SendEmailRequestDto{
			To:      request.Email,
			Subject: "Forgot Password",
			Body:    fmt.Sprintf("Your OTP code is %s", otpCode),
		})
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return nil
}

func (s JWTAuthService) ValidateOTP(request dto.ValidateOTPRequestDto) (*dto.ValidateOTPResponseDto, error) {
	s.logger.Info("Getting user by email: ", request.Email)
	user, err := s.userRepository.GetByEmail(request.Email)
	if err != nil {
		return nil, httperror.NewHttpError(constants.ERROR_USER_NOT_FOUND, "", http.StatusNotFound)
	}

	s.logger.Info("Getting OTP data by code: ", request.Code)
	otp, err := s.otpService.GetByCode(int(user.ID), request.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httperror.NewHttpError("OTP code not found", "", http.StatusNotFound)
		}
		return nil, httperror.NewHttpError(constants.ERROR_GETTING_OTP_BY_CODE, "", http.StatusInternalServerError)
	}

	// update otp data
	s.logger.Info("Updating OTP data")
	err = s.otpService.UpdateById(otp.ID, true)
	if err != nil {
		return nil, httperror.NewHttpError(constants.ERROR_UPDATING_OTP, "", http.StatusInternalServerError)
	}

	s.logger.Info("Generating reset token")
	resetToken, err := s.CreateToken(user, constants.TypeResetToken)
	if err != nil {
		return nil, httperror.NewHttpError("Error creating reset token", "", http.StatusInternalServerError)
	}

	// saving token to db
	var tokenStore = models.TokenStore{
		Token: resetToken,
	}

	s.logger.Info("Saving reset token to DB")
	err = s.tokenStoreService.CreateToken(tokenStore)
	if err != nil {
		return nil, httperror.NewHttpError("Error saving reset token", "", http.StatusInternalServerError)
	}

	response := &dto.ValidateOTPResponseDto{
		ResetToken: resetToken,
	}

	return response, nil
}

func (s JWTAuthService) ResetPassword(request dto.AuthResetPasswordDto) error {
	s.logger.Info("Getting user by user_id: ", request.UserId)
	user, err := s.userService.GetOneUserById(request.UserId)
	if err != nil {
		return httperror.NewHttpError(constants.ERROR_USER_NOT_FOUND, "", http.StatusNotFound)
	}

	// validate reset token exist
	s.logger.Info("Validating reset token")
	err = s.tokenStoreService.ValidateToken(request.ResetToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httperror.NewHttpError("Invalid token", "", http.StatusUnauthorized)
		}
		return err
	}

	s.logger.Info("Updating user password")
	err = s.userService.UpdateUser(user.ID, dto.CreateUserRequest{
		Email:    user.Email,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	s.logger.Info("Deleting reset token")
	err = s.tokenStoreService.DeleteToken(request.ResetToken)
	if err != nil {
		return err
	}

	return nil
}

func (s JWTAuthService) Logout(tokenString string, refreshToken string) error {
	err := s.tokenStoreService.DeleteToken(tokenString)
	if err != nil {
		s.logger.Error("Failed to delete")
		return err
	}
	return s.tokenStoreService.DeleteToken(refreshToken)
}
