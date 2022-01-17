package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "fmt"
  "io/ioutil"
 
  
)

type Users struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
  }
  

var db *sql.DB
var err error

const (
	username string = "root"
	password string = ""
	database string = "user_manag"
)

var (
	dsn = fmt.Sprintf("%v:%v@/%v", username, password, database)
)


func main() {
  
 

  db, err = sql.Open("mysql",dsn)
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()

  router := mux.NewRouter()
  router.HandleFunc("/user", getUsers).Methods("GET")
  router.HandleFunc("/user", createUser).Methods("POST")
  router.HandleFunc("/user/{limit}/{offset}", getPagination).Methods("GET")
  router.HandleFunc("/user/{id}", getUser).Methods("GET")
  router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
  router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
  http.ListenAndServe(":8000", router)

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []Users
	result, err := db.Query("SELECT id, username, password, name from users")
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
	  var user Users
	  err := result.Scan(&user.ID, &user.Username, &user.Password, &user.Name)
	  if err != nil {
		panic(err.Error())
	  }
	  users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
  }

  func createUser(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO users(username, password, name ) VALUES(?,?,?)")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	password := keyVal["password"]
	name := keyVal["name"]
	_,err = stmt.Exec(username,password,name)
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
  }

  func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, username, password, name FROM users WHERE id = ?", params["id"])
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var user Users
	for result.Next() {
	  err := result.Scan(&user.ID, &user.Username, &user.Password, &user.Name)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(user)
  }

  func getPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, username, password, name FROM users LIMIT ?", params["limit"],"OFFSET",params["offset"]  )
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var user Users
	for result.Next() {
	  err := result.Scan(&user.ID, &user.Username, &user.Password, &user.Name)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(user)
  }


  func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE users SET username = ?,password = ?,name = ? WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	newusername := keyVal["username"]
	newpassword := keyVal["password"]
	newname := keyVal["name"]

	_, err = stmt.Exec(newusername,newpassword,newname, params["id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
  }


  func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
   if err != nil {
	  panic(err.Error())
	}
  fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
  }