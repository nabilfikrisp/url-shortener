package url

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/common/helpers"
	"github.com/nabilfikrisp/url-shortener/internal/common/response"
)

type URLHandler interface {
	Create(c *fiber.Ctx) error
}
type urlHandler struct {
	service URLService
}

func NewURLHandler(service URLService) URLHandler {
	return &urlHandler{
		service: service,
	}
}

type shortenPostRequest struct {
	Url string `json:"url" validate:"required,url"`
}

func validateShortenRequest(c *fiber.Ctx, req *shortenPostRequest) error {
	if err := c.BodyParser(req); err != nil {
		return errors.New("invalid request body: " + err.Error())
	}

	if req.Url == "" {
		return errors.New("URL is required")
	}

	return nil
}
func validateShortenRule(c *fiber.Ctx, req *shortenPostRequest) error {
	if !govalidator.IsURL(req.Url) {
		return errors.New("the provided URL is not valid")
	}

	isOurDomain, err := helpers.OurDomainValidator(c, req.Url)
	if err != nil {
		return errors.New("failed to validate URL: " + err.Error())
	}
	if isOurDomain {
		return errors.New("cannot shorten URLs from our own domain")
	}

	return nil
}

func (h *urlHandler) Create(c *fiber.Ctx) error {
	req := new(shortenPostRequest)
	if err := validateShortenRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid request",
				Err:     err.Error(),
			}),
		)
	}

	if err := validateShortenRule(c, req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid request",
				Err:     err.Error(),
			}),
		)
	}

	url, err := h.service.CreateShortToken(req.Url)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Failed to create short URL",
				Err:     err.Error(),
			}),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessPayload(response.SuccessPayloadParams{
		Message: "Short URL created successfully",
		Data:    url,
	}))
}
