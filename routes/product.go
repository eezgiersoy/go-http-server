package routes

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber}
}

func CreateProduct(ctx *fiber.Ctx) error {
	var product models.Product

	if err := ctx.BodyParser(&product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	db.Database.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)

	return ctx.Status(200).JSON(responseProduct)
}

func GetProducts(ctx *fiber.Ctx) error {
	var products []models.Product
	var responseProducts []Product

	db.Database.Db.Find(&products)

	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return ctx.Status(200).JSON(responseProducts)
}

func findProduct(id int, product *models.Product) error {
	db.Database.Db.Find(&product, "id = ?", id)
	if product.ID == 0 {
		return errors.New("product does not exist")
	}

	return nil
}

func GetProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var product models.Product

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that ID is an integer")
	}

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)

	return ctx.Status(200).JSON(responseProduct)
}

func UpdateProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var product models.Product

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that ID is an integer")
	}

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	type UpdateProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	var updateData UpdateProduct

	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(500).JSON(err.Error())
	}

	product.Name = updateData.Name
	product.SerialNumber = updateData.SerialNumber

	db.Database.Db.Save(&product)

	responseProduct := CreateResponseProduct(product)
	return ctx.Status(200).JSON(responseProduct)
}

func DeleteProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var product models.Product

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	if err := db.Database.Db.Delete(&product).Error; err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	return ctx.Status(200).JSON("Successfully deleted product")
}
