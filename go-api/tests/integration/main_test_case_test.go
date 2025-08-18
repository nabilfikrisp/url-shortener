package integration

import (
	"encoding/json"
	"io"
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

			var urlModel url.URLModel
			err = json.Unmarshal(dataBytes, &urlModel)
			if err != nil {
				t.Fatal(err)
			}

			// Verify body contains a short token
			assert.NotEmpty(t, urlModel.ShortToken, "Response should contain a short token")
			assert.NotEmpty(t, urlModel.Original, "Response should contain the original URL")
			assert.Equal(t, longUrlInput, urlModel.Original)

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

		t.Run("Negative", func(t *testing.T) {

		})
	})

}
