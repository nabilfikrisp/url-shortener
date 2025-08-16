package url

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/common/helpers"
	"github.com/nabilfikrisp/url-shortener/internal/common/response"
)

type URLHandler struct {
	service *URLService
}

func NewURLHandler(service *URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type ShortenPostRequest struct {
	Original string `json:"original" validate:"required,url"`
}

func validateShortenRequest(c *fiber.Ctx, req *ShortenPostRequest) error {
	if err := c.BodyParser(req); err != nil {
		return errors.New("invalid request body: " + err.Error())
	}

	if req.Original == "" {
		return errors.New("original URL is required")
	}

	if !govalidator.IsURL(req.Original) {
		return errors.New("the provided URL is not valid")
	}

	isOurDomain, err := helpers.OurDomainValidator(c, req.Original)
	if err != nil {
		return errors.New("failed to validate URL: " + err.Error())
	}
	if isOurDomain {
		return errors.New("cannot shorten URLs from our own domain")
	}

	return nil
}

func (h *URLHandler) Create(c *fiber.Ctx) error {
	req := new(ShortenPostRequest)
	if err := validateShortenRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid request",
				Err:     err.Error(),
			}),
		)
	}

	url, err := h.service.CreateShortToken(req.Original)
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
