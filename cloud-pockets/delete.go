package cloud_pockets

import (
	"net/http"
	"strconv"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Resp struct {
	Message string `json:"message"`
}

func (h handler) DeleteCloudPocket(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")
	rowID, err := strconv.Atoi(id)

	if err != nil {
		logger.Error("id should be int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "id should be int ", err.Error())
	}

	var balance float64
	err = h.db.QueryRowContext(ctx, "SELECT balance FROM cloud_pockets WHERE id = $1", rowID).Scan(&balance)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Resp{Message: "Cloud pocket not found"})
	}

	if balance <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, Resp{Message: "Unable to delete this Cloud Pocket. there is amount left in this Cloud Pocket, please move money out and try again"})
	}

	err = h.db.QueryRowContext(ctx, "DELETE FROM cloud_pockets WHERE id = $1", rowID).Scan()

	if err != nil {
		logger.Error("can't execute delete statement", zap.Error(err))
	}

	logger.Info("delete successfully")

	return c.JSON(http.StatusOK, Resp{Message: "Cloud pocket deleted successfully"})

}
