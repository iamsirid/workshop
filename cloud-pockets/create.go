package cloud_pockets

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) CreateCloudPocket(c echo.Context) error {
	cloudPocket := new(CreatePocket)
	tableName := "cloud_pockets"
	if err := c.Bind(cloudPocket); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	row := h.db.QueryRow(
		`INSERT INTO `+tableName+` (name, category, currency, balance) VALUES ($1, $2, $3, $4) RETURNING id`,
		cloudPocket.Name, cloudPocket.Category, cloudPocket.Currency, cloudPocket.Balance,
	)
	err := row.Scan(&cloudPocket.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, CreatePocket{
		ID:       cloudPocket.ID,
		Name:     cloudPocket.Name,
		Category: cloudPocket.Category,
		Currency: cloudPocket.Currency,
		Balance:  cloudPocket.Balance,
	})
}
