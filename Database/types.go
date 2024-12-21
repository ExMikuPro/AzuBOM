package Database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Begin struct { // 启动函数组
}

type Add struct { // 添加信息函数组
}

type AdminFunction struct { // 数据管理处理函数组
}

type Get struct { // 数据获取函数组
}

type DBService struct { // 数据库操作函数组
	Client *mongo.Client
}
