package url

import (
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

func (h *URLHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Original string `json:"original" validate:"required,url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid request body",
				Err:     err.Error(),
			}),
		)
	}

	if !govalidator.IsURL(req.Original) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid URL",
				Err:     "The provided URL is not valid",
			}),
		)
	}

	if helpers.IsOurDomain(c, req.Original) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid URL",
				Err:     "The provided URL is from our own domain",
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
