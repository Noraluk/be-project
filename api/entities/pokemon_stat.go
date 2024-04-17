package entities

const (
	PokemonStatTableName = "pokemon_stats"
)

type PokemonStat struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	BaseStat  int
	PokemonID int
}
