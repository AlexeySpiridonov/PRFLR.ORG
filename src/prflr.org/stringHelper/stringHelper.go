package stringHelper

import (
    "crypto/rand"
    "math/big"
)

func RandomString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

    symbols := big.NewInt(int64(len(alphanum)))
    states  := big.NewInt(0)
    states.Exp(symbols, big.NewInt(int64(n)), nil)

    r, err := rand.Int(rand.Reader, states)
    if err != nil {
        return ""
    }

    var bytes = make([]byte, n)

    r2     := big.NewInt(0)
    symbol := big.NewInt(0)
    for i := range bytes {
        r2.DivMod(r, symbols, symbol)
        r, r2 = r2, r
        bytes[i] = alphanum[symbol.Int64()]
    }

    return string(bytes)
}

func GetCappedCollectionNameForApiKey(apiKey string) string {
    collectionName := "timers_" + apiKey

    if len(collectionName) > 125 {
        collectionName = collectionName[:125] // NO More than 125 chars!
    }

    return collectionName
}