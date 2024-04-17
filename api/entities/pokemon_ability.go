package entities

const (
	PokemonAbilityTableName = "pokemon_abilities"
)

type PokemonAbility struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	PokemonID int
}
