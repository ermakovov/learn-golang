package webserver

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	CreateLinkRequest struct {
		ExtLink string `json:"external"`
		IntLink string `json:"internal"`
	}

	GetLinkResponse struct {
		IntLink string `json:"internal"`
	}
)

func StartURLExchangerServer() {
	webApp := fiber.New()

	linkHandler := &LinkHandler{
		storage: &LinkStorage{
			links: make(map[string]string),
		},
	}

	webApp.Post("/links", linkHandler.CreateLink)
	webApp.Get("/links/:extLink", linkHandler.GetLink)

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}

type LinkCreatorGetter interface {
	CreateLink(extLink, intLink string) error
	GetLink(extLink string) (string, error)
}

type LinkHandler struct {
	storage LinkCreatorGetter
}

func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var req CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	h.storage.CreateLink(req.ExtLink, req.IntLink)

	return c.SendStatus(fiber.StatusOK)
}

func (h *LinkHandler) GetLink(c *fiber.Ctx) error {
	extLink, err := url.QueryUnescape(c.Params("extLink"))
	if err != nil {
		return fmt.Errorf("link escaping: %w", err)
	}

	intLink, err := h.storage.GetLink(extLink)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Link not found")
	}

	return c.JSON(GetLinkResponse{IntLink: intLink})
}

// Storage
type LinkStorage struct {
	links map[string]string
}

var errLinkNotCreated = "link not created"

func (ls *LinkStorage) CreateLink(extLink, intLink string) error {
	ls.links[extLink] = intLink
	if ls.links[extLink] != intLink {
		return fmt.Errorf(errLinkNotCreated)
	}

	return nil
}

var errInternalLinkNotFound = "internal link not found"

func (ls *LinkStorage) GetLink(extLink string) (string, error) {
	intLink, ok := ls.links[extLink]
	if !ok {
		return "", fmt.Errorf(errInternalLinkNotFound)
	}

	return intLink, nil
}
