package models

// Item is a struct that represents an item
type Item struct {
	Model
	Name  string  `db:"name" json:"name"`
	Icon  string  `db:"icon" json:"icon"`
	Price float64 `db:"price" json:"price"`
}

// OrderedItem is when an item is ordered
type OrderedItem struct {
	Model
	OrderedItem Item `db:"item" json:"item"`
	Quantity    int  `db:"qty" json:"quantity"`
}
