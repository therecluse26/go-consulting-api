package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../database"
)

func GetCourses(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT c.id, c.code, c.name, c.description 
									FROM Courses c` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetCourse(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: "SELECT c.id, c.code, c.name, c.description FROM Courses c WHERE c.id = {{course_id}}", Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result[0])

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

	sql := database.Statement{ Sql: `SELECT cs.id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.id = {{course_id}} ORDER BY cs.session_number`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetCourseSession(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cs.id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Courses c 
									INNER JOIN Course_Sessions cs ON cs.course_id = c.id
									WHERE c.id = {{course_id}} and cs.id = {{session_id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

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

	sql := database.Statement{ Sql: `SELECT c.id, c.code, c.name, c.description,
										(SELECT id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM Course_Sessions cs 
    									WHERE cs.course_id = c.id ORDER BY cs.session_number FOR JSON PATH) as sessions
									FROM Courses c` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetCourseRegistrations(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cr.student_id, u.first_name, u.last_name, u.username, si.major
    									FROM Course_Registrations cr
    									INNER JOIN Student_Info si on si.user_id = cr.student_id
										INNER JOIN Users u on u.id = cr.student_id
									WHERE cr.course_id = {{course_id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}