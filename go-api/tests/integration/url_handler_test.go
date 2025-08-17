package integration

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"

	// "strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/common/response"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestApp(t *testing.T) *fiber.App {
	db := SetupTestDB(t)

	// init handler + register routes
	handler := url.InitURLHandler(db)
	app := fiber.New()
	url.RegisterRoutes(app, handler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	return app
}

func TestURLHandler(t *testing.T) {

	t.Run("GET /", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req, 1)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("POST /shorten", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{"url":"https://www.google.com/"}`

			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		})

		t.Run("Invalid URL", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{"url":"invalid-url"}`
			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		})

		t.Run("Missing URL", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{}`
			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Invalid JSON", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{"url":invalid-url}`
			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Shorten URL for our domain", func(t *testing.T) {
			app := setupTestApp(t)

			body := `{"url":"http://127.0.0.1:3001/"}`
			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// simulate request coming from localhost:3001
			req.Host = "127.0.0.1:3001"

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		})

		t.Run("POST /shorten - Service Error", func(t *testing.T) {
			// build app with mock service
			app := fiber.New()
			mockService := new(MockURLService)
			mockService.On("CreateShortToken", mock.Anything).Return(nil, errors.New("db insert failed"))
			h := url.NewURLHandler(mockService)
			app.Post("/shorten", h.Create)

			body := `{"url":"https://www.google.com/"}`

			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		})

	})

	t.Run("GET /stats/:shortToken", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{"url":"https://www.google.com/"}`

			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			respBody, _ := io.ReadAll(resp.Body)

			// Decode into response.Success first
			var successResp response.Success
			if err := json.Unmarshal(respBody, &successResp); err != nil {
				t.Fatal(err)
			}

			// Convert Data(any) into URLModel
			dataBytes, err := json.Marshal(successResp.Data)
			if err != nil {
				t.Fatal(err)
			}

			var created url.URLModel
			if err := json.Unmarshal(dataBytes, &created); err != nil {
				t.Fatal(err)
			}

			req = httptest.NewRequest("GET", "/stats/"+created.ShortToken, nil)
			resp, err = app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("Short Token Not Found", func(t *testing.T) {
			app := setupTestApp(t)

			req := httptest.NewRequest("GET", "/stats/nonexistent", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		})

		t.Run("Service Error", func(t *testing.T) {
			// Setup app with mock service
			app := fiber.New()
			mockService := new(MockURLService)
			mockService.On("FindByShortToken", "error-token").Return(nil, errors.New("database connection failed"))
			h := url.NewURLHandler(mockService)
			app.Get("/stats/:shortToken", h.FindByShortToken)

			req := httptest.NewRequest("GET", "/stats/error-token", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		})

	})

	t.Run("GET /:shortToken - REDIRECT", func(t *testing.T) {
		t.Run("Success - Redirects to original URL", func(t *testing.T) {
			app := setupTestApp(t)
			body := `{"url":"https://www.google.com/"}`
			req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			respBody, _ := io.ReadAll(resp.Body)
			var successResp response.Success
			if err := json.Unmarshal(respBody, &successResp); err != nil {
				t.Fatal(err)
			}

			dataBytes, _ := json.Marshal(successResp.Data)
			var created url.URLModel
			if err := json.Unmarshal(dataBytes, &created); err != nil {
				t.Fatal(err)
			}

			req = httptest.NewRequest("GET", "/"+created.ShortToken, nil)
			resp, err = app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusFound, resp.StatusCode)
			assert.Equal(t, "https://www.google.com/", resp.Header.Get("Location"))
		})

		t.Run("Failure - Short token not found", func(t *testing.T) {
			app := setupTestApp(t)
			req := httptest.NewRequest("GET", "/nonexistent", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)

			var errResp response.Error
			if err := json.Unmarshal(respBody, &errResp); err != nil {
				t.Fatal(err)
			}

			assert.Contains(t, errResp.Message, "not found")
		})
	})

}
