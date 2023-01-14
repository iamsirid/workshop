package cloud_pockets

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) GetCloudPocketById(c echo.Context) error {
	id := c.Param("id")
	cp := CloudPocket{}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	stmt, err := h.db.Prepare("SELECT id, name, category, currency, balance FROM cloud_pockets WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Prepare statement is error",
		})
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&cp.PocketID, &cp.PocketName, &cp.Category, &cp.Currency, &cp.Balance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Scan is error")
	}

	return c.JSON(http.StatusOK, cp)

}
