package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db *sql.DB

func getEnv(key, defaultValue string) string {

	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
func initDB() {
	var username string
	var password string
	var host string
	var port string
	var database string
	var driver string
	var err error

	//get env with default value
	username = getEnv("DB_USERNAME", "root")
	password = getEnv("DB_PASSWORD", "")
	host = getEnv("DB_HOST", "localhost")
	port = getEnv("DB_PORT", "3306")
	database = getEnv("DB_DATABASE", "database")
	driver = getEnv("DB_DRIVER", "mysql")

	if driver == "mysql" {
		db, err = sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database)
		if err != nil {
			log.Fatal(err)
		}
	} else if driver == "postgres" {
		db, err = sql.Open("postgres", "postgresql://"+username+":"+password+"@"+host+"/"+database)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Database driver not supported")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

// cek tabel exist
func checkTableExist() {

	_, err := db.Query("SELECT 1 FROM provinsis LIMIT 1")
	if err != nil {
		_, err = db.Query("SELECT 1 FROM kab_kotas LIMIT 1")
		if err != nil {
			_, err = db.Query("SELECT 1 FROM kecamatans LIMIT 1")
			if err != nil {
				_, err = db.Query("SELECT 1 FROM kelurahan_desas LIMIT 1")
				if err != nil {
					log.Fatal("Table not exist")
				}
			}
		}
	}
}

func info(c *gin.Context) {
	// jumlah provinsi, kabupaten, kecamatan, kelurahan
	var jumlahProvinsi int
	var jumlahKabupaten int
	var jumlahKecamatan int
	var jumlahKelurahan int

	rowCountProvinsi := db.QueryRow("SELECT COUNT(*) FROM provinsis").Scan(&jumlahProvinsi)
	rowCountKabupaten := db.QueryRow("SELECT COUNT(*) FROM kab_kotas").Scan(&jumlahKabupaten)
	rowCountKecamatan := db.QueryRow("SELECT COUNT(*) FROM kecamatans").Scan(&jumlahKecamatan)
	rowCountKelurahan := db.QueryRow("SELECT COUNT(*) FROM kelurahan_desas").Scan(&jumlahKelurahan)

	if rowCountProvinsi != nil {
		log.Fatal(rowCountProvinsi)
	}
	if rowCountKabupaten != nil {
		log.Fatal(rowCountKabupaten)
	}
	if rowCountKecamatan != nil {
		log.Fatal(rowCountKecamatan)
	}
	if rowCountKelurahan != nil {
		log.Fatal(rowCountKelurahan)
	}

	c.JSON(200, gin.H{
		"jumlah_provinsi":  jumlahProvinsi,
		"jumlah_kabupaten": jumlahKabupaten,
		"jumlah_kecamatan": jumlahKecamatan,
		"jumlah_kelurahan": jumlahKelurahan,
	})
}

func getProvinsi(c *gin.Context) {
	searchQuery := c.Query("search")

	var rows *sql.Rows
	var err error

	if searchQuery != "" {
		// Jika ada parameter search, lakukan pencarian berdasarkan nama
		rows, err = db.Query("SELECT id, nama_provinsi FROM provinsis WHERE nama_provinsi LIKE ? ORDER BY id ASC", "%"+searchQuery+"%")
	} else {
		// Jika tidak ada parameter search, ambil semua data provinsi
		rows, err = db.Query("SELECT id, nama_provinsi FROM provinsis ORDER BY id ASC")
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var provinsi []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_provinsi string
		if err := rows.Scan(&id, &nama_provinsi); err != nil {
			log.Fatal(err)
		}
		provinsi = append(provinsi, map[string]interface{}{"id": id, "nama": nama_provinsi})
	}

	c.JSON(200, provinsi)
}

// get detail provinsi
// hitung jumlah kabupaten kode, kecamatan dan kelurahan
func getDetailProvinsi(c *gin.Context) {
	idProvinsi := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_provinsi FROM provinsis WHERE id = ?", idProvinsi)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var provinsi map[string]interface{}
	for rows.Next() {
		var id int
		var nama_provinsi string
		if err := rows.Scan(&id, &nama_provinsi); err != nil {
			log.Fatal(err)
		}
		var jumlahKabupaten int
		var jumlahKecamatan int
		var jumlahKelurahan int
		rowCountKabupaten := db.QueryRow("SELECT COUNT(*) FROM kab_kotas WHERE provinsi_id = ?", id).Scan(&jumlahKabupaten)
		rowCountKecamatan := db.QueryRow("SELECT COUNT(*) FROM kecamatans LEFT JOIN kab_kotas ON kecamatans.kab_kota_id = kab_kotas.id WHERE kab_kotas.provinsi_id = ?", id).Scan(&jumlahKecamatan)
		rowCountKelurahan := db.QueryRow("SELECT COUNT(*) FROM kelurahan_desas LEFT JOIN kecamatans ON kelurahan_desas.kecamatan_id = kecamatans.id LEFT JOIN kab_kotas ON kecamatans.kab_kota_id = kab_kotas.id WHERE kab_kotas.provinsi_id = ?", id).Scan(&jumlahKelurahan)
		if rowCountKabupaten != nil {
			log.Fatal(rowCountKabupaten)
		}
		if rowCountKecamatan != nil {
			log.Fatal(rowCountKecamatan)
		}
		if rowCountKelurahan != nil {
			log.Fatal(rowCountKelurahan)
		}

		provinsi = map[string]interface{}{
			"id":               id,
			"nama_provinsi":    nama_provinsi,
			"jumlah_kabupaten": jumlahKabupaten,
			"jumlah_kecamatan": jumlahKecamatan,
			"jumlah_kelurahan": jumlahKelurahan,
		}
	}

	c.JSON(200, provinsi)
}

// get all kabupaten with pagination and limit, search by name
func getKabupatenAll(c *gin.Context) {
	searchQuery := c.Query("search")
	// page and limit default value is 1 and 10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	var rows *sql.Rows
	var err error

	if searchQuery != "" {
		// Jika ada parameter search, lakukan pencarian berdasarkan nama
		rows, err = db.Query("SELECT id, nama_kab_kota FROM kab_kotas WHERE nama_kab_kota LIKE ?  ORDER BY id ASC LIMIT ? OFFSET ?", "%"+searchQuery+"%", limit, offset)
	} else {
		// Jika tidak ada parameter search, ambil semua data provinsi
		rows, err = db.Query("SELECT id, nama_kab_kota FROM kab_kotas ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kabupaten []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kab_kota string
		if err := rows.Scan(&id, &nama_kab_kota); err != nil {
			log.Fatal(err)
		}
		kabupaten = append(kabupaten, map[string]interface{}{"id": id, "nama": nama_kab_kota})
	}

	c.JSON(200, kabupaten)
}

// get kabupaten by provinsi
func getKabupaten(c *gin.Context) {
	idProvinsi := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kab_kota FROM kab_kotas WHERE provinsi_id = ? ORDER BY id ASC", idProvinsi)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kabupaten []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kab_kota string
		if err := rows.Scan(&id, &nama_kab_kota); err != nil {
			log.Fatal(err)
		}
		kabupaten = append(kabupaten, map[string]interface{}{"id": id, "nama": nama_kab_kota})
	}

	c.JSON(200, kabupaten)
}

// get detail kabupaten
// hitung jumlah kecamatan dan kelurahan
func getDetailKabupaten(c *gin.Context) {
	idKabupaten := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kab_kota FROM kab_kotas WHERE id = ?", idKabupaten)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kabupaten map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kab_kota string
		if err := rows.Scan(&id, &nama_kab_kota); err != nil {
			log.Fatal(err)
		}
		var jumlahKecamatan int
		var jumlahKelurahan int
		rowCountKecamatan := db.QueryRow("SELECT COUNT(*) FROM kecamatans WHERE kab_kota_id = ?", id).Scan(&jumlahKecamatan)
		rowCountKelurahan := db.QueryRow("SELECT COUNT(*) FROM kelurahan_desas LEFT JOIN kecamatans ON kelurahan_desas.kecamatan_id = kecamatans.id WHERE kecamatans.kab_kota_id = ?", id).Scan(&jumlahKelurahan)
		if rowCountKecamatan != nil {
			log.Fatal(rowCountKecamatan)
		}
		if rowCountKelurahan != nil {
			log.Fatal(rowCountKelurahan)
		}

		kabupaten = map[string]interface{}{
			"id":               id,
			"nama_kab_kota":    nama_kab_kota,
			"jumlah_kecamatan": jumlahKecamatan,
			"jumlah_kelurahan": jumlahKelurahan,
		}
	}

	c.JSON(200, kabupaten)
}

// get all kecamatan with pagination and limit, search by name
func getKecamatanAll(c *gin.Context) {
	searchQuery := c.Query("search")
	// page and limit default value is 1 and 10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	var rows *sql.Rows
	var err error

	if searchQuery != "" {
		// Jika ada parameter search, lakukan pencarian berdasarkan nama
		rows, err = db.Query("SELECT id, nama_kecamatan FROM kecamatans WHERE nama_kecamatan LIKE ?  ORDER BY id ASC LIMIT ? OFFSET ?", "%"+searchQuery+"%", limit, offset)
	} else {
		// Jika tidak ada parameter search, ambil semua data provinsi
		rows, err = db.Query("SELECT id, nama_kecamatan FROM kecamatans ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kecamatan []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kecamatan string
		if err := rows.Scan(&id, &nama_kecamatan); err != nil {
			log.Fatal(err)
		}
		kecamatan = append(kecamatan, map[string]interface{}{"id": id, "nama": nama_kecamatan})
	}

	c.JSON(200, kecamatan)
}

// get kecamatan by kabupaten
func getKecamatan(c *gin.Context) {
	idKabupaten := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kecamatan FROM kecamatans WHERE kab_kota_id = ? ORDER BY id ASC", idKabupaten)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kecamatan []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kecamatan string
		if err := rows.Scan(&id, &nama_kecamatan); err != nil {
			log.Fatal(err)
		}
		kecamatan = append(kecamatan, map[string]interface{}{"id": id, "nama": nama_kecamatan})
	}

	c.JSON(200, kecamatan)
}

// get detail kecamatan
// hitung jumlah kelurahan
func getDetailKecamatan(c *gin.Context) {
	idKecamatan := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kecamatan FROM kecamatans WHERE id = ?", idKecamatan)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kecamatan map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kecamatan string
		if err := rows.Scan(&id, &nama_kecamatan); err != nil {
			log.Fatal(err)
		}
		var jumlahKelurahan int
		rowCountKelurahan := db.QueryRow("SELECT COUNT(*) FROM kelurahan_desas WHERE kecamatan_id = ?", id).Scan(&jumlahKelurahan)
		if rowCountKelurahan != nil {
			log.Fatal(rowCountKelurahan)
		}

		kecamatan = map[string]interface{}{
			"id":               id,
			"nama_kecamatan":   nama_kecamatan,
			"jumlah_kelurahan": jumlahKelurahan,
		}
	}

	c.JSON(200, kecamatan)
}

