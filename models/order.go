package models

// Order is an order
type Order struct {
	Model
	Items          []Item  `db:"Item"`
	PaidByLydia    float64 `db:"PaidByLydia"`
	PaidByCash     float64 `db:"PaidByCash"`
	GivenBackLydia float64 `db:"GivenBackLydia"`
	GivenBackCash  float64 `db:"GivenBackCash"`
}
