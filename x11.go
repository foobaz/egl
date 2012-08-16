package egl

/*
#cgo LDFLAGS: -lEGL -lX11 -lGLEW

#include <GL/glew.h>
#include <EGL/egl.h>
#include <stdlib.h>
#include <strings.h>
#include <X11/Xlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Display struct {
	xDisplay *C.Display
	eglDisplay C.EGLDisplay
	majorVersion, minorVersion int
}

func Init() {
//	C.glewExperimental = C.GL_TRUE
//	gl.Init()
}

func GetDisplayWithX11Name(name string) (*Display, error) {
	display := new(Display)
	runtime.SetFinalizer(display, terminateDisplay)

	displayName := C.CString(name)
	display.xDisplay = C.XOpenDisplay(displayName)
	C.free(unsafe.Pointer(displayName))
	if display.xDisplay == nil {
		return nil, getError()
	}

	display.eglDisplay = C.eglGetDisplay(C.EGLNativeDisplayType(display.xDisplay))
	if display.eglDisplay == nil {
		return nil, getError()
	}

	return display, nil
}

/*
func (display *Display) FlushX11() {
	C.XFlush(display.xDisplay)
}
*/

func (surface *Surface) CopyBuffers() (*RGBAImage, error) {
	width, widthErr := surface.Query(Width)
	if widthErr != nil {
		return nil, widthErr
	}

	height, heightErr := surface.Query(Height)
	if heightErr != nil {
		return nil, heightErr
	}

	xDisplay := surface.display.xDisplay
	pixmap := C.XCreatePixmap(xDisplay, C.Drawable(C.XDefaultRootWindow(xDisplay)), C.uint(width), C.uint(height), 32)

	C.eglCopyBuffers(surface.display.eglDisplay, surface.eglSurface, C.EGLNativePixmapType(pixmap))
	xImage := C.XGetImage(xDisplay, C.Drawable(pixmap), 0, 0, C.uint(width), C.uint(height), 0, C.ZPixmap)
//	xImage := C.XGetImage(xDisplay, pixmap, 0, 0, width, height, C.AllPlanes, C.XYPixmap)

	image := new(RGBAImage)
	length := 4 * xImage.width * xImage.height
	image.pixels = make([]byte, length)
	C.bcopy(unsafe.Pointer(xImage.data), unsafe.Pointer(&image.pixels[0]), C.size_t(length))
	image.width = int(xImage.width)
	image.height = int(xImage.height)
	image.rowBytes = int(xImage.bytes_per_line)

	C.XFreePixmap(xDisplay, pixmap)

	return image, nil
}
