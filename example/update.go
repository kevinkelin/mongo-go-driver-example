package example

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateDocOrInsert(ctx context.Context, db *mongo.Database) error {
	/*
		更新文档，如果没有则新添加
		ModifiedCount, 如果没有更新的，比如说要更新的内容和搜索出来的文件是一致的，则不更新
	*/
	c := db.Collection("yyxtest")
	filter := bson.M{
		"name": "lisi",
		// "age":  20,
	}
	updates := bson.M{
		"$set": bson.M{
			"age": 21,
		},
	}
	upsert := true
	ops := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := c.UpdateOne(ctx, filter, updates, &ops)
	if err != nil {
		return err
	}
	fmt.Println("MatchedCount:", result.MatchedCount)
	fmt.Println("ModifiedCount:", result.ModifiedCount)
	fmt.Println("UpsertedCount:", result.UpsertedCount)
	return nil

}
