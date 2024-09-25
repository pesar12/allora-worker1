package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

type Kline struct {
	OpenTime  time.Time
	CloseTime time.Time
	Interval  string
	Symbol    string
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	Closed    bool
}

func main() {

	cfg := &envConfig{
		APIKey: os.Getenv("CMC_APIKEY"), // تغییر به CoinMarketCap APIKey
		RPC:    os.Getenv("RPC"),
	}

	fmt.Println("CMC_API_KEY: ", cfg.APIKey)
	fmt.Println("RPC: ", cfg.RPC)

	router := gin.Default()

	router.GET("/inference/:token", func(c *gin.Context) {
		token := c.Param("token")
		if token == "MEME" {
			handleMemeRequest(c, cfg)
			return
		}

		symbol := fmt.Sprintf("%s", token)

		price, err := getCryptoPrice(symbol, cfg.APIKey)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching price: %v", err))
			return
		}

		c.String(http.StatusOK, strconv.FormatFloat(price, 'f', 2, 64))
	})

	router.Run(":8000")
}

func handleMemeRequest(c *gin.Context, cfg *envConfig) {

	if cfg.APIKey == "" {
		c.String(400, "need api key")
	}

	if cfg.RPC == "" {
		panic("Invalid env.json file")
	}

	lb, err := getLatestBlock(cfg.RPC)
	if err != nil {
		fmt.Println(err)
		return
	}

	meme, err := getMemeOracleData(lb, cfg.APIKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	mp, err := getMemePrice(meme.Data.Platform, meme.Data.Address)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nBlockHeight: \"%s\", Meme: \"%s\", Platform: \"%s\", Price: \"%s\"\n\n",
		lb, meme.Data.TokenSymbol, meme.Data.Platform, mp)

	mpf, _ := strconv.ParseFloat(mp, 64)

	c.String(http.StatusOK, strconv.FormatFloat(random(mpf), 'g', -1, 64))
}

// تابع جدید برای دریافت قیمت از کوین مارکت کپ
func getCryptoPrice(symbol string, apiKey string) (float64, error) {
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s&convert=USD", symbol)
	req,
