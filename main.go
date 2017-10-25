package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/guregodevo/discomotionslack/misc"
	"github.com/guregodevo/discomotionslack/server"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	showVersion = flag.Bool("version", false, "print version string")
	configFile  = flag.String("config", "dev.config.yaml", "Config file to use")

	hostname = "unknown"

	// set by Makefile
	Version   = "unknown"
	Build     = "unknown"
	BuildTime = "unknown"
)

func main() {
	if *showVersion {
		fmt.Printf("discomotion v%s (%s) (%s)\n", Version, Build, BuildTime)
		return
	}

	flag.Parse()

	api := slack.New("xoxp-158375954754-260390656065-261049069733-553e910bf231963d14cda02cd8a8cee5")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)

	var mutex = &sync.Mutex{}

	conf, err := misc.LoadConf(*configFile)
	if err != nil {
		log.WithField("error", err).Error("Failed to load config")
		os.Exit(-1)
	}

	setupLogging(&conf.Log)

	log.WithFields(log.Fields{"version": Version, "build": Build, "built": BuildTime}).Info("Starting up discomotion")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	if host, err := os.Hostname(); err == nil {
		hostname = host
	}

	s := startHTTPServer(&conf.Http, termChan, api, mutex, conf.CoreURL)

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:

		case <-termChan:
			shutdown(ticker, s)
		}
	}
}

func setupLogging(conf *misc.Log) {
	lvl, err := log.ParseLevel(conf.Level)
	if err == nil {
		log.SetLevel(lvl)
	}

	if conf.Json {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	var lj *lumberjack.Logger
	if conf.Filename != "" {
		lj = &lumberjack.Logger{
			Filename:   conf.Filename,
			MaxSize:    conf.MaxSizeMB,
			MaxBackups: conf.MaxBackups,
			MaxAge:     conf.MaxAgeDays,
		}
	}

	if conf.WriteStdout {
		if lj == nil {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(io.MultiWriter(os.Stdout, lj))
		}
	} else {
		if lj == nil {
			log.SetOutput(os.Stdout)

			log.Warn("Logging filename cannot be empty, using stdout")
		} else {
			log.SetOutput(lj)
		}
	}
}

func startHTTPServer(conf *misc.Http, termChan chan os.Signal, api *slack.Client, mutex *sync.Mutex, url string) *server.Server {

	s := &server.Server{
		Info: &server.ServerInfo{
			Server:    "discomotion",
			Version:   Version,
			Build:     Build,
			BuildTime: BuildTime,
			Hostname:  hostname,
		},
		Uptime:  time.Now(),
		Api:     api,
		Mutex:   mutex,
		BaseURL: url,
	}

	go func() {
		if err := s.Run(conf); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err).Error("Failed to start HTTP server")
		}

		termChan <- syscall.SIGTERM
	}()

	return s
}

func shutdown(ticker *time.Ticker, s *server.Server) {

	if ticker != nil {
		ticker.Stop()
	}

	if s != nil {
		s.Close()
	}

	log.Info("Shutdown complete.")
	os.Exit(0)
}
