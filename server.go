package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
)

func main() {
	connStr, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Println("Couldn't get database url. Exiting")
		return
	} else {
		log.Println("Database url: ", connStr)
	}
	log.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
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

	// Disease Type Methods
	app.Get("/diseasetype", func(c *fiber.Ctx) error {
		return getDiseaseTHandler(c, db)
	})

	// Login methods
	app.Get("/login", func(c *fiber.Ctx) error {
		return getLoginHandler(c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return postLoginHandler(c, db)
	})

	app.Static("/", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func getDiseaseTHandler(c *fiber.Ctx, db *sql.DB) error {
	// check cookie
	email := c.Cookies("email")
	if email == "" {
		log.Printf("No cookie found")
		return c.Redirect("/login")
	}

	type Record struct {
		cname          string
		disease_code   string
		total_deaths   int
		total_patients int
	}
	var diseaseTypeList []Record
	var r Record
	rows, err := db.Query("select cname, disease_code, total_deaths, total_patients from record")
	defer rows.Close()
	if err != nil {
		log.Println(err)
		c.JSON("Internal error")
		return err
	}
	for rows.Next() {
		if err = rows.Scan(&r.cname, &r.disease_code, &r.total_deaths, &r.total_patients); err != nil {
			log.Println(err)
			break
		}

		diseaseTypeList = append(diseaseTypeList, r)
	}
	return c.Render("disease-types/index", fiber.Map{
		"DiseaseTypes": diseaseTypeList,
	})
}

func getLoginHandler(c *fiber.Ctx) error {
	return c.Render("login/index", fiber.Map{})
}

func postLoginHandler(c *fiber.Ctx, db *sql.DB) error {
	type response struct {
		Email string
	}

	r := response{}

	if err := c.BodyParser(&r); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	// fmt.Printf("%v", r.Email)
	var email string
	err := db.QueryRow("select email from users where email=$1", r.Email).Scan(&email)
	if err != nil {
		log.Printf("query error: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).Render("login/index", fiber.Map{
			"content": "Wrong email",
		})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "email"
	cookie.Value = r.Email
	cookie.Expires = time.Now().Add(5 * time.Minute)

	c.Cookie(cookie)

	return c.Redirect("/diseasetype")
}
