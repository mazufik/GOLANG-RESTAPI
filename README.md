# MEMBUAT REST API DENGAN GOLANG, ECHO, GORM DAN MYSQL

## Menjalankan Program

- Install GO 1.21.4 or newer
- Create Database with fields id_item, nama_item, unit, stok, and harga_satuan
- Clone github repository `https://github.com/mazufik/GOLANG-RESTAPI.git`
- Copy environtment `cp config.example.json config.json`
- run `go mod tidy`
- run `go run main.go`

## Tutorial

## Tahap Pertama

Langkah pertama buat sebuah database terlebih dahulu, misalnya kita buat
database **_eshop_**. Kemudian buat satu buah tabel dengan nama **_item_**
dengan field sebagai berikut.

| PK | id_item int autoincrement |
| -- | ------------------------- |
|    | nama_item varchar(100)    |
|    | unit enum('kg','pcs')     |
|    | stok int                  |
|    | harga_satuan float        |
|    |                           |

## Tahap Kedua

Buat sebuah folder project dengan nama **_golang-restapi_**, lalu lakukan
inisialisasi project golang dengan perintah `go mod init golang-restapi`.

Setelah dilakukan inisialisasi project golang, buka project tersebut dengan text
editor atau IDE favorit teman-teman. disini saya akan menggunakan text editor
_vscode_. Lalu buka terminal dari text editor install depedensi yang diperlukan
yaitu gorm, viper dan echo dengan cara `go get gorm.io.gorm`,
`go get github.com/spf13/viper` dan `go get github.com/labstack/echo`.

Selanjutnya pertama-tama kita akan membuat config untuk databasenya, buat file
_config.json_ di root project dan buat sebuah folder di dalam project kita
dengan nama **_config_**.

```json
{
  "database": {
    "host": "hostname",
    "port": "portnya",
    "dbname": "nama database",
    "username": "nama user database",
    "password": "password database"
  }
}
```

Lalu didalam folder tersebut buat file dengan nama **_database.go_**.

- file _database.go_

```go
package config

import (
 "github.com/spf13/viper"
 "gorm.io/driver/mysql"
 "gorm.io/gorm"
)

func InitDB() *gorm.DB {
 viper.SetConfigFile("config.json")
 viper.AddConfigPath(".")

 err := viper.ReadInConfig()
 if err != nil {
  panic(err)
 }

 host := viper.GetString("database.host")
 port := viper.GetString("database.port")
 dbname := viper.GetString("database.dbname")
 username := viper.GetString("database.username")
 password := viper.GetString("database.password")

 dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname

 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

 if err != nil {
  panic("Can't connect to database")
 }

 return db
}
```

## Tahap Ketiga

Pada project kita nanti, kita akan menerapkan repository pattern tapi akan saya
sesuaikan dengan kebiasaan saya ya patternya.

Pertama buat folder dengan nama **_domain_** di dalam root project kita. Lalu di
dalam folder **_domain_** buat folder **_item_** dan di dalamnya buat folder
**_models_**, **_repositories_**, **_controllers_**, dan **_services_**.

Kita akan mulai dari models ya, buat sebuah fil baru di dalam folder
**_models_** dengan nama _item.go_.

- file _/domain/item/models/item.go_

```go
package models

type Item struct {
 IdItem      int     `json:"id_item" gorm:"column:id_item;primaryKey;autoIncrement"`
 NamaItem    string  `json:"nama_item" gorm:"column:nama_item"`
 Unit        string  `json:"unit" gorm:"column:unit"`
 Stok        int     `json:"stok" gorm:"column:stok"`
 HargaSatuan float64 `json:"harga_satuan" gorm:"column:harga_satuan"`
}

func (Item) TableName() string {
 return "item"
}
```

## Tahap Keempat

Setelah kita membuat model, kita lanjutkan membuat repositorynya. Untuk membuat
repositorynya buat sebuah file dengan nama _item_repository.go_ di dalam folder
**_/domain/item/repositories/_**.

- file _/domain/item/repositories/item_repository.go_

