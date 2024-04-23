package constants

var (
	PokemonWeakness = map[string][]string{
		"normal":   {"fighting"},
		"fire":     {"water", "ground", "rock"},
		"water":    {"electric", "grass"},
		"electric": {"ground"},
		"grass":    {"fire", "ice", "poison", "flying", "bug"},
		"ice":      {"fire", "fighting", "rock", "steel"},
		"fighting": {"fighting", "psychic", "fairy"},
		"poison":   {"ground", "psychic"},
		"ground":   {"water", "grass", "ice"},
		"flying":   {"electric", "ice", "rock"},
		"psychic":  {"bug", "ghost", "dark"},
		"bug":      {"fire", "flying", "rock"},
		"rock":     {"water", "grass", "fighting", "ground", "steel"},
		"ghost":    {"ghost", "dark"},
		"dragon":   {"ice", "dragon", "fairy"},
		"dark":     {"fighting", "bug", "fairy"},
		"steel":    {"fire", "fighting", "ground"},
		"fairy":    {"poison", "steel"},
	}
)
