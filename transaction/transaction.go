package transaction

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"fmt"
	// "math/big"
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
}

type CloudPocket struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}


type handler struct {
	db  *sql.DB
}


func New(db *sql.DB) *handler {
	return &handler{db}
}
const (
	cStmt             = "SELECT id,source_cloud_pocket_id,destination_cloud_pocket_id,amount,description,datetime FROM transaction WHERE id=$1;"
	CreateQuery       = "INSERT INTO transaction (source_cloud_pocket_id, destination_cloud_pocket_id, amount, datetime, description) VALUES($1,$2, $3,$4,$5);"
	FindPocketQuery   = "SELECT id, name, category, currency, balance FROM cloud_pockets WHERE id = $1;"
	UpdatePocketQuery = "UPDATE public.cloud_pockets	SET balance=$2 WHERE id=$1 RETURNING id;"
	InsertTransationQuery = "INSERT INTO transaction (source_cloud_pocket_id, destination_cloud_pocket_id, amount, datetime,description) VALUES($1,$2,$3,$4,$5) RETURNING id;"
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
		err = rows.Scan(&t.ID, &t.Soucre_Cloud_Pocket_ID, &t.Destination_Cloud_Pocket_ID, &t.Amount, &t.Description, &t.Date)
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
	SoucrePocket, err := FindidCloudPocket(h, tn.Soucre_Cloud_Pocket_ID)
	if err != nil {
		logger.Error("SoucrePocket not Found", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "SoucrePocket not Found", err.Error())
	}
	DestinationPocket, err := FindidCloudPocket(h, tn.Destination_Cloud_Pocket_ID)
	if err != nil {
		logger.Error("DestinationPocket not Found", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "SoucrePocket not Found", err.Error())
	}
	if SoucrePocket.Balance < tn.Amount {
		logger.Error("Your wallet has insufficient funds for this transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "SoucrePocket balance  not Enough", err.Error())
	}
	err=UpdateCloudPocket(h,SoucrePocket,DestinationPocket,tn.Amount)
	if err != nil {
		logger.Error("UpdateCloudPocket error", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "UpdateCloudPocket error", err.Error())
	}
    fmt.Println("%s",DestinationPocket)
	var lastInsertId int64
	tn.Date = time.Now()
	err = h.db.QueryRowContext(ctx, InsertTransationQuery, tn.Soucre_Cloud_Pocket_ID, tn.Destination_Cloud_Pocket_ID, tn.Amount, tn.Date, tn.Description).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	tn.ID = lastInsertId
	return c.JSON(http.StatusCreated, tn)
}
func FindidCloudPocket(h handler, ID int64) (CloudPocket, error) {
	cp := CloudPocket{}
	row := h.db.QueryRow(FindPocketQuery, ID)
	err := row.Scan(&cp.ID, &cp.Name, &cp.Category, &cp.Currency, &cp.Balance)
	if err != nil {
		return cp, err
	}
	return cp, nil
}
func UpdateCloudPocket(h handler,Soucre CloudPocket,Des CloudPocket,balance float64) ( error) {
	cp := CloudPocket{}
	NewSocureBlance := Soucre.Balance - balance
	NewDesBlance := Des.Balance + balance
	rowSrc := h.db.QueryRow(UpdatePocketQuery, Soucre.ID,NewSocureBlance)
	rowDesc := h.db.QueryRow(UpdatePocketQuery,Des.ID,NewDesBlance)
	errSrc  := rowSrc.Scan(&cp.ID)
	if errSrc != nil {
		return  errSrc
	}
	errDesc := rowDesc.Scan(&cp.ID)
	if errDesc != nil {
		return  errDesc
	}
	return  nil
}
