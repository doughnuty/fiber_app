package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Get("/diseasetype", func(c *fiber.Ctx) error {
		return getHandler(c, db)
	})

	app.Static("/diseasetype", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func getHandler(c *fiber.Ctx, db *sql.DB) error {
	var description string
	var diseaseTypeList []string
	rows, err := db.Query("select description from diseasetype")
	defer rows.Close()
	if err != nil {
		log.Println(err)
		c.JSON("Internal error")
		return err
	}
	for rows.Next() {
		rows.Scan(&description)
		diseaseTypeList = append(diseaseTypeList, description)
	}
	return c.Render("index", fiber.Map{
		"DiseaseTypes": diseaseTypeList,
	})
}
