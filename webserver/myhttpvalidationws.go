package webserver

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      int64
	Email   string
	Age     int
	Country string
}

var users = map[int64]User{}

type (
	CreateUserRequest struct {
		// BEGIN (write your solution here)
		ID      int64  `json:"id" validate:"required,min=0"`
		Email   string `json:"email" validate:"required,email"`
		Age     int    `json:"age" validate:"required,gte=18,lte=130"`
		Country string `json:"country" validate:"required,allowable_country"`
		// END
	}
)

func (req *CreateUserRequest) toUser() User {
	return User{
		ID:      req.ID,
		Email:   req.Email,
		Age:     req.Age,
		Country: req.Country,
	}
}

func StartHTTPValidationServer() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		usersList := make([]User, 0, len(users))
		for _, user := range users {
			usersList = append(usersList, user)
		}

		return c.JSON(usersList)
	})

	// BEGIN (write your solution here) (write your solution here)
	var allowedCountries = []string{"USA", "Germany", "France"}
	validate := validator.New()
	vErr := validate.RegisterValidation("allowable_country", func(fl validator.FieldLevel) bool {
		country := fl.Field().String()
		for _, allowedCountry := range allowedCountries {
			if country == allowedCountry {
				return true
			}
		}

		return false
	})

	if vErr != nil {
		logrus.Fatal("register validation")
	}

	webApp.Post("/users", func(ctx *fiber.Ctx) error {
		var req CreateUserRequest
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		err := validate.Struct(req)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		users[req.ID] = req.toUser()

		return ctx.SendStatus(fiber.StatusOK)
	})
	// END
	logrus.Fatal(webApp.Listen(":8080"))
}
