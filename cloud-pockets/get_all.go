package cloud_pockets

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CloudPocket struct {
	PocketID   int     `json:"pocketID"`
	PocketName string  `json:"pocketName"`
	Category   string  `json:"category"`
	Currency   string  `json:"currency"`
	Balance    float64 `json:"balance"`
	// AccountID  int     `json:"accountID"`
}

func (h handler) GetAllCloudPocket(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	stmt, err := h.db.Prepare("SELECT id, name, category, currency, balance FROM cloud_pockets")
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
		err = rows.Scan(&cp.PocketID, &cp.PocketName, &cp.Category, &cp.Currency, &cp.Balance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Scan is error")
		}

		cloudPockets = append(cloudPockets, cp)
	}

	return c.JSON(http.StatusOK, cloudPockets)

}
