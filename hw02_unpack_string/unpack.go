package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrInvalidNumber = errors.New("invalid number")
)

const (
	tNumber int = iota
	tLetter
	tShield
	tOther
)

func Unpack(in string) (string, error) {
	var (
		isShielded = false
		count      int
		result     strings.Builder
		err        error
	)

	for i := 0; i < len(in); i++ {
		symbol := in[i]
		symbolType := getSymbolType(symbol)

		if !isShielded && symbolType == tNumber {
			return "", ErrInvalidString
		}
		if len(in) <= i+1 {
			result.WriteByte(symbol)
			break
		}

		nextSymbol := in[i+1]
		nextSymbolType := getSymbolType(nextSymbol)

		switch {
		case !isShielded && symbolType == tShield && (nextSymbolType == tLetter || nextSymbolType == tOther):
			return "", ErrInvalidString
		case !isShielded && symbolType == tShield:
			isShielded = true
			continue
		case nextSymbolType == tNumber:
			count, err = strconv.Atoi(string(nextSymbol))
			if err != nil {
				return "", ErrInvalidNumber
			}
			i++
		default:
			count = 1
		}

		result.WriteString(strings.Repeat(string(symbol), count))
		isShielded = false
	}
	return result.String(), nil
}

func getSymbolType(s byte) int {
	switch {
	case s == '\\':
		return tShield
	case s >= 'a' && s <= 'z':
		return tLetter
	case s >= '0' && s <= '9':
		return tNumber
	default:
		return tOther
	}
}
