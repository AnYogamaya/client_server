package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type employee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	BAL  string `json:"bal"`
}

type JsonResponse struct {
	Type    string     `json:"type"`
	Data    []employee `json:"data"`
	Message string     `json:"message"`
}
type JsonResponse2 struct {
	Idd  int64  `json:"id2"`
	Type string `json:"status"`
}
type JsonResponse3 struct {
	Type string   `json:"status"`
	Data employee `json:"data2"`
}

func setupDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Ayogamaya@28@tcp(127.0.0.1:3306)/students")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	rout := mux.NewRouter()
	rout.HandleFunc("/", show).Methods("GET")
	rout.HandleFunc("/show_selc", show_selective).Methods("GET")
	rout.HandleFunc("/insert", insert).Methods("POST")
	rout.HandleFunc("/delete", delete).Methods("POST")
	rout.HandleFunc("/update", update).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", rout))
}

//........................... Show all ...................................

func show(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	rows, err := db.Query("SELECT * FROM transaction")

	if err != nil {
		panic(err)
	}

	var emp []employee

	for rows.Next() {
		var Id string
		var Nam string
		var bl string

		err = rows.Scan(&Id, &Nam, &bl)

		if err != nil {
			panic(err)
		}

		emp = append(emp, employee{ID: Id, Name: Nam, BAL: bl})
	}

	var response = JsonResponse{Type: "success", Data: emp}

	json.NewEncoder(w).Encode(response)
}

//....................................... Show Selective ..........................

func show_selective(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	id1 := keyVal["id"]

	result, err := db.Query("SELECT * FROM transaction WHERE id = ?", id1)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	var emp employee
	for result.Next() {
		err := result.Scan(&emp.ID, &emp.Name, &emp.BAL)
		if err != nil {
			panic(err.Error())
		}
	}
	if emp.ID == "" {
		num, _ := strconv.Atoi(id1)
		var response = JsonResponse2{Type: "Fail!!.", Idd: int64(num)}
		json.NewEncoder(w).Encode(response)
	} else {
		var response = JsonResponse3{Type: "Sucess", Data: emp}
		json.NewEncoder(w).Encode(response)
	}
}

//....................................... Insert ...........................

func insert(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	//var response = JsonResponse{}

	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO transaction(id,name,balance) VALUES(?,?,?)")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	json.Unmarshal(body, &keyVal)

	if keyVal["id"] == "" && keyVal["name"] == "" && keyVal["bal"] == "" {
		var response = JsonResponse{Type: "error", Message: "Missing Data"}
		json.NewEncoder(w).Encode(response)
	} else {
		_, err = stmt.Exec(keyVal["id"], keyVal["name"], keyVal["bal"])
	}

	if err != nil {
		panic(err.Error())
	}
	num, _ := strconv.Atoi(keyVal["id"])
	var response2 = JsonResponse2{Type: "success", Idd: int64(num)}

	json.NewEncoder(w).Encode(response2)
}

//................................. Delete .........................

func delete(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	var response = JsonResponse{}
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("DELETE FROM transaction WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	rows, err := db.Query("SELECT id FROM transaction")
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &keyVal)
	var Id string
	for rows.Next() {
		err = rows.Scan(&Id)
		if err != nil {
			panic(err)
		}
	}

	num, _ := strconv.Atoi(Id)
	num2, _ := strconv.Atoi(keyVal["id"])

	if num < num2 {
		var response2 = JsonResponse2{Type: "Error, id not found", Idd: int64(num2)}
		json.NewEncoder(w).Encode(response2)
	} else {
		json.Unmarshal(body, &keyVal)
		del_id := keyVal["id"]
		_, err = stmt.Exec(del_id)
		if err != nil {
			panic(err.Error())
		}
		response = JsonResponse{Type: "success", Message: "Record deleted successfully!"}
		json.NewEncoder(w).Encode(response)
	}
}

//........... Update ...........................

func update(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	var response = JsonResponse{}
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("UPDATE transaction SET balance = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)

	rows, err := db.Query("SELECT id FROM transaction")
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &keyVal)
	var Id string
	for rows.Next() {
		err = rows.Scan(&Id)
		if err != nil {
			panic(err)
		}
	}

	num, _ := strconv.Atoi(Id)
	num2, _ := strconv.Atoi(keyVal["id"])
	if num < num2 {
		var response2 = JsonResponse2{Type: "Error, id not found", Idd: int64(num2)}
		json.NewEncoder(w).Encode(response2)
	} else {

		newbal := keyVal["bal"]
		id1 := keyVal["id"]
		_, err = stmt.Exec(newbal, id1)
		if err != nil {
			panic(err.Error())
		}
		response = JsonResponse{Type: "success", Message: "Balance Updated successfully."}
		json.NewEncoder(w).Encode(response)
	}
}
