package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"tiny-steam-client/go-client/internal/accounts"
)

const cmURL = "https://api.steampowered.com/ISteamDirectory/GetCMListForConnect/v1/?cellid=0&format=json"

type cmResp struct {
	Response struct {
		Serverlist []string `json:"serverlist"`
	} `json:"response"`
}

type Manager struct {
	accounts []accounts.Account
	servers  []string
}

func NewManager(acc []accounts.Account) *Manager {
	return &Manager{accounts: acc}
}

func (m *Manager) FetchCMServers(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cmURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected cm status: %s", resp.Status)
	}
	var cr cmResp
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return err
	}
	if len(cr.Response.Serverlist) == 0 {
		return fmt.Errorf("steam returned empty CM list")
	}
	m.servers = cr.Response.Serverlist
	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	if len(m.servers) == 0 {
		if err := m.FetchCMServers(ctx); err != nil {
			return err
		}
	}
	var wg sync.WaitGroup
	for _, acc := range m.accounts {
		acc := acc
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.runAccountLoop(ctx, acc)
		}()
	}
	wg.Wait()
	return nil
}

func (m *Manager) runAccountLoop(ctx context.Context, account accounts.Account) {
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		server := m.servers[idx%len(m.servers)]
		idx++
		dialer := net.Dialer{Timeout: 6 * time.Second}
		conn, err := dialer.DialContext(ctx, "tcp", server)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
		_, _ = conn.Write([]byte{})
		_ = conn.Close()
		time.Sleep(2 * time.Second)
	}
}
