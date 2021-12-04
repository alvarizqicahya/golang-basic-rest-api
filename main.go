package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mahasiswa struct {
	// gorm.Model
	Id        uint64         `gorm:"primary_key:auto_increment" json:"id"`
	Nama      string         `gorm:"type:varchar(255)" json:"nama"`
	NoTelp    string         `gorm:"type:varchar(255)" json:"no_telp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type MahasiswaHandler interface{}

var db *gorm.DB

func main() {
	DatabaseConnect()
	defer DatabaseClose()
	Route()
}

func DatabaseConnect() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/test?parseTime=True&loc=Asia%2FJakarta&charset=utf8"), &gorm.Config{})
	if err != nil {
		log.Panic("Koneksi database error, " + err.Error())
	}

	db.AutoMigrate(&Mahasiswa{})

	log.Println("Koneksi database berhasil")
}

func DatabaseClose() {
	dbConn, _ := db.DB()

	dbConn.Close()

	log.Println("Koneksi database tertutup")
}

func Route() {
	r := gin.Default()
	defer r.Run(":8080")

	r.GET("/", func(c *gin.Context) {
		res := gin.H{
			"status":  http.StatusOK,
			"message": "REST-Api",
		}
		c.JSON(http.StatusOK, res)
	})

	person := r.Group("/mahasiswa")
	{
		person.GET("/", ShowAllMahasiswa)
		person.GET("/:id", ShowMahasiswa)
		person.POST("/create", CreateMahasiswa)
		person.PUT("/update/:id", UpdateMahasiswa)
		person.DELETE("/delete/:id", DeleteMahasiswa)
	}
}

func ShowAllMahasiswa(c *gin.Context) {
	var mhs []Mahasiswa
	var message string

	if err := db.Find(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menampilkan semua data mahasiswa",
			"error":   err.Error(),
		})
		return
	}

	count := len(mhs)

	if count > 0 {
		message = "suskes menampilkan semua data mahasiswa"
	} else {
		message = "data mahasiswa kosong"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    mhs,
		"count":   count,
	})

}

func ShowMahasiswa(c *gin.Context) {
	id := c.Params.ByName("id")
	var mhs Mahasiswa

	if err := db.Where("id = ?", id).First(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menampilkan data mahasiswa",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "suskes menampilkan data mahasiswa",
		"data":    mhs,
	})
}

func CreateMahasiswa(c *gin.Context) {
	var mhs Mahasiswa
	if err := c.BindJSON(&mhs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menambah data",
			"error":   err.Error(),
		})
		return
	}

	db.Create(&mhs)

	c.JSON(http.StatusOK, gin.H{
		"message": "sukses menambah data",
		"data":    mhs,
	})
}

func UpdateMahasiswa(c *gin.Context) {
	id := c.Params.ByName("id")
	var mhs Mahasiswa

	if err := db.Where("id = ?", id).First(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal memperbarui data",
			"error":   err.Error(),
		})
		return
	}

	c.BindJSON(&mhs)
	db.Save(&mhs)

	c.JSON(http.StatusOK, gin.H{
		"message": "sukses memperbarui data",
		"data":    mhs,
	})
}

func DeleteMahasiswa(c *gin.Context) {
	id := c.Params.ByName("id")
	var mhs Mahasiswa

	if err := db.Where("id = ?", id).Delete(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menghapus data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sukses menghapus data",
		"id":      id,
	})
}
