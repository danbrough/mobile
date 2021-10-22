// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"io"
	"io/ioutil"
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


	pkgName := pkgs[0].Name

	if buildO == "" {
		buildO = "."
	}
	buildO, _ = filepath.Abs(buildO)

	println("buildO:", buildO, "pkgName:", pkgName)


	// Run gobind to generate the bindings
	cmd := exec.Command(
		gobind,
		"-lang=go,java",
		"-outdir="+tmpdir,
	)

	var javaHome *string

	for _, s := range os.Environ() {
		//println("ENV:",s)
		if strings.HasPrefix(s, "JAVA_HOME") {
			println("FOUND ", s)
			javaHome = &strings.Split(s, "=")[1]
			println("FOUND JAVA_HOME:", *javaHome)
			//cmd.Env = append(cmd.Env, "CGO_CFLAGS=-I"+*javaHome+"/include -I"+*javaHome+"/include/linux")
			//cmd.Env = append(cmd.Env, "CFLAGS=-I"+*javaHome+"/include -I"+*javaHome+"/include/linux")
		} else if strings.HasPrefix(s, "CFLAGS") {
			cmd.Env = append(cmd.Env, s)
		} else if strings.HasPrefix(s, "CGO_CFLAGS") {
			println("ADDING", s)
			cmd.Env = append(cmd.Env, s)
		} else if strings.HasPrefix(s, "CGO_LDFLAGS") {
			println("ADDING", s)
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

	//linuxDir := filepath.Join(tmpdir, "linux")

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

		//toolchain := ndk.Toolchain(t.arch)
		err := goBuildAt(
			filepath.Join(tmpdir, "src"),
			"./gobind",
			env,
			"-buildmode=c-shared",
			"-o="+filepath.Join(buildO, "libs",t.arch,"libgojni.so"),
		)
		if err != nil {
			return err
		}
	}


	jsrc := filepath.Join(tmpdir, "java")
	/*if err := buildAAR(jsrc, linuxDir, pkgs, targets); err != nil {
		return err
	}*/

	/*	jarw := zip.NewWriter(out)
		jarwcreate := func(name string) (io.Writer, error) {
			if buildV {
				fmt.Fprintf(os.Stderr, "aar: %s\n", name)
			}
			return jarw.Create(name)
		}*/
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
	println("buildLinuxJar() srcDir:",srcDir)
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
	//	"-target", javacTargetVer,
	//	"-bootclasspath", bClspath,
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
