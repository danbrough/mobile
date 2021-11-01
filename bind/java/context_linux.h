// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <jni.h>
<<<<<<< HEAD:bind/java/context_linux.c

#include "seq_linux.h"

=======
#include "seq_linux.h"
>>>>>>> dev:bind/java/context_linux.h
#include "_cgo_export.h"

JNIEXPORT void JNICALL
Java_go_Seq_setContext(JNIEnv* env, jclass clazz, jobject ctx) {
    printf("Java_go_Seq_setContext()\n");
	JavaVM* vm;

		if ((*env)->GetJavaVM(env, &vm) != 0) {
		#ifdef __GOBIND_ANDROID__
    		LOG_FATAL("failed to get JavaVM");
    		#else
    		printf("Failed to get JavaVM\n");
    		#endif
    	}
	setContext(vm, (*env)->NewGlobalRef(env, ctx));
}
