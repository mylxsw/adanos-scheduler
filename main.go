package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/mylxsw/adanos-scheduler/pattern"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/sshx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	client, err := sshx.NewClient(
		"192.168.1.223:22",
		sshx.Credential{User: "root", PrivateKeyPath: filepath.Join(homeDir, ".ssh/id_rsa")},
		sshx.SetLogger(log.Module("ssh")),
		sshx.SetEstablishTimeout(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	if err := client.Handle(func(sub sshx.Client) error {
		processlists, err := sub.Command(context.TODO(), "ps -ef")
		if err != nil {
			return err
		}

		res, err := pattern.Eval(`JQ(".name")`, string(processlists))
		if err != nil {
			panic(err)
		}

		log.Debugf("processlist: %s", res)
		return nil
	}); err != nil {
		panic(err)
	}
}
