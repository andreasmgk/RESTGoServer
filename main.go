package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Article struct {
    gorm.Model
    Key string `json:"Key"`
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}


// our initial migration function
func initialMigration() {
    db, err = gorm.Open("mysql", "andreasmg:HsnTO57$@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
    if err != nil {
        fmt.Println(err.Error())
        panic("failed to connect database")
    }
    //defer db.Close()

    // Migrate the schema
    db.AutoMigrate(&Article{})
}

// handles all requests to the root URL
func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
    // creates a new instance of a mux router
    myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET") // add the /articles route and map it to the returnAllArticles function
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST") // add the /articles route and map it to the createNewArticle function
    myRouter.HandleFunc("/article/{key}", updateArticle).Methods("PATCH") // add the /articles/{id} route and map it to the updateArticle function
    myRouter.HandleFunc("/article/{key}", deleteArticle).Methods("DELETE") // add the /articles/{id} route and map it to the deleteArticle function
    myRouter.HandleFunc("/article/{key}", returnSingleArticle).Methods("GET") // add the /articles/{id} route and map it to the returnSingleArticle function
    log.Fatal(http.ListenAndServe(":10000", myRouter)) // the API starts up on port 10000 if it’s not already been taken by another process
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")

    articles := []Article{}
    db.Find(&articles)
    fmt.Println("{}", articles)

    // return all the articles encoded as JSON
    json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnSingleArticle")

    // parse the path parameters
    vars := mux.Vars(r)

    // extract the `id` of the article 
    // to be returned
    k := vars["key"]

    articles := []Article{}
    db.Find(&articles)
    
    for _, article := range articles {
        // string to int
        if article.Key == k {
            fmt.Println(article)
            fmt.Println("Endpoint Hit: Article No:",k)
            json.NewEncoder(w).Encode(article)
        }
     }
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: createNewArticle")

    // get the body of the POST request
    // unmarshal this into a new Article struct
    // append this to the Articles array
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article 
    json.Unmarshal(reqBody, &article)

    //id := article.Id
    //title := article.Title
    //desc := article.Desc
    //content := article.Content

    //Article{Id: id, Title: title, Desc: desc, Content: content}
    db.Create(&article)

    fmt.Fprintf(w, "New Article Successfully Created")
    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: deleteArticle")

    // parse the path parameters
    vars := mux.Vars(r)

    // extract the `id` of the article 
    // to be deleted
    k := vars["key"]
    /*
    var article Article
    db.Where("Key = ?", k).Find(&article)

    db.Delete(&article)
    */
    articles := []Article{}
    db.Find(&articles)

    for _, a := range articles {
        if a.Key == k {
            db.Delete(&a)
            fmt.Fprintf(w, "Successfully Deleted Article")
        }
    }
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: updateArticle")

    // get the body of the POST request
    // unmarshal this into a new Article struct
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article 
    json.Unmarshal(reqBody, &article)

    //db.Where("Key = ?", k).Find(&article)

    articles := []Article{}
    db.Find(&articles)

    for _, a := range articles {
        // string to int
        if a.Key == article.Key {
            a.Title = article.Title
            a.Desc = article.Desc
            a.Content = article.Content
            db.Save(&a)
            json.NewEncoder(w).Encode(a)
            fmt.Fprintf(w, "Successfully Updated Article")
        }
     }
}

func main() {
    initialMigration()

    handleRequests()
}