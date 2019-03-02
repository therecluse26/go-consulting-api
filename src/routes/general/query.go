package general

import (
	"net/http"
	"../../database"
)

type FuncLoader struct {}

type QueryMeta struct {
	sqlQuery string
	sqlParams string
}

func (t *FuncLoader) SelectQuery (w http.ResponseWriter, r *http.Request) {

	sql := database.Statement{ Sql: `SELECT
  										p.id, p.title, p.description, p.start_date, p.due_date, ps.status_display,
										b.name as business_name, c.name as category_name
									FROM Consulting.Projects p
  										INNER JOIN Accounts.Businesses b on b.id = p.business_id
  										INNER JOIN Consulting.Categories c on c.id = p.category_id
  										INNER JOIN Consulting.Project_Status ps on ps.id = p.status_id` }

	database.SelectAndWriteJsonResponse(sql, w)
}
