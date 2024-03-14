package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Mr-LvGJ/jota/log"

	"github.com/Mr-LvGJ/gobase/pkg/common/setting"
	"github.com/Mr-LvGJ/gobase/pkg/gobase"
)

func Run(ctx context.Context) error {
	gin.SetMode(setting.C().RunMode)
	g := gin.New()
	gobase.LoadRouter(g)
	insecureServer := &http.Server{
		Addr:    setting.C().Addr,
		Handler: g,
	}
	log.Info(ctx, "Start to listening the incoming request on http address: %s", "addr", setting.C().Addr)
	go func() {
		if err := insecureServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error(ctx, "Listen: %s\n", "err", err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		if err := pingServer(ctx); err != nil {
			log.Error(ctx, err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, "Shutting down server...")

	if err := insecureServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Insecure server fored to shutdown")
		return err
	}
	return nil

}

func pingServer(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", setting.C().Addr)
	bind := strings.Split(setting.C().Addr, ":")[0]
	if bind == "" || bind == "0.0.0.0" {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(setting.C().Addr, ":")[1])
	}
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		log.Info(ctx, "Wait for router, retry in 1 second.")
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			return fmt.Errorf("can not ping server within the specified time interval")
		default:
		}
	}
}
