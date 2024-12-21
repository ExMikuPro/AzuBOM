package Database

import (
	"AzuBOM/Function"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://" + Function.GetEvn("DB_USER") + ":" + Function.GetEvn("DB_PASSWD") + "@" + Function.GetEvn("DB_HOST") + ":" + Function.GetEvn("DB_PORT")).SetTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return client
}

func (s *DBService) ReadAllDB(Collection string) ([]gin.H, error) {
	collection := s.Client.Database(Function.GetEvn("DB_DATA_BASE")).Collection(Collection)

	cur, err := collection.Find(context.Background(), bson.D{}, options.Find())
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cur.Close(context.Background())

	results := []gin.H{}
	for cur.Next(context.Background()) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		results = append(results, gin.H(result))
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

func (s *DBService) ReadOneDB(Collection string, filter bson.D) (gin.H, error) {
	collection := s.Client.Database(Function.GetEvn("DB_DATA_BASE")).Collection(Collection)

	var result gin.H
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DBService) WriteOneDB(Collection string, Data any) error {
	collection := s.Client.Database(Function.GetEvn("DB_DATA_BASE")).Collection(Collection)

	_, err := collection.InsertOne(context.Background(), Data)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	return nil
}

func (s *DBService) UpdateOneDB(collectionName string, filter bson.D, update bson.D) error {
	// 创建上下文并设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取集合
	collection := s.Client.Database(Function.GetEvn("DB_DATA_BASE")).Collection(collectionName)

	// 执行更新操作
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update document in collection %s: %v", collectionName, err)
	}

	// 检查匹配和更新的文档数量
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document matched the filter in collection %s", collectionName)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("document matched but no modifications were made in collection %s", collectionName)
	}

	return nil
}

func (s *DBService) DeleteOneDB(Collection string, filter bson.D) error {
	collection := s.Client.Database(Function.GetEvn("DB_DATA_BASE")).Collection(Collection)

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	return nil
}
