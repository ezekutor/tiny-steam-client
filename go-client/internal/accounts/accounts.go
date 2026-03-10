package accounts

import (
	"encoding/json"
	"fmt"
	"os"
)

type Account struct {
	User string
	Pass string
	TFC  string
	AC   string
}

type accountsFile struct {
	Accounts map[string]struct {
		User string `json:"user"`
		Pass string `json:"passwd"`
	} `json:"accounts"`
}

func FromSingle(user, pass, tfc, ac string) []Account {
	return []Account{{User: user, Pass: pass, TFC: tfc, AC: ac}}
}

func FromFile(path string) ([]Account, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read accounts file: %w", err)
	}
	var af accountsFile
	if err := json.Unmarshal(b, &af); err != nil {
		return nil, fmt.Errorf("parse accounts json: %w", err)
	}
	out := make([]Account, 0, len(af.Accounts))
	for _, a := range af.Accounts {
		if a.User == "" || a.Pass == "" {
			continue
		}
		out = append(out, Account{User: a.User, Pass: a.Pass})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid accounts found")
	}
	return out, nil
}
