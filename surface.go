package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>

const EGLSurface kNoSurface = EGL_NO_SURFACE;
*/
import "C"

var noSurface C.EGLSurface = C.kNoSurface

func destroySurface(surface *Surface) {
	surface.Destroy()
}

func (surface *Surface) Query(name Attrib) (Attrib, error) {
	var value Attrib
	success := C.eglQuerySurface(surface.Display.eglDisplay, surface.eglSurface, C.EGLint(name), (*C.EGLint)(&value))
	if success == C.EGL_FALSE {
		return None, getError()
	}

	return value, nil
}

func (surface *Surface) SwapBuffers() error {
	success := C.eglSwapBuffers(surface.Display.eglDisplay, surface.eglSurface)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

