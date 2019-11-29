package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oxodao/caisseoverflow/dal"
	"github.com/oxodao/caisseoverflow/services"
)

// ListItemsRoute lists all items
func ListItemsRoute(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := dal.GetItems(prv)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			return
		}

		it, err := json.Marshal(items)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(it)
	}
}
