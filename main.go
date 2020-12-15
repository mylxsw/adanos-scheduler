package main

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/pattern"
	"github.com/mylxsw/sshx"
)

type Server struct {
	pattern.Helpers
	Host       string
	Credential *sshx.Credential
	Labels     []string
}

var servers = []Server{
	{Host: "192.168.1.223", Labels: []string{"server-223", "office-agent"}},
	{Host: "192.168.1.225", Labels: []string{"server-225", "test1"}},
	{Host: "192.168.1.226", Labels: []string{"server-226", "test1"}},
	{Host: "192.168.1.231", Labels: []string{"server-231", "dev"}},
	{Host: "192.168.1.170", Labels: []string{"server-170", "docker"}, Credential: &sshx.Credential{
		User:           "root",
		PrivateKeyPath: filepath.Join(getHomeDir(), ".ssh/id_rsa_prometheus"),
	}},
}

func main() {
	log.DefaultLogFormatter(formatter.NewDefaultCleanFormatter(true))

	credential := sshx.Credential{User: "root", PrivateKeyPath: filepath.Join(getHomeDir(), ".ssh/id_rsa")}

	matcher, err := pattern.NewMatcher(`any(Labels, {# == "test1"})`, Server{})
	if err != nil {
		panic(err)
	}

	hostRegexp, _ := regexp.Compile(`:\d+$`)
	coll.MustNew(servers).Map(func(server Server) Server {
		if server.Credential == nil {
			server.Credential = &credential
		}
		if !hostRegexp.MatchString(server.Host) {
			server.Host += ":22"
		}

		return server
	}).Filter(func(server Server) bool {
		matched, err := matcher.Match(server)
		if err != nil {
			log.Errorf("match failed: %v", err)
		}
		return matched
	}).Each(func(server Server) {
		client, err := sshx.NewClient(server.Host, *server.Credential, sshx.SetLogger(log.Module("ssh")), sshx.SetEstablishTimeout(3*time.Second))
		if err != nil {
			log.Errorf("can not connect to %s: %w", server.Host, err)
			return
		}

		resp, err := client.Command(context.TODO(), "ip addr | grep 192.168")
		if err != nil {
			log.Errorf("execute command on remote server %s failed: %w", server.Host, err)
		}

		log.Debugf("%s -> %s", server.Host, string(resp))
	})
}

func getHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func TestServerSSH() {
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
