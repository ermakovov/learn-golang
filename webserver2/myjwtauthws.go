package webserver2

import (
	"errors"
	"fmt"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const contextKeyUser = "user"

func StartJWTAuthServer() {
	webApp := fiber.New()

	authHandler := &AuthHandler{storage: &AuthStorage{users: map[string]User{}}}

	publicGroup := webApp.Group("")
	publicGroup.Post("/register", authHandler.CreateUser)
	publicGroup.Post("/login", authHandler.AuthUser)

	authorizedGroup := webApp.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: jwtSecretKey},
		ContextKey: contextKeyUser,
	}))
	authorizedGroup.Get("/profile", authHandler.GetUserData)

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}

type (
	AuthHandler struct {
		storage *AuthStorage
	}

	// In-memory storage of created users
	AuthStorage struct {
		users map[string]User
	}

	User struct {
		Email    string
		Name     string
		password string
	}
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	req := CreateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if _, exists := h.storage.users[req.Email]; exists {
		return errors.New("User with provided email already exists")
	}

	h.storage.users[req.Email] = User{
		Email:    req.Email,
		Name:     req.Name,
		password: req.Password,
	}

	return c.SendStatus(fiber.StatusCreated)
}

type (
	AuthUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthUserResponse struct {
		AccessToken string `json:"access_token"`
	}
)

var (
	errBadCredentials = errors.New("email or password is incorrect")
	jwtSecretKey      = []byte("secret-phrase")
)

func (h *AuthHandler) AuthUser(c *fiber.Ctx) error {
	req := AuthUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	user, exists := h.storage.users[req.Email]
	if !exists {
		return errBadCredentials
	}
	if user.password != req.Password {
		return errBadCredentials
	}

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logrus.WithError(err).Error("JWT signing")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(AuthUserResponse{AccessToken: signedToken})
}

type GetUserDataResponse struct {
	Email string `json:"email"`
	Name  string `json:"Name"`
}

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := c.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Context().Value(contextKeyUser),
		}).Error("wrong type of token in context")

		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of token claims")
		return nil, false
	}

	return payload, true
}

func (h *AuthHandler) GetUserData(c *fiber.Ctx) error {
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userData, ok := h.storage.users[jwtPayload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return c.JSON(GetUserDataResponse{
		Email: userData.Email,
		Name:  userData.Name,
	})
}
