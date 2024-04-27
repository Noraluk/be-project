package routes

import "github.com/gofiber/fiber/v2"

func NewRoutes(app *fiber.App) {
	handler := NewHandler()

	apiGroup := app.Group("/api")

	pokemonGroup := apiGroup.Group("/pokemons")
	pokemonItemGroup := pokemonGroup.Group("/items")

	// pokemon
	pokemonGroup.Get("", handler.pokemon.GetPokemons)
	pokemonGroup.Post("", handler.pokemon.CreatePokemon)
	pokemonGroup.Get("/:id", handler.pokemon.GetPokemon)
	pokemonGroup.Delete("/:id", handler.pokemon.DeletePokemon)

	// pokemon item
	pokemonItemGroup.Get("", handler.pokemon.GetPokemonItems)

	apiGroup.Post("/register", handler.auth.Register)
	apiGroup.Post("/login", handler.auth.Login)
}
