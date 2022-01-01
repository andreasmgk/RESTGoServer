package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
    Id string `json:"Id"`
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}

// a global Articles array
// that will then be populated in the main function
// to simulate a database
var Articles []Article

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
    myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PATCH") // add the /articles/{id} route and map it to the updateArticle function
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE") // add the /articles/{id} route and map it to the deleteArticle function
    myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET") // add the /articles/{id} route and map it to the returnSingleArticle function
    log.Fatal(http.ListenAndServe(":10000", myRouter)) // the API start up on port 10000 if it’s not already been taken by another process
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")

    // return all the articles encoded as JSON
    json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnSingleArticle")

    // parse the path parameters
    vars := mux.Vars(r)

    // extract the `id` of the article 
    // to be returned
    key := vars["id"]

    // Loop over all of the Articles
    // if the article.Id equals the key passed in
    // return the article encoded as JSON
    for _, article := range Articles {
        if article.Id == key {
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
    
    // update the global Articles array to include
    // the new Article
    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: deleteArticle")

    // parse the path parameters
    vars := mux.Vars(r)

    // extract the `id` of the article 
    // to be delete
    id := vars["id"]

    // loop through all the articles
    for index, article := range Articles {
        // if the id path parameter matches one of the
        // articles
        if article.Id == id {
            // update the Articles array to remove the 
            // article
            Articles = append(Articles[:index], Articles[index+1:]...)
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
    
    // loop through all the articles
    for index, a := range Articles {
        // if the id path parameter matches one of the
        // articles
        if a.Id == article.Id {
            // update the specific article with the
            // new parameter values
            Articles[index].Title = article.Title
            Articles[index].Desc  = article.Desc
            Articles[index].Content = article.Content
        }
    }

    json.NewEncoder(w).Encode(article)
}

func main() {
    Articles = []Article{
        {Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        {Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }
    handleRequests()
}