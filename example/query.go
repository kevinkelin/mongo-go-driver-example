package example

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}

func QueryOneByM(ctx context.Context, db *mongo.Database) error {
	/*
		FindOne 方法返回的是第一个匹配到的文档,如果没有匹配到则返回一个error
		类型为  mongo.ErrNoDocuments
	*/
	c := db.Collection("yyxtest")
	filter := bson.M{
		"name": "yyx",
	}
	// SingleResult
	result := c.FindOne(ctx, filter)
	// 将result 解码成对应的数据
	if result.Err() != nil {
		return result.Err()
	}
	var data map[string]interface{}
	result.Decode(&data)
	fmt.Println(data)
	return nil
}

func QueryOneByProjection(ctx context.Context, db *mongo.Database) error {
	/*
		返回指定的字段, Projection 中设置的字段
	*/
	c := db.Collection("yyxtest")
	filter := bson.M{
		"Bname": "yyy",
	}
	projectopt := options.FindOneOptions{
		Projection: bson.M{
			"age":      1,
			"livecity": 1,
		},
	}
	result := c.FindOne(ctx, filter, &projectopt)
	if result.Err() != nil {
		return result.Err()
	}
	var data map[string]interface{}
	result.Decode(&data)
	fmt.Println(data)

	return nil

}

func QueryOneToM(ctx context.Context, db *mongo.Database) error {
	/*
		将查询结果Decode到bson.M
	*/
	type data struct {
		HasTime bool      `bson:"hasTime"`
		Name    string    `bson:"name"`
		Time    time.Time `bson:"time"`
	}
	c := db.Collection("yyxtest")
	filter := bson.M{
		"hasTime": true,
	}
	result := c.FindOne(ctx, filter)
	var dat bson.M
	err := result.Decode(&dat)
	fmt.Println(dat)
	return err

}

func QueryManyByM(ctx context.Context, db *mongo.Database) error {
	/*
		查询多条数据
	*/
	c := db.Collection("yyxtest")
	filter := bson.M{
		"name": "yyx",
	}
	result, err := c.Find(ctx, filter)
	if err != nil {
		return err
	}
	// 通过遍历result
	var datas []bson.M
	for result.Next(ctx) {
		var data bson.M
		err = result.Decode(&data)
		if err != nil {
			return err
		}
		datas = append(datas, data)
	}
	fmt.Println(datas)
	return nil

}

func QueryManyToAll(ctx context.Context, db *mongo.Database) error {
	/*
		查询多条数据
	*/
	c := db.Collection("yyxtest")
	filter := bson.M{
		"name": "yyx",
	}
	result, err := c.Find(ctx, filter)
	if err != nil {
		return err
	}
	// 通过使用All 方法来一次性将数据转换到go
	var datas []bson.M
	result.All(ctx, &datas)
	fmt.Println(datas)
	return nil

}

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

func QueryByIn(ctx context.Context, db *mongo.Database) error {
	/*
		in 查询, name 为 ["yyx", "zhangsan"]
	*/
	c := db.Collection("yyxtest")
	result, err := c.Find(ctx, bson.M{
		"name": bson.M{"$in": []string{"yyx", "zhangsan"}},
	})
	if err != nil {
		return err
	}
	var data []bson.M
	err = result.All(ctx, &data)
	fmt.Println(data)
	return err
}

func QueryByOr(ctx context.Context, db *mongo.Database) error {
	/*
		or 查询, or 查询的条件为一个array, 所以需要使用bson.A
	*/
	c := db.Collection("yyxtest")
	result, err := c.Find(ctx, bson.M{
		"$or": bson.A{
			bson.M{"name": "yyx"},
			bson.M{"Bname": "yyy"},
		},
	})
	if err != nil {
		return err
	}
	var data []bson.M
	err = result.All(ctx, &data)
	fmt.Println(data)
	return err
}

func QueryByAnd(ctx context.Context, db *mongo.Database) error {
	/*
		and 查询, and 查询的条件为一个array, 所以需要使用bson.A
		使用and 意义不大
	*/
	c := db.Collection("yyxtest")
	result, err := c.Find(ctx, bson.M{
		"$and": bson.A{
			bson.M{"name": "lisi"},
			bson.M{"age": 20},
		},
	})
	if err != nil {
		return err
	}
	var data []bson.M
	err = result.All(ctx, &data)
	fmt.Println(data)
	return err
}

func QueryByRange(ctx context.Context, db *mongo.Database) error {
	/*
		范围查询，
		查询年龄在18-20的
	*/
	c := db.Collection("yyxtest")
	option := options.FindOptions{
		Projection: bson.M{
			"age": 1,
		},
	}
	filter := bson.M{"age": bson.M{"$gte": 18, "$lte": 20}}
	result, err := c.Find(ctx, filter, &option)
	if err != nil {
		return err
	}
	var data []bson.M
	err = result.All(ctx, &data)
	PrettyPrint(data)
	return err
}

func QueryByRangeOr(ctx context.Context, db *mongo.Database) error {
	/*
		范围查询，
		多个查询条件
	*/
	c := db.Collection("yyxtest")
	option := options.FindOptions{
		Projection: bson.M{
			"age": 1,
		},
	}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"age": bson.M{"$gte": 18, "$lte": 20}},
			bson.M{"time": bson.M{"$lte": time.Now()}},
		},
	}
	result, err := c.Find(ctx, filter, &option)
	if err != nil {
		return err
	}
	var data []bson.M
	err = result.All(ctx, &data)
	PrettyPrint(data)
	return err
}
