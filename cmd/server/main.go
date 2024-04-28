package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/derkora/master/common/genproto/taskmaster"
	"github.com/derkora/master/common/repository"
	"github.com/derkora/master/common/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func init() {

	args := os.Args[1:]
	var configname string = "default-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}
	log.Printf("loading config file %s.yml \n", configname)

	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s.yml \n", err))
	}

}

func main() {

	log.Println("Starting up GRPC server...")

	log.Println("Creating connection to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("app.mongodb.uri")))
	if err != nil {
		log.Fatalf("%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db := client.Database(viper.GetString("app.mongodb.database"))
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Connected to database...")

	server := grpc.NewServer()

	urepo := repository.NewTaskRepo(db)
	usvc := service.NewTaskService(urepo)
	taskmaster.RegisterTaskApiServer(server, usvc)

	port := ":" + viper.GetString("app.grpc.port")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", port, err)
	}

	panic(server.Serve(listener))
}
