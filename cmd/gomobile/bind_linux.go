// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func goLinuxBind(gobind string, pkgs []*packages.Package, targets []targetInfo) error {

	println("goLinuxBind() gobind:", gobind, "os:", runtime.GOOS, "arch:", runtime.GOARCH)
	for _, t := range targets {
		println("target: platform:", t.platform, "arch:", t.arch)
	}

	// Run gobind to generate the bindings
	cmd := exec.Command(
		gobind,
		"-lang=go,java",
		"-outdir="+tmpdir,
	)

	var javaHome *string

	for _, s := range os.Environ() {
		if strings.HasPrefix(s, "JAVA_HOME") {
			println("FOUND ", s)
			javaHome = &strings.Split(s, "=")[1]
			println("FOUND JAVA_HOME:", *javaHome)
			cmd.Env = append(cmd.Env, "CGO_CFLAGS=-I"+*javaHome+"/include -I"+*javaHome+"/include/linux")
			break
		} else if strings.HasPrefix(s, "CFLAGS") {
			cmd.Env = append(cmd.Env, s)
		}
	}
	if javaHome == nil {
		println("NO JAVA_HOME FOUND!!!!!!!!!!!!")
	}

	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	if len(buildTags) > 0 {
		cmd.Args = append(cmd.Args, "-tags="+strings.Join(buildTags, ","))
	}
	if bindJavaPkg != "" {
		cmd.Args = append(cmd.Args, "-javapkg="+bindJavaPkg)
	}
	if bindClasspath != "" {
		cmd.Args = append(cmd.Args, "-classpath="+bindClasspath)
	}
	if bindBootClasspath != "" {
		cmd.Args = append(cmd.Args, "-bootclasspath="+bindBootClasspath)
	}
	for _, p := range pkgs {
		cmd.Args = append(cmd.Args, p.PkgPath)
	}
	if err := runCmd(cmd); err != nil {
		return err
	}

	linuxDir := filepath.Join(tmpdir, "linux")

	modulesUsed, err := areGoModulesUsed()
	if err != nil {
		return err
	}

	// Generate binding code and java source code only when processing the first package.
	for _, t := range targets {
		if err := writeGoMod(tmpdir, "linux", t.arch); err != nil {
			return err
		}

		env := androidEnv[t.arch]
		// Add the generated packages to GOPATH for reverse bindings.
		gopath := fmt.Sprintf("GOPATH=%s%c%s", tmpdir, filepath.ListSeparator, goEnv("GOPATH"))
		env = append(env, gopath)

		// Run `go mod tidy` to force to create go.sum.
		// Without go.sum, `go build` fails as of Go 1.16.
		if modulesUsed {
			if err := goModTidyAt(filepath.Join(tmpdir, "src"), env); err != nil {
				return err
			}
		}

		toolchain := ndk.Toolchain(t.arch)
		err := goBuildAt(
			filepath.Join(tmpdir, "src"),
			"./gobind",
			env,
			"-buildmode=c-shared",
			"-o="+filepath.Join(linuxDir, "libs/"+toolchain.abi+"/libgojni.so"),
		)
		if err != nil {
			return err
		}
	}

	jsrc := filepath.Join(tmpdir, "java")
	/*if err := buildAAR(jsrc, linuxDir, pkgs, targets); err != nil {
		return err
	}*/
	return buildSrcJar(jsrc)
}
