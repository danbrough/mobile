// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package java // import "github.com/danbrough/mobile/bind/java"

//
//#include <jni.h>
import "C"

import (
	"unsafe"
	"github.com/danbrough/internal/mobileinit"
)

//export setContext
func setContext(vm *C.JavaVM, ctx C.jobject) {
	printf("%s:%d setContext()\n",__LINE__,__FILE__);
	mobileinit.SetCurrentContext(unsafe.Pointer(vm), uintptr(ctx))
}
