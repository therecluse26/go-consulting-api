package consulting

import (
	"github.com/therecluse26/fortisure-api/src/database"
	"github.com/gorilla/mux"
	"net/http"
)

func GetAllProjects(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT
  										p.id, p.title, p.description, p.start_date, p.due_date, ps.status_display,
										b.name as business_name, c.name as category_name
									FROM Consulting.Projects p
  										INNER JOIN Accounts.Businesses b on b.id = p.business_id
  										INNER JOIN Consulting.Categories c on c.id = p.category_id
  										INNER JOIN Consulting.Project_Status ps on ps.id = p.status_id` }

	database.SelectAndWriteJsonResponse(sql, w)

}

func GetProject(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT
  										p.id, p.title, p.description, p.start_date, p.due_date, ps.status_display,
										b.name as business_name, c.name as category_name, t.name as team_name,
  										co.first_name + ' ' + co.last_name as contact_name, co.email as contact_email
									FROM Consulting.Projects p
  										INNER JOIN Accounts.Businesses b on b.id = p.business_id
  										INNER JOIN Consulting.Categories c on c.id = p.category_id
  										INNER JOIN Consulting.Project_Status ps on ps.id = p.status_id
										INNER JOIN Consulting.Teams t on t.id = p.team_assigned
  										INNER JOIN Accounts.Contacts co on co.id = p.primary_contact
									WHERE p.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndWriteJsonResponse(sql, w)

}