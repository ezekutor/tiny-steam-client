package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"tiny-steam-client/go-client/internal/accounts"
)

type Client struct {
	Addr string
}

func (c Client) SendHeartbeat(ctx context.Context, accs []accounts.Account) error {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", c.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := struct {
		Type      string   `json:"type"`
		Accounts  []string `json:"accounts"`
		Timestamp int64    `json:"timestamp"`
	}{
		Type:      "tiny-steam-client-go-heartbeat",
		Accounts:  make([]string, 0, len(accs)),
		Timestamp: time.Now().Unix(),
	}
	for _, a := range accs {
		payload.Accounts = append(payload.Accounts, a.User)
	}
	b, _ := json.Marshal(payload)
	if _, err := conn.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("send auth heartbeat: %w", err)
	}
	return nil
}