// get all kelurahan with pagination and limit, search by name
func getKelurahanAll(c *gin.Context) {
	searchQuery := c.Query("search")
	// page and limit default value is 1 and 10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	var rows *sql.Rows
	var err error

	if searchQuery != "" {
		// Jika ada parameter search, lakukan pencarian berdasarkan nama
		rows, err = db.Query("SELECT id, nama_kelurahan_desa FROM kelurahan_desas WHERE nama_kelurahan_desa LIKE ?  ORDER BY id ASC LIMIT ? OFFSET ?", "%"+searchQuery+"%", limit, offset)
	} else {
		// Jika tidak ada parameter search, ambil semua data provinsi
		rows, err = db.Query("SELECT id, nama_kelurahan_desa FROM kelurahan_desas ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kelurahan []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kelurahan_desa string
		if err := rows.Scan(&id, &nama_kelurahan_desa); err != nil {
			log.Fatal(err)
		}
		kelurahan = append(kelurahan, map[string]interface{}{"id": id, "nama": nama_kelurahan_desa})
	}

	c.JSON(200, kelurahan)
}

// get kelurahan by kecamatan
func getKelurahan(c *gin.Context) {
	idKecamatan := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kelurahan_desa FROM kelurahan_desas WHERE kecamatan_id = ? ORDER BY id ASC", idKecamatan)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kelurahan []map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kelurahan_desa string
		if err := rows.Scan(&id, &nama_kelurahan_desa); err != nil {
			log.Fatal(err)
		}
		kelurahan = append(kelurahan, map[string]interface{}{"id": id, "nama": nama_kelurahan_desa})
	}

	c.JSON(200, kelurahan)
}

