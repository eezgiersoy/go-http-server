package routes

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

func CreateUser(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	db.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return ctx.Status(200).JSON(responseUser)
}

func GetUsers(ctx *fiber.Ctx) error {
	var users []models.User
	var responseUsers []User

	db.Database.Db.Find(&users)

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return ctx.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	db.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func GetUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var user models.User

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that ID is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return ctx.Status(200).JSON(responseUser)
}

func UpdateUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var user models.User

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that ID is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateData UpdateUser

	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(500).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	db.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return ctx.Status(200).JSON(responseUser)
}

func DeleteUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var user models.User

	if err != nil {
		return ctx.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	if err := db.Database.Db.Delete(&user).Error; err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	return ctx.Status(200).JSON("Successfully deleted user")
}
