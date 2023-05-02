package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/keyauth/v2"
	. "github.com/lnzx/strnx/internal"
	"github.com/lnzx/strnx/tools"
	"log"
	"os"
)

const (
	MaxAge = 3600 // 1 hour
)

var envKey = tools.IfThen(os.Getenv("KEY"), "123456")
var hashToken string

func init() {
	hashKey := sha256.Sum256([]byte(envKey))
	hashToken = hex.EncodeToString(hashKey[:32])
	StartAsync()
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Static("/", "./dist", fiber.Static{
		MaxAge: MaxAge,
	})

	api := app.Group("/api", keyauth.New(keyauth.Config{
		Filter: func(c *fiber.Ctx) bool {
			return c.Path() == "/api/login"
		},
		KeyLookup: "header:token",
		Validator: validator,
	}))

	api.Post("/login", login)

	api.Get("/wallets", GetWallets)
	api.Post("/wallets", AddWallet)
	api.Delete("/wallets", DelWallets)

	api.Get("/nodes", GetNodes)

	api.Get("/summary", Summary)

	// fix vue history router 404
	app.Static("*", "./dist/index.html")

	log.Fatal(app.Listen(":8080"))
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return fiber.ErrBadRequest
	}
	if err := tools.ValidateStruct(user); err != nil {
		return fiber.ErrBadRequest
	}
	if user.Password != envKey {
		return fiber.ErrUnauthorized
	}

	return c.JSON(fiber.Map{
		"username": user.Username,
		"token":    hashToken,
	})
}

func validator(_ *fiber.Ctx, token string) (bool, error) {
	if token != "" && token == hashToken {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
