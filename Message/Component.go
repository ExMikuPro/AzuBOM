package Message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

// AddComponent @Title 元件
// @Tags 元件
// @Summary	添加新元件
// @Description 添加新元件
// @Produce	json
// @Param component_type query string true "类型"
// @Param component_values query string true "数值"
// @Param component_deviation query string true "误差"
// @Param component_size query string true "封装"
// @Param component_PCS query string true "数量"
// @Router		/v1/component [POST]
func AddComponent(ctx *gin.Context) {
	// 创建新组件并直接赋值
	newComponent := Component{
		ID:        primitive.NewObjectID(),
		Type:      ctx.Query("component_type"),
		Values:    ctx.Query("component_values"),
		Deviation: ctx.Query("component_deviation"),
		Size:      ctx.Query("component_size"),
		PCS:       ctx.Query("component_PCS"),
	}

	// 生成随机 PIN（时间戳 + 随机三位数）
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	componentPIN := fmt.Sprintf("%03d", r.Intn(1000)) // 生成三位随机数并补零
	newComponent.PIN = fmt.Sprintf("%d%s", time.Now().Unix(), componentPIN)

	// 写入数据库
	if err := DataBase.WriteOneDB("component", newComponent); err != nil {
		// 如果写入失败，返回错误响应
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "failed to write component to database",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, NewComponentJson{
		Code: 0,
		Data: newComponent,
	})
}

// UpdateComponent @Title 元件
// @Tags 元件
// @Summary 更新元件
// @Description 更新元件
// @Param component_id query string true "元件ID"
// @Param component_type query string false "类型"
// @Param component_values query string false "数值"
// @Param component_deviation query string false "误差"
// @Param component_size query string false "封装"
// @Param component_PCS query string false "数量"
// @Produce json
// @Router /v1/component [PUT]
func UpdateComponent(ctx *gin.Context) {
	// 获取元件 ID
	componentID := ctx.Query("component_id")
	if componentID == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "component_id is required",
		})
		return
	}

	// 构建过滤条件
	cid, _ := primitive.ObjectIDFromHex(componentID)
	filter := bson.D{{Key: "_id", Value: cid}}

	// 构建更新字段
	updateFields := bson.D{}
	if componentType := ctx.Query("component_type"); componentType != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: componentType})
	}
	if componentValues := ctx.Query("component_values"); componentValues != "" {
		updateFields = append(updateFields, bson.E{Key: "values", Value: componentValues})
	}
	if componentDeviation := ctx.Query("component_deviation"); componentDeviation != "" {
		updateFields = append(updateFields, bson.E{Key: "component_deviation", Value: componentDeviation})
	}
	if componentSize := ctx.Query("component_size"); componentSize != "" {
		updateFields = append(updateFields, bson.E{Key: "size", Value: componentSize})
	}
	if componentPCS := ctx.Query("component_PCS"); componentPCS != "" {
		updateFields = append(updateFields, bson.E{Key: "PCS", Value: componentPCS})
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

	fmt.Println(filter)

	// 更新数据库
	err := DataBase.UpdateOneDB("component", filter, update)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to update component",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, ComponentJson{
		Code: 0,
		Data: []gin.H{gin.H{"_id": componentID}},
	})
}

// GetComponent @Title 元件
// @Tags 元件
// @Summary	获取元件
// @Param component_id query string false "元件ID(留空返回所有)"
// @Description 获取元件列表
// @Produce	json
// @Router		/v1/component [GET]
func GetComponent(ctx *gin.Context) {
	// 从查询参数中获取 component_id
	componentID := ctx.Query("component_id") // 如果不存在，返回空字符串 ""

	var data []gin.H
	var err error

	if componentID == "" {
		// 如果没有提供 component_id，查询所有数据
		data, err = DataBase.ReadAllDB("component")
	} else {
		// 如果提供了 component_id，根据 ID 查询单条数据
		var oneData gin.H
		cid, _ := primitive.ObjectIDFromHex(componentID)
		filter := bson.D{
			{Key: "_id", Value: cid},
		}
		oneData, err = DataBase.ReadOneDB("component", filter)
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
	ctx.JSON(200, ComponentJson{
		Code: 0,
		Data: data,
	})
}

// DeleteComponent @Title 元件
// @Tags 元件
// @Summary 删除元件
// @Description 删除元件
// @Param component_id query string true "元件ID"
// @Produce json
// @Router /v1/component [DELETE]
func DeleteComponent(ctx *gin.Context) {
	// 获取 component_id
	componentID := ctx.Query("component_id")
	if componentID == "" {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "component_id is required",
		})
		return
	}

	// 转换 component_id 为 ObjectID
	cid, err := primitive.ObjectIDFromHex(componentID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "invalid component_id format",
			"error":   err.Error(),
		})
		return
	}

	// 构建过滤条件
	filter := bson.D{
		{Key: "_id", Value: cid},
	}

	// 执行删除操作
	err = DataBase.DeleteOneDB("component", filter)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "failed to delete component",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(200, ComponentJson{
		Code: 0,
		Data: []gin.H{gin.H{"_id": componentID}},
	})
}