```go
package repositories

import (
 "github.com/mazufik/GOLANG-RESTAPI/domain/item/models"
 "gorm.io/gorm"
)

type dbItem struct {
 Conn *gorm.DB
}

// Create implements ItemRepository.
func (db *dbItem) Create(item models.Item) error {
 return db.Conn.Create(item).Error
}

// Delete implements ItemRepository.
func (db *dbItem) Delete(idItem int) error {
 return db.Conn.Delete(&models.Item{IdItem: idItem}).Error
}

// GetAll implements ItemRepository.
func (db *dbItem) GetAll() ([]models.Item, error) {
 var data []models.Item
 result := db.Conn.Find(&data)
 return data, result.Error
}

// GetById implements ItemRepository.
func (db *dbItem) GetById(idItem int) (models.Item, error) {
 var data models.Item
 result := db.Conn.Where("id_item", idItem).First(&data)
 return data, result.Error
}

// Update implements ItemRepository.
func (db *dbItem) Update(idItem int, item models.Item) error {
 return db.Conn.Where("id_item", idItem).Updates(item).Error
}

type ItemRepository interface {
 Create(item models.Item) error
 Update(idItem int, item models.Item) error
 Delete(idItem int) error
 GetById(idItem int) (models.Item, error)
 GetAll() ([]models.Item, error)
}

func NewItemRepository(Conn *gorm.DB) ItemRepository {
 return &dbItem{Conn: Conn}
}
```

## Tahap Kelima

Setelah kita selesai dengan repository, tahap selanjutnya kita akan membuat
service. di service ini kita akan membuat logic dari aplikasi yang kita bangun.
Buat sebuah file dengan nama _item_service.go_ didalam folder
**_/domain/item/services/_**.

Sebelum kita lanjut ke item_service, kita buat dulu sebuah folder di root
project dengan nama **_helpers_** folder ini berfungsi untuk menyimpan
fungsi-fungsi custom yang kita buat. didalam folder **_helpers_** buat file
dengan nama _response.go_

- file _/helpers/response.go_

```go
package helpers

type Response struct {
 Status   int         `json:"status"`
 Messages string      `json:"messages"`
 Data     interface{} `json:"data"`
}
```

Kita lanjutkan ke file _item_service.go_

- file _/domain/item/services/item_service.go_

```go
package services

import (
 "fmt"

 "github.com/mazufik/GOLANG-RESTAPI/domain/item/models"
 "github.com/mazufik/GOLANG-RESTAPI/domain/item/repositories"
 "github.com/mazufik/GOLANG-RESTAPI/helpers"
 "gorm.io/gorm"
)

type itemService struct {
 itemRepo repositories.ItemRepository
}

// Create implements ItemService.
func (service *itemService) Create(item models.Item) helpers.Response {
 var response helpers.Response
 if err := service.itemRepo.Create(item); err != nil {
  response.Status = 400
  response.Messages = "Failed to create a new item"
 } else {
  response.Status = 200
  response.Messages = "Success to create a new item"
 }
 return response
}

// Delete implements ItemService.
func (service *itemService) Delete(idItem int) helpers.Response {
 var response helpers.Response
 if err := service.itemRepo.Delete(idItem); err != nil {
  response.Status = 400
  response.Messages = fmt.Sprint("Failed to delete item: ", idItem)
 } else {
  response.Status = 200
  response.Messages = "Success to delete item"
 }
 return response
}

// GetAll implements ItemService.
func (service *itemService) GetAll() helpers.Response {
 var response helpers.Response
 data, err := service.itemRepo.GetAll()
 if err != nil {
  response.Status = 400
  response.Messages = "Failed to get all item"
 } else {
  response.Status = 200
  response.Messages = "Success to get all item"
  response.Data = data
 }
 return response
}

// GetById implements ItemService.
func (service *itemService) GetById(idItem int) helpers.Response {
 var response helpers.Response
 data, err := service.itemRepo.GetById(idItem)
 if err != nil {
  response.Status = 404
  response.Messages = fmt.Sprintf("Item %d not found", idItem)
 } else {
  response.Status = 200
  response.Messages = "Success to get all item"
  response.Data = data
 }
 return response
}

// Update implements ItemService.
func (service *itemService) Update(idItem int, item models.Item) helpers.Response {
 var response helpers.Response
 if err := service.itemRepo.Update(idItem, item); err != nil {
  response.Status = 400
  response.Messages = fmt.Sprint("Failed to update item: ", idItem)
 } else {
  response.Status = 200
  response.Messages = "Success to update item"
 }
 return response
}

type ItemService interface {
 Create(item models.Item) helpers.Response
 Update(idItem int, item models.Item) helpers.Response
 Delete(idItem int) helpers.Response
 GetById(idItem int) helpers.Response
 GetAll() helpers.Response
}

func NewItemService(db *gorm.DB) ItemService {
 return &itemService{itemRepo: repositories.NewItemRepository(db)}
}
```

