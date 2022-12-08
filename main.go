package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB

func init() {

	host := os.Getenv("PGHOST")         //google host for pgadmin
	po := os.Getenv("PGPORT")           //port 5432
	user := os.Getenv("PGUSER")         // my username must be username in env
	password := os.Getenv("PGPASSWORD") // my password
	dbname := os.Getenv("PGDATABASE")   //database name

	port, _ := strconv.Atoi(po)

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error

	db, err = sql.Open("postgres", sqlInfo) //fix this
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil { // fix this
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/items", GetArmy).Methods("GET")
	r.HandleFunc("/items/{id}", GetItem).Methods("GET")
	r.HandleFunc("/items", AddArmy).Methods("POST")
	r.HandleFunc("/items/{id}", UpdateArmy).Methods("PUT")
	r.HandleFunc("/items/{id}", DeleteArmy).Methods("DELETE")
	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe("localhost:8080", r))

}

type Army struct {
	ID          int    `json:"id"`
	Name        string `json:"model_name"`
	UnitSize    int    `json:"unit_size"`
	ModelsOwned string `json:"models_owned"`
	Points      int    `json:"price"`
}

func GetArmy(w http.ResponseWriter, r *http.Request) {
	var models []Army
	stmt, err := db.Prepare("SELECT * FROM nighthaunt;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := stmt.Query()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = stmt.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for result.Next() {
		var model Army
		err := result.Scan(&model.ID, &model.Name, &model.Points, &model.ModelsOwned, &model.UnitSize)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		models = append(models, model)
	}

	err = json.NewEncoder(w).Encode(models)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	intid, err := strconv.Atoi(stringid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if intid > 100 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stmt, err := db.Prepare("SELECT * FROM nighthaunt WHERE ID =?;") // need to fix this
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := stmt.QueryRow(intid)
	err = stmt.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var army Army
	err = result.Scan(&army.ID, &army.Name, &army.Points)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(army)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func AddArmy(w http.ResponseWriter, r *http.Request) {
	request, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newModel := Army{
		ID:          0,
		Name:        "",
		Points:      0,
		ModelsOwned: "",
		UnitSize:    0,
	}

	err = json.Unmarshal(request, &newModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare("INSERT INTO army (id, name, owned, points,) VALUES (?,?,?,?,);")
	_, err = stmt.Exec(newModel.Name, newModel.ModelsOwned, newModel.ID, newModel.Points, newModel.UnitSize)
	w.WriteHeader(http.StatusCreated)

}

func DeleteArmy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	id, err := strconv.Atoi(stringid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if id > 100 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stmt, err := db.Prepare("DELETE FROM army WHERE ID =?") // need to fix this
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateArmy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	id, err := strconv.Atoi(stringid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if id > 100 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newModel := Army{
		ID:          0,
		Name:        "",
		Points:      0,
		ModelsOwned: "",
		UnitSize:    0,
	}

	err = json.Unmarshal(request, &newModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE army SET Name =?, Models_owned =?, WHERE ID =?") // need to fix this
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(newModel.Name, newModel.ModelsOwned, newModel.ID, newModel.Points, newModel.UnitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
