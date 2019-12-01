package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

type orderedItem struct {
	ID        *int64 `db:"id" json:"id"`
	ItemID    *int64 `db:"item_id" json:"item_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	ItemName  string `db:"item_name" json:"item_name"`
	ItemIcon  string `db:"item_icon" json:"item_icon"`
	ItemPrice int    `db:"item_price" json:"item_price"`
}

type operationItem struct {
	ID            *int64 `db:"id" json:"id"`
	OperationType int    `db:"operation_type" json:"operation_type"`
	Amount        int    `db:"amt" json:"amount"`
}

type order struct {
	ID         int64           `db:"id" json:"id"`
	CreatedAt  *time.Time      `db:"created_at" json:"created_at"`
	Items      []orderedItem   `json:"items"`
	Operations []operationItem `json:"operations"`
}

// ListOrderRoute lists all order made
func ListOrderRoute(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rq := `SELECT * FROM ORDERS`
		res, err := prv.Db.Queryx(rq)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var orders []order = []order{}
		for res.Next() {
			var currOrder order
			res.StructScan(&currOrder)

			rqCurr := ` SELECT oi.ID, oi.ITEM_ID, oi.QUANTITY, i.NAME AS ITEM_NAME, i.ICON AS ITEM_ICON, i.PRICE AS ITEM_PRICE
					    FROM ORDER_ITEM oi
							   INNER JOIN ITEMS i ON i.ID = oi.ITEM_ID
						WHERE oi.ORDER_ID = $1`

			currOrder.Items = []orderedItem{}
			resCurr, err := prv.Db.Queryx(rqCurr, currOrder.ID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for resCurr.Next() {
				var currOrderedItem orderedItem
				resCurr.StructScan(&currOrderedItem)
				currOrder.Items = append(currOrder.Items, currOrderedItem)
			}

			rqCurr = `  SELECT ID, OPERATION_TYPE, AMT
					    FROM ORDER_OPERATION
						WHERE ORDER_ID = $1`

			currOrder.Operations = []operationItem{}
			resCurr, err = prv.Db.Queryx(rqCurr, currOrder.ID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for resCurr.Next() {
				var currOp operationItem
				resCurr.StructScan(&currOp)
				currOrder.Operations = append(currOrder.Operations, currOp)
			}

			orders = append(orders, currOrder)
		}

		jsonTx, err := json.Marshal(orders)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonTx)

	}
}
