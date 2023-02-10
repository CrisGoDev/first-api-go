package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
)

func main() {
	globals := map[string]interface{}{}
	bgColorBlue := true

	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".tmpl").
			AddFunc("args", func(values ...interface{}) ([]interface{}, error) { return values, nil }).
			AddFunc("add", func(a, b int) int { return a + b }).
			AddFunc("odd", func(x int) bool { return x%2 == 1 }).
			AddFunc("neg", func(x int) int { return -x }).
			AddFunc("currency", func(x int) string {
				if x == 0 {
					return "0"
				}
				if x < 0 {
					x = -x
					return fmt.Sprintf("($%v.%02d)", x/100, x%100)
				}
				return fmt.Sprintf("$%v.%02d", x/100, x%100)
			}).
			AddFunc("bgColor", func() string {
				bgColorBlue = !bgColorBlue
				if !bgColorBlue {
					return "bg-blue"
				}
				return "bg-white"
			}).
			AddFunc("getGlobal", func(k string) (interface{}, error) {
				v, ok := globals[k]
				if !ok {
					return "", fmt.Errorf("Global \"%v\" not defined", k)
				}
				return v, nil
			}).
			AddFunc("setGlobal", func(k string, v interface{}) (interface{}, error) {
				globals[k] = v
				return "", nil
			}).
			AddFunc("incGlobal", func(k string) (interface{}, error) {
				v, ok := globals[k].(int)
				if !ok {
					return "", fmt.Errorf("global \"%v\" is not a int.", k)
				}

				globals[k] = v + 1

				return "", nil
			}),
		JSONDecoder: func(data []byte, v interface{}) error {
			jd := json.NewDecoder(bytes.NewReader(data))
			jd.DisallowUnknownFields()
			return jd.Decode(v)
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	api := app.Group("/api")

	api.Use(cors.New())

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, from cors 👋!")
	})

	app.Listen(":3000")
}
