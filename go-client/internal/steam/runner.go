package steam

import (
	"context"

	"tiny-steam-client/go-client/internal/accounts"
)

type Runner interface {
	FetchCMServers(context.Context) error
	Run(context.Context) error
}

type RunnerOptions struct {
	PreferSteamLib bool
}

func NewRunner(accs []accounts.Account, opts RunnerOptions) Runner {
	if opts.PreferSteamLib {
		if r := newSteamLibRunner(accs); r != nil {
			return r
		}
	}
	return NewManager(accs)
}
