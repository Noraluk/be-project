package dtos

type PokemonList struct {
	Pokemon
	PokemonTypes []pokemonType `gorm:"foreignKey:PokemonID;references:PokemonID" json:"pokemon_types"`
}

type PokemonDetail struct {
	ID                                   int               `json:"id"`
	PokemonID                            int               `json:"pokemon_id"`
	Name                                 string            `json:"name"`
	SpriteFrontDefaultOfficialArtworkURL string            `json:"sprite_front_default_official_artwork_url"`
	PokemonTypes                         []pokemonType     `gorm:"foreignKey:PokemonID;references:PokemonID" json:"pokemon_types"`
	PokemonAbilities                     []pokemonAbility  `gorm:"foreignKey:PokemonID;references:PokemonID" json:"pokemon_abilities"`
	Height                               float64           `json:"height"`
	Weight                               float64           `json:"weight"`
	BaseExperience                       int               `json:"base_experience"`
	PokemonWeaknesses                    []pokemonWeakness `gorm:"foreignKey:PokemonID;references:PokemonID" json:"pokemon_weaknesses"`
	PokemonStats                         []PokemonStat     `gorm:"foreignKey:PokemonID;references:PokemonID" json:"pokemon_stats"`
	BasePokemonID                        int               `json:"-"`
	EvolvedPokemonID                     int               `json:"-"`
	EvolvedPokemon                       *evolvedPokemon   `gorm:"foreignKey:BasePokemonID" json:"evolved_pokemon"`
	NextPokemon                          *Pokemon          `json:"next_pokemon" gorm:"-"`
	PrevPokemon                          *Pokemon          `json:"prev_pokemon" gorm:"-"`
}

type pokemonType struct {
	Name      string `json:"name"`
	PokemonID int    `json:"-"`
}

type pokemonAbility struct {
	Name      string `json:"name"`
	PokemonID int    `json:"-"`
}

type pokemonWeakness struct {
	Name      string `json:"name"`
	PokemonID int    `json:"-"`
}

type PokemonStat struct {
	Name      string `json:"name"`
	BaseStat  int    `json:"base_stat"`
	PokemonID int    `json:"-"`
}

type evolvedPokemon struct {
	ID                                   int             `json:"id" gorm:"primaryKey"`
	SpriteFrontDefaultOfficialArtworkURL string          `json:"sprite_front_default_official_artwork_url"`
	MinimumLevel                         int             `json:"minumum_level"`
	EvolvedPokemonID                     int             `json:"-"`
	EvolvedPokemon                       *evolvedPokemon `json:"evolved_pokemon" gorm:"foreignKey:EvolvedPokemonID"`
}

type Pokemon struct {
	ID                            int    `json:"id"`
	PokemonID                     int    `json:"pokemon_id"`
	Name                          string `json:"name"`
	SpriteFrontDefaultShowdownURL string `json:"sprite_front_default_showdown_url"`
}
