package supply

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/canxium/supply-information/config"
)

type handlers struct {
	client *ethclient.Client
	cfg    *config.Config
	db     *gorm.DB
}

func NewHandlers(cfg *config.Config, db *gorm.DB) *handlers {
	client, err := ethclient.Dial(cfg.Supply.RpcApi)
	if err != nil {
		log.Fatal(err)
	}

	return &handlers{cfg: cfg, client: client, db: db}
}

func (h *handlers) getBalance(addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := h.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (h *handlers) getTotalSupply() (*big.Int, error) {
	var sum string
	// sum all address balances from blockscout explorer database
	err := h.db.Table("addresses").Select("sum(fetched_coin_balance)").Row().Scan(&sum)
	supply, _ := new(big.Int).SetString(sum, 10)
	return supply, err
}

func (h *handlers) GetSupplyInfo(c echo.Context) error {
	switch c.QueryParam("q") {
	case "totalSupply":
		totalSupply, err := h.getTotalSupply()
		if err != nil {
			fmt.Printf("%v\n", err)
			return c.JSON(http.StatusInternalServerError, "")
		}

		tmp, _ := new(big.Float).SetString(totalSupply.String())
		supply, _ := new(big.Float).Quo(tmp, big.NewFloat(math.Pow10(18))).Float64()
		return c.JSON(http.StatusOK, supply)
	default:
		totalSupply, err := h.getTotalSupply()
		if err != nil {
			fmt.Printf("%v\n", err)
			return c.JSON(http.StatusInternalServerError, "")
		}

		tmp, _ := new(big.Float).SetString(totalSupply.String())
		supply, _ := new(big.Float).Quo(tmp, big.NewFloat(math.Pow10(18))).Float64()

		// sub foundation funds balances
		accounts := strings.Split(h.cfg.Supply.Addresses, ":")
		for _, account := range accounts {
			balance, err := h.getBalance(account)
			if err != nil {
				return err
			}

			totalSupply.Sub(totalSupply, balance)
		}

		balance, _ := new(big.Float).SetString(totalSupply.String())
		circulating := new(big.Float).Quo(balance, big.NewFloat(math.Pow10(18)))
		circulatingSupply, _ := circulating.Float64()

		if c.QueryParam("q") == "" {
			return c.JSON(http.StatusOK, map[string]float64{
				"total_supply":       supply,
				"circulating_supply": circulatingSupply,
			})
		}

		return c.JSON(http.StatusOK, circulatingSupply)
	}
}
