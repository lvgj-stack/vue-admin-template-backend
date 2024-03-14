package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/Mr-LvGJ/jota/log"

	"github.com/Mr-LvGJ/gobase/cmd/gobase/bootstrap"
	"github.com/Mr-LvGJ/gobase/pkg/common/setting"
	"github.com/Mr-LvGJ/gobase/pkg/common/token"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

var (
	configPath  = flag.String("config", "/etc/goserver/gobase.yaml", "The gobase config file default user home config")
	pprofAddr   = flag.String("pprof-addr", "", "The pprof addr")
	showVersion = flag.Bool("version", false, "")
)

func main() {
	ctx := context.Background()
	flag.Parse()
	if pprofAddr != nil && len(*pprofAddr) > 0 {
		go func() {
			if err := http.ListenAndServe(*pprofAddr, nil); err != nil {
				log.Info(ctx, "fail to start pprof", "err:", err)
			}
		}()
	}
	setting.InitConfig(*configPath)
	if err := log.NewGlobal(setting.C().Log); err != nil {
		panic(err)
	}
	token.Init(setting.C().Jwt.Key, setting.C().Jwt.IdentityKey)
	if _, err := store.Setup(ctx); err != nil {
		log.Error(ctx, "database init fail", "err", err)
	}
	bootstrap.Run(ctx)
}
