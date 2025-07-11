package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jung-kurt/gofpdf"
)

var db *sql.DB

func init() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/upgris"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

// Halaman login
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Proses login
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user string
	err := db.QueryRow("SELECT username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Username atau password salah"})
			return
		}
		log.Println("Login error:", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Server error"})
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

// Halaman home menampilkan data elektronik
func HomePage(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama, deskripsi, harga, stok FROM elektronik")
	if err != nil {
		log.Println("Query error:", err)
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"error": "Gagal ambil data"})
		return
	}
	defer rows.Close()

	type Elektronik struct {
		ID        int
		Nama      string
		Deskripsi string
		Harga     int
		Stok      int
	}

	var produks []Elektronik
	for rows.Next() {
		var p Elektronik
		err := rows.Scan(&p.ID, &p.Nama, &p.Deskripsi, &p.Harga, &p.Stok)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		produks = append(produks, p)
	}

	c.HTML(http.StatusOK, "home.html", gin.H{"produks": produks})
}

// Tambah data elektronik
func TambahElektronik(c *gin.Context) {
	nama := c.PostForm("nama")
	deskripsi := c.PostForm("deskripsi")
	harga := c.PostForm("harga")
	stok := c.PostForm("stok")

	_, err := db.Exec("INSERT INTO elektronik (nama, deskripsi, harga, stok) VALUES (?, ?, ?, ?)", nama, deskripsi, harga, stok)
	if err != nil {
		log.Println("Insert error:", err)
	}

	c.Redirect(http.StatusFound, "/home")
}

// Edit data elektronik
func EditElektronik(c *gin.Context) {
	id := c.PostForm("id")
	nama := c.PostForm("nama")
	deskripsi := c.PostForm("deskripsi")
	harga := c.PostForm("harga")
	stok := c.PostForm("stok")

	_, err := db.Exec("UPDATE elektronik SET nama=?, deskripsi=?, harga=?, stok=? WHERE id=?", nama, deskripsi, harga, stok, id)
	if err != nil {
		log.Println("Update error:", err)
	}

	c.Redirect(http.StatusFound, "/home")
}

// Hapus data elektronik
func HapusElektronik(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM elektronik WHERE id = ?", id)
	if err != nil {
		log.Println("Delete error:", err)
	}
	c.Redirect(http.StatusFound, "/home")
}

// Export PDF
func ExportPDF(c *gin.Context) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Times", "B", 16)
	pdf.CellFormat(0, 20, "Laporan Data Produk Elektronik", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(10, 10, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, 10, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(75, 10, "Deskripsi", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Harga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "Stok", "1", 1, "C", false, 0, "")

	rows, err := db.Query("SELECT nama, deskripsi, harga, stok FROM elektronik")
	if err != nil {
		log.Println("Gagal query data elektronik:", err)
		c.String(http.StatusInternalServerError, "Gagal mengambil data")
		return
	}
	defer rows.Close()

	pdf.SetFont("Times", "", 12)
	no := 1
	for rows.Next() {
		var nama, deskripsi string
		var harga, stok int

		if err := rows.Scan(&nama, &deskripsi, &harga, &stok); err != nil {
			log.Println("Gagal scan:", err)
			continue
		}

		if len(nama) > 35 {
			nama = nama[:35] + "..."
		}
		if len(deskripsi) > 80 {
			deskripsi = deskripsi[:80] + "..."
		}

		pdf.CellFormat(10, 10, fmt.Sprintf("%d", no), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 10, nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(75, 10, deskripsi, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("Rp %d", harga), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", stok), "1", 1, "C", false, 0, "")

		no++
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		log.Println("Gagal generate PDF:", err)
		c.String(http.StatusInternalServerError, "Gagal membuat PDF")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `attachment; filename="produk_elektronik.pdf"`)
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
