package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) InsertMongoDB(json []byte, table_name string) error {
	log.Println("InsertMongoDB is called.")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	bsonMap := bson.M{}
	// Convert json to bson, which is mongodb document.
	err := bson.UnmarshalExtJSON([]byte(json), false, &bsonMap)
	if err != nil {
		log.Println("error from InsertMongo(), bson.UnmarshalExtJSON")
		log.Println(err)
		return err
	}

	collection := db.Client.Database(Dbname).Collection(table_name)
	//log.Println("Mongo DB name:", db.Client.Database(Dbname).Name())
	if table_name == Colname { //table_name==Job
		var episodesFiltered JsonJob
		filter := bson.D{{"company", bsonMap["company"]}, {"title", bsonMap["title"]}}
		err := collection.FindOne(context.Background(), filter).Decode(&episodesFiltered)
		if err != nil {
			//log.Println("Error from collection.Find.")
			return err
		}

		//log.Println("episodesFiltered:", episodesFiltered)
		if len(episodesFiltered.Company) > 0 {
			//log.Println("there already exists:", bsonMap["company"])
		} else {
			log.Println(bsonMap["company"], "will be inserted.")
		}
	}

	if table_name == ColnameUser {
		var result User
		var results []User
		readOne, _ := collection.Find(context.Background(), bson.D{{"email", bsonMap["email"]}})
		readOne.Decode(&result)
		readOne.All(context.Background(), &results)
		log.Println("results:", len(results))
		if len(results) > 0 {
			log.Println("This user is already registered.")
			return fmt.Errorf("%s", "This user is already registered.")
		}
	}

	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *DB) ReadMongo(user_iput ...string) []JsonJob {
	log.Println("ReadMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.Client.Database(Dbname).Collection(Colname)

	findOptions := options.Find()
	// Sort by `date` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})
	currentTime := time.Now()
	lastMonth := time.Date(currentTime.Year(), currentTime.Month()-1, currentTime.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")

	cur, err := collection.Find(context.Background(), bson.D{{"dateadded", bson.D{{"$gt", lastMonth}}}}, findOptions)
	if err != nil {
		log.Println("err from collection.Find()")
		return nil
	}

	if len(user_iput) > 0 {
		cur, err = collection.Find(context.Background(), bson.M{"company": user_iput[0]}, findOptions)
		if err != nil {
			log.Println("err from user input:", err)
			return nil
		}
	}

	var jobs []JsonJob
	var doc JsonJob
	for cur.Next(context.Background()) {
		//var doc JsonJob
		err := cur.Decode(&doc)
		if err != nil {
			log.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
	}
	return jobs
}
