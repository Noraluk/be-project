package routes

import "github.com/gofiber/fiber/v2"

func NewRoutes(app *fiber.App) {
	handler := NewHandler()

	apiGroup := app.Group("/api")

	// pokemon
	pokemonGroup := apiGroup.Group("/pokemons")
	pokemonGroup.Get("", handler.pokemon.GetPokemons)
	pokemonGroup.Get("/:id", handler.pokemon.GetPokemon)
}
