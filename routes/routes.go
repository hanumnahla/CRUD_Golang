package routes

import (
	"E-commerce/handlers"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	// Load template HTML dari folder templates
	r.LoadHTMLGlob("templates/*")

	// Routing untuk Login
	r.GET("/", handlers.LoginPage)
	r.POST("/login", handlers.LoginHandler)

	// Routing untuk halaman utama dan laporan
	r.GET("/home", handlers.HomePage)
	r.GET("/export-pdf", handlers.ExportPDF)

	// Routing untuk CRUD data elektronik
	r.POST("/elektronik/tambah", handlers.TambahElektronik)  // Tambah data
	r.GET("/elektronik/edit/:id", handlers.EditElektronik)   // TAMPILKAN FORM EDIT
	r.POST("/elektronik/edit", handlers.EditElektronik)      // SIMPAN PERUBAHAN
	r.GET("/elektronik/hapus/:id", handlers.HapusElektronik) // Hapus data berdasarkan ID
}
