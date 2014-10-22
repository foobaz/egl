package egl

/*
#cgo pkg-config: egl

#include <EGL/egl.h>

// These variables are necessary because EGL_DEFAULT_DISPLAY and EGL_NO_DISPLAY
// are pointer constants, and cgo doesn't translate them correctly.
const EGLNativeDisplayType kDefaultDisplay = EGL_DEFAULT_DISPLAY;
const EGLDisplay kNoDisplay = EGL_NO_DISPLAY;
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	//"runtime"
)

var defaultDisplay C.EGLNativeDisplayType = C.kDefaultDisplay
var noDisplay C.EGLDisplay = C.kNoDisplay

func OpenDisplay() (*Display, error) {
	display := new(Display)

	display.eglDisplay = C.eglGetDisplay(defaultDisplay)
	return display, nil
}

func (display *Display) Initialize() error {
	if display.eglDisplay == noDisplay {
		return getError()
	}

	var major, minor C.EGLint
	success := C.eglInitialize(display.eglDisplay, &major, &minor)
//fmt.Printf("display == %v, version == %d.%d, success == %d\n", display.eglDisplay, major, minor, success)
	if success == C.EGL_FALSE {
		return getError()
	}
	display.majorVersion = int(major)
	display.minorVersion = int(minor)

	return nil
}

func (display *Display) GetVersion() (major, minor int) {
	return display.majorVersion, display.minorVersion
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

	if configCount <= 0 {
		configCount = 10
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
	if configCount <= 0 {
		return nil, errors.New("eglGetConfigs() returned zero configs")
	}

	return configurations[:configCount], nil
}

func (display *Display) ChooseConfig(attribList []Attrib) ([]Config, error) {
	var eglAttribs *C.EGLint
	if attribList != nil {
		eglAttribs = (*C.EGLint)(&(attribList[0]))
	}

	var configCount C.EGLint
	success := C.eglChooseConfig(display.eglDisplay, eglAttribs, nil, 0, &configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}

	if configCount <= 0 {
		configCount = 10
	}
	configurations := make([]Config, configCount)

	success = C.eglChooseConfig(
		display.eglDisplay,
		eglAttribs,
		(*C.EGLConfig)(&(configurations[0])),
		configCount,
		&configCount)
//fmt.Printf("success == %d, configCount == %d\n", success, configCount)
	if success == C.EGL_FALSE {
		return nil, getError()
	}
	if configCount <= 0 {
		return nil, errors.New("eglGetConfigs() returned zero configs")
	}

	return configurations[:configCount], nil
}

func (display *Display) GetConfigAttrib(config Config, name Attrib) (Attrib, error) {
	var value Attrib
	success := C.eglGetConfigAttrib(display.eglDisplay, C.EGLConfig(config), C.EGLint(name), (*C.EGLint)(&value))
	if success == C.EGL_FALSE {
		return None, getError()
	}

	return value, nil
}

func (display *Display) PrintAllConfigs() {
	configs, configErr := display.GetConfigs()
	if configErr != nil {
		fmt.Printf("could not get configs: %v\n", configErr)
	}

	fmt.Printf("got %d configs:\n", len(configs))
	for _, name := range(AllConfigAttribNames) {
		fmt.Printf("\t%v:\n\t", name)
		for _, config := range(configs) {
			value, err := display.GetConfigAttrib(config, name)
			if err != nil {
				fmt.Printf("? ")
				continue
			}
			fmt.Printf("%d ", value)
		}
		fmt.Printf("\n")
	}
}

func (display *Display) PrintConfigs(configs []Config) {
	fmt.Printf("got %d configs:\n", len(configs))
	for _, name := range(AllConfigAttribNames) {
		fmt.Printf("\t%v:\n\t", name)
		for _, config := range(configs) {
			value, err := display.GetConfigAttrib(config, name)
			if err != nil {
				fmt.Print("? ")
				continue
			}
			fmt.Printf("%d ", value)
		}
		fmt.Printf("\n")
	}
}

func (display *Display) PrintConfig(config Config) {
	for _, name := range(AllConfigAttribNames) {
		fmt.Printf("%v: ", name)
		value, err := display.GetConfigAttrib(config, name)
		if err != nil {
			fmt.Print("?\n")
			continue
		}
		fmt.Printf("%d\n", value)
	}
	fmt.Printf("\n")
}

func (display *Display) CreatePbufferSurface(config Config, attribList []Attrib) (*Surface, error) {
	var eglAttribs *C.EGLint
	if attribList != nil {
		eglAttribs = (*C.EGLint)(&(attribList[0]))
	}

	eglSurface := C.eglCreatePbufferSurface(display.eglDisplay, C.EGLConfig(config), eglAttribs)
	if eglSurface == noSurface {
		return nil, getError()
	}

	surface := new(Surface)
	//runtime.SetFinalizer(surface, destroySurface)
	surface.Display = display
	surface.eglSurface = eglSurface
	return surface, nil
}

func (display *Display) CreateContext(config Config, shareContext *Context, attribList []Attrib) (*Context, error) {
	var eglShareContext C.EGLContext
	if shareContext != nil {
		eglShareContext = shareContext.eglContext
	}

	var eglAttribs *C.EGLint
	if attribList != nil {
		eglAttribs = (*C.EGLint)(&(attribList[0]))
	}

	eglContext := C.eglCreateContext(display.eglDisplay, C.EGLConfig(config), eglShareContext, eglAttribs)
	if eglContext == noContext {
		return nil, getError()
	}

	context := new(Context)
	//runtime.SetFinalizer(context, destroyContext)
	context.eglContext = eglContext
	context.Display = display
	return context, nil
}

func (display *Display) String() string {
	var buffer bytes.Buffer

	vendor, vendorErr := display.QueryString(Vendor)
	var vendorStr string
	if vendorErr != nil {
		vendorStr = "couldn't query vendor\n"
	} else {
		vendorStr = fmt.Sprintf("vendor is:\n\t%v\n", vendor)
	}
	buffer.WriteString(vendorStr)

	version, versionErr := display.QueryString(Version)
	var versionStr string
	if versionErr != nil {
		versionStr = "couldn't query version\n"
	} else {
		versionStr = fmt.Sprintf("version is:\n\t%v\n", version)
	}
	buffer.WriteString(versionStr)

	extensions, extensionsErr := display.QueryString(Extensions)
	var extensionsStr string
	if extensionsErr != nil {
		extensionsStr = "couldn't query extensions\n"
	} else {
		extensionsStr = fmt.Sprintf("extensions are:\n\t%v\n", extensions)
	}
	buffer.WriteString(extensionsStr)

	apis, apisErr := display.QueryString(ClientAPIs)
	var apisStr string
	if apisErr != nil {
		apisStr = "couldn't query APIs\n"
	} else {
		apisStr = fmt.Sprintf("client APIs are:\n\t%v\n", apis)
	}
	buffer.WriteString(apisStr)

	return buffer.String()
}

/*
 * WARNING: this method does not work with Mesa as of
 * version 8.x. Their implementation violates the EGL
 * specification and does not allow you to release a context.
 */
func (display *Display) ReleaseCurrentContext() error {
	success := C.eglMakeCurrent(display.eglDisplay, noSurface, noSurface, noContext)
	if success == C.EGL_FALSE {
		return getError()
	}
	return nil
}

