package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	bolt "go.etcd.io/bbolt"
)

type LinkForm struct {
	Link string `json:"link"`
}

const BUCKET_NAME = "SHORTEN_LINKS"

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Please add base url as first argument to program")
	}

	baseUrl := os.Args[1]

	app := fiber.New()
	db, err := bolt.Open("store", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Send Post request with { 'url': '' } as request to receive { 'link': '/shorten' } link",
		})
	})

	app.Get("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")

		var link string
		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(BUCKET_NAME))
			link = string(bucket.Get([]byte(code)))
			return nil
		})

		if len(link) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Code invalid",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"link":    link,
		})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		linkForm := new(LinkForm)

		if err := c.BodyParser(linkForm); err != nil {
			return err
		}

		link := linkForm.Link

		if len(link) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Please add link",
			})
		}

		code := getRandCode(6)
		db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(BUCKET_NAME))
			err := bucket.Put([]byte(code), []byte(link))
			return err
		})

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"link":    fmt.Sprintf("%s/%s", baseUrl, code),
		})

	})

	app.Listen(":80")
}

func getRandCode(length int) string {
	data := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxy0123456789"
	const LEN_OF_DATA = 61
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(data[rand.Int63()%LEN_OF_DATA])
	}
	return sb.String()
}
