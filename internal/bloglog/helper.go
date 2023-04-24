// Copyright 2022 Innkeeper Belm(梁广庆) &lt;138521257@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/guangqingliang/blog

package bloglog

import (
	"blog/internal/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	// 定义放置blog服务配置的默认目录
	recommendedHomeDir = ".blog"

	// 指定服务的默认配置文件名
	defaultConfigName = "blog.yaml"
)

// initConfig 设置需要读取的配置文件名、环境变量，并读取配置文件内容到viper中
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找用户主目录
		home, err := os.UserHomeDir()
		// 如果获取用户主目录失败，打印错误日志并退出
		cobra.CheckErr(err)

		// 用 `$HOME/<recommendedHomeDir>` 目录加入到配置文件的搜索路径中
		viper.AddConfigPath(filepath.Join(home) + "/GolandProjects/blog/configs/")

		// 把目前目录加入到配置文件的搜索路径中
		viper.AddConfigPath(".")

		// 设置配置文件格式为YAML(YAML 格式清晰易读，并且支持复杂的配置结果)
		viper.SetConfigType("yaml")

		// 配置文件名称(没有文件拓展名)
		viper.SetConfigName(defaultConfigName)
	}

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为blog，将自动转变为大写
	viper.SetEnvPrefix("BLOG")

	// 以下2行，将viper.Get(key) key 字符串中'.' 和 '-'替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("failed to read viper configuration file", "err", err)
	}
	// 打印 viper 当前使用的配置文件 方便debug
	log.Infow("Using config file", "file", viper.ConfigFileUsed())
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
