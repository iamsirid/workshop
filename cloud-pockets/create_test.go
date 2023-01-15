//go:build unit

package cloud_pockets

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocket(t *testing.T) {
	body := `{"pocketName":"test-pocket", "currency":"THB", "balance":100.00}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectPrepare("INSERT INTO").ExpectQuery().WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	hCloudPocket := New(db)

	if assert.NoError(t, hCloudPocket.CreateCloudPocket(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "test-pocket", response["pocketName"])
		assert.Equal(t, "THB", response["currency"])
		assert.Equal(t, 100.0, response["balance"])
	}
}
