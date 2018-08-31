package routes

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"../database"
	"github.com/auth0/go-jwt-middleware"
)

func SetCourseRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	// Course Routes
	router.HandleFunc("/courses/all", GetAllCourseData).Methods("GET")
	//router.Handle("/courses", middleware.Handler(GetCourses)).Methods("GET")

	router.HandleFunc("/courses", GetCourses).Methods("GET")

	router.HandleFunc("/courses", NewCourse).Methods("PUT")
	router.HandleFunc("/courses/{course_id}", GetCourse).Methods("GET")
	router.HandleFunc("/courses/{course_id}/grades", GetCourseGrades).Methods("GET")

	router.HandleFunc("/courses/{course_id}/assignments", GetCourseAssignments).Methods("GET")


	//router.HandleFunc("/courses/{code}", UpdateCourse).Methods("POST")
	//router.HandleFunc("/courses/{code}", DeleteCourse).Methods("DELETE")
	router.HandleFunc("/courses/{course_id}/registrants", GetCourseRegistrants).Methods("GET")

	// Course Session Routes
	router.HandleFunc("/courses/{course_id}/sessions", GetCourseSessions).Methods("GET")
	router.HandleFunc("/courses/{course_id}/sessions", NewCourseSession).Methods("PUT")
	router.HandleFunc("/courses/{course_id}/sessions/{session_id}", GetCourseSession).Methods("GET")
	router.HandleFunc("/courses/{course_id}/sessions/{session_id}/assignments", GetSessionAssignments).Methods("GET")
	//router.HandleFunc("/courses/{code}/sessions/{session_number}", UpdateCourseSession).Methods("POST")
	//router.HandleFunc("/courses/{code}/sessions/{session_number}", DeleteCourseSession).Methods("DELETE")
}


func GetCourses(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT c.id, c.code, c.name, c.description 
									FROM School.Courses c` }

	database.SelectAndReturnJson(sql, w)

}

func GetCourse(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: "SELECT c.id, c.code, c.name, c.description FROM School.Courses c WHERE c.id = {{course_id}}", Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func NewCourse(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: `INSERT INTO School.Courses (code, name, description) VALUES ('{{code}}', '{{name}}', '{{description}}')`, Params: courseData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))
	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))	}

}

func GetCourseSessions(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cs.id, cs.session_number, cs.title, cs.description, cs.start_datetime 
										FROM School.Courses c 
										INNER JOIN School.Course_Sessions cs ON cs.course_id = c.id
									WHERE c.id = {{course_id}} ORDER BY cs.session_number`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetCourseSession(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cs.id, cs.session_number, cs.title, cs.description, cs.start_datetime 
										FROM School.Courses c 
										INNER JOIN School.Course_Sessions cs ON cs.course_id = c.id
									WHERE c.id = {{course_id}} and cs.id = {{session_id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)
}

func GetSessionAssignments(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT ca.id, cs.id as session_id, cs.session_number, ca.type, ca.name, ca.description, ca.weight
										FROM School.Courses c 
										INNER JOIN School.Course_Sessions cs ON cs.course_id = c.id
										INNER JOIN School.Course_Assignments ca ON ca.course_session_id = cs.id 
									WHERE c.id = {{course_id}} and cs.id = {{session_id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)
}


func NewCourseSession(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseSessionData := map[string]string{ "course_id": r.Form.Get("course_id"),  "session_number": r.Form.Get("session_number"), "title": r.Form.Get("title"), "description": r.Form.Get("description"), "start_datetime": r.Form.Get("start_datetime") }

	sql := database.Statement{ Sql: `INSERT INTO School.Course_Sessions (course_id, session_number, title, description, start_datetime) VALUES ('{{course_id}}', '{{session_number}}', '{{title}}', '{{description}}', '{{start_datetime}}')`, Params: courseSessionData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}
	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))
	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))	}
}


func GetCourseAssignments(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cs.id, cs.id as session_id, cs.session_number, cs.title, cs.description, cs.start_datetime 
										FROM School.Courses c 
										INNER JOIN School.Course_Sessions cs ON cs.course_id = c.id
									WHERE c.id = {{course_id}} ORDER BY cs.session_number`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}


func GetAllCourseData(w http.ResponseWriter, r *http.Request) {

	sql := database.Statement{Sql: `SELECT c.id, c.code, c.name, c.description,
										(SELECT id, cs.session_number, cs.title, cs.description, cs.start_datetime FROM School.Course_Sessions cs 
    									WHERE cs.course_id = c.id ORDER BY cs.session_number FOR JSON PATH, WITHOUT_ARRAY_WRAPPER) as sessions
									FROM School.Courses c`}

	database.SelectAndReturnJson(sql, w)

}

func GetCourseRegistrants(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT cr.student_id, u.first_name, u.last_name, u.username, si.major
    									FROM School.Course_Registrations cr
    									INNER JOIN School.Student_Info si on si.user_id = cr.student_id
										INNER JOIN People.Users u on u.id = cr.student_id
									WHERE cr.course_id = {{course_id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetCourseGrades(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT
    									g.course_id, g.name, g.student_id, g.first_name, g.last_name, g.username, 
										cast(g.final_percentage as float) as final_percentage, g.final_grade 
									FROM School.f_GetCourseGrades({{course_id}}) g`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}


