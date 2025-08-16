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
	FindByShortToken(c *fiber.Ctx) error
	RedirectToOriginal(c *fiber.Ctx) error
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
		return errors.New("request body is not valid JSON: " + err.Error())
	}

	if req.Url == "" {
		return errors.New("URL field is required")
	}

	return nil
}
func validateShortenRule(c *fiber.Ctx, req *shortenPostRequest) error {
	if !govalidator.IsURL(req.Url) {
		return errors.New("please provide a valid URL")
	}

	isOurDomain, err := helpers.OurDomainValidator(c, req.Url)
	if err != nil {
		return errors.New("unable to validate URL domain: " + err.Error())
	}
	if isOurDomain {
		return errors.New("cannot create short URLs for this domain")
	}

	return nil
}

func (h *urlHandler) Create(c *fiber.Ctx) error {
	req := new(shortenPostRequest)
	if err := validateShortenRequest(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Invalid request format",
				Err:     err.Error(),
			}),
		)
	}

	if err := validateShortenRule(c, req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "URL validation failed",
				Err:     err.Error(),
			}),
		)
	}

	url, err := h.service.CreateShortToken(req.Url)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Unable to create short URL",
				Err:     err.Error(),
			}),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessPayload(response.SuccessPayloadParams{
		Message: "Short URL created successfully",
		Data:    url,
	}))
}

func (h *urlHandler) FindByShortToken(c *fiber.Ctx) error {
	shortToken := c.Params("shortToken")
	if shortToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Short token parameter is required",
			}),
		)
	}

	url, err := h.service.FindByShortToken(shortToken)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Short URL not found",
				Err:     err.Error(),
			}),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessPayload(response.SuccessPayloadParams{
		Message: "URL retrieved successfully",
		Data:    url,
	}))
}

func (h *urlHandler) RedirectToOriginal(c *fiber.Ctx) error {
	shortToken := c.Params("shortToken")
	if shortToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Short token parameter is required",
			}),
		)
	}

	url, err := h.service.RedirectService(shortToken)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "Short URL not found",
				Err:     err.Error(),
			}),
		)
	}
	if url == nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.ErrorPayload(response.ErrorResponseParams{
				Message: "The requested short URL does not exist",
			}),
		)
	}

	return c.Redirect(url.Original, fiber.StatusFound)
}
