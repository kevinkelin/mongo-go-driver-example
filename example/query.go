package example

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func QueryOneByNothing(ctx context.Context, db *mongo.Database) error {
	/*
		没有任何搜索条件， 但是也要写一个, 默认会查询数据库中的第一条记录
	*/
	c := db.Collection("yyxtest")
	result := c.FindOne(ctx, bson.M{})
	var data map[string]interface{}
	err := result.Decode(&data)
	fmt.Println(data)
	return err
}

func QueryOneByNothingDesc(ctx context.Context, db *mongo.Database) error {
	/*
		倒序返回最后一个, 使用_id 进行排序
	*/
	c := db.Collection("yyxtest")
	findOpt := options.FindOneOptions{Sort: bson.M{"_id": -1}}
	result := c.FindOne(ctx, bson.M{}, &findOpt)
	var data map[string]interface{}
	err := result.Decode(&data)
	fmt.Println(data)
	return err
}

func QueryOneByNothingMulSort(ctx context.Context, db *mongo.Database) error {
	/*
		倒序返回最后一个, 这里使用多个排序因子，这里就要使用bson.D了，先按name升序，再按age做降序
	*/
	c := db.Collection("yyxtest")
	findOpt := options.FindOneOptions{Sort: bson.D{
		{Key: "name", Value: 1},
		{Key: "age", Value: -1},
	}}
	result := c.FindOne(ctx, bson.M{}, &findOpt)
	var data map[string]interface{}
	err := result.Decode(&data)
	fmt.Println(data)
	return err
}
