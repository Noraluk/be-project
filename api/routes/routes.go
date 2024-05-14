package routes

import (
	"be-project/api/middlewares"
	"be-project/api/models/response"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewRoutes(app *fiber.App) {
	handler := NewHandler()

	apiGroup := app.Group("/api")
	apiGroup.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(response.Response{Status: fiber.StatusOK})
	})
	apiGroup.Post("/register", handler.auth.Register)
	apiGroup.Post("/login", handler.auth.Login)

	protectedGroup := apiGroup.Group("", middlewares.Protected(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		c.Locals("id", claims["id"])
		return c.Next()
	})

	pokemonGroup := protectedGroup.Group("/pokemons")
	pokemonItemGroup := protectedGroup.Group("/pokemon-items")

	// pokemon
	pokemonGroup.Get("", handler.pokemon.GetPokemons)
	pokemonGroup.Post("", handler.pokemon.CreatePokemon)
	pokemonGroup.Get("/:id", handler.pokemon.GetPokemon)
	pokemonGroup.Delete("/:id", handler.pokemon.DeletePokemon)

	// pokemon item
	pokemonItemGroup.Get("", handler.pokemon.GetPokemonItems)

	wsGroup := app.Group("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			c.Locals("username", string(c.Query("username")))
			c.Locals("target", string(c.Query("target")))
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	go handler.chat.CreateConnection()
	wsGroup.Get("", websocket.New(handler.chat.Broadcast))
}
