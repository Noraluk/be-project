package services

import (
	"be-project/api/dtos"
	"be-project/api/entities"
	"be-project/api/utils"
	"be-project/pkg/base"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PokemonService interface {
	GetPokemons(c *fiber.Ctx) ([]dtos.PokemonList, int64, error)
	GetPokemon(pokemonID int) (dtos.PokemonDetail, error)
	GetPokemonItems(c *fiber.Ctx) ([]entities.PokemonItem, int64, error)
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

	db := s.repository.Table(entities.PokemonTableName).
		Joins("left join pokemon_types on pokemon_types.pokemon_id = pokemons.pokemon_id").
		Where("pokemons.name LIKE ?", fmt.Sprintf("%s%%", c.Query("name"))).
		Group("pokemons.pokemon_id")

	if len(c.Query("pokemon_type")) > 0 {
		db = db.Where("pokemon_types.name = ?", c.Query("pokemon_type"))
	}

	var totalRecords int64
	err := db.Count(&totalRecords).Error()
	if err != nil {
		return nil, 0, err
	}

	err = db.
		Preload("PokemonTypes", func(db *gorm.DB) *gorm.DB {
			return db.Table(entities.PokemonTypeTableName)
		}).
		Scopes(utils.Paginate(c)).
		Order("id asc").
		Find(&pokemons).
		Error()
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

func (s pokemonService) GetPokemonItems(c *fiber.Ctx) ([]entities.PokemonItem, int64, error) {
	var items []entities.PokemonItem
	err := s.repository.Scopes(utils.Paginate(c)).Order("id asc").Find(&items).Error()
	if err != nil {
		return nil, 0, err
	}

	var totalRecords int64
	err = s.repository.Table(entities.PokemonItemTableName).Count(&totalRecords).Error()
	if err != nil {
		return nil, 0, err
	}

	return items, totalRecords, nil
}
