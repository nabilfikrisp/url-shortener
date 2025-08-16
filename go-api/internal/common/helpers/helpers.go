package helpers

import (
	"errors"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func OurDomainValidator(c *fiber.Ctx, inputURL string) (bool, error) {
	if inputURL == "" {
		return false, errors.New("URL cannot be empty")
	}

	parsed, err := url.Parse(inputURL)
	if err != nil {
		return false, errors.New("failed to parse URL")
	}

	if parsed.Hostname() == "" {
		return false, errors.New("URL must have a hostname")
	}

	appHost := strings.ToLower(strings.Split(c.Hostname(), ":")[0])
	urlHost := strings.ToLower(parsed.Hostname())

	if urlHost == "localhost" {
		return true, nil
	}

	return appHost == urlHost, nil
}
