package web3

import (
	"math/big"
	"net/http"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/shopspring/decimal"
)

func handleAPIErrorResponse(res *http.Response) error {
	switch res.StatusCode {
	case 400:
		return &APIError{}
	case 403:
		return &APIError{}
	case 404:
		return &APIError{}
	case 429:
		return &APIError{}
	case 418:
		return &APIError{}
	case 500:
		return &APIError{}
	default:
		return &APIError{}
	}
}

func (c *apiClient) handleRequest(req *http.Request) (res *http.Response, err error) {
	req.Header.Add("Content-Type", "application/json")
	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return
	}
	res, err = c.client.Do(rreq)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return res, handleAPIErrorResponse(res)
	}

	return res, nil
}

func convHexToDec(hexa string) (string, error) {
	hexa = hexa[2:]
	decimal, err := strconv.ParseInt(hexa, 16, 64)

	if err != nil {
		return "0", err
	}
	stringValue := strconv.FormatInt(decimal, 10)
	return stringValue, nil
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)

	default:
		return false
	}
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}
