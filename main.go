package main

import (
	"AzuBOM/Api"
	"AzuBOM/Database"
	"AzuBOM/Function"
	"AzuBOM/Message"
	_ "AzuBOM/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AzuBOM
// @version 1.0.1233
// @description 一只轻量级BOM管理器
// @host 127.0.0.1:8080
// @BasePath /
func main() {
	err := Function.BeginEvn()
	if err != nil {
		return
	}

	Message.DataBase = Database.DBService{ // 初始化数据库函数
		Client: Database.InitDB(),
	}

	router := gin.Default()

	Api.Router(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(Function.GetEvn("WEBSITE_PORT"))
	if err != nil {

		return
	}

	//cs, _ := code.Create_Code(code.Capacitance_Tage_Tpye, "123", "10")
	//file, _ := os.Create("qr2.png")
	//defer file.Close()
	//
	//qrCode, _ := barcode.Scale(cs, 350, 70)
	//png.Encode(file, qrCode)

}
