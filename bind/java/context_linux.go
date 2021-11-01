// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package java // import "github.com/danbrough/mobile/bind/java"

//
//#include <jni.h>
import "C"

import (
	"unsafe"
<<<<<<< HEAD:bind/java/context_linux.go
	"github.com/danbrough/internal/mobileinit"
=======

	"github.com/danbrough/mobile/internal/mobileinit"
>>>>>>> dev:bind/java/context_android.go
)

//export setContext
func setContext(vm *C.JavaVM, ctx C.jobject) {
	printf("%s:%d setContext()\n",__LINE__,__FILE__);
	mobileinit.SetCurrentContext(unsafe.Pointer(vm), uintptr(ctx))
}
