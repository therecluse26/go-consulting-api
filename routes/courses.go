package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../database"
)

func GetCourses(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cast(c.id AS char(36)) AS id, c.code, c.name, c.description 
									FROM Courses c` }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetCourse(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: "SELECT cast(id as char(36)) AS id, code, name, description FROM Courses WHERE code = '{{code}}'", Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func NewCourse(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: `INSERT INTO Courses (code, name, description) VALUES ('{{code}}', '{{name}}', '{{description}}')`, Params: courseData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)
	}

}

func GetCourseSessions(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.code = '{{code}}' ORDER BY cs.session_number`, Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetCourseSession(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.code = '{{code}}' and cs.session_number = {{session_number}}`, Params: mux.Vars(r) }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func NewCourseSession(w http.ResponseWriter, r *http.Request){

	r.ParseForm()


	courseSessionData := map[string]string{ "course_id": r.Form.Get("course_id"),  "session_number": r.Form.Get("session_number"), "title": r.Form.Get("title"), "description": r.Form.Get("description"), "start_datetime": r.Form.Get("start_datetime") }

	sql := database.Statement{ Sql: `INSERT INTO Course_Sessions (course_id, session_number, title, description, start_datetime) VALUES ('{{course_id}}', '{{session_number}}', '{{title}}', '{{description}}', '{{start_datetime}}')`, Params: courseSessionData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}
	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)
	}
}

func GetAllCourseData(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cast(c.id AS char(36)) AS id, c.code, c.name, c.description,
										(SELECT CAST(cs.id as char(36)) AS id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Course_Sessions cs 
    									WHERE cs.course_id = c.id ORDER BY cs.session_number FOR JSON PATH) as sessions
									FROM Courses c` }

	result := database.DbSelect(database.Dbconn, sql)

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

