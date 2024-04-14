package handlers

import (
	"be-project/api/models"
	"be-project/api/services"

	"github.com/gofiber/fiber/v2"
)

type PokemonHandler interface {
	GetPokemons(c *fiber.Ctx) error
}

type pokemonHandler struct {
	pokemonService services.PokemonService
}

func NewPokemonHandler(
	pokemonService services.PokemonService,
) PokemonHandler {
	return &pokemonHandler{
		pokemonService: pokemonService,
	}
}

func (h pokemonHandler) GetPokemons(c *fiber.Ctx) error {
	pokemons, totalRecords, err := h.pokemonService.GetPokemons(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Data: pokemons,
	}.ToPagination(c, totalRecords))
}
