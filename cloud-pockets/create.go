package cloud_pockets

import "github.com/labstack/echo/v4"

type CreatePocket struct {
	ID       int     `json:"id"`
	name     string  `json:"name"`
	category string  `json:"category"`
	currency string  `json:"currency"`
	balance  float64 `json:"balance"`
}

func CreateCloudPocket(c echo.Context) error {

	//cp := CreatePocket{}

}
