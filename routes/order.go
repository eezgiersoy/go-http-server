package routes

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Order struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"order_date"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

func CreateOrder(ctx *fiber.Ctx) error {
	var order models.Order

	if err := ctx.BodyParser(&order); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	db.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return ctx.Status(200).JSON(responseOrder)
}

func GetOrders(ctx *fiber.Ctx) error {
	var orders []models.Order
	var responseOrders []Order
	db.Database.Db.Find(&orders)

	for _, order := range orders {
		var user models.User
		var product models.Product
		db.Database.Db.Find(&user, "id = ?", order.UserRefer)
		db.Database.Db.Find(&product, "id = ?", order.ProductRefer)
		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)
	}

	return ctx.Status(200).JSON(responseOrders)
}

func FindOrder(id int, order *models.Order) error {
	db.Database.Db.Find(&order, "id = ?", id)
	if order.ID == 0 {
		return errors.New("order does not exist")
	}

	return nil
}

func GetOrder(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	var order models.Order

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindOrder(id, &order); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	db.Database.Db.First(&user, order.UserRefer)
	db.Database.Db.First(&product, order.ProductRefer)
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)

	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return ctx.Status(200).JSON(responseOrder)
}

func UpdateOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	var order models.Order
	if err := FindOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateOrder struct {
		ProductRefer int `json:"product_id"`
	}

	var updateData UpdateOrder
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	order.ProductRefer = updateData.ProductRefer
	db.Database.Db.Save(&order)

	responseUser := CreateResponseUser(order.User)
	responseProduct := CreateResponseProduct(order.Product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}

func DeleteOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	var order models.Order
	if err := FindOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Database.Db.Delete(&order)

	return c.Status(204).Send(nil)
}
