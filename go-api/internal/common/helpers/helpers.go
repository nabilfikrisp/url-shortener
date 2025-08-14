package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CHECKPOINT BENERIN LOCALHOST MASIH LOLOS
func IsOurDomain(c *fiber.Ctx, inputURL string) bool {
	parsed, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	appHost := strings.Split(c.Hostname(), ":")[0]
	urlHost := strings.Split(parsed.Host, ":")[0]

	fmt.Print(appHost + "APP HOST")
	fmt.Print(urlHost + "URL HOST")

	return strings.EqualFold(urlHost, appHost)
}