## Tahap Keenam

Pada tahap ini kita akan masuk ke controllers, buat file dengan nama
_item_controller.go_ di dalam folder **_/domain/item/controllers_**.

- file _/domain/item/controllers/item_controller.go_

```go
package controllers

import (
 "net/http"
 "strconv"

 "github.com/labstack/echo/v4"
 "github.com/mazufik/GOLANG-RESTAPI/domain/item/models"
 "github.com/mazufik/GOLANG-RESTAPI/domain/item/services"
)

type ItemController struct {
 itemService services.ItemService
}

func (controller ItemController) Create(c echo.Context) error {
 type payload struct {
  NamaItem    string  `json:"nama_item" validate:"required"`
  Unit        string  `json:"unit" validate:"required"`
  Stok        int     `json:"stok" validate:"required"`
  HargaSatuan float64 `json:"harga_satuan" validate:"required"`
 }

 payloadValidator := new(payload)

 if err := c.Bind(payloadValidator); err != nil {
  return err
 }

 result := controller.itemService.Create(models.Item{NamaItem: payloadValidator.NamaItem, Unit: payloadValidator.Unit, Stok: payloadValidator.Stok, HargaSatuan: payloadValidator.HargaSatuan})

 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": result,
 })
}

func (controller ItemController) Update(c echo.Context) error {
 type payload struct {
  NamaItem    string  `json:"nama_item" validate:"required"`
  Unit        string  `json:"unit" validate:"required"`
  Stok        int     `json:"stok" validate:"required"`
  HargaSatuan float64 `json:"harga_satuan" validate:"required"`
 }

 payloadValidator := new(payload)

 if err := c.Bind(payloadValidator); err != nil {
  return err
 }

 idItem, _ := strconv.Atoi(c.Param("id_item"))
 result := controller.itemService.Update(idItem, models.Item{NamaItem: payloadValidator.NamaItem, Unit: payloadValidator.Unit, Stok: payloadValidator.Stok, HargaSatuan: payloadValidator.HargaSatuan})

 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": result,
 })
}

func (controller ItemController) Delete(c echo.Context) error {
 idItem, _ := strconv.Atoi(c.Param("id_item"))
 result := controller.itemService.Delete(idItem)

 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": result,
 })
}

func (controller ItemController) GetAll(c echo.Context) error {
 result := controller.itemService.GetAll()

 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": result,
 })
}

func (controller ItemController) GetById(c echo.Context) error {
 idItem, _ := strconv.Atoi(c.Param("id_item"))
 result := controller.itemService.GetById(idItem)

 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": result,
 })
}

func NewItemController(db *gorm.DB) ItemController {
 service := services.NewItemService(db)
 controller := ItemController{
  itemService: service,
 }

 return controller
}
```

## Tahap Ketujuh

Pada tahap ini, kita akan membuat file _main.go_ nya di dalam folder root
project kita.

- file _main.go_

```go
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
 apiV1.GET("item/detail/:id_item", itemController.GetById)

 route.Logger.Print("Starting ", viper.GetString("server.appName"))
 route.Logger.Fatal(route.Start(":" + viper.GetString("server.appPort")))
}
```
