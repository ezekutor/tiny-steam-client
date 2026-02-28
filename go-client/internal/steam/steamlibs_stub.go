//go:build !steamlibs

package steam

import "tiny-steam-client/go-client/internal/accounts"

func newSteamLibRunner(_ []accounts.Account) Runner { return nil }
