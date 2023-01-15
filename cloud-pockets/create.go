package cloud_pockets

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) CreateCloudPocket(c echo.Context) error {
	cloudPocket := new(CloudPocket)
	if err := c.Bind(cloudPocket); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	stmt, err := h.db.Prepare("INSERT INTO cloud_pockets (name, category, currency, balance) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(cloudPocket.PocketName, cloudPocket.Category, cloudPocket.Currency, cloudPocket.Balance)
	err = row.Scan(&cloudPocket.PocketID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, CloudPocket{
		PocketID:   cloudPocket.PocketID,
		PocketName: cloudPocket.PocketName,
		Category:   cloudPocket.Category,
		Currency:   cloudPocket.Currency,
		Balance:    cloudPocket.Balance,
	})
}
