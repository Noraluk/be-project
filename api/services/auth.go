package services

import (
	"be-project/api/entities"
	"be-project/api/models/request"
	"be-project/api/models/response"
	"be-project/pkg/base"
	"be-project/pkg/config"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(c *fiber.Ctx, req request.RegisterRequest) error
	Login(c *fiber.Ctx, req request.LoginRequest) (response.LoginResponse, error)
}

type authService struct {
	repository base.BaseRepository[any]
	config     config.Config
}

func NewAuthService(repository base.BaseRepository[any]) AuthService {
	return &authService{
		repository: repository,
		config:     config.GetConfig(),
	}
}

func (s authService) Register(c *fiber.Ctx, req request.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return err
	}

	auth := entities.Auth{
		Username: req.Username,
		Password: string(passwordHash),
	}
	db := s.repository.Where("username = ?", req.Username).FirstOrCreate(&auth)
	if err := db.Error(); err != nil {
		return err
	}

	if db.RowsAffected() == 0 {
		return errors.New("user was created")
	}

	return nil
}

func (s authService) Login(c *fiber.Ctx, req request.LoginRequest) (response.LoginResponse, error) {
	var auth entities.Auth
	err := s.repository.First(&auth, "username = ?", req.Username).Error()
	if err != nil {
		return response.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password))
	if err != nil {
		return response.LoginResponse{}, err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = auth.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(s.config.Auth.Secret))
	if err != nil {
		return response.LoginResponse{}, err
	}

	return response.LoginResponse{
		Token: t,
	}, nil
}
