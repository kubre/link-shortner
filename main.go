package main

import (
	"github.com/gofiber/fiber/v2"
)

type LinkForm struct {
	Link string `json:"link"`
}

func main() {

	app := fiber.New()

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

		return c.JSON(fiber.Map{
			"success": true,
			"link":    link,
		})

	})

	app.Listen(":80")
}
