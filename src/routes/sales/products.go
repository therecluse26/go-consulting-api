package sales

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/therecluse26/fortisure-api/src/database"
	"github.com/therecluse26/fortisure-api/src/util"
)

func GetProducts(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT p.id, cast(p.price as varchar) as price, p.name, p.category, p.description FROM Sales.Products p` }

	database.SelectAndWriteJsonResponse(sql, w)

}

func GetProduct(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT p.id, cast(p.price as varchar) as price, p.name, p.category, p.description FROM Sales.Products p WHERE p.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndWriteJsonResponse(sql, w)

}

func NewProduct(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	productData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: ``, Params: productData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		util.ErrorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))	}

}


func GetOrders(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT o.id, o.user_id, u.username, o.status, sum(oli.count) as item_count,
     									 cast(sum(oli.count * p.price) as varchar) as order_total, o.notes, o.created_on, o.updated_on
									FROM Sales.Orders o
  										INNER JOIN Sales.Order_Line_Items oli on o.id = oli.order_id
  										INNER JOIN Sales.Products p on oli.product_id = p.id
										INNER JOIN People.Users u on u.id = o.user_id
									GROUP BY o.id, o.user_id, u.username, o.status, o.notes, o.created_on, o.updated_on` }

	database.SelectAndWriteJsonResponse(sql, w)

}

func GetOrder(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT o.id, o.user_id, u.username, o.status, sum(oli.count) as item_count,
     									 cast(sum(oli.count * p.price) as varchar) as order_total, o.notes, o.created_on, o.updated_on
									FROM Sales.Orders o
  										INNER JOIN Sales.Order_Line_Items oli on o.id = oli.order_id
  										INNER JOIN Sales.Products p on oli.product_id = p.id
										INNER JOIN People.Users u on u.id = o.user_id
										WHERE o.id = {{id}}
									GROUP BY o.id, o.user_id, u.username, o.status, o.notes, o.created_on, o.updated_on`, Params: mux.Vars(r) }

	database.SelectAndWriteJsonResponse(sql, w)

}

func GetOrderDetails(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT o.id, o.user_id, u.username, o.status, o.notes, o.created_on,
    									(SELECT p.id as product_id, p.name as product_name, oli.count, p.price
      										FROM Sales.Orders o
												INNER JOIN Sales.Order_Line_Items oli on o.id = oli.order_id
          										INNER JOIN Sales.Products p on oli.product_id = p.id
      										WHERE o.id = {{id}}
    										FOR JSON PATH, WITHOUT_ARRAY_WRAPPER) as line_items, o.updated_on
										FROM Sales.Orders o
   											INNER JOIN People.Users u on u.id = o.user_id
  										WHERE o.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndWriteJsonResponse(sql, w)

}