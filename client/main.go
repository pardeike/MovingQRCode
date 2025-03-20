package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"strconv"
	"time"
)

// computeToken calculates the HMAC-SHA256 token.
func computeToken(secret []byte, sessionID, timestamp string) string {
	data := sessionID + ":" + timestamp
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	sessionID := flag.String("session", "", "Session ID")
	secretHex := flag.String("secret", "", "Secret in hex")
	flag.Parse()
	if *sessionID == "" || *secretHex == "" {
		fmt.Println("Both -session and -secret must be provided.")
		return
	}
	secret, err := hex.DecodeString(*secretHex)
	if err != nil {
		fmt.Println("Error decoding secret:", err)
		return
	}

	fmt.Println("Generating 10 QR code tokens (one per second):")
	for i := 0; i < 10; i++ {
		// For each iteration, generate a token with a timestamp offset.
		ts := time.Now().Unix() + int64(i)
		tsStr := strconv.FormatInt(ts, 10)
		token := computeToken(secret, *sessionID, tsStr)
		// The QR code content is a comma-separated string: sessionID,timestamp,token.
		qrContent := *sessionID + "," + tsStr + "," + token
		fmt.Printf("Iteration %d: %s\n", i+1, qrContent)
		// time.Sleep(1 * time.Second)
	}
}
