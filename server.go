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

type Record struct {
	cname          string
	disease_code   string
	total_deaths   int
	total_patients int
}

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

	app.Get("/", func(c *fiber.Ctx) error {
		// check cookie
		if email := authorize(c, db); email == "" {
			log.Printf("No cookie found")
			return c.Redirect("/login")
		}
		return c.Redirect("/record")
	})

	// Record Methods
	app.Get("/record", func(c *fiber.Ctx) error {
		return getRecordHandler(c, db)
	})

	app.Post("/record", func(c *fiber.Ctx) error {
		return postRecordHandler(c, db)
	})

	app.Put("/record", func(c *fiber.Ctx) error {
		return putRecordHandler(c, db)
	})

	app.Delete("/record", func(c *fiber.Ctx) error {
		return deleteRecordHandler(c, db)
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

func getRecordHandler(c *fiber.Ctx, db *sql.DB) error {
	// check cookie
	email := authorize(c, db)
	if email == "" {
		return c.Redirect("/login")
	}

	var recordList []Record
	var r Record
	rows, err := db.Query("select cname, disease_code, total_deaths, total_patients from record where email=$1", email)
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

		recordList = append(recordList, r)
	}
	return c.Render("record/index", fiber.Map{
		"Records": recordList,
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
	var email string

	if err := c.BodyParser(&r); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}

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

func postRecordHandler(c *fiber.Ctx, db *sql.DB) error {
	email := authorize(c, db)
	if email == "" {
		return c.Redirect("/login")
	}

	newRecord := Record{}
	if err := c.BodyParser(&newRecord); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newRecord)
	if newRecord.cname != "" && newRecord.disease_code != "" {
		_, err := db.Exec("INSERT into record VALUES ($1, $2, $3, $4, $5)", email, newRecord.cname, newRecord.disease_code, newRecord.total_deaths, newRecord.total_patients)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/record")
}

func putRecordHandler(c *fiber.Ctx, db *sql.DB) error {
	email := authorize(c, db)
	if email == "" {
		return c.Redirect("/login")
	}

	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)

	return c.Redirect("/record")
}

func deleteRecordHandler(c *fiber.Ctx, db *sql.DB) error {
	email := authorize(c, db)
	if email == "" {
		return c.Redirect("/login")
	}

	todoToDelete := c.Query("item")
	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)

	return c.Redirect("/record")
}

func authorize(c *fiber.Ctx, db *sql.DB) string {
	email := c.Cookies("email")
	if email == "" {
		log.Printf("No cookie found")
		return ""
	}
	return email
}
