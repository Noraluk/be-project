package services

import (
	"be-project/api/dtos"
	"be-project/api/entities"
	"be-project/api/utils"
	"be-project/pkg/base"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PokemonService interface {
	GetPokemons(c *fiber.Ctx) ([]dtos.PokemonList, int64, error)
}

type pokemonService struct {
	repository base.BaseRepository[any]
}

func NewPokemonService(repository base.BaseRepository[any]) PokemonService {
	return &pokemonService{
		repository: repository,
	}
}

func (s pokemonService) GetPokemons(c *fiber.Ctx) ([]dtos.PokemonList, int64, error) {
	var pokemons []dtos.PokemonList

	err := s.repository.Table(entities.PokemonTableName).Scopes(utils.Paginate(c)).Preload("PokemonTypes", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonTypeTableName)
	}).Order("ID asc").Find(&pokemons).Error()
	if err != nil {
		return nil, 0, err
	}

	var totalRecords int64
	err = s.repository.Table(entities.PokemonTableName).Count(&totalRecords).Error()
	if err != nil {
		return nil, 0, err
	}

	return pokemons, totalRecords, nil
}
