// Copyright 2022 Innkeeper Belm(梁广庆) &lt;138521257@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/guangqingliang/blog

package bloglog

import (
	"blog/internal/pkg/log"
	"blog/pkg/version/verflag"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var cfgFile string

func NewBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 制定命名的名字，该名字会出现在帮助信息中
		Use: "blog",
		// 命令的简短描述
		Short: "A good Go practical project",
		// 命令的详细描述
		Long: `A good Go practical project, used to create user with basic information`,
		// 命令出错时, 不打印帮助信息。
		SilenceUsage: true,
		// 指定调用cmd.Execute()时，执行的Run函数,函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果 `--version = true`,则打印版本并退出
			verflag.PrintAndExitIfRequested()
			// 初始化日志
			log.Init(logOptions())
			// 将缓存中的日志异步刷新到磁盘文件中
			defer log.Sync()

			return run()
		},
		// 这里设置命令运行时,不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					log.Errorw("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
	// 以下配置使initConfig函数在每个命令行运行时都会调用以读取配置
	cobra.OnInitialize(initConfig)
	// cobra 支持持久性表示，该标志可以用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the blog configuration file. Empty string for no configuration file.")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 添加--version标志
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

// run函数实际是业务代码入口函数
func run() error {
	// 设置gin模式
	gin.SetMode(viper.GetString("runmode"))
	// 创建gin引擎
	g := gin.New()

	// 注册404 handler
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 10003,
			"message": "Page not found"})
	})

	// 注册/health handler
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// 创建http server实例
	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}
	// 运行http服务器
	// 打印一条日志 用来提示服务已经起来，方便排除障碍
	log.Infow("start to listening the incoming requests on http address", viper.GetString("addr"))
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw(err.Error())
	}
	return nil
}
