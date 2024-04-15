package dtos

type PokemonList struct {
	ID                            int    `json:"id"`
	Name                          string `json:"name"`
	SpriteFrontDefaultShowdownURL string `json:"sprite_front_default_showdown_url"`
	PokemonTypes                  []struct {
		Name      string `json:"name"`
		PokemonID int    `json:"-"`
	} `gorm:"foreignKey:PokemonID" json:"pokemon_types"`
}
