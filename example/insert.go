package example

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//  添加文档操作

func InsertOneByM(ctx context.Context, db *mongo.Database) (string, error) {
	// 添加单个文档
	c := db.Collection("yyxtest")
	// bson.M 是顺序无关的，在顺序不敏感的环境下使用 bson.M{"foo": "bar", "hello": "world", "pi": 3.14159}
	/*
		InsertOneResult, 有一个InsertedID, 它是一个interface{},
		使用的时候需要转换为primitive.ObjectID, 如果要转换为string, 需要使用ObjectID.Hex() 方法
	*/
	iresult, err := c.InsertOne(ctx, bson.M{
		"name": "yyx",
		"age":  18,
	})
	if err != nil {
		return "", err
	}
	if oid, ok := iresult.InsertedID.(primitive.ObjectID); !ok {
		return "", errors.New("插入文档失败,转换objectid失败")
	} else {
		return oid.Hex(), nil
	}
}

func InsertOneByD(ctx context.Context, db *mongo.Database) (string, error) {
	c := db.Collection("yyxtest")
	// bson.D 是顺序的，在顺序敏感的环境下使用，如果不敏感，还是推荐使用bson.M,
	// 在插入的场景下，其实使用M更适合
	// bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
	iresult, err := c.InsertOne(ctx, bson.D{
		{Key: "name", Value: "zhangsan"},
		{Key: "age", Value: 19},
	})
	if err != nil {
		return "", err
	}
	if oid, ok := iresult.InsertedID.(primitive.ObjectID); !ok {
		return "", errors.New("插入文档失败,转换objectid失败")
	} else {
		return oid.Hex(), nil
	}
}

func InsertOneByMap(ctx context.Context, db *mongo.Database) (string, error) {
	c := db.Collection("yyxtest")
	//使用map
	iresult, err := c.InsertOne(ctx, map[string]interface{}{
		"name":    "yyx",
		"type":    "map",
		"hasTime": true,
		"time":    time.Now(),
	})
	if err != nil {
		return "", err
	}
	if oid, ok := iresult.InsertedID.(primitive.ObjectID); !ok {
		return "", errors.New("插入文档失败,转换objectid失败")
	} else {
		return oid.Hex(), nil
	}
}

type Person struct {
	Name     string `json:"name" bson:"Bname"`
	Age      int    `json:"age"`
	LiveCity string `json:"livecity111"`
	sex      string
}

func InsertOneByStruct(ctx context.Context, db *mongo.Database) (string, error) {
	c := db.Collection("yyxtest")
	/*
		使用struct,
		1. struct 的字段中如果定义了bson tag, 则使用bson, 如果没有定义bson,则使用字段名的小写方式，
		2. 只有公开的字段可以写
	*/
	p := Person{"yyy", 21, "bj", "male"}
	iresult, err := c.InsertOne(ctx, &p)
	if err != nil {
		return "", err
	}
	if oid, ok := iresult.InsertedID.(primitive.ObjectID); !ok {
		return "", errors.New("插入文档失败,转换objectid失败")
	} else {
		return oid.Hex(), nil
	}
}

func InsertManyByM(ctx context.Context, db *mongo.Database) ([]string, error) {
	c := db.Collection("yyxtest")
	datas := []interface{}{
		bson.M{"name": "yang1", "age": 18},
		bson.M{"name": "yang2", "age": 19},
		bson.M{"name": "yang3", "age": 20},
	}
	result, err := c.InsertMany(ctx, datas)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range result.InsertedIDs {
		if oid, ok := item.(primitive.ObjectID); !ok {
			fmt.Println("解析失败")
		} else {
			ids = append(ids, oid.Hex())
		}
	}
	return ids, nil
}
