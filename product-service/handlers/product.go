package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"product-service/models"
	"product-service/pkg/database"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func paginate(r *http.Request, query *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	query := paginate(r, database.DB)
	if err := query.Preload("Images").Find(&products).Error; err != nil {
		log.Println("Error retrieving products:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, product := range products {
		var productType models.ProductType
		if err := database.DB.First(&productType, product.ProductTypeID).Error; err != nil {
			log.Println("Error retrieving product type:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products[i].ProductType = productType
	}

	var total int64
	if err := database.DB.Model(&models.Product{}).Count(&total).Error; err != nil {
		log.Println("Error counting products:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"total": total,
		"data":  products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.Preload("Images").First(&product, params["id"]).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var productType models.ProductType
	if err := database.DB.First(&productType, product.ProductTypeID).Error; err != nil {
		http.Error(w, "Product type not found", http.StatusNotFound)
		return
	}

	product.ProductType = productType

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func GetAllProductsByProductType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var productType models.ProductType
	if err := database.DB.First(&productType, params["id"]).Error; err != nil {
		http.Error(w, "Product type not found", http.StatusNotFound)
		return
	}

	var products []models.Product
	query := database.DB.Where("product_type_id = ?", productType.ID).Preload("Images")
	query = paginate(r, query)

	if err := query.Find(&products).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var total int64
	if err := database.DB.Model(&models.Product{}).Where("product_type_id = ?", productType.ID).Count(&total).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"total": total,
		"data":  products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProductsByIDs(w http.ResponseWriter, r *http.Request) {
	idsQuery := r.URL.Query().Get("product_ids")
	if idsQuery == "" {
		http.Error(w, "Product IDs are required", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsQuery, ",")

	log.Println("Product IDs:", ids)

	var products []models.Product
	if err := database.DB.Where("id IN ?", ids).Preload("Images").Find(&products).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, product := range products {
		var productType models.ProductType
		if err := database.DB.First(&productType, product.ProductTypeID).Error; err != nil {
			http.Error(w, "Product type not found", http.StatusNotFound)
			return
		}
		products[i].ProductType = productType
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// check the product type is exists
	var productType models.ProductType
	if err := database.DB.First(&productType, product.ProductTypeID).Error; err != nil {
		http.Error(w, "Product type not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	required_quantity := r.URL.Query().Get("request_quantity")
	if required_quantity == "" {
		http.Error(w, "Stock is required", http.StatusBadRequest)
		return
	}

	required_quantity_int, err := strconv.Atoi(required_quantity)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	if required_quantity_int > product.Stock {
		http.Error(w, "Insufficient stock", http.StatusBadRequest)
		return
	}

	product.Stock -= required_quantity_int
	database.DB.Save(&product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// check the product type is exists
	var productType models.ProductType
	if err := database.DB.First(&productType, product.ProductTypeID).Error; err != nil {
		http.Error(w, "Product type not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Save(&product).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete existing images
	if err := database.DB.Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add new images
	for _, image := range product.Images {
		image.ProductID = product.ID
		// Ensure the ID is not set manually
		image.ID = 0
		if err := database.DB.Create(&image).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.Preload("Images").First(&product, params["id"]).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Collect all file URLs
	var fileUrls []string
	for _, image := range product.Images {
		fileUrls = append(fileUrls, image.URL)
	}

	// Delete the images associated with the product
	database.DB.Where("product_id = ?", product.ID).Delete(&models.ProductImage{})
	database.DB.Delete(&product)

	response := map[string]interface{}{
		"message":   "Product deleted successfully",
		"file_urls": fileUrls,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	query := database.DB

	if name := r.URL.Query().Get("name"); name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}

	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")

	if minPrice != "" && maxPrice != "" {
		query = query.Where("price BETWEEN ? AND ?", minPrice, maxPrice)
	} else if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	} else if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	if inStock := r.URL.Query().Get("in_stock"); inStock != "" && inStock == "true" {
		query = query.Where("stock > 0")
	}

	if productTypeID := r.URL.Query().Get("product_type_id"); productTypeID != "" {
		query = query.Where("product_type_id = ?", productTypeID)
	}

	var total int64
	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		log.Println("Error counting products:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query = paginate(r, query)
	if err := query.Preload("Images").Find(&products).Error; err != nil {
		log.Println("Error retrieving products:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"total": total,
		"data":  products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllProductTypes(w http.ResponseWriter, r *http.Request) {
	var productTypes []models.ProductType
	if err := database.DB.Find(&productTypes).Error; err != nil {
		log.Println("Error retrieving product types:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productTypes)
}

func CreateProductType(w http.ResponseWriter, r *http.Request) {
	var productType models.ProductType
	if err := json.NewDecoder(r.Body).Decode(&productType); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	database.DB.Create(&productType)
	json.NewEncoder(w).Encode(productType)
}

func DeleteProductType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var productType models.ProductType
	if err := database.DB.First(&productType, params["id"]).Error; err != nil {
		http.Error(w, "Product type not found", http.StatusNotFound)
		return
	}

	// delete the products associated with the product type
	database.DB.Model(&productType).Association("Products").Clear()

	database.DB.Delete(&productType)
	json.NewEncoder(w).Encode(productType)
}

func GetProductStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":    product.ID,
		"stock": product.Stock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
