package entities

const (
	PokemonTypeTableName = "pokemon_types"
)

type PokemonType struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	PokemonID int
}
