package transaction

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	ID                          int64     `json:"id"`
	SOURCE_CLOUD_POCKET_ID      int64     `json:"source_cloud_pocket_id"`
	DESTINATION_CLOUD_POCKET_ID int64     `json:"destination_cloud_pocket_id"`
	AMOUNT                      float64   `json:"amount"`
	DESCRIPTION                 string    `json:"description"`
	DATE                        time.Time `json:"date"`
	STATUS                      string    `json:"status"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

const (
	cStmt = "SELECT id,source_cloud_pocket_id,destination_cloud_pocket_id,amount,description,datetime,status FROM transaction WHERE id=$1;"
)

func (h handler) GetTransactionbyAccountid(c echo.Context) error {
	id := c.Param("id")
	logger := mlog.L(c)
	rowID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("id should be int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "id should be int ", err.Error())
	}
	tn := []Transaction{}
	rows, err := h.db.Query(cStmt, rowID)
	if err != nil {
		logger.Error("qury row error", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.ID, &t.SOURCE_CLOUD_POCKET_ID, &t.DESTINATION_CLOUD_POCKET_ID, &t.AMOUNT, &t.DESCRIPTION, &t.DATE, &t.STATUS)
		if err != nil {
			logger.Error("quey row error", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		tn = append(tn, t)
	}
	logger.Info("send successfully", zap.String("id", id))
	return c.JSON(http.StatusOK, tn)
}
