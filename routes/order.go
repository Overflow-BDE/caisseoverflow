package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/oxodao/caisseoverflow/services"
)

type addOrderRequest struct {
	OrderedItems []struct {
		ID       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"orderedItems"`
	Operations []struct {
		Type   int `json:"type"`
		Amount int `json:"amt"`
	} `json:"operations"`
}

// AddOrderRoute add an order to the database
func AddOrderRoute(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var rc addOrderRequest

		err = json.Unmarshal(body, &rc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// For some reason, we can't insert empty parameters (Everything is generated)
		rq := `
			INSERT INTO ORDERS(CREATED_AT) VALUES (CURRENT_TIMESTAMP) RETURNING ID	
		`
		var id int64
		err = prv.Db.Get(&id, rq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tx, err := prv.Db.Beginx()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i := 0; i < len(rc.OrderedItems); i++ {
			rq := `INSERT INTO ORDER_ITEM (ORDER_ID, ITEM_ID, QUANTITY) VALUES ($1, $2, $3)`
			tx.MustExec(rq, id, rc.OrderedItems[i].ID, rc.OrderedItems[i].Quantity)
		}

		for i := 0; i < len(rc.Operations); i++ {
			rq := `INSERT INTO ORDER_OPERATION (ORDER_ID, OPERATION_TYPE, AMT) VALUES ($1, $2, $3)`
			tx.MustExec(rq, id, rc.Operations[i].Type, rc.Operations[i].Amount)
		}

		err = tx.Commit()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			rq := `DELETE FROM ORDERS WHERE ID = $1`
			prv.Db.Exec(rq, id)

			return
		}

		w.WriteHeader(201)
	}
}

// ListOrderRoute lists all order made
func ListOrderRoute(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//rq := ``
	}
}
