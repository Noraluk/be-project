package request

type CreatedPokemon struct {
	PokemonID                            int      `json:"pokemon_id"`
	Name                                 string   `json:"name"`
	SpriteFrontDefaultOfficialArtworkURL string   `json:"sprite_front_default_official_artwork_url"`
	SpriteFrontDefaultShowdownURL        string   `json:"sprite_front_default_showdown_url"`
	PokemonTypes                         []string `json:"pokemon_types"`
	PokemonAbilities                     []string `json:"pokemon_abilities"`
	Height                               float64  `json:"height"`
	Weight                               float64  `json:"weight"`
	BaseExperience                       int      `json:"base_experience"`
	PokemonWeaknesses                    []string `json:"pokemon_weaknesses"`
	PokemonStats                         []struct {
		Name     string `json:"name"`
		BaseStat int    `json:"base_stat"`
	} `json:"pokemon_stats"`
	EvolvedPokemonID *int `json:"evolved_pokemon_id"`
	MinimumLevel     int  `json:"minimum_level"`
	BasePokemonID    int  `json:"base_pokemon_id"`
}
