package main

import (
	"os"
	"path"

	microplugin "github.com/jakexks/netdata-collector/micro"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/cli"
	"github.com/netdata/go-orchestrator/pkg/multipath"
)

var (
	cd, _         = os.Getwd()
	netdataConfig = multipath.New(
		os.Getenv("NETDATA_USER_CONFIG_DIR"),
		os.Getenv("NETDATA_STOCK_CONFIG_DIR"),
		path.Join(cd, "/../../../../etc/netdata"),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d"),
	)
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.debug.collector"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	microplugin.New().WithClient(service.Client()).Register()

	netdata := orchestrator.New()
	netdata.Name = "micro.d"
	netdata.Option = &cli.Option{
		UpdateEvery: 1,
		Debug:       true,
		Module:      "all",
		ConfigDir:   netdataConfig,
		Version:     false,
	}
	netdata.ConfigPath = netdataConfig

	if !netdata.Setup() {
		log.Fatal("Netdata failed to Setup()")
	}
	go netdata.Serve()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
