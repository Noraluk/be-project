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
	GetPokemon(pokemonID int) (dtos.PokemonDetail, error)
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
	}).Order("id asc").Find(&pokemons).Error()
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

func (s pokemonService) GetPokemon(id int) (dtos.PokemonDetail, error) {
	var pokemonDetail dtos.PokemonDetail
	err := s.repository.Table(entities.PokemonTableName).
		Preload("PokemonTypes", func(db *gorm.DB) *gorm.DB {
			return db.Table(entities.PokemonTypeTableName)
		}).Preload("PokemonAbilities", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonAbilityTableName)
	}).Preload("PokemonWeaknesses", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonWeaknessTableName)
	}).Preload("PokemonStats", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonStatTableName)
	}).Preload("EvolvedPokemon", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonTableName)
	}).Preload("EvolvedPokemon.EvolvedPokemon", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonTableName)
	}).Preload("EvolvedPokemon.EvolvedPokemon.EvolvedPokemon", func(db *gorm.DB) *gorm.DB {
		return db.Table(entities.PokemonTableName)
	}).Where("id = ?", id).First(&pokemonDetail).Error()
	if err != nil {
		return dtos.PokemonDetail{}, err
	}

	var nextPokemon entities.Pokemon
	err = s.repository.Where("id = ?", id+1).First(&nextPokemon).Error()
	if err != nil && err != gorm.ErrRecordNotFound {
		return dtos.PokemonDetail{}, err
	}
	pokemonDetail.NextPokemon = &dtos.Pokemon{
		ID:                            nextPokemon.ID,
		PokemonID:                     nextPokemon.PokemonID,
		Name:                          nextPokemon.Name,
		SpriteFrontDefaultShowdownURL: nextPokemon.SpriteFrontDefaultShowdownURL,
	}

	if id > 1 {
		var prevPokemon entities.Pokemon
		err = s.repository.Where("id = ?", id-1).First(&prevPokemon).Error()
		if err != nil && err != gorm.ErrRecordNotFound {
			return dtos.PokemonDetail{}, err
		}
		pokemonDetail.PrevPokemon = &dtos.Pokemon{
			ID:                            prevPokemon.ID,
			PokemonID:                     prevPokemon.PokemonID,
			Name:                          prevPokemon.Name,
			SpriteFrontDefaultShowdownURL: prevPokemon.SpriteFrontDefaultShowdownURL,
		}
	}

	var totalStat int = 0
	for _, v := range pokemonDetail.PokemonStats {
		totalStat += v.BaseStat
	}
	pokemonDetail.PokemonStats = append(pokemonDetail.PokemonStats, dtos.PokemonStat{
		Name:     "total",
		BaseStat: totalStat,
	})

	return pokemonDetail, nil
}
