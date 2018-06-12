package main

import (
	"github.com/urfave/cli"
	"github.com/leekchan/accounting"
	"os"
	"log"
	"fmt"
	"strings"
)

const VERSION  = "0.1.0"

const Reset = "\033[0m"
const Red ="\033[0;31m"          // Red
const Green ="\033[0;32m"        // Green

// Bold
const BRed="'\033[1;31m'"         // Red
const BGreen="'\033[1;32m'"       // Green

// Underline
const URed="'\033[4;31m'"         // Red
const UGreen="'\033[4;32m'"       // Green

// Background
const On_Red="'\033[41m'"         // Red
const On_Green="'\033[42m'"       // Green

// High Intensity
const IRed="'\033[0;91m'"         // Red
const IGreen="'\033[0;92m'"       // Green

// Bold High Intensity
const BIRed="'\033[1;91m'"        // Red
const BIGreen="'\033[1;92m'"      // Green

// High Intensity backgrounds
const On_IRed="'\033[0;101m'"     // Red
const On_IGreen="'\033[0;102m'"   // Green

func main() {
	market := Market{}

	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "qoin"
	app.Usage = "CLI bitcoin market tools"
	app.Description = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "coin, c",
			Usage: "Coin symbol e.g 'BTC', 'XRP','ETH'",
		},
	}

	app.Action = func(c *cli.Context) error {
		var coin string
		if len(c.String("c")) == 0 {
			coin = c.Args().Get(0)
		} else {
			coin = c.String("c")
		}
		coin = strings.ToUpper(coin)
		ticker := market.ticker()

		money := accounting.Accounting{Symbol: "$", Precision: 0}
		supply := accounting.Accounting{Symbol: ""};

		for _, data := range ticker {
			if data.Symbol != coin {
				continue
			}

			var changes1h string
			if data.Quotes["USD"].PercentChange1h > 0 {
				changes1h = Green
			} else {
				changes1h = Red
			}

			var changes24h string
			if data.Quotes["USD"].PercentChange24h > 0 {
				changes24h = Green
			} else {
				changes24h = Red
			}

			var changes7d string
			if data.Quotes["USD"].PercentChange7d > 0 {
				changes7d = Green
			} else {
				changes7d = Red
			}

			fmt.Printf("Name: %s\n", data.Name)
			fmt.Printf("Rank: %d\n", data.Rank)
			fmt.Printf("Total Supply: %s\n", supply.FormatMoney(FloatToInt(data.TotalSupply)))
			fmt.Printf("Max Supply: %s\n", supply.FormatMoney(FloatToInt(data.MaxSupply)))
			fmt.Printf("Market Capacity: %s\n", money.FormatMoney(FloatToInt(data.Quotes["USD"].MarketCap)))
			fmt.Printf("Volume (24H): %s\n", money.FormatMoney(FloatToInt(data.Quotes["USD"].Volume24h)))
			fmt.Printf("Changes in 1 Hour: %s%.2f%%%s\n", changes1h, data.Quotes["USD"].PercentChange1h, Reset)
			fmt.Printf("Changes in 24 Hours: %s%.2f%%%s\n", changes24h, data.Quotes["USD"].PercentChange24h, Reset)
			fmt.Printf("Changes in 7 Days: %s%.2f%%%s\n", changes7d, data.Quotes["USD"].PercentChange7d, Reset)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Show all coins name and symbol",
			Action: func(c *cli.Context) error {
				fmt.Println("Symbol\tName")
				response := market.listing()
				for _, data := range response.Data {
					fmt.Printf("%s\t%s\n", data.Symbol, data.Name)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
