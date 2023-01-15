package cloud_pockets

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCloudPocketsById(t *testing.T) {
	//Arrange
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets/:id", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
		AddRow(1, "test1", "testc", "THB", 200)
	db, mock, err := sqlmock.New()
	mock.ExpectPrepare("SELECT id, name, category, currency, balance FROM cloud_pockets").ExpectQuery().WillReturnRows(newsMockRows)
	// fmt.Println(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	h := New(db)
	c := echo.New().NewContext(req, rec)
	c.SetPath("/cloud-pockets/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	expected := "[{\"pocketID\":1,\"pocketName\":\"test1\",\"category\":\"testc\",\"currency\":\"THB\",\"balance\":200}]"

	//Act
	err = h.GetAllCloudPocket(c)
	//Arrange
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}
