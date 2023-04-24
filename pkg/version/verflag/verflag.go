// Copyright 2022 Innkeeper Belm(梁广庆) &lt;138521257@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/guangqingliang/blog

package verflag

import (
	"blog/pkg/version"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strconv"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const (
	strRawVersion   = "raw"
	versionFlagName = "version"
)

var versionFlag = Version(versionFlagName, VersionFalse, "print version information and quit")

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw

		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

func (v *versionValue) Type() string {
	return "version"
}

func VersionVal(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	pflag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVal(p, name, value, usage)

	return p
}

func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}
