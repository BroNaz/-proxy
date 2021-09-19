package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	errConfLog "log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/BroNaz/proxy/internal/config"
	"github.com/BroNaz/proxy/internal/logger"
	"github.com/BroNaz/proxy/internal/proxy"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

var (
	// versionInfo populated at compile time
	versionInfo string

	configPath  = flag.String("config", "", "path to config file")
	showVersion = flag.Bool("version", false, "display the proxy version and exit")
)

func init() {
	expvar.NewString("startedAt").Set(time.Now().Format(time.RFC3339Nano))
	expvar.NewString("versionInfo").Set(versionInfo)
	expvar.Publish("goroutines", expvar.Func(goroutines))

}
func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(versionInfo)
		return
	}
	if *configPath == "" {
		errConfLog.Fatalln("Usage: proxy --config=<path_to_config>")
		return
	}

	var conf config.TomlConfig
	if _, err := toml.DecodeFile(*configPath, &conf); err != nil {
		errConfLog.Fatalln(err.Error())
		return
	}

	// setup the global logger
	logger.SetupLogging(&conf.Log)

	log.Info().Msg("Proxy service started")

	server := proxy.NewServer(conf.Settings, conf.Addr)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatal().Err(err).Msg("server ended")
		}
	}()

	// When receiving signals, we give the current requests a few seconds to complete.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to gently stop the server")
	}
}
