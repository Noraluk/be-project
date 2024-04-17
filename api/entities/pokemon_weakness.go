package entities

const (
	PokemonWeaknessTableName = "pokemon_weaknesses"
)

type PokemonWeakness struct {
	ID            int `gorm:"primaryKey"`
	Name          string
	PokemonID     int
	PokemonTypeID int
}
