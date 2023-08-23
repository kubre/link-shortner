package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gofiber/fiber/v2"
	bolt "go.etcd.io/bbolt"
)

type LinkForm struct {
	Link string `json:"link"`
}

const BUCKET_NAME = "SHORTEN_LINKS"

func main() {

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

	app.Get("/:code?", func(c *fiber.Ctx) error {
		code := c.Params("code")

		if code == "" {
			return c.SendString("Send Post request with { 'url': '' } as request to receive { 'link': '/shorten' } link")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"link":    code,
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

		return c.JSON(fiber.Map{
			"success": true,
			"code":    code, // Fix: provide full link
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
