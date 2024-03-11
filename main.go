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
	ProductID    int
	ProductName  string
	Price        int
	SupplierName string
}

type AddProductWithSupplier struct {
	ProductName  string
	Price        int
	SupplierName string
}

type Supplier struct {
	// define your supplier fields
	Name string
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

	app.Get("/product/:id", getProductHandler)

	app.Listen(":8080")

	// productAdd := Product{
	// 	Name:  "Example Product",
	// 	Price: 100,
	// }

	// supplier := Supplier{
	// 	Name: "Example Supplier",
	// }

	// err = addProductAndSupplier(productAdd, supplier)
	// if err != nil {
	// 	log.Fatalf("Error adding product and supplier: %v", err)
	// } else {
	// 	fmt.Println("Product and supplier added successfully")
	// }

	// err = createProduct(&Product{Name: "Go Product3", Price: 333})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Retrieve all products
	// productAll, err := getProducts()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Products:")
	// for _, p := range productAll {
	// 	fmt.Printf("ID: %d, Name: %s, Price: %d\n", p.ID, p.Name, p.Price)
	// }

	// getProductById
	// product, err := getProductById(3)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Product with ID %d: %+v\n", product.ID, product)

	// update ProductById
	// productUpdate, err := updateProduct(3, &Product{Name: "Go update1", Price: 123})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Update product SuccessFul !", productUpdate)

	// delete productById
	// err = deleteProduct(3)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Delete product SuccessFul !")

	// Join Data getProductsAndSuppliers
	// ProductsAndSuppliers, err := getProductsAndSuppliers()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("ProductsAndSuppliers:")
	// for _, p := range ProductsAndSuppliers {
	// 	fmt.Printf("ID: %d, Name: %s, Price: %d, SupplierName: %s\n", p.ProductID, p.ProductName, p.Price, p.SupplierName)
	// }
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
