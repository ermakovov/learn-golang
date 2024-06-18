package webserver

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	CreateOrderRequest struct {
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}

	CreateOrderResponse struct {
		ID string `json:"id"`
	}

	GetOrderResponse struct {
		ID         string  `json:"id"`
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}
)

func StartSimpleStorageServer() {
	webApp := fiber.New()

	orderHandler := &OrderHandler{
		storage: &OrderStorage{
			orders: make(map[string]Order),
		},
	}

	webApp.Post("/orders", orderHandler.CreateOrder)
	webApp.Get("/orders/:id", orderHandler.GetOrder)

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}

type OrderCreatorGetter interface {
	CreateOrder(order Order) (string, error)
	GetOrder(orderID string) (Order, error)
}

type OrderHandler struct {
	storage OrderCreatorGetter
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parsing: %w", err)
	}

	order := Order{
		ID:         uuid.NewString(),
		UserID:     req.UserID,
		ProductIDs: req.ProductIDs,
	}
	orderID, err := h.storage.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("order creation: %w", err)
	}

	return c.JSON(CreateOrderResponse{ID: orderID})
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")

	order, err := h.storage.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	return c.JSON(GetOrderResponse(order))
}

// Order model
type Order struct {
	ID         string
	UserID     int64
	ProductIDs []int64
}

// Storage
type OrderStorage struct {
	mu     sync.Mutex
	orders map[string]Order
}

func (o *OrderStorage) CreateOrder(order Order) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.orders[order.ID] = order

	return order.ID, nil
}

var errOrderNotFound = errors.New("order not found")

func (o *OrderStorage) GetOrder(orderID string) (Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, ok := o.orders[orderID]
	if !ok {
		return Order{}, errOrderNotFound
	}

	return order, nil
}
