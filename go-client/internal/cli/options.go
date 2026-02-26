package cli

import (
	"errors"
	"flag"
)

type Options struct {
	User         string
	Pass         string
	TFC          string
	AC           string
	ServerIP     string
	ServerPort   int
	AccountsFile string
}

func Parse() (Options, error) {
	var o Options
	flag.StringVar(&o.User, "user", "", "Steam account username")
	flag.StringVar(&o.Pass, "pw", "", "Steam account password")
	flag.StringVar(&o.TFC, "tfc", "", "Steam two factor code")
	flag.StringVar(&o.AC, "ac", "", "Steam auth code")
	flag.StringVar(&o.ServerIP, "sip", "", "Tiny csgo server ip")
	flag.IntVar(&o.ServerPort, "sport", 0, "Tiny csgo server port")
	flag.StringVar(&o.AccountsFile, "acfile", "", "Accounts file path")
	flag.Parse()

	if o.TFC != "" && o.AC != "" {
		return o, errors.New("you can only provide one of -tfc or -ac")
	}
	if o.AccountsFile == "" && (o.User == "" || o.Pass == "") {
		return o, errors.New("if -acfile is not provided, -user and -pw are required")
	}
	if (o.ServerIP == "") != (o.ServerPort == 0) {
		return o, errors.New("-sip and -sport must be provided together")
	}
	return o, nil
}
