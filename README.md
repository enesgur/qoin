## What's the Qoin?  
It is an application which provides bitcoin market data to you through your terminal. It is written with golang 1.9. It gets all datas from coinmarketcap.com continuously. All operations you are done by using datas provided by the application are your own responsibility.
  
## Installation
	  go get -u github.com/enesgur/qoin

## Usage

List all coins:

	$ qoin list

Get price and market datas:

	$ qoin -c symbol-name
	$ qoin -c btc
