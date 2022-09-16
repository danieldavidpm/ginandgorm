package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB
var err error

type Cliente struct {
	Id      int    `json:"id"`
	Nombres string `json:"nombre"`
	Email   string `json:"email"`
}

func main() {
	db, _ = gorm.Open("mssql", "sqlserver://sa:12345@localhost:1433?database=dbPba")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&Cliente{})
	r := gin.Default()
	r.GET("/clientes/", GetClientes)
	r.GET("/clientes/:id", GetCliente)
	r.POST("/clientes/", CreateCliente)
	r.PUT("/clientes/:id", UpdateCliente)
	r.DELETE("/clientes/:id", DeleteCliente)
	r.Run(":8080")
}

func DeleteCliente(c *gin.Context) {
	id := c.Params.ByName("id")
	var cliente Cliente
	d := db.Where("id = ?", id).Delete(&cliente)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateCliente(c *gin.Context) {
	var cliente Cliente
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&cliente).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&cliente)
	db.Save(&cliente)
	c.JSON(200, cliente)
}

func CreateCliente(c *gin.Context) {
	var cliente Cliente
	c.BindJSON(&cliente)
	db.Create(&cliente)
	c.JSON(200, cliente)
}

func GetCliente(c *gin.Context) {
	id := c.Params.ByName("id")
	var cliente Cliente
	if err := db.Where("id = ?", id).First(&cliente).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, cliente)
	}
}

func GetClientes(c *gin.Context) {
	var clientes []Cliente
	if err := db.Find(&clientes).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, clientes)
	}
}
