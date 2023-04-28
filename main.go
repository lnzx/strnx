package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/keyauth/v2"
	. "github.com/lnzx/strnx/internal"
	"github.com/lnzx/strnx/tools"
	"log"
	"os"
)

var EnvKey = tools.IfThen(os.Getenv("KEY"), "123456")
var hashedAPIKey = sha256.Sum256([]byte(EnvKey))

var tokens = make(map[string]string)

func init() {
	StartAsync()
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Static("/", "./dist")

	api := app.Group("/api", keyauth.New(keyauth.Config{
		Filter: func(c *fiber.Ctx) bool {
			return c.Path() == "/api/login"
		},
		KeyLookup: "header:token",
		Validator: validator,
	}))

	api.Post("/login", login)
	api.Post("/logout", logout)

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
	if user.Password != EnvKey {
		return fiber.ErrUnauthorized
	}

	return c.JSON(fiber.Map{
		"username": user.Username,
		"token":    hashedAPIKey,
	})
}

func logout(c *fiber.Ctx) error {
	token := c.Get("token")
	delete(tokens, token)
	log.Println("logout", token)
	return nil
}

func validator(_ *fiber.Ctx, token string) (bool, error) {
	if token != "" {
		hashedKey := sha256.Sum256([]byte(token))
		if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
			return true, nil
		}
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
