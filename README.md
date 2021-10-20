# Go support for Mobile devices

## Hacked to run on jvm.

This fork of [GoMobile](https://pkg.go.dev/golang.org/x/mobile) has been modified to enable libraries to be built for the JVM  as well as android.

It's currently working on linux (amd64 but arm64, arm, x86 should work also).

I'm not a mac developer but it shouldn't too hard to get this running on darwin,
if the original project doesn't do that already.

To get started have a look at the [demos/basic](demos/basic/README.md) demonstration.

Feel free to make any contributionsS.


[![Go Reference](https://pkg.go.dev/badge/golang.org/x/mobile.svg)](https://pkg.go.dev/golang.org/x/mobile)

The Go mobile repository holds packages and build tools for using Go on mobile platforms.

Package documentation as a starting point:

- [Building all-Go apps](https://golang.org/x/mobile/app)
- [Building libraries for SDK apps](https://golang.org/x/mobile/cmd/gobind)

![Caution image](doc/caution.png)

The Go Mobile project is experimental. Use this at your own risk.
While we are working hard to improve it, neither Google nor the Go
team can provide end-user support.

This is early work and installing the build system requires Go 1.5.
Follow the instructions on
[golang.org/wiki/Mobile](https://golang.org/wiki/Mobile)
to install the gomobile command, build the
[basic](https://golang.org/x/mobile/example/basic)
and the [bind](https://golang.org/x/mobile/example/bind) example apps.

--

Contributions to Go are appreciated. See https://golang.org/doc/contribute.html.

* Bugs can be filed at the [Go issue tracker](https://golang.org/issue/new?title=x/mobile:+).
* Feature requests should preliminary be discussed on
[golang-nuts](https://groups.google.com/forum/#!forum/golang-nuts)
mailing list.
