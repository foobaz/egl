package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>

// These variables are necessary because EGL_DEFAULT_DISPLAY and EGL_NO_DISPLAY
// are pointer constants, and cgo doesn't translate them correctly.
const EGLNativeDisplayType kDefaultDisplay = EGL_DEFAULT_DISPLAY;
const EGLDisplay kNoDisplay = EGL_NO_DISPLAY;
const EGLSurface kNoSurface = EGL_NO_SURFACE;
const EGLContext kNoContext = EGL_NO_CONTEXT;
*/
import "C"

import (
	"fmt"
	"runtime"
)

type Surface struct {
	eglSurface C.EGLSurface
	display *Display
}

type Context struct {
	eglContext C.EGLContext
	display *Display
}

type Config C.EGLConfig
type Attrib C.EGLint
type NativeDisplay C.EGLNativeDisplayType
type NativePixmap C.EGLNativePixmapType

func GetDisplay() (*Display, error) {
	display := new(Display)
	runtime.SetFinalizer(display, terminateDisplay)

	display.eglDisplay = C.eglGetDisplay(C.kDefaultDisplay)
	return display, nil
}

func terminateDisplay(display *Display) {
	fmt.Printf("terminating display == %v\n", display)
	display.Terminate()
}

func (display *Display) Initialize() error {
	if display.eglDisplay == C.kNoDisplay {
		return getError()
	}

	var major, minor C.EGLint
	success := C.eglInitialize(display.eglDisplay, &major, &minor)
	if success == C.EGL_FALSE {
		return getError()
	}
	display.majorVersion = int(major)
	display.minorVersion = int(minor)

	return nil
}

func GetVersion(display *Display) (major, minor int) {
	return display.majorVersion, display.minorVersion
}

func (display *Display) Terminate() error {
	if display.xDisplay != nil {
		C.XCloseDisplay(display.xDisplay)
	}

	if display.eglDisplay != nil {
		success := C.eglTerminate(C.EGLDisplay(display.eglDisplay))
		if success == C.EGL_FALSE {
			return getError()
		}
	}

	return nil
}

func (display *Display) QueryString(name int) (string, error) {
	cString := C.eglQueryString(display.eglDisplay, C.EGLint(name))
	if cString == nil {
		return "", getError()
	}
	return C.GoString(cString), nil
}

func (display *Display) GetConfigs() ([]Config, error) {
	var configCount C.EGLint
	success := C.eglGetConfigs(display.eglDisplay, nil, 0, &configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}

	configurations := make([]Config, configCount)
	success = C.eglGetConfigs(
		display.eglDisplay,
		(*C.EGLConfig)(&(configurations[0])),
		configCount,
		&configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}

	return configurations, nil
}

func (display *Display) ChooseConfig(attribList []Attrib) ([]Config, error) {
	var configCount C.EGLint
	success := C.eglChooseConfig(
		display.eglDisplay,
		(*C.EGLint)(&(attribList[0])),
		nil,
		0,
		&configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}

	configurations := make([]Config, configCount)
	if configCount <= 0 {
		return configurations, nil
	}

	success = C.eglChooseConfig(
		display.eglDisplay,
		(*C.EGLint)(&(attribList[0])),
		(*C.EGLConfig)(&(configurations[0])),
		configCount,
		&configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}

	return configurations, nil
}

func (display *Display) GetConfigAttrib(config Config, name Attrib) (Attrib, error) {
	var value Attrib
	success := C.eglGetConfigAttrib(display.eglDisplay, C.EGLConfig(config), C.EGLint(name), (*C.EGLint)(&value))
	if success == C.EGL_FALSE {
		return None, getError()
	}

	return value, nil
}

func (display *Display) CreatePbufferSurface(config Config, attribList []Attrib) (*Surface, error) {
	eglSurface := C.eglCreatePbufferSurface(display.eglDisplay, C.EGLConfig(config), (*C.EGLint)(&(attribList[0])))
	if eglSurface == C.kNoSurface {
		return nil, getError()
	}

	surface := new(Surface)
	runtime.SetFinalizer(surface, destroySurface)
	surface.display = display
	surface.eglSurface = eglSurface
	return surface, nil
}

func destroySurface(surface *Surface) {
	fmt.Printf("destroying surface == %v\n", surface)
	surface.Destroy()
}

func (surface *Surface) Destroy() error {
	success := C.eglDestroySurface(surface.display.eglDisplay, surface.eglSurface)
	if success == C.EGL_FALSE {
		return getError()
	}

	return nil
}

func (surface *Surface) Query(name Attrib) (Attrib, error) {
	var value Attrib
	success := C.eglQuerySurface(surface.display.eglDisplay, surface.eglSurface, C.EGLint(name), (*C.EGLint)(&value))
	if success == C.EGL_FALSE {
		return None, getError()
	}

	return value, nil
}

func (surface *Surface) SwapBuffers() {
	C.eglSwapBuffers(surface.display.eglDisplay, surface.eglSurface)
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

func (display *Display) CreateContext(config Config, shareContext *Context, attribList []Attrib) (*Context, error) {
	var eglShareContext C.EGLContext
	if shareContext == nil {
		eglShareContext = C.kNoContext
	} else {
		eglShareContext = shareContext.eglContext
	}

	var eglAttribs *C.EGLint
	if attribList == nil {
		eglAttribs = nil
	} else {
		eglAttribs = (*C.EGLint)(&(attribList[0]))
	}

	eglContext := C.eglCreateContext(display.eglDisplay, C.EGLConfig(config), eglShareContext, eglAttribs)
	if eglContext == C.kNoContext {
		return nil, getError()
	}

	context := new(Context)
	runtime.SetFinalizer(context, destroyContext)
	context.eglContext = eglContext
	context.display = display
	return context, nil
}

func destroyContext(context *Context) {
//	fmt.Printf("destroying context == %v\n", context)
	context.Destroy()
}

func (context *Context) Destroy() error {
	success := C.eglDestroyContext(context.display.eglDisplay, context.eglContext)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

func (context *Context) MakeCurrent(draw *Surface, read *Surface) error {
	success := C.eglMakeCurrent(context.display.eglDisplay, draw.eglSurface, read.eglSurface, context.eglContext)
	fmt.Printf("draw == %v, read == %v, success == %d\n", draw, read, success)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

func WaitClient() error {
	success := C.eglWaitClient()
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}
