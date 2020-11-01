package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB構造体へInsert用のメソッドを定義
// JSONファイルから読み込んだバイトスライスを渡し、MongoDBへInsert
func (db *DB) InsertMongoDB(json []byte, table_name string) {
	// 3秒でタイムアウトするコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	bsonMap := bson.M{}
	// JSONのバイトスライスをMongoDBのドキュメント型であるbsonへマップ
	err := bson.UnmarshalExtJSON([]byte(json), false, &bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Insert先のコレクション名からクライアント作成
	collection := db.Client.Database(Dbname).Collection(table_name) //colname is the parameter for table_name
	fmt.Println("bsonMap:", bsonMap)
	//fmt.Println("ctx:", ctx)

	if table_name == Colname {
		readOne, _ := collection.Find(context.Background(), bson.D{{"url", bsonMap["url"]}})
		if readOne != nil {
			fmt.Println("there already exists:", bsonMap["company"])
			return
		}
	}

	if table_name == ColnameUser {
		readOne, _ := collection.Find(context.Background(), bson.D{{"email", bsonMap["email"]}})
		if readOne != nil {
			//This user is already registered.
			return
		}
	}

	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (db *DB) ReadMongo(user_iput ...string) []JsonJob {
	log.Println("ReadMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.Client.Database(Dbname).Collection(Colname)

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})

	cur, err := collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
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
			fmt.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
		log.Println("searched company:", doc.Company)
	}
	return jobs
}
