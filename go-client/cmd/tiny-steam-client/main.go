package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tiny-steam-client/go-client/internal/accounts"
	"tiny-steam-client/go-client/internal/auth"
	"tiny-steam-client/go-client/internal/cli"
	"tiny-steam-client/go-client/internal/steam"
)

func main() {
	opts, err := cli.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var accs []accounts.Account
	if opts.AccountsFile != "" {
		accs, err = accounts.FromFile(opts.AccountsFile)
	} else {
		accs = accounts.FromSingle(opts.User, opts.Pass, opts.TFC, opts.AC)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mgr := steam.NewManager(accs)
	if err := mgr.FetchCMServers(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "fetch cm servers:", err)
		os.Exit(1)
	}

	if opts.ServerIP != "" && opts.ServerPort > 0 {
		go func() {
			tick := time.NewTicker(20 * time.Second)
			defer tick.Stop()
			aclient := auth.Client{Addr: fmt.Sprintf("%s:%d", opts.ServerIP, opts.ServerPort)}
			for {
				if err := aclient.SendHeartbeat(ctx, accs); err != nil {
					fmt.Fprintln(os.Stderr, "auth client warning:", err)
				}
				select {
				case <-ctx.Done():
					return
				case <-tick.C:
				}
			}
		}()
	}

	if err := mgr.Run(ctx); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "steam manager error:", err)
		os.Exit(1)
	}
}
