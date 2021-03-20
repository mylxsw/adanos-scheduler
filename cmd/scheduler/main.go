package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/mylxsw/adanos-scheduler/api"
	"github.com/mylxsw/adanos-scheduler/config"
	"github.com/mylxsw/adanos-scheduler/pubsub"
	repoMock "github.com/mylxsw/adanos-scheduler/repo/mock"
	"github.com/mylxsw/adanos-scheduler/scheduler"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	log.DefaultLogFormatter(formatter.NewDefaultCleanFormatter(true))

	app := application.Create(fmt.Sprintf("%s(%s)", Version, GitCommit))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "listen",
		Usage: "http listen addr",
		Value: ":15777",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "log_path",
		Usage: "日志文件输出目录（非文件名），默认为空，输出到标准输出",
	}))

	app.BeforeServerStart(func(cc container.Container) error {
		stackWriter := writer.NewStackWriter()
		cc.MustResolve(func(c infra.FlagContext) {
			logPath := c.String("log_path")
			if logPath == "" {
				stackWriter.PushWithLevels(writer.NewStdoutWriter())
				return
			}

			log.All().LogFormatter(formatter.NewJSONWithTimeFormatter())
			stackWriter.PushWithLevels(writer.NewDefaultRotatingFileWriter(func(le level.Level, module string) string {
				return filepath.Join(logPath, fmt.Sprintf("server-%s.%s.log", le.GetLevelName(), time.Now().Format("20060102")))
			}))
		})

		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(app.Container()),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	app.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{

		}
	})

	app.Provider(repoMock.Provider{})

	app.Provider(api.Provider{})
	app.Provider(service.Provider{})
	app.Provider(pubsub.Provider{})
	app.Provider(scheduler.Provider{})

	app.Main(func(conf *config.Config, publisher event.Publisher) {
		rand.Seed(time.Now().Unix())

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"config": conf,
			}).Debug("configuration")
		}

		publisher.Publish(pubsub.SystemUpDownEvent{
			Up:        true,
			CreatedAt: time.Now(),
		})
	})
	app.BeforeServerStop(func(cc infra.Resolver) error {
		return cc.Resolve(func(publisher event.Publisher) {
			publisher.Publish(pubsub.SystemUpDownEvent{
				Up:        false,
				CreatedAt: time.Now(),
			})
		})
	})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %v", err)
	}
}

type ErrorCollectorWriter struct {
	cc container.Container
}

func NewErrorCollectorWriter(cc container.Container) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{cc: cc}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return nil
}

func (e *ErrorCollectorWriter) ReOpen() error {
	return nil
}

func (e *ErrorCollectorWriter) Close() error {
	return nil
}
