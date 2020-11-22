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

func InternalHandler(w http.ResponseWriter, r *http.Request) {
	//get user niput first
	email := r.FormValue("email")
	password := r.FormValue("password")
	//debug
	log.Println("input email in InternalHandler:", email)
	log.Println("input PW in InternalHandler:", password)
	mongoClient, _ := ConnectMongoDB()
	collection := mongoClient.Client.Database(Dbname).Collection(ColnameUser)
	cur, _ := collection.Find(context.Background(), bson.M{"email": email, "password": password})

	log.Println("cur.Current:", cur.Current)
	if cur.Current == nil {
		//if cur.Next(context.Background()) {
		var error Error
		errorInResponse(w, http.StatusBadRequest, error)
		log.Println("Login failed.")
		return
	} else {
		//http Redirect
		target := "http://" + r.Host + "/userpost"
		log.Println("http redirect to ", target)
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Working Directory
	/*
		wd := os.Getenv("WORKINGDIRECTLY")
		if wd == "" {
			wd,_ = os.Getwd()
		}*/
	wd, err := os.Getwd() //need to switch to without wd when we work with EKS
	log.Println("LoginHandler called. wd := ", wd)
	t, err := template.ParseFiles(wd + "/app/view/login.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	/*
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
			//target := "https://" + r.Host + "/userpost"
			//http.Redirect(w, r, target, http.StatusFound)
		}
	*/
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

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
