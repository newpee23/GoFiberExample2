package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"   // or the Docker service name if running in another container
	port     = 5432          // default PostgreSQL port
	user     = "postgres"    // as defined in docker-compose.yml
	password = "123456789"   // as defined in docker-compose.yml
	dbname   = "mikelopster" // as defined in docker-compose.yml
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductWithSupplier struct {
	ProductID    int    `json:"productID"`
	ProductName  string `json:"productName"`
	Price        int    `json:"price"`
	SupplierName string `json:"supplierName"`
}

type AddProductWithSupplier struct {
	ProductName  string `json:"productName"`
	Price        int    `json:"price"`
	SupplierName string `json:"supplierName"`
}

var dbConnect *sql.DB

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	var err error
	dbConnect, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	// defer จะทำงานเมื่อทุกอย่างเสร็จ
	defer dbConnect.Close()

	// Check the connection
	err = dbConnect.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")

	app := fiber.New()

	app.Get("/productsSupplier", getAllProductAndSupplierHandler)
	app.Get("/products", getAllProductHandler)
	app.Get("/product/:id", getProductHandler)
	app.Post("/product", createProductHandler)
	app.Post("/productSupplier", createProductAndSupplierHandler)
	app.Put("/product/:id", updateProductHandler)
	app.Delete("/product/:id", deleteProductHandler)

	app.Listen(":8080")
}

func getProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	product, err := getProductById(productId)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(product)
}

func deleteProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = deleteProduct(productId)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Deleted product with ID: %d", productId)})
}

func getAllProductHandler(c *fiber.Ctx) error {
	product, err := getProducts()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(product)
}

func getAllProductAndSupplierHandler(c *fiber.Ctx) error {
	products, err := getProductsAndSuppliers()
	if err != nil {
		// return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(products)
}

func createProductHandler(c *fiber.Ctx) error {
	p := new(Product)
	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := createProduct(p)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(p)
}

func createProductAndSupplierHandler(c *fiber.Ctx) error {
	p := new(ProductWithSupplier)
	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := addProductAndSupplier(p)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(p)
}

func updateProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	p := new(Product)

	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	product, err := updateProduct(productId, p)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(product)
}
