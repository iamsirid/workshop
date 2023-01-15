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
	Soucre_Cloud_Pocket_ID      int64     `json:"source_cloud_pocket_id"`
	Destination_Cloud_Pocket_ID int64     `json:"destination_cloud_pocket_id"`
	Amount                      float64   `json:"amount"`
	Description                 string    `json:"description"`
	Date                        time.Time `json:"date"`
	Status                      string    `json:"status"`
}

type CloudPocket struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

const (
	cStmt = "SELECT id,source_cloud_pocket_id,destination_cloud_pocket_id,amount,description,datetime,status FROM transaction WHERE id=$1;",
	CreateQuery="INSERT INTO transaction (source_cloud_pocket_id, destination_cloud_pocket_id, amount, datetime, description, status) VALUES($1,$2, $3,$4,$5,$6);
	",
	FindPocketQuery="SELECT id, \"name\", category, currency, balance
	FROM public.cloud_pockets where id=$1;"
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
		err = rows.Scan(&t.ID, &t.Soucre_Cloud_Pocket_ID, &t.Destination_Cloud_Pocket_ID, &t.Amount, &t.Description, &t.Date, &t.Status)
		if err != nil {
			logger.Error("row scan error", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		tn = append(tn, t)
	}
	logger.Info("send successfully", zap.String("id", id))
	return c.JSON(http.StatusOK, tn)
}

func (h handler) CreateTransaction(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var tn Transaction
	err := c.Bind(&tn)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	SoucrePocket,err:=FindidCloudPocket(h,tn.Soucre_Cloud_Pocket_ID)
    if err !=nil{
		logger.Error("SoucrePocket not Found", zap.Error(err))
		return echo.NewHTTPError(StatusBadRequest"SoucrePocket not Found", err.Error())
	}
	DestinationPocket:=FindidCloudPocket(h,tn.Destination_Cloud_Pocket_ID)
	if err !=nil{
		logger.Error("DestinationPocket not Found", zap.Error(err))
		return echo.NewHTTPError(StatusBadRequest"SoucrePocket not Found", err.Error())
	}
	if SoucrePocket.balance<tn.amount{
		logger.Error("SoucrePocket balance  not Enough", zap.Error(err))
		return echo.NewHTTPError(StatusBadRequest"SoucrePocket balance  not Enough", err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, ac.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	ac.ID = lastInsertId
	return c.JSON(http.StatusCreated, ac)
}
func FindidCloudPocket(h handler,int64 id) (CloudPocket,error){
	var lastInsertId int64
	cp:=CloudPocket{}
	row,err = h.db.Query(FindPocketQuery,id)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return  nil,err
	}
	err=row.Scan($cp.ID,$cp.Name,$cp.Category,$cp.Currency,$cp.Balance)
	if err != nil {
		logger.Error("row scan error", zap.Error(err))
		return nil,err
	}
	return cp,nil
}
func UpdateCloudPocket(h handler,int64 id) (CloudPocket,error){
	var lastInsertId int64
	cp:=CloudPocket{}
	row,err = h.db.Query(FindPocketQuery,id)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return  nil,err
	}
	err=row.Scan($cp.ID,$cp.Name,$cp.Category,$cp.Currency,$cp.Balance)
	if err != nil {
		logger.Error("row scan error", zap.Error(err))
		return nil,err
	}
	return cp,nil
}
