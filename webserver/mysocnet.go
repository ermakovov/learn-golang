package webserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var postLikes = map[string]int64{}

const postIdUnknown = "unknown"

func StartSocialNetworkServer() {
	webApp := fiber.New()

	webApp.Get("/likes/:post_id?", func(c *fiber.Ctx) error {
		postId := c.Params("post_id", postIdUnknown)
		if postId == postIdUnknown {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		likes, ok := postLikes[postId]
		if !ok {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.SendString(strconv.FormatInt(likes, 10))
	})

	webApp.Post("likes/:post_id?", func(c *fiber.Ctx) error {
		postId := c.Params("post_id", postIdUnknown)
		if postId == postIdUnknown {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, ok := postLikes[postId]
		if !ok {
			postLikes[postId] = 1
			createdLikes := postLikes[postId]

			return c.Status(fiber.StatusCreated).SendString(strconv.FormatInt(createdLikes, 10))
		}

		postLikes[postId]++

		return c.SendString(strconv.FormatInt(postLikes[postId], 10))
	})

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}
