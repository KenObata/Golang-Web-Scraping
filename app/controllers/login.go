package controllers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/login.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	//get user niput first
	email := r.FormValue("email")
	password := r.FormValue("password")

	//we ganna check if User exists in MongoDB
	//var error Error //temporarily comment out

	mongoClient, _ := ConnectMongoDB()
	collection := mongoClient.Client.Database(Dbname).Collection(ColnameUser)
	cur, err := collection.Find(context.Background(), bson.M{"email": email, "password": password})

	if cur == nil {
		//if cur.Next(context.Background()) {
		var error Error
		errorInResponse(w, http.StatusBadRequest, error)
		log.Println("Login failed.")
		return
	} else {
		//http Redirect
		target := "https://" + r.Host + "/userpost"
		http.Redirect(w, r, target, http.StatusFound)
	}

}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	//target := "https://" + r.Host + "/userpost"
	//http.Redirect(w, r, target, http.StatusFound)

	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/userpost.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	//get user niput first
	var job JsonJob
	companyName := r.FormValue("companyName")
	jobTitle := r.FormValue("jobTitle")
	jobURL := r.FormValue("jobURL")

	//append to JSON object
	if companyName != "" {
		job.Company = companyName
	}
	if jobTitle != "" {
		job.Title = jobTitle
	}
	if jobURL != "" {
		job.URL = jobURL
	}
	//get dateadded column
	currentTime := time.Now()
	job.DateAdded = currentTime.Format("2006-01-02")

	jsonJobJSON, err := json.Marshal(job)
	//insert my JSON object into mongoDB
	mongoClient, _ := ConnectMongoDB()
	mongoClient.InsertMongoDB(jsonJobJSON, Colname)

}
