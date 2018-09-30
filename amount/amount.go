package amount

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

const Zero = Amount("00000000")

var (
	ErrorNumString      = errors.New("Error Num")
	ErrorFractionDigits = errors.New("Error Fraction Digits")
)

var bigBase = new(big.Int)

type Amount string

func (a Amount) String() string {
	return a.FloatString()
}

func (a Amount) Float64() float64 {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	return float64(bigA.Int64()) / 1e8
}

func (a Amount) CmpFloat(x float64) int {
	xTmp, _ := FromFloat64(x)
	return a.Cmp(xTmp)
}

func (a Amount) Cmp(x Amount) int {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	bigX := new(big.Int)
	bigX.SetString(string(x), 10)
	return bigA.Cmp(bigX)
}

func (a Amount) DivFloat64(x float64) Amount {
	amountX, _ := FromFloat64(x)
	return a.Div(amountX)
}

func (a Amount) Div(x Amount) Amount {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	bigX := new(big.Int)
	bigX.SetString(string(x), 10)
	bigA.Mul(bigA, bigBase)
	bigA.Div(bigA, bigX)
	numBigStr := formatBigInt28Digits(bigA)
	return Amount(numBigStr)
}

func (a Amount) Mul(x Amount) Amount {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	bigX := new(big.Int)
	bigX.SetString(string(x), 10)
	bigA.Mul(bigA, bigX)
	bigA.Div(bigA, bigBase)
	numBigStr := formatBigInt28Digits(bigA)
	return Amount(numBigStr)
}

func (a Amount) Sub(x Amount) Amount {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	bigX := new(big.Int)
	bigX.SetString(string(x), 10)
	bigA.Sub(bigA, bigX)
	numBigStr := formatBigInt28Digits(bigA)
	return Amount(numBigStr)
}

func (a Amount) Add(x Amount) Amount {
	bigA := new(big.Int)
	bigA.SetString(string(a), 10)
	bigX := new(big.Int)
	bigX.SetString(string(x), 10)
	bigA.Add(bigA, bigX)
	numBigStr := formatBigInt28Digits(bigA)
	return Amount(numBigStr)
}

func (a Amount) Minus() Amount {
	tmp := string(a)
	if strings.HasPrefix(tmp, "-") {
		return a
	}
	if a == Zero {
		return a
	}
	return "-" + a

}

func (a Amount) FloatStringShort() string {
	if a == "" {
		return "0.0"
	}
	if a == Zero {
		return "0.0"
	}
	amountLen := len(a)
	intPart := ""
	numPartLen := 0
	if amountLen == 8 {
		intPart = "0"
	} else if amountLen == 9 && a[0] == '-' {
		intPart = "-0"
		numPartLen = 1
	} else {
		numPartLen = amountLen - 8
		intPart = string(a[:numPartLen])
	}
	fracPart := string(a[numPartLen:])
	fracLastNotZeroIndex := strings.LastIndexFunc(fracPart,
		func(c rune) bool {
			if c != '0' {
				return true
			}
			return false
		})
	if intPart == "0" && fracLastNotZeroIndex == -1 {
		return "0.0"
	} else if fracLastNotZeroIndex == -1 {
		return intPart + ".0"
	} else {
		return intPart + "." + fracPart[:fracLastNotZeroIndex+1]
	}
}

func (a Amount) FloatString() string {
	if a == "" {
		return "0.0"
	}
	if a == Zero {
		return "0.0"
	}
	strLen := len(a)
	if strLen == 8 {
		return "0." + string(a)
	}
	if strLen == 9 && a[0] == '-' {
		return "-0." + string(a[1:])
	}
	numPartLen := strLen - 8
	fracPart := string(a[numPartLen:])
	if fracPart == "00000000" {
		return string(a[:numPartLen]) + ".0"
	}
	return string(a[:numPartLen]) + "." + string(a[numPartLen:])
}

func ValidateAmountStr(value string) error {
	numParts := strings.Split(value, ".")
	numPartLen := len(numParts)
	if numPartLen > 2 {
		return ErrorNumString
	}
	if numPartLen == 2 && len(numParts[1]) > 8 {
		return ErrorFractionDigits
	}
	numBig := new(big.Int)
	_, ok := numBig.SetString(numParts[0], 10)
	if !ok && numParts[0] != "" && numParts[0] != "-" {
		return ErrorNumString
	}
	return nil
}

func FromNumString(value string) (Amount, error) {
	if value == "" {
		return Zero, nil
	}
	numParts := strings.Split(value, ".")
	numPartLen := len(numParts)
	if numPartLen > 2 {
		return Zero, ErrorNumString
	}
	if numPartLen == 2 && len(numParts[1]) > 8 {
		return Zero, ErrorFractionDigits
	}
	numBig := new(big.Int)
	_, ok := numBig.SetString(numParts[0], 10)
	if !ok && numParts[0] != "" && numParts[0] != "-" {
		return Zero, ErrorNumString
	}
	numBig.Mul(numBig, bigBase)
	if numPartLen == 2 {
		fracPartLen := len(numParts[1])
		append0Nums := 8 - fracPartLen
		fracPartStr := numParts[1] + strings.Repeat("0", append0Nums)
		fracBig := new(big.Int)
		_, ok := fracBig.SetString(fracPartStr, 10)
		if !ok {
			return Zero, ErrorNumString
		}
		isMinus := strings.HasPrefix(value, "-")
		if isMinus {
			numBig.Sub(numBig, fracBig)
		} else {
			numBig.Add(numBig, fracBig)
		}

	}
	numBigStr := formatBigInt28Digits(numBig)
	return Amount(numBigStr), nil
}

func formatBigInt28Digits(numBig *big.Int) string {
	numBigStr := numBig.String()
	isMinus := strings.HasPrefix(numBigStr, "-")
	if isMinus {
		numBigStr = strings.Replace(numBigStr, "-", "", -1)
	}
	numBigStrLen := len(numBigStr)
	if numBigStrLen < 8 {
		append0Nums := 8 - numBigStrLen
		numBigStr = strings.Repeat("0", append0Nums) + numBigStr
	}
	if isMinus {
		numBigStr = "-" + numBigStr
	}
	return numBigStr
}

func FromFloat64(value float64) (Amount, error) {
	return FromNumString(strconv.FormatFloat(value, 'f', -1, 64))
}

// func New() Amount {
// return Zero
// }

func init() {
	bigBase.SetInt64(1e8)
}
