//go:build steamlibs

package steam

import (
	"context"
	"fmt"

	_ "github.com/SteamDatabase/Protobufs"
	steamlib "github.com/paralin/go-steam"

	"tiny-steam-client/go-client/internal/accounts"
)

type steamLibRunner struct {
	accounts []accounts.Account
}

func newSteamLibRunner(accs []accounts.Account) Runner {
	return &steamLibRunner{accounts: accs}
}

func (r *steamLibRunner) FetchCMServers(_ context.Context) error {
	return nil
}

func (r *steamLibRunner) Run(ctx context.Context) error {
	for _, acc := range r.accounts {
		client := steamlib.NewClient()
		if err := client.Connect(); err != nil {
			return fmt.Errorf("connect steam for %s: %w", acc.User, err)
		}
		go func(c *steamlib.Client, a accounts.Account) {
			<-ctx.Done()
			c.Disconnect()
		}(client, acc)
	}
	<-ctx.Done()
	return nil
}
