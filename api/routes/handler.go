package routes

import (
	"be-project/api/handlers"
	"be-project/api/services"
	"be-project/pkg/base"
)

type handler struct {
	pokemon handlers.PokemonHandler
}

func NewHandler() handler {
	// repository
	repository := base.NewBaseRepository[any]()

	// service
	pokemonService := services.NewPokemonService(repository)

	return handler{
		pokemon: handlers.NewPokemonHandler(pokemonService),
	}
}
