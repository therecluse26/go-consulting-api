package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../database"
)

func Courses(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cast(c.id AS char(36)) AS id, c.code, c.name, c.description 
									FROM Courses c` }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func Course(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: "SELECT cast(id as char(36)) AS id, code, name, description FROM Courses WHERE code = '{{code}}'", Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func CourseSessions(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.code = '{{code}}' ORDER BY cs.session_number`, Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func CourseSession(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.code = '{{code}}' and cs.session_number = {{session_number}}`, Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func AllCourseData(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cast(c.id AS char(36)) AS id, c.code, c.name, c.description,
										(SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Course_Sessions cs 
    									WHERE cs.course_id = c.id ORDER BY cs.session_number FOR JSON PATH) as sessions
									FROM Courses c` }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

