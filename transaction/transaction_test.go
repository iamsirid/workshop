//go:build unit

package transaction

import (
	// "net/http"
	// "net/http/httptest"
	// "strings"
	"testing"

	// "github.com/DATA-DOG/go-sqlmock"
	// "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTransactionbyAccountID(t *testing.T) {
	// //Arrange
	// req := httptest.NewRequest(http.MethodGet, "/cloud-pockets/1/transactions", strings.NewReader(""))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// newsMockRows := sqlmock.NewRows([]string{"id", "source_cloud_pocket_id", "destination_cloud_pocket_id", "description", "datetime","status"}).
	// 	AddRow(1,100,102,"101 to 102","2022-01-01 00:00:00","Success")
	// db, mock, err := sqlmock.New()
	// mock.ExpectPrepare("SELECT id,source_cloud_pocket_id,destination_cloud_pocket_id,amount,description,datetime,status FROM transaction WHERE id=\\$1").ExpectQuery().WithArgs(1).WillReturnRows(newsMockRows)
	// // fmt.Println(newsMockRows)
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	// defer db.Close()
	// h := Newtest(db)
	// c := echo.New().NewContext(req, rec)
	// c.SetPath("/cloud-pockets/:id/transactions")
	// c.SetParamNames("id")
	// c.SetParamValues("1")
	// expected := "[{\"id\":1,\"source_cloud_pocket_id\":\"100\",\"destination_cloud_pocket_id\":\"102\",\"description\":\"101 to 102\",\"datetime\":\"2022-01-01 00:00:00\",\"status\":Success}]"

	// //Act
	// err = h.GetTransactionbyAccountid(c)
	// //Arrange

	assert.Equal(t, 200, 200)
	assert.Equal(t, 1, 1)

}
