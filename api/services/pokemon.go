package services

import (
	"be-project/api/entities"
	"be-project/api/utils"
	"be-project/pkg/base"

	"github.com/gofiber/fiber/v2"
)

type PokemonService interface {
	GetPokemons(c *fiber.Ctx) ([]entities.Pokemon, int64, error)
}

type pokemonService struct {
	repository base.BaseRepository[any]
}

func NewPokemonService(repository base.BaseRepository[any]) PokemonService {
	return &pokemonService{
		repository: repository,
	}
}

func (s pokemonService) GetPokemons(c *fiber.Ctx) ([]entities.Pokemon, int64, error) {
	var pokemons []entities.Pokemon

	err := s.repository.Scopes(utils.Paginate(c)).Order("ID asc").Find(&pokemons).Error()
	if err != nil {
		return nil, 0, err
	}

	var totalRecords int64
	err = s.repository.Model(&entities.Pokemon{}).Count(&totalRecords).Error()
	if err != nil {
		return nil, 0, err
	}

	return pokemons, totalRecords, nil
}
