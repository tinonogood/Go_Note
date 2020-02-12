package main

import (
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    "html/template"
    "fmt"
    "encoding/csv"
    "os"
)

type User struct {
    Id  int64 // csv int64
    Name        string
    Age int64
}

var allUsers = []User {
}


func getUsersHTML(w http.ResponseWriter, r *http.Request){
    tpl, _ := template.ParseFiles("users.html")
    tpl.Execute(w, allUsers)
}


func getUser(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    idURL := vars["id"]
    //id, _ := strconv.Atoi(idURL)
    id, _ := strconv.ParseInt(idURL,0,0)
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

func getNewUserID() int64 {
    newID := int64(0)
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
    writeCSV()
}

func updateUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    //id, _ := strconv.Atoi(idURL)
    id, _ := strconv.ParseInt(idURL,0,0)
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
    writeCSV()
}


func deleteUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    //id, _ := strconv.Atoi(idURL)
    id, _ := strconv.ParseInt(idURL,0,0)
    for index, originalUser := range allUsers{
        if originalUser.Id == id{
            allUsers = append(allUsers[:index], allUsers[index+1:]...)
            break
        }
    }
    writeCSV()
}

func writeCSV() error {
    // write csv
    file, err := os.Create("users.csv")
    if err != nil{
        return err
    }

    write := csv.NewWriter(file)
    for _,user := range allUsers {
        line := []string{strconv.FormatInt(user.Id,10), user.Name, strconv.FormatInt(user.Age,10)}
        err := write.Write(line)
        if err != nil {
            return err
        }
    }
    write.Flush()
    return file.Close()
}

func main() {

    // open file
    file, err := os.Open("users.csv")
    if err != nil {
        panic(err)
    }

    // read csv
    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1
    record, err := reader.ReadAll()
    if err != nil {
        panic(err)
    }

    for _, item := range record {
        var user User
        user.Id, _ = strconv.ParseInt(item[0],0,0)
        user.Name = item[1]
        user.Age, _ = strconv.ParseInt(item[2],0,0)
        allUsers = append(allUsers, user)
        fmt.Printf("id: %d, name: %s, age: %d\n", user.Id, user.Name, user.Age)
    }

    file.Close()

    r := mux.NewRouter()
    r.HandleFunc("/users", getUsersHTML).Methods("GET")
    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    r.HandleFunc("/users/create", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

    http.ListenAndServe(":8080", r)

}
      	
