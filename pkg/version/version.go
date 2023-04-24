// Copyright 2022 Innkeeper Belm(梁广庆) &lt;138521257@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/guangqingliang/blog

package version

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"runtime"
)

var (
	GitVersion = "v0.0.0-master+$Format:%h$"

	BuildDate = "1970-01-01T00:00:00Z"

	GitCommit = "$Format:%H$"

	GitTreeState = ""
)

// Info 包含了版本信息
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// string 返回人性化的版本信息字符串
func (info Info) String() string {
	s, err := info.Text()
	if err == nil {
		return string(s)
	}
	return info.GitVersion
}

func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)
	return string(s)
}

func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GoVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.Bytes(), nil
}

// Get 返回详尽的代码库版本信息，用来标明二进制文件由那个版本的代码构建
func Get() Info {
	// 一下变量通常由 -ldflags 进行设置
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitVersion,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
