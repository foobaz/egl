package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>
*/
import "C"

type Config C.EGLConfig
type Attrib C.EGLint
type NativeDisplay C.EGLNativeDisplayType
type NativePixmap C.EGLNativePixmapType

func WaitClient() error {
	success := C.eglWaitClient()
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

func BindAPI(api int) error {
	success := C.eglBindAPI(C.EGLenum(api))
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

func QueryAPI() int {
	return int(C.eglQueryAPI())
}
