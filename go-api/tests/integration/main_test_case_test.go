package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nabilfikrisp/url-shortener/internal/common/response"
	"github.com/nabilfikrisp/url-shortener/internal/features/url"
	"github.com/stretchr/testify/assert"
)

func TestMainTestCase(t *testing.T) {
	t.Run("TC-1: Generate Short URL", func(t *testing.T) {
		// You submit a long URL, expecting a short token in response.
		t.Run("Positive", func(t *testing.T) {
			app := setupTestApp(t)
			longUrlInput := "https://tanstack.com/query/latest/docs/framework/react/guides/queries"
			var urlModel url.URLModel

			t.Run("Response: 201 Created", func(t *testing.T) {
				body := `{"url":"` + longUrlInput + `"}`

				req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				// Verify response status is 201 Created
				assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var successResp response.Success
				err = json.Unmarshal(respBody, &successResp)
				if err != nil {
					t.Fatal(err)
				}

				dataBytes, err := json.Marshal(successResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(dataBytes, &urlModel)
				if err != nil {
					t.Fatal(err)
				}
			})

			t.Run("Body contains a short token", func(t *testing.T) {
				assert.NotEmpty(t, urlModel.ShortToken, "Response should contain a short token")
				assert.NotEmpty(t, urlModel.Original, "Response should contain the original URL")
				assert.Equal(t, longUrlInput, urlModel.Original)
			})

			t.Run("Database stores mapping of token to original URL", func(t *testing.T) {
				statsReq := httptest.NewRequest("GET", "/stats/"+urlModel.ShortToken, nil)
				statsResp, err := app.Test(statsReq, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer statsResp.Body.Close()

				assert.Equal(t, fiber.StatusOK, statsResp.StatusCode)

				statsBody, err := io.ReadAll(statsResp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var statsSuccessResp response.Success
				err = json.Unmarshal(statsBody, &statsSuccessResp)
				if err != nil {
					t.Fatal(err)
				}

				statsDataBytes, err := json.Marshal(statsSuccessResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				var storedUrlModel url.URLModel
				err = json.Unmarshal(statsDataBytes, &storedUrlModel)
				if err != nil {
					t.Fatal(err)
				}

				// Verify database stores correct mapping
				assert.Equal(t, urlModel.ShortToken, storedUrlModel.ShortToken)
				assert.Equal(t, longUrlInput, storedUrlModel.Original)
			})
		})

		t.Run("Negative", func(t *testing.T) {
			t.Run("Missing URL - 400 Bad Request", func(t *testing.T) {
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

			t.Run("Invalid URL format - 422 Unprocessable Entity", func(t *testing.T) {
				app := setupTestApp(t)
				body := `{"url":"invalid-url-format"}`

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
	})

	t.Run("TC-2: Redirect To Long URL", func(t *testing.T) {
		// Visiting a short URL should redirect to the long one.
		t.Run("Positive", func(t *testing.T) {
			app := setupTestApp(t)
			longUrlInput := "https://docs.gofiber.io/guide/routing/"
			var urlModel url.URLModel
			var redirectResp *http.Response

			// First create a short URL
			t.Run("Setup - Create short URL", func(t *testing.T) {
				body := `{"url":"` + longUrlInput + `"}`

				req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var successResp response.Success
				err = json.Unmarshal(respBody, &successResp)
				if err != nil {
					t.Fatal(err)
				}

				dataBytes, err := json.Marshal(successResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(dataBytes, &urlModel)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, urlModel.ShortToken)
			})

			t.Run("HTTP status 302 Found", func(t *testing.T) {
				req := httptest.NewRequest("GET", "/"+urlModel.ShortToken, nil)
				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				redirectResp = resp

				assert.True(t,
					resp.StatusCode == fiber.StatusFound,
					"Expected status 302, got %d", resp.StatusCode)
			})

			t.Run("Location header points to the original URL", func(t *testing.T) {
				location := redirectResp.Header.Get("Location")
				assert.Equal(t, longUrlInput, location, "Location header should point to original URL")
			})
		})

		t.Run("Negative", func(t *testing.T) {
			t.Run("Unknown token - 404 Not Found", func(t *testing.T) {
				app := setupTestApp(t)

				req := httptest.NewRequest("GET", "/nonexistent-token", nil)
				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			})
		})
	})

	t.Run("TC-3: Idempotent Token Generation", func(t *testing.T) {
		//  The same long URL generates the same token consistently.
		t.Run("Positive", func(t *testing.T) {
			app := setupTestApp(t)
			longUrlInput := "https://docs.gofiber.io/guide/routing/"
			var firstUrlModel url.URLModel
			var secondUrlModel url.URLModel

			t.Run("First request - Submit URL", func(t *testing.T) {
				body := `{"url":"` + longUrlInput + `"}`

				req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var successResp response.Success
				err = json.Unmarshal(respBody, &successResp)
				if err != nil {
					t.Fatal(err)
				}

				dataBytes, err := json.Marshal(successResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(dataBytes, &firstUrlModel)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, firstUrlModel.ShortToken)
				assert.Equal(t, longUrlInput, firstUrlModel.Original)
			})

			t.Run("Second request - Submit same URL", func(t *testing.T) {
				body := `{"url":"` + longUrlInput + `"}`

				req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var successResp response.Success
				err = json.Unmarshal(respBody, &successResp)
				if err != nil {
					t.Fatal(err)
				}

				dataBytes, err := json.Marshal(successResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(dataBytes, &secondUrlModel)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, secondUrlModel.ShortToken)
				assert.Equal(t, longUrlInput, secondUrlModel.Original)
			})

			t.Run("Returns identical short tokens for both requests", func(t *testing.T) {
				assert.Equal(t, firstUrlModel.ShortToken, secondUrlModel.ShortToken, "Same URL should return identical tokens")
				t.Logf("Idempotent behavior confirmed: same token '%s' returned for identical URL", firstUrlModel.ShortToken)

				req := httptest.NewRequest("GET", "/"+firstUrlModel.ShortToken, nil)
				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()
				assert.Equal(t, fiber.StatusFound, resp.StatusCode)
				assert.Equal(t, longUrlInput, resp.Header.Get("Location"))
			})
		})
	})

	t.Run("TC-4: Click Analytics Tracking", func(t *testing.T) {
		// Track the number of times a short URL has been accessed.
		t.Run("Positive", func(t *testing.T) {
			app := setupTestApp(t)
			longUrlInput := "https://docs.gofiber.io/guide/routing/"
			var urlModel url.URLModel
			numberOfClicks := 3

			t.Run("Setup - Create short URL", func(t *testing.T) {
				body := `{"url":"` + longUrlInput + `"}`

				req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var successResp response.Success
				err = json.Unmarshal(respBody, &successResp)
				if err != nil {
					t.Fatal(err)
				}

				dataBytes, err := json.Marshal(successResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(dataBytes, &urlModel)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, urlModel.ShortToken)
			})

			t.Run("Click counter increments per call", func(t *testing.T) {

				for i := 0; i < numberOfClicks; i++ {
					req := httptest.NewRequest("GET", "/"+urlModel.ShortToken, nil)
					resp, err := app.Test(req, -1)
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					assert.Equal(t, fiber.StatusFound, resp.StatusCode)
					assert.Equal(t, longUrlInput, resp.Header.Get("Location"))
				}
			})

			t.Run("200 OK response with accurate hit count and metadata", func(t *testing.T) {
				statsReq := httptest.NewRequest("GET", "/stats/"+urlModel.ShortToken, nil)
				statsResp, err := app.Test(statsReq, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer statsResp.Body.Close()

				assert.Equal(t, fiber.StatusOK, statsResp.StatusCode)

				statsBody, err := io.ReadAll(statsResp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var statsSuccessResp response.Success
				err = json.Unmarshal(statsBody, &statsSuccessResp)
				if err != nil {
					t.Fatal(err)
				}

				statsDataBytes, err := json.Marshal(statsSuccessResp.Data)
				if err != nil {
					t.Fatal(err)
				}

				var statsUrlModel url.URLModel
				err = json.Unmarshal(statsDataBytes, &statsUrlModel)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, urlModel.ShortToken, statsUrlModel.ShortToken)
				assert.Equal(t, longUrlInput, statsUrlModel.Original)
				assert.Equal(t, numberOfClicks, statsUrlModel.ClickCount, "Click count should match number of accesses")
				assert.NotEmpty(t, statsUrlModel.CreatedAt, "Should have creation timestamp")

				t.Logf("Analytics tracking confirmed: %d clicks recorded for token '%s'",
					statsUrlModel.ClickCount, statsUrlModel.ShortToken)
			})
		})

		t.Run("Negative", func(t *testing.T) {
			t.Run("Stats request for invalid token - 404 Not Found", func(t *testing.T) {
				app := setupTestApp(t)

				statsReq := httptest.NewRequest("GET", "/stats/invalid-nonexistent-token", nil)
				statsResp, err := app.Test(statsReq, -1)
				if err != nil {
					t.Fatal(err)
				}
				defer statsResp.Body.Close()

				assert.Equal(t, fiber.StatusNotFound, statsResp.StatusCode)
			})
		})
	})
}
