package webserver

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	BinarySearchRequest struct {
		Numbers []int `json:"numbers"`
		Target  int   `json:"target"`
	}

	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)

const targetNotFound = -1

func StartArrayFinderServer() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go to /search")
	})

	webApp.Post("/search", func(c *fiber.Ctx) error {
		var req BinarySearchRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Invalid JSON",
			})
		}

		targetIndex := sort.SearchInts(req.Numbers, req.Target)
		if targetIndex >= len(req.Numbers) || targetIndex < len(req.Numbers) && req.Target != req.Numbers[targetIndex] {
			return c.Status(fiber.StatusNotFound).JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Target was not found",
			})
		}

		return c.JSON(BinarySearchResponse{
			TargetIndex: targetIndex,
		})
	})

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}
