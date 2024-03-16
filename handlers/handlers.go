package handlers

import (
	"fmt"
	"time"

	"github.com/bwhitney2439/muzz/database"
	"github.com/bwhitney2439/muzz/models"
	"github.com/bwhitney2439/muzz/types"
	"github.com/bwhitney2439/muzz/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pioz/faker"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	// users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   true,
	})
}

// UserCreate registers a user
func CreateUser(c *fiber.Ctx) error {

	user := new(models.User)
	var userInput struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		Name        string `json:"name"`
		Gender      string `json:"gender"`
		DateOfBirth string `json:"date_of_birth"`
	}

	userInput.Email = faker.SafeEmail()
	userInput.Password = faker.Username()
	userInput.Name = faker.FullName()
	genders := []string{"Male", "Female", "Other"}

	userInput.Gender = genders[faker.IntInRange(0, len(genders)-1)]

	// Generate a random date within the range
	userInput.DateOfBirth = faker.Time().Format("02/01/2006")

	var err error
	var exists bool
	_, exists, err = database.GetUserByEmail(userInput.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error checking if user exists"})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already exists"})
	}

	layout := "02/01/2006"
	fmt.Println(userInput.DateOfBirth)
	t, err := time.Parse(layout, userInput.DateOfBirth)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse date of birth"})
	}
	age := utils.CalculateAge(t)

	user = &models.User{Email: userInput.Email, Password: userInput.Password, Name: userInput.Name, Gender: userInput.Gender, Age: uint8(age)}
	err = database.InsertUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot insert user"})
	}

	return c.JSON(fiber.Map{"result": user})
}

func LoginUser(c *fiber.Ctx) error {

	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	user, exists, err := database.GetUserByEmail(userInput.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error checking if user exists"})
	}
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not exist"})
	}

	if user.Password != userInput.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password"})

	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(utils.GoDotEnvVariable("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})

}

func Discover(c *fiber.Ctx) error {

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(float64)

	usersResponse := new([]types.UserResponse)

	err := database.GetUsers(usersResponse, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot get users"})
	}

	return c.JSON(fiber.Map{"results": usersResponse})
}

func Swipe(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(float64)

	var swipeInput struct {
		TargetUserID float64 `json:"targetUserId"`
		Preference   string  `json:"preference"`
	}

	if err := c.BodyParser(&swipeInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	matched, matchedUser, err := database.SwipeAction(uint(user_id), uint(swipeInput.TargetUserID), swipeInput.Preference)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot insert swipe"})
	}
	if matched {
		return c.JSON(fiber.Map{"results": fiber.Map{
			"matched": matched,
			"matchID": matchedUser.ID,
		}})
	}

	return c.JSON(fiber.Map{"results": fiber.Map{
		"matched": matched,
	}})
}
