ackage main

import (
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    "html/template"
    _ "github.com/lib/pq"
    "database/sql"
    "fmt"
)


const (
    host        = "127.0.0.1"
    port        = 5432
    user        = "user"
    password    = "password"
    dbname      = "dbname"
)

var Db *sql.DB

type User struct {
    Id  int
    Name        string
    Age int
}

func init() {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    Db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
}

func getUsersHTML(w http.ResponseWriter, r *http.Request){
    users,_ := getUsers(10)
    tpl, _ := template.ParseFiles("users.html")
    tpl.Execute(w, users)
}

func getUser(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    var user User
    var err error
    err = Db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Age)
}

func getUsers(limit int) (users []User, err error) {
    rows, err := Db.Query("SELECT id, name, age FROM users limit $1", limit)
    if err != nil {
        return
    }
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Name, &user.Age)
        if err != nil {
            return
        }
        users = append(users, user)
    }
    rows.Close()
    return
}


func createUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    var user User
    _ = json.NewDecoder(r.Body).Decode(&user)
    var err error
    _, err = Db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
    if err != nil{
        fmt.Println("create wrongly")
    }
    js, err := json.Marshal(user)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(js)
    return
}

func updateUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    var updateUser User
    _ = json.NewDecoder(r.Body).Decode(&updateUser)
    updateUser.Id = id
    var err error
    _, err = Db.Exec("UPDATE users SET name = $2, age = $3 WHERE id = $1", updateUser.Id, updateUser.Name, updateUser.Age)
    js, err :=json.Marshal(updateUser)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(js)
    return
}

func deleteUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    var err error
    _, err = Db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


func main() {
    r := mux.NewRouter()
    r.HandleFunc("/users", getUsersHTML).Methods("GET")
    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    r.HandleFunc("/users/create", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

    http.ListenAndServe(":8080", r)
}

/*
postgres=# \d users

                            Table "public.users"
 Column |  Type   | Collation | Nullable |              Default              
--------+---------+-----------+----------+-----------------------------------
 id     | integer |           | not null | nextval('users_id_seq'::regclass)
 name   | text    |           | not null | 
 age    | integer |           | not null | 
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)

*/
