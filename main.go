package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mazufik/GOLANG-RESTAPI/config"
	"github.com/mazufik/GOLANG-RESTAPI/domain/item/controllers"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.json")
	viper.AddConfigPath(".")
	db := config.InitDB()

	route := echo.New()
	apiV1 := route.Group("api/v1/")

	itemController := controllers.NewItemController(db)
	apiV1.POST("item/create", itemController.Create)
	apiV1.PUT("item/update/:id_item", itemController.Update)
	apiV1.DELETE("item/delete/:id_item", itemController.Delete)
	apiV1.GET("item/get_all", itemController.GetAll)
	apiV1.GET("item/detail", itemController.GetById)

	route.Logger.Print("Starting ", viper.GetString("server.appName"))
	route.Logger.Fatal(route.Start(":" + viper.GetString("server.appPort")))
}
