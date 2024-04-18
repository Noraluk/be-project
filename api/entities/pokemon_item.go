package entities

const (
	PokemonItemTableName = "pokemon_items"
)

type PokemonItem struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	ItemID    int    `json:"item_id"`
	Name      string `json:"name"`
	Cost      int    `json:"cost"`
	SpriteURL string `json:"sprite_url"`
}
