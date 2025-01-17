package router

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/account"
	cloud_pockets "github.com/kkgo-software-engineering/workshop/cloud-pockets"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/kkgo-software-engineering/workshop/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(db)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hAccount := account.New(cfg.FeatureFlag, db)
	e.POST("/accounts", hAccount.Create)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	hCloudPocket := cloud_pockets.New(db)
	e.POST("/cloud-pockets", hCloudPocket.CreateCloudPocket)
	hTransaction := transaction.New(db)
	e.GET("/cloud-pockets/:id/transactions", hTransaction.GetTransactionbyAccountid)
	e.GET("/cloud-pockets", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "cloud-pocket")
	})

	e.GET("/cloud-pockets", hCloudPocket.GetAllCloudPocket)
	e.GET("/cloud-pockets/:id", hCloudPocket.GetCloudPocketById)
	e.DELETE("/cloud-pockets/:id", hCloudPocket.DeleteCloudPocket)
	e.POST("/cloud-pockets/transfer",hTransaction.CreateTransaction)
    e.GET("/cloud-pockets/:id/csv",hTransaction.GetCsv)
	return e
}
