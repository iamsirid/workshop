//go:build integration

package transaction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloud_pockets "github.com/kkgo-software-engineering/workshop/cloud-pockets"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTransactionbyAccountID(t *testing.T) {
	cfg := config.New().All()
	db, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	cloudPocket1 := new(CloudPocket)
	cloudPocket2 := new(CloudPocket)
	db.Exec("DROP TABLE IF EXISTS cloud_pockets")
	db.Exec("CREATE TABLE IF NOT EXISTS cloud_pockets (id SERIAL PRIMARY KEY, name VARCHAR(255), category VARCHAR(255), currency VARCHAR(255), balance NUMERIC(10,2))")
	row1 := db.QueryRow("INSERT INTO cloud_pockets (name, category, currency, balance) VALUES ('test-pocket-1', 'test', 'THB', 100.00) RETURNING id")
	row2 := db.QueryRow("INSERT INTO cloud_pockets (name, category, currency, balance) VALUES ('test-pocket-2', 'test', 'THB', 100.00) RETURNING id")
	err = row1.Scan(&cloudPocket1.ID)
	if err != nil {
		t.Error(err)
	}
	err = row2.Scan(&cloudPocket2.ID)
	if err != nil {
		t.Error(err)
	}

	e := echo.New()
	body := fmt.Sprintf(`{"source_cloud_pocket_id":%d, "destination_cloud_pocket_id":%d, "amount": 10.25, "description": "test"}`, cloudPocket1.ID, cloudPocket2.ID)
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	hTransaction := New(db)
	if assert.NoError(t, hTransaction.CreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "test", response["description"])
		assert.Equal(t, float64(10.25), response["amount"])
	}

	e = echo.New()
	req = httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.SetBasicAuth("admin", "secret")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	hCloudPocket := cloud_pockets.New(db)
	if assert.NoError(t, hCloudPocket.GetAllCloudPocket(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		if response[0]["pocketID"] == cloudPocket1.ID {
			assert.Equal(t, float64(89.75), response[0]["balance"])
		} else {
			assert.Equal(t, float64(110.25), response[1]["balance"])
		}
	}
}
