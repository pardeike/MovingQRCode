package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	SessionID string
	Secret    []byte
	CreatedAt time.Time
}

var session *Session

// generateRandomBytes returns securely generated random bytes.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// generateSession creates a new session with a random session ID and secret.
func generateSession() (*Session, error) {
	sessionIDBytes, err := generateRandomBytes(16)
	if err != nil {
		return nil, err
	}
	secretBytes, err := generateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	return &Session{
		SessionID: hex.EncodeToString(sessionIDBytes),
		Secret:    secretBytes,
		CreatedAt: time.Now(),
	}, nil
}

// computeToken calculates the HMAC-SHA256 token using the secret and the message "sessionID:timestamp".
func computeToken(secret []byte, sessionID, timestamp string) string {
	data := sessionID + ":" + timestamp
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// parseTimestamp converts a Unix timestamp string into a time.Time value.
func parseTimestamp(ts string) (time.Time, error) {
	seconds, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(seconds, 0), nil
}

func main() {
	var err error
	session, err = generateSession()
	if err != nil {
		fmt.Println("Error generating session:", err)
		return
	}
	// Display the seed that is sent to the third party.
	fmt.Println("Session generated. Seed sent to third party:")
	fmt.Println("Session ID:", session.SessionID)
	fmt.Println("Secret:", hex.EncodeToString(session.Secret))
	fmt.Println("Waiting for QR code input (format: sessionID,timestamp,token):")
	time.Sleep(5)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)
	parts := strings.Split(input, ",")
	if len(parts) != 3 {
		fmt.Println("Invalid input format. Expected sessionID,timestamp,token")
		return
	}
	receivedSessionID := parts[0]
	timestampStr := parts[1]
	receivedToken := parts[2]

	// Invalidate the session immediately after any QR content is received.
	defer func() { session = nil }()

	// Verify session ID.
	if receivedSessionID != session.SessionID {
		fmt.Println("Session ID does not match.")
		return
	}

	// Parse the provided timestamp.
	ts, err := parseTimestamp(timestampStr)
	if err != nil {
		fmt.Println("Invalid timestamp:", err)
		return
	}

	// Check that the timestamp is within a 10-second window of the current time.
	now := time.Now()
	diff := now.Sub(ts)
	if diff < 0 {
		diff = -diff
	}
	if diff > 10*time.Second {
		fmt.Println("Timestamp is not within the valid 10-second window.")
		return
	}

	expectedToken := computeToken(session.Secret, session.SessionID, timestampStr)
	if hmac.Equal([]byte(receivedToken), []byte(expectedToken)) {
		fmt.Println("Valid token!")
	} else {
		fmt.Println("Invalid token!")
	}
}
