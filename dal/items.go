package dal

import (
	"github.com/oxodao/caisseoverflow/models"
	"github.com/oxodao/caisseoverflow/services"
)

// GetItems returns the list of sold items
func GetItems(prv *services.Provider) ([]models.Item, error) {
	rq := `SELECT ID, CREATED_AT, NAME, ICON, PRICE
		   FROM ITEMS`

	var items []models.Item = []models.Item{}
	rows, err := prv.Db.Queryx(rq)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		var item models.Item
		rows.StructScan(&item)

		items = append(items, item)
	}

	return items, err
}
