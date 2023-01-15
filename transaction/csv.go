package transaction		
import (
	"encoding/csv"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"os"
	"log"
	"fmt"
)

func (h handler) GetCsv(c echo.Context) error {
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
	file, err := os.Create("records.csv")
    if err != nil {
        log.Fatalln("failed to open file", err)
    }
    w := csv.NewWriter(file)
    defer w.Flush()
	
    // Using WriteAll
    var data [][]string
    for _, t := range tn {
        row := []string{strconv.FormatInt(t.ID,10),strconv.FormatInt(t.Soucre_Cloud_Pocket_ID,10),strconv.FormatInt(t.Destination_Cloud_Pocket_ID,10),fmt.Sprintf("%v", t.Amount),t.Description,t.Date.String()}
        data = append(data, row)
    }
    w.WriteAll(data)
	return c.File("records.csv")
}

