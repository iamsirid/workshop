//go:build integration

package cloud_pockets

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocketIT(t *testing.T) {
	body := `{"pocketName":"test-pocket", "currency":"THB", "balance":100.00}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	hCloudPocket := New(sql)

	if assert.NoError(t, hCloudPocket.CreateCloudPocket(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "test-pocket", response["pocketName"])
		assert.Equal(t, "THB", response["currency"])
		assert.Equal(t, 100.0, response["balance"])
		assert.NotEmpty(t, response["pocketID"])
	}
}
