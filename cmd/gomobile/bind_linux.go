// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || windows || android
// +build linux windows android

package main

import (
	"fmt"
	"github.com/danbrough/mobile/klog"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func goWindowsBind(gobind string, pkgs []*packages.Package, targets []targetInfo) error {
	klog.KLog.Info("goWindowsBind() gobind:%s", gobind)
	var jdkDir string
	if jdkDir = os.Getenv("JAVA_HOME"); jdkDir == "" {
		return fmt.Errorf("this command requires JAVA_HOME environment variable (path to the Java SDK)")
	}

	// Run gobind to generate the bindings
	cmd := exec.Command(
		gobind,
		"-lang=go,java",
		"-outdir="+tmpdir,
	)
	cmd.Env = append(cmd.Env, "GOOS=windows")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	cmd.Env = append(cmd.Env, "CC=/usr/bin/x86_64-w64-mingw32-cc")
	cmd.Env = append(cmd.Env, "CXX=/usr/bin/x86_64-w64-mingw32-c++")

	jdkIncludes := " -I" + filepath.Join(jdkDir, "include")
	w32Includes := filepath.Join(jdkDir, "include", "win32")
	if _, err := os.Stat(w32Includes); err == nil {
		jdkIncludes = jdkIncludes + " -I" + w32Includes
	} else {
		jdkIncludes = jdkIncludes + " -I" + filepath.Join(jdkIncludes,"linux")
	}
	cmd.Env = append(cmd.Env, "CGO_CFLAGS="+os.Getenv("CGO_CFLAGS")+jdkIncludes)
	cmd.Env = append(cmd.Env, "CGO_LDFLAGS=-static -fPIC "+os.Getenv("CGO_LDFLAGS"))

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

	buildDir, _ := filepath.Abs(buildO)
	pkgName := pkgs[0].Name
	modulesUsed, err := areGoModulesUsed()

	// Generate binding code and java source code only when processing the first package.
	for _, t := range targets {
		if err := writeGoMod(tmpdir, "linux", t.arch); err != nil {
			return err
		}

		//env := androidEnv[t.arch]

		// Add the generated packages to GOPATH for reverse bindings.
		gopath := fmt.Sprintf("GOPATH=%s%c%s", tmpdir, filepath.ListSeparator, goEnv("GOPATH"))
		cmd.Env = append(cmd.Env, gopath)

		// Run `go mod tidy` to force to create go.sum.
		// Without go.sum, `go build` fails as of Go 1.16.
		if modulesUsed {
			if err := goModTidyAt(filepath.Join(tmpdir, "src"), cmd.Env); err != nil {
				return err
			}
		}

		//toolchain := ndk.Toolchain(t.arch)
		klog.KLog.Warn("calling goBuildAt()")
		err := goBuildAt(
			filepath.Join(tmpdir, "src"),
			"./gobind",
			cmd.Env,
			"-buildmode=c-shared",
			"-o="+filepath.Join(buildDir, "libs", t.arch, "libgojni.dll"),
		)
		if err != nil {
			return err
		}
	}

	jsrc := filepath.Join(tmpdir, "java")

	err = buildLinuxSrcJar(jsrc, filepath.Join(buildO, pkgName+"-sources.jar"))
	if err != nil {
		return err
	}

	var out io.Writer = ioutil.Discard
	if !buildN {
		f, err := os.Create(filepath.Join(buildO, pkgName+".jar"))
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		out = f
	}
	return buildLinuxJar(out, jsrc)

}

func goLinuxBind(gobind string, pkgs []*packages.Package, targets []targetInfo) error {
	klog.KLog.Info("goLinuxBind() gobind:%s", gobind)
	var jdkDir string
	if jdkDir = os.Getenv("JAVA_HOME"); jdkDir == "" {
		return fmt.Errorf("this command requires JAVA_HOME environment variable (path to the Java SDK)")
	}

	// Run gobind to generate the bindings
	cmd := exec.Command(
		gobind,
		"-lang=go,java",
		"-outdir="+tmpdir,
	)
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	cmd.Env = append(cmd.Env, "CGO_CFLAGS="+os.Getenv("CGO_CFLAGS")+" -I"+filepath.Join(jdkDir, "include")+" -I"+filepath.Join(jdkDir, "include", "linux"))
	cmd.Env = append(cmd.Env, "CGO_LDFLAGS=-fPIC "+os.Getenv("CGO_LDFLAGS"))

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

	buildDir, _ := filepath.Abs(buildO)
	pkgName := pkgs[0].Name
	modulesUsed, err := areGoModulesUsed()

	// Generate binding code and java source code only when processing the first package.
	for _, t := range targets {
		if err := writeGoMod(tmpdir, "linux", t.arch); err != nil {
			return err
		}

		//env := androidEnv[t.arch]

		// Add the generated packages to GOPATH for reverse bindings.
		gopath := fmt.Sprintf("GOPATH=%s%c%s", tmpdir, filepath.ListSeparator, goEnv("GOPATH"))
		cmd.Env = append(cmd.Env, gopath)

		// Run `go mod tidy` to force to create go.sum.
		// Without go.sum, `go build` fails as of Go 1.16.
		if modulesUsed {
			if err := goModTidyAt(filepath.Join(tmpdir, "src"), cmd.Env); err != nil {
				return err
			}
		}

		//toolchain := ndk.Toolchain(t.arch)
		klog.KLog.Warn("calling goBuildAt()")
		err := goBuildAt(
			filepath.Join(tmpdir, "src"),
			"./gobind",
			cmd.Env,
			"-buildmode=c-shared",
			"-o="+filepath.Join(buildDir, "libs", t.arch, "libgojni.so"),
		)
		if err != nil {
			return err
		}
	}

	jsrc := filepath.Join(tmpdir, "java")

	err = buildLinuxSrcJar(jsrc, filepath.Join(buildO, pkgName+"-sources.jar"))
	if err != nil {
		return err
	}

	var out io.Writer = ioutil.Discard
	if !buildN {
		f, err := os.Create(filepath.Join(buildO, pkgName+".jar"))
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		out = f
	}
	return buildLinuxJar(out, jsrc)

}

func buildLinuxJar(w io.Writer, srcDir string) error {
	println("buildLinuxJar() srcDir:", srcDir)
	var srcFiles []string
	if buildN {
		srcFiles = []string{"*.java"}
	} else {
		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".java" {
				srcFiles = append(srcFiles, filepath.Join(".", path[len(srcDir):]))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	dst := filepath.Join(tmpdir, "javac-output")
	if !buildN {
		if err := os.MkdirAll(dst, 0700); err != nil {
			return err
		}
	}

	//bClspath := bindBootClasspath

	args := []string{
		"-d", dst,
		//"-source", javacTargetVer,
		//      "-target", javacTargetVer,
		//      "-bootclasspath", bClspath,
	}
	if bindClasspath != "" {
		args = append(args, "-classpath", bindClasspath)
	}

	args = append(args, srcFiles...)

	javac := exec.Command("javac", args...)
	javac.Dir = srcDir
	if err := runCmd(javac); err != nil {
		return err
	}

	if buildX {
		printcmd("jar c -C %s .", dst)
	}
	return writeJar(w, dst)
}

func buildLinuxSrcJar(src string, output string) error {
	println("buildLinuxSrcJar() src:", src, "output:", output)
	var out io.Writer = ioutil.Discard
	if !buildN {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		out = f
	}

	return writeJar(out, src)
}
