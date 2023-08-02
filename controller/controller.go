package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Asad2730/DynoDb_Crud/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var (
	ctx       = context.TODO()
	tableName = "Your Table name"
)

func Create(client *dynamodb.Client) {

	item := model.Item{
		Id:        uuid.New().String(),
		Name:      "Asad Sajjad",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)

	if err != nil {
		log.Fatalf("Error unmarchling: %v", err.Error())
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		log.Fatalf("Error creating item: %v", err.Error())
	}

	fmt.Println("Item created successfully!")

}

func Read(client *dynamodb.Client) {

	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	rs, err := client.Scan(ctx, input)

	if err != nil {
		fmt.Println("failed to scan items:", err)
	}

	var items []model.Item

	for _, itemMap := range rs.Items {
		var i model.Item
		if err := attributevalue.UnmarshalMap(itemMap, &i); err != nil {
			fmt.Println("failed to unmarshal item:", err.Error())
		}
		items = append(items, i)
	}

	for i := 0; i < len(items); i++ {
		fmt.Println(items[i])
	}

}

func ReadById(client *dynamodb.Client, id string) {

	key := map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberS{
			Value: id,
		},
	}

	res, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})

	if err != nil {
		log.Fatalf("Error reading item: %v", err)
	}

	if res != nil {
		var read model.Item
		err = attributevalue.UnmarshalMap(res.Item, &read)
		if err != nil {
			log.Fatalf("Error unmarshaling item: %v", err)
		}
		fmt.Printf("Item read: %+v\n", read)
	} else {
		fmt.Println("Item not found")
	}
}

func Update(client *dynamodb.Client, id, name string) {

	key := map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberS{
			Value: id,
		},
	}

	updateExp := aws.String("SET #name = :n,#updatedAt = :u")

	expNames := map[string]string{
		"#name":      "Name",
		"#updatedAt": "UpdatedAt",
	}

	expValues := map[string]types.AttributeValue{
		":n": &types.AttributeValueMemberS{Value: name},
		":u": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          updateExp,
		ExpressionAttributeNames:  expNames,
		ExpressionAttributeValues: expValues,
		ReturnValues:              types.ReturnValueAllNew,
	}

	_, err := client.UpdateItem(ctx, input)

	if err != nil {
		fmt.Println("Error", err.Error())
	}

	fmt.Println("Updated!")
}

func Delete(client *dynamodb.Client, id string) {

	key := map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberS{
			Value: id,
		},
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	_, err := client.DeleteItem(ctx, input)

	if err != nil {
		fmt.Println("Error", err.Error())
	}

	fmt.Println("Deleted!")
}
