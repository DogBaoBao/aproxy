package main

import (
	"aproxy/pkg"
	"aproxy/pkg/logger"
	"aproxy/pkg/proxy"
	"github.com/urfave/cli"
	"runtime"
)

var (
	flagToLogLevel = map[string]string{
		"trace":    "TRACE",
		"debug":    "DEBUG",
		"info":     "INFO",
		"warning":  "WARN",
		"error":    "ERROR",
		"critical": "FATAL",
	}

	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start dubbogo proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "DUBBOGO_PROXY_CONFIG",
				Value:  "configs/proxy_config.json",
			},
			cli.StringFlag{
				Name:   "log-level, l",
				Usage:  "dubbogo proxy log level, trace|debug|info|warning|error|critical",
				EnvVar: "LOG_LEVEL",
			},
			cli.StringFlag{
				Name:  "log-format, lf",
				Usage: "dubbogo proxy log format, currently useless",
			},
			cli.StringFlag{
				Name:  "limit-cpus, limc",
				Usage: "dubbogo proxy schedule threads count",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			flagLogLevel := c.String("log-level")

			bootstrap := pkg.Load(configPath)
			if logLevel, ok := flagToLogLevel[flagLogLevel]; ok {
				logger.SetLoggerLevel(logLevel)
			}

			limitCpus := c.Int("limit-cpus")
			if limitCpus <= 0 {
				runtime.GOMAXPROCS(runtime.NumCPU())
			} else {
				runtime.GOMAXPROCS(limitCpus)
			}

			proxy.Start(bootstrap)
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop dubbogo proxy",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	cmdReload = cli.Command{
		Name:  "reload",
		Usage: "reconfiguration",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)
