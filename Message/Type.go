package Message

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddType @Title 类型
// @Tags 类型
// @Summary	添加类型
// @Description 添加类型
// @Param type_name query string true "名称"
// @Produce	json
// @Router		/v1/type [POST]
func AddType(ctx *gin.Context) {
	// 获取类型名称
	typeName := ctx.Query("type_name")
	if typeName == "" {
		// 如果 type_name 为空，返回错误响应
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "type_name is required",
		})
		return
	}

	// 构建类型对象
	newType := Type{
		ID:   primitive.NewObjectID(),
		Name: typeName,
	}

	// 写入数据库
	if err := DataBase.WriteOneDB("type", newType); err != nil {
		// 如果写入失败，返回错误响应
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "failed to write type to database",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, NewTypeJson{
		Code: 0,
		Data: newType,
	})
}

// UpdateType @Title 类型
// @Tags 类型
// @Summary	更新类型
// @Description 更新类型
// @Param type_id query string true "类型ID"
// @Param type_name query string false "名称"
// @Produce	json
// @Router		/v1/type [PUT]
func UpdateType(ctx *gin.Context) {
	typeID := ctx.Query("type_id")
	if typeID == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "type_id is required",
		})
		return
	}

	// 构建过滤条件
	cid, _ := primitive.ObjectIDFromHex(typeID)
	filter := bson.D{{Key: "_id", Value: cid}}

	// 构建更新字段
	updateFields := bson.D{}
	if typeName := ctx.Query("type_name"); typeName != "" {
		updateFields = append(updateFields, bson.E{Key: "name", Value: typeName})
	}

	// 检查是否有更新字段
	if len(updateFields) == 0 {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "No fields to update",
		})
		return
	}

	// 构建更新操作
	update := bson.D{{Key: "$set", Value: updateFields}}

	// 更新数据库
	err := DataBase.UpdateOneDB("type", filter, update)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to update type",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, ComponentJson{
		Code: 0,
		Data: []gin.H{gin.H{"_id": typeID}},
	})
}

// GetType @Title 类型
// @Tags 类型
// @Summary	获取类型
// @Description 获取类型
// @Param type_id query string false "类型ID(留空返回所有)"
// @Produce	json
// @Router		/v1/type [GET]
func GetType(ctx *gin.Context) {
	typeID := ctx.Query("type_id") // 如果不存在，返回空字符串 ""

	var data []gin.H
	var err error

	if typeID == "" {
		// 如果没有提供 component_id，查询所有数据
		data, err = DataBase.ReadAllDB("type")
	} else {
		// 如果提供了 component_id，根据 ID 查询单条数据
		var oneData gin.H
		cid, _ := primitive.ObjectIDFromHex(typeID)
		filter := bson.D{
			{Key: "_id", Value: cid},
		}
		oneData, err = DataBase.ReadOneDB("type", filter)
		data = append(data, oneData)
	}

	// 错误处理
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	// 返回结果
	ctx.JSON(200, TypeJson{
		Code: 0,
		Data: data,
	})
}

// DeleteType @Title 类型
// @Tags 类型
// @Summary	删除类型
// @Description 删除类型
// @Param type_id query string true "类型ID"
// @Produce	json
// @Router		/v1/type [DELETE]
func DeleteType(ctx *gin.Context) {
	// 获取 component_id
	typeID := ctx.Query("type_id")
	if typeID == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "type_id is required",
		})
		return
	}

	// 转换 component_id 为 ObjectID
	cid, err := primitive.ObjectIDFromHex(typeID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "invalid type_id format",
			"error":   err.Error(),
		})
		return
	}

	// 构建过滤条件
	filter := bson.D{
		{Key: "_id", Value: cid},
	}

	// 执行删除操作
	err = DataBase.DeleteOneDB("type", filter)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "failed to delete type",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, ComponentJson{
		Code: 0,
		Data: []gin.H{gin.H{"_id": typeID}},
	})
}
