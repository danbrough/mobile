// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func goLinuxBind(gobind string, pkgs []*packages.Package, targets []targetInfo) error {
	if sdkDir := os.Getenv("ANDROID_HOME"); sdkDir == "" {
		return fmt.Errorf("this command requires ANDROID_HOME environment variable (path to the Android SDK)")
	}

	// Run gobind to generate the bindings
	cmd := exec.Command(
		gobind,
		"-lang=go,java",
		"-outdir="+tmpdir,
	)
	cmd.Env = append(cmd.Env, "GOOS=android")
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

	androidDir := filepath.Join(tmpdir, "android")

	modulesUsed, err := areGoModulesUsed()
	if err != nil {
		return err
	}

	// Generate binding code and java source code only when processing the first package.
	for _, t := range targets {
		if err := writeGoMod(tmpdir, "android", t.arch); err != nil {
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
			"-o="+filepath.Join(androidDir, "src/main/jniLibs/"+toolchain.abi+"/libgojni.so"),
		)
		if err != nil {
			return err
		}
	}

	jsrc := filepath.Join(tmpdir, "java")
	if err := buildAAR(jsrc, androidDir, pkgs, targets); err != nil {
		return err
	}
	return buildSrcJar(jsrc)
}
