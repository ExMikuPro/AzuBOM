package Message

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
import "AzuBOM/Database"

var DataBase Database.DBService

type ComponentJson struct {
	Code int     `json:"code"`
	Data []gin.H `json:"data"`
}

type NewComponentJson struct {
	Code int       `json:"code"`
	Data Component `json:"data"`
}

type SizeJson struct {
	Code int  `json:"code"`
	Data Size `json:"data"`
}

type TypeJson struct {
	Code int     `json:"code"`
	Data []gin.H `json:"data"`
}

type NewTypeJson struct {
	Code int  `json:"code"`
	Data Type `json:"data"`
}

// Component 元件数据表类型
type Component struct { // 元件表
	ID        primitive.ObjectID `bson:"_id" json:"id"` // 主键，自增
	Type      string             `bson:"type" json:"Type"`
	Values    string             `bson:"values" json:"value"`
	Deviation string             `bson:"deviation" json:"deviation"`
	Size      string             `bson:"size" json:"size"`
	PCS       string             `bson:"PCS" json:"PCS"`
	PIN       string             `bson:"PIN" json:"PIN"`
}

///

type Type struct { // 类型表
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type Size struct { // 封装表
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Values string             `bson:"values" json:"values"`
	Types  string             `bson:"types" json:"types"`
}
