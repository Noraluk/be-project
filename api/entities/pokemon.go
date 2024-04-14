package entities

type Pokemon struct {
	ID                                   int     `gorm:"primaryKey" json:"id"`
	Name                                 string  `json:"name"`
	SpriteFrontDefaultShowdownURL        string  `json:"sprite_front_default_showdown_url"`
	SpriteFrontDefaultOfficialArtworkURL string  `json:"sprite_front_default_official_artwork_url"`
	Height                               float64 `json:"height"`
	Weight                               float64 `json:"weight"`
	BaseExperience                       int     `json:"base_experience"`
	MinimumLevel                         int     `json:"minumum_level"`
	EvolvedPokemonID                     *int    `json:"evlved_pokemon_id"`
}
