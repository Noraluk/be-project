package services

import (
	"be-project/api/constants"
	"be-project/api/dtos"
	"be-project/api/entities"
	"be-project/api/models/request"
	"be-project/api/utils"
	"be-project/pkg/base"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PokemonService interface {
	GetPokemons(c *fiber.Ctx) ([]dtos.PokemonList, int64, error)
	GetPokemon(pokemonID int) (dtos.PokemonDetail, error)
	CreatePokemon(req request.CreatedPokemon) error
	DeletePokemon(c *fiber.Ctx, id int) error
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
		Group("pokemons.id")

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
		Order(fmt.Sprintf("%s %s", c.Query("sort_by"), c.Query("sort_order"))).
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

func (s pokemonService) CreatePokemon(req request.CreatedPokemon) error {
	var pokemonTypes []entities.PokemonType
	var pokemonAbilities []entities.PokemonAbility
	var pokemonWeaknesses []entities.PokemonWeakness
	var pokemonStats []entities.PokemonStat

	channelWg := &sync.WaitGroup{}
	stopCh := make(chan struct{}, 1)
	pokemonTypeCh := make(chan entities.PokemonType, 1)
	pokemonAbilityCh := make(chan entities.PokemonAbility, 1)
	pokemonWeaknessCh := make(chan []entities.PokemonWeakness)
	pokemonStatCh := make(chan entities.PokemonStat, 1)

	channelWg.Add(1)
	go func() {
		defer channelWg.Done()

		for {
			select {
			case pt := <-pokemonTypeCh:
				pokemonTypes = append(pokemonTypes, pt)
			case pa := <-pokemonAbilityCh:
				pokemonAbilities = append(pokemonAbilities, pa)
			case pw := <-pokemonWeaknessCh:
				pokemonWeaknesses = append(pokemonWeaknesses, pw...)
			case ps := <-pokemonStatCh:
				pokemonStats = append(pokemonStats, ps)
			case <-stopCh:
				return
			}
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()

		pwMap := make(map[string]bool)
		for _, pokemonType := range req.PokemonTypes {
			pokemonTypeCh <- entities.PokemonType{
				Name: pokemonType,
			}

			pws := []entities.PokemonWeakness{}
			for _, v := range constants.PokemonWeakness[pokemonType] {
				if !pwMap[v] {
					pws = append(pws, entities.PokemonWeakness{
						Name: v,
					})
				}
				pwMap[v] = true
			}
			pokemonWeaknessCh <- pws
		}
	}()

	go func() {
		defer wg.Done()
		for _, ability := range req.PokemonAbilities {
			pokemonAbilityCh <- entities.PokemonAbility{
				Name: ability,
			}
		}
	}()

	go func() {
		defer wg.Done()
		for _, stat := range req.PokemonStats {
			pokemonStatCh <- entities.PokemonStat{
				BaseStat: stat.BaseStat,
				Name:     stat.Name,
			}
		}
	}()

	wg.Wait()
	close(stopCh)
	channelWg.Wait()

	var maxID int
	err := s.repository.Table(entities.PokemonTableName).Select("max(id)").Scan(&maxID).Error()
	if err != nil {
		return err
	}

	pokemon := entities.Pokemon{
		ID:                                   maxID + 1,
		PokemonID:                            req.PokemonID,
		Name:                                 req.Name,
		SpriteFrontDefaultShowdownURL:        req.SpriteFrontDefaultShowdownURL,
		SpriteFrontDefaultOfficialArtworkURL: req.SpriteFrontDefaultOfficialArtworkURL,
		Height:                               req.Height,
		Weight:                               req.Weight,
		BaseExperience:                       req.BaseExperience,
		MinimumLevel:                         req.MinimumLevel,
		EvolvedPokemonID:                     req.EvolvedPokemonID,
		PokemonTypes:                         pokemonTypes,
		PokemonAbilities:                     pokemonAbilities,
		PokemonWeaknesses:                    pokemonWeaknesses,
		PokemonStats:                         pokemonStats,
	}

	err = s.repository.Create(&pokemon).Error()
	if err != nil {
		return err
	}

	return nil
}

func (s pokemonService) DeletePokemon(c *fiber.Ctx, id int) error {
	var pokemon entities.Pokemon
	err := s.repository.First(&pokemon, id).Error()
	if err != nil {
		return err
	}

	err = s.repository.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&entities.PokemonAbility{}, "pokemon_id = ?", pokemon.PokemonID).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&entities.PokemonStat{}, "pokemon_id = ?", pokemon.PokemonID).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&entities.PokemonType{}, "pokemon_id = ?", pokemon.PokemonID).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&entities.PokemonWeakness{}, "pokemon_id = ?", pokemon.PokemonID).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&entities.Pokemon{}, "id = ?", id).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
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
