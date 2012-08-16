package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>
*/
import "C"

import (
	"errors"
	"fmt"
)

func getError() error {
	errorCode := C.eglGetError()
	switch errorCode {
		case C.EGL_SUCCESS:
			return errors.New("EGL succeeded.")
		case C.EGL_NOT_INITIALIZED:
			return errors.New("EGL is not initialized, or could not be initialized, for the specified display.")
		case C.EGL_BAD_ACCESS:
			return errors.New("EGL cannot access a requested resource.")
		case C.EGL_BAD_ALLOC:
			return errors.New("EGL failed to allocate resources for the requested operation.")
		case C.EGL_BAD_ATTRIBUTE:
			return errors.New("An unrecognized attribute or attribute value was passed in an attribute list.")
		case C.EGL_BAD_CONTEXT:
			return errors.New("An EGLContext argument does not name a valid EGLContext.")
		case C.EGL_BAD_CONFIG:
			return errors.New("An EGLConfig argument does not name a valid EGLConfig.")
		case C.EGL_BAD_CURRENT_SURFACE:
			return errors.New("The current surface of the calling thread is a window, pbuffer, or pixmap that is no longer valid.")
		case C.EGL_BAD_DISPLAY:
			return errors.New("An EGLDisplay argument does not name a valid EGLDisplay.")
		case C.EGL_BAD_SURFACE:
			return errors.New("An EGLSurface argument does not name a valid surface (window, pbuffer, or pixmap) configured for rendering.")
		case C.EGL_BAD_MATCH:
			return errors.New("Arguments are inconsistent; for example, an otherwise valid context requires buffers (e.g. depth or stencil) not allocated by an otherwise valid surface.")
		case C.EGL_BAD_PARAMETER:
			return errors.New("One or more argument values are invalid.")
		case C.EGL_BAD_NATIVE_PIXMAP:
			return errors.New("An EGLNativePixmapType argument does not refer to a valid native pixmap.")
		case C.EGL_BAD_NATIVE_WINDOW:
			return errors.New("An EGLNativeWindowType argument does not refer to a valid native window.")
		case C.EGL_CONTEXT_LOST:
			return errors.New("A power management event has occurred.")
	}
	return fmt.Errorf("EGL error code %d.", errorCode)
}
