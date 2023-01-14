package cloudpocket

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CloudPocket struct {
	PocketID   int     `json:"pocketID"`
	PocketName string  `json:"pocketName"`
	Balance    float64 `json:"balance"`
	AccountID  int     `json:"accountID"`
}

var db *sql.DB

func getAllCloudPocket(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	stmt, err := db.Prepare("SELECT pocket_id, pocket_name, balance, account_id FROM cloud_pocket")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Prepare statement is error",
		})
	}
	cloudPockets := []CloudPocket{}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Query is error")
	}

	for rows.Next() {
		cp := CloudPocket{}
		err = rows.Scan(&cp.PocketID, &cp.PocketName, &cp.Balance, &cp.AccountID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Scan is error")
		}

		cloudPockets = append(cloudPockets, cp)
	}

	return c.JSON(http.StatusOK, cloudPockets)

}
