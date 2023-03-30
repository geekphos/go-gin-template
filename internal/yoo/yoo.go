package yoo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"phos.cc/yoo/internal/pkg/known"
	"phos.cc/yoo/internal/pkg/log"
	"phos.cc/yoo/internal/pkg/middleware"
	"phos.cc/yoo/internal/pkg/validator"
	"phos.cc/yoo/pkg/token"
	"phos.cc/yoo/pkg/version/verflag"
)

var cfgFile string

func NewYooCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "yoo",
		Short:        "yoo server",
		Long:         "yoo server for front end",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			log.Init(logOptions())
			defer log.Sync()

			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/yoo.yaml)")

	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io:8080
// @BasePath /v1
func run() error {

	// 初始化数据库
	if err := initStore(); err != nil {
		return err
	}

	// 初始化翻译器
	if err := validator.InitTrans("zh"); err != nil {
		return err
	}

	// 初始化 JWT
	token.Init(viper.GetString("jwt-secret"), known.XEmailKey)

	gin.SetMode(viper.GetString("runmode"))

	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.RequestID())

	if err := installRouters(r); err != nil {
		return err
	}

	httpsrv := &http.Server{Addr: ":8080", Handler: r}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))

	go func() {

		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}

	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Infow("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Fatalw("Server Shutdown:", err)
		return err
	}

	log.Infow("Server exiting")

	return nil
}
