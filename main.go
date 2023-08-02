package main

import (
	"context"
	"fmt"

	"github.com/Asad2730/DynoDb_Crud/controller"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		fmt.Println("Error loading AWS config", err)
	}

	db := dynamodb.NewFromConfig(cfg)

	controller.Create(db)
}
