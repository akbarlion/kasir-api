package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	_ "kasir-api/docs" // Import generated docs

	"github.com/spf13/viper"
)

// @title Kasir API
// @version 1.0
// @description API untuk sistem kasir dengan manajemen produk dan kategori
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host kasir-api-production-8d59.up.railway.app
// @BasePath /api

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// setup database nya
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Println("Failed to initialize database:", err)
	}
	defer db.Close()

	// cek server hidup apa engga di http://localhost:8080
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "Application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "server is running well",
		})
	})

	fmt.Println("Server running di http://localhost:" + config.Port)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Register routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.ProductByID)
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.CategoryByID)
	http.HandleFunc("/api/checkout/", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleReport)
	http.HandleFunc("/api/report", reportHandler.HandleReport)

	// Swagger documentation routes
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs/"))))
	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Root route redirect to swagger
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Gagal running server:", err)
	}
}
