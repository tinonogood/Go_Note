package main

import (
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
)


type User struct {
    Id  int
    Name        string
    Age int
}

var allUsers = []User {
    {Id: 1,Name: "john", Age: 20},
    {Id: 2,Name: "mary", Age: 18},
}

func getUser(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for _, user := range allUsers {
        if user.Id == id {
            js, err :=json.Marshal(user)
            if err != nil{
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type","application/json")
            w.Write(js)
        }
    }
}

func getNewUserID() int {
    newID := 0
    for _, user := range allUsers {
        if (newID <= user.Id) {
            newID = user.Id + 1
        }
    }
    return newID
}

func createUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    var user User
    _ = json.NewDecoder(r.Body).Decode(&user)
    user.Id = getNewUserID()
    allUsers = append(allUsers, user)
    js, err := json.Marshal(user)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(js)
}

func updateUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for index, originalUser := range allUsers{
        if originalUser.Id == id{
            allUsers = append(allUsers[:index], allUsers[index+1:]...)
            var updateUser User
            _ = json.NewDecoder(r.Body).Decode(&updateUser)
            updateUser.Id = id
            allUsers = append(allUsers, updateUser)
            js, err := json.Marshal(updateUser)
            if err != nil{
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Write(js)
          }
    }
}

func deleteUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for index, originalUser := range allUsers{
        if originalUser.Id == id{
            allUsers = append(allUsers[:index], allUsers[index+1:]...)
            break
        }
    }
}


func main() {
    r := mux.NewRouter()

    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    r.HandleFunc("/users/create", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

    http.ListenAndServe(":8080", r)
}
