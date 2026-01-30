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

	"github.com/spf13/viper"
)

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

	// Register routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.ProductByID)
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.CategoryByID)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal running server")
	}

}

// func handleProduk(w http.ResponseWriter, r *http.Request) {
// 	// Cek apakah ada ID di path
// 	if strings.HasPrefix(r.URL.Path, "/api/produk/") {
// 		// Handle /api/produk/{id}
// 		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
// 			return
// 		}

// 		switch r.Method {
// 		case "GET":
// 			for _, p := range produk {
// 				if p.ID == id {
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(p)
// 					return
// 				}
// 			}
// 			http.Error(w, "Produk belum ada", http.StatusNotFound)
// 		case "PUT":
// 			var updateProduk Produk
// 			err = json.NewDecoder(r.Body).Decode(&updateProduk)
// 			if err != nil {
// 				http.Error(w, "Invalid Request", http.StatusBadRequest)
// 				return
// 			}
// 			for i := range produk {
// 				if produk[i].ID == id {
// 					updateProduk.ID = id
// 					produk[i] = updateProduk
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(produk[i])
// 					return
// 				}
// 			}
// 			http.Error(w, "Produk Belum Ada", http.StatusNotFound)
// 		case "DELETE":
// 			for i := range produk {
// 				if produk[i].ID == id {
// 					produk = append(produk[:i], produk[i+1:]...)
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(map[string]string{
// 						"message": "sukses delete",
// 					})
// 					return
// 				}
// 			}
// 			http.Error(w, "Produk belum ada", http.StatusNotFound)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	} else {
// 		// Handle /api/produk
// 		switch r.Method {
// 		case "GET":
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(produk)
// 		case "POST":
// 			var produkBaru Produk
// 			err := json.NewDecoder(r.Body).Decode(&produkBaru)
// 			if err != nil {
// 				http.Error(w, "Invalid request", http.StatusBadRequest)
// 				return
// 			}
// 			produkBaru.ID = len(produk) + 1
// 			produk = append(produk, produkBaru)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusCreated)
// 			json.NewEncoder(w).Encode(produkBaru)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

// func handleCategories(w http.ResponseWriter, r *http.Request) {
// 	// Cek apakah ada ID di path
// 	if strings.HasPrefix(r.URL.Path, "/api/categories/") {
// 		// Handle /api/categories/{id}
// 		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			http.Error(w, "Invalid Category ID", http.StatusBadRequest)
// 			return
// 		}

// 		switch r.Method {
// 		case "GET":
// 			for _, c := range categories {
// 				if c.ID == id {
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(c)
// 					return
// 				}
// 			}
// 			http.Error(w, "Category belum ada", http.StatusNotFound)
// 		case "PUT":
// 			var updateCategory Category
// 			err := json.NewDecoder(r.Body).Decode(&updateCategory)
// 			if err != nil {
// 				http.Error(w, "Invalid Request", http.StatusBadRequest)
// 				return
// 			}
// 			for i := range categories {
// 				if categories[i].ID == id {
// 					updateCategory.ID = id
// 					categories[i] = updateCategory
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(updateCategory)
// 					return
// 				}
// 			}
// 			http.Error(w, "Category Belum Ada", http.StatusNotFound)
// 		case "DELETE":
// 			for i, c := range categories {
// 				if c.ID == id {
// 					categories = append(categories[:i], categories[i+1:]...)
// 					w.Header().Set("Content-Type", "application/json")
// 					json.NewEncoder(w).Encode(map[string]string{
// 						"message": "sukses delete category",
// 					})
// 					return
// 				}
// 			}
// 			http.Error(w, "Category belum ada", http.StatusNotFound)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	} else {
// 		// Handle /api/categories
// 		switch r.Method {
// 		case "GET":
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(categories)
// 		case "POST":
// 			var newCategory Category
// 			err := json.NewDecoder(r.Body).Decode(&newCategory)
// 			if err != nil {
// 				http.Error(w, "Invalid Request", http.StatusBadRequest)
// 				return
// 			}
// 			newCategory.ID = len(categories) + 1
// 			categories = append(categories, newCategory)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusCreated)
// 			json.NewEncoder(w).Encode(newCategory)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }
