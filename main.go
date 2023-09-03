package main

import (
	"context"
	"fmt"
	"log"
	"sample-go-crud/app/config"
	"sample-go-crud/app/database"
	"sample-go-crud/app/inject"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()
	ctx := context.Background()
	config.InitConfig()
	mongoClient := database.ConnectMongoDB(ctx)
	defer ctx.Done()
	db := mongoClient.Database(viper.GetString("MongoDatabase"))
	DI := inject.InitDI(db)
	DI.Inject(ctx, router)
	port := viper.GetString("PORT")
	e := router.Run(fmt.Sprintf(":%s", port))
	if e != nil {
		log.Fatal("Error while starting server")
	}
}