// get detail kelurahan
func getDetailKelurahan(c *gin.Context) {
	idKelurahan := c.Param("id")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, nama_kelurahan_desa FROM kelurahan_desas WHERE id = ?", idKelurahan)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var kelurahan map[string]interface{}
	for rows.Next() {
		var id int
		var nama_kelurahan_desa string
		if err := rows.Scan(&id, &nama_kelurahan_desa); err != nil {
			log.Fatal(err)
		}

		kelurahan = map[string]interface{}{
			"id":                  id,
			"nama_kelurahan_desa": nama_kelurahan_desa,
		}
	}

	c.JSON(200, kelurahan)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	initDB()
	checkTableExist()
	defer db.Close()

	// Mulai membuat API menggunakan Gin
	router := gin.New()

	router.GET("/", info)

	// Tambahkan endpoint-endpoint API di sini
	router.GET("/provinsi", getProvinsi)
	router.GET("/provinsi/:id", getDetailProvinsi)

	router.GET("/kota", getKabupatenAll)
	router.GET("/provinsi/:id/kota", getKabupaten)
	router.GET("/kota/:id", getDetailKabupaten)

	router.GET("/kecamatan", getKecamatanAll)
	router.GET("/kota/:id/kecamatan", getKecamatan)
	router.GET("/kecamatan/:id", getDetailKecamatan)

	router.GET("/kelurahan", getKelurahanAll)
	router.GET("/kecamatan/:id/kelurahan", getKelurahan)
	router.GET("/kelurahan/:id", getDetailKelurahan)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
