package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>

const EGLContext kNoContext = EGL_NO_CONTEXT;
*/
import "C"

var noContext C.EGLContext = C.kNoContext

type Context struct {
	eglContext C.EGLContext
	Display *Display
}

func destroyContext(context *Context) {
//	fmt.Printf("destroying context == %v\n", context)
	context.Destroy()
}

func (context *Context) Destroy() error {
	success := C.eglDestroyContext(context.Display.eglDisplay, context.eglContext)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

func (context *Context) MakeCurrent(draw *Surface, read *Surface) error {
	var eglDraw C.EGLSurface
	if draw == nil {
		eglDraw = noSurface
	} else {
		eglDraw = draw.eglSurface
	}

	var eglRead C.EGLSurface
	if read == nil {
		eglRead = noSurface
	} else {
		eglRead = read.eglSurface
	}

	success := C.eglMakeCurrent(context.Display.eglDisplay, eglDraw, eglRead, context.eglContext)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}
