package webserver

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var exchangeRate = map[string]float64{
	"USD/EUR": 0.8,
	"EUR/USD": 1.25,
	"USD/GBP": 0.7,
	"GBP/USD": 1.43,
	"USD/JPY": 110,
	"JPY/USD": 0.0091,
}

func StartCurrExchangeServer() {
	currUnknown := "unknown"
	webApp := fiber.New()

	webApp.Get("/convert", func(c *fiber.Ctx) error {
		from := c.Query("from", currUnknown)
		to := c.Query("to", currUnknown)

		if from == "" || to == "" || from == currUnknown || to == currUnknown {
			return c.SendStatus(fiber.StatusNotFound)
		}

		currPair := from + "/" + to

		currRate, ok := exchangeRate[currPair]
		if !ok {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.SendString(fmt.Sprintf("%.2f", currRate))
	})

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}
