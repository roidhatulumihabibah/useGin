package main

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MySQLData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MongoDBData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Koneksi MySQL
	mysqlDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mahasiswa")
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlDB.Close()

	// Koneksi MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Inisialisasi router Gin
	router := gin.Default()

	// Mengatur pengaturan proxy
	router.ForwardedByClientIP = true

	// Handler untuk membaca data dari MySQL
	router.GET("/mysql", func(c *gin.Context) {
		var data []MySQLData

		rows, err := mysqlDB.Query("SELECT * FROM person")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch data from MySQL"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var d MySQLData
			err := rows.Scan(&d.ID, &d.Name)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to scan data from MySQL"})
				return
			}
			data = append(data, d)
		}

		c.JSON(200, data)
	})

	// Handler untuk membaca data dari MongoDB
	router.GET("/mongodb", func(c *gin.Context) {
		var data []MongoDBData

		collection := mongoClient.Database("dosen").Collection("person")
		cur, err := collection.Find(context.Background(), bson.D{{}})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch data from MongoDB"})
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			var d MongoDBData
			err := cur.Decode(&d)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to decode data from MongoDB"})
				return
			}
			data = append(data, d)
		}

		c.JSON(200, data)
	})

	// Menjalankan server
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
