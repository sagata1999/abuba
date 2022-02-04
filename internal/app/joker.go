package app

import (
	"awesomeProject/internal/pkg/config"
	"awesomeProject/internal/pkg/version"
	"fmt"
	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func initLogger(logLevel string, dsn, filename string, tags map[string]string, maxAge, maxSize, maxBackups int) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("invalid log level %s, error level will be used", logLevel)
		level = log.ErrorLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})

	hook, err := logrus_sentry.NewAsyncWithTagsSentryHook(dsn, tags, []log.Level{
		log.ErrorLevel,
		log.PanicLevel,
		log.FatalLevel,
	})
	if err != nil {
		log.Errorf("connect to sentry: %s", err)

		return
	}

	hook.StacktraceConfiguration.Enable = true

	log.AddHook(hook)
}

// Run ...
func Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Recovered in main", r)
		}
	}()

	application := kingpin.New("Joker app. ", "Awesome application to get jokes by category or random.")
	configPath := application.Flag(
		"config", "").Short('c').Default("./configs/main.yml").String()

	log.Debug("Joker app start\nTo interrupt execution press Ctrl+C...")
	application.Version(fmt.Sprintf("Joker version: %s\nCopyright (C) 2021  AwesomeCo\n\n"+
		"Compiler: %s\nOn: %s\nSystem: %s\n", version.Version,
		strings.ReplaceAll(version.Compiler, "_", " "),
		strings.ReplaceAll(version.BuildTime, "_", " "),
		strings.ReplaceAll(version.System, "_", " ")))

	_, err := application.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	config.InitConfig(*configPath)

	initLogger(viper.GetString("log.level"), viper.GetString("log.dsn"),
		viper.GetString("log.file"), viper.GetStringMapString("log.tags"), viper.GetInt("log.age"), viper.GetInt("log.size"), viper.GetInt("log.backups"))

	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	sig := <-gracefulStop
	log.Warnf("Caught sig: %+v", sig)
	log.Info("Wait to finish processing")

	log.Debug("Joker app stop")
}
