package egl

/*
#cgo pkg-config: egl x11

#include <EGL/egl.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"os"
	//"runtime"
	"strconv"
	"unsafe"
)

type Display struct {
	xDisplay *C.Display
	eglDisplay C.EGLDisplay
	majorVersion, minorVersion int
}

type Surface struct {
	eglSurface C.EGLSurface
	Display *Display
	xPixmap C.Pixmap
}

/*
const pixelByteCount = 2
const pixelBitCount = pixelByteCount * 8
const pixelMask = 0x0000FFFF // 32 bits, for 32-bit ARGB
*/

const pixelByteCount = 4
const pixelBitCount = pixelByteCount * 8
const pixelMask = 0xFFFFFFFF // 32 bits, for 32-bit ARGB

func init() {
	C.XInitThreads()
}

func openXDisplayWithCString(displayName *C.char) (*Display, error) {
	xDisplay := C.XOpenDisplay(displayName)
	if xDisplay == nil {
		return nil, fmt.Errorf(
			"XOpenDisplay returned nil, display %v may be in use by another user",
			C.GoString(displayName),
		)
	}

	eglDisplay := C.eglGetDisplay(C.EGLNativeDisplayType(xDisplay))
	if eglDisplay == nil {
		return nil, getError()
	}

	display := new(Display)
	display.xDisplay = xDisplay
	display.eglDisplay = eglDisplay

	return display, nil
}

func OpenMainXDisplay() (*Display, error) {
	return openXDisplayWithCString(nil)
}

func OpenXDisplay(name string) (*Display, error) {
	displayName := C.CString(name)
	display, displayErr := openXDisplayWithCString(displayName)
	C.free(unsafe.Pointer(displayName))

	return display, displayErr
}

func OpenAllDisplaysOnXServer(name string) ([]*Display, error) {
	serverName := C.CString(name)
	xDisplay := C.XOpenDisplay(serverName)
	C.free(unsafe.Pointer(serverName))
	if xDisplay == nil {
		return nil, fmt.Errorf(
			"XOpenDisplay returned nil, screen %v may be in use by another user",
			name,
		)
	}

	count := int(C.XScreenCount(xDisplay))
	C.XCloseDisplay(xDisplay)
	if count <= 0 {
		return nil, fmt.Errorf("XScreenCount returned %d", count)
	}

	baseName := name + "."
	var lastError error
	var allDisplays []*Display
	for i := 0; i < count; i++ {
		displayName := baseName + strconv.Itoa(i)
		display, displayErr := OpenXDisplay(displayName)
		if displayErr != nil {
			fmt.Printf("failed on display %v with error %v\n", displayName, displayErr)
			lastError = displayErr
			continue
		}

		//fmt.Printf("found X display named %v\n", displayName)
		allDisplays = append(allDisplays, display)
	}
	if allDisplays == nil {
		return nil, lastError
	}

	return allDisplays, nil
}

func OpenAllDisplaysOnAllXServers() ([]*Display, error) {
	servers, serversErr := allXServers()
	if serversErr != nil {
		// sane default
		servers = []string{":0"}
	}

	var allDisplays []*Display
	var lastError error
	for _, oneServer := range servers {
		theseDisplays, displaysErr := OpenAllDisplaysOnXServer(oneServer)
		if displaysErr != nil {
			lastError = displaysErr
			continue
		}

		allDisplays = append(allDisplays, theseDisplays...)
	}
	if allDisplays == nil {
		return nil, lastError
	}

	return allDisplays, nil
}

func allXServers() ([]string, error) {
	xDir, openErr := os.Open("/tmp/.X11-unix")
	if openErr != nil {
		return nil, openErr
	}

	const maxFileCount = 1000
	files, readErr := xDir.Readdir(maxFileCount)
	xDir.Close()
	if readErr != nil {
		return nil, readErr
	}

	var allNames []string
	for _, oneFile := range files {
		name := oneFile.Name()
		if len(name) < 2 {
			continue
		}

		allNames = append(allNames, ":" + name[1:])
	}
	if allNames == nil {
		return nil, errors.New("No X sockets found in /tmp/.X11-unix")
	}

	return allNames, nil
}

func (display *Display) Close() error {
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

func (display *Display) CreatePixmapSurface(config Config, attribList []Attrib, width, height int) (*Surface, error) {
	rootWindow := C.XDefaultRootWindow(display.xDisplay)
//	fmt.Printf("got root window == %d\n", rootWindow)

	if width < 0 {
		width = -width
	}
	if height < 0 {
		height = -height
	}
	pixmap := C.XCreatePixmap(display.xDisplay, C.Drawable(rootWindow), C.uint(width), C.uint(height), pixelBitCount)
//	fmt.Printf("created pixmap == %d\n", pixmap)

/*
	// fills magenta for debugging
	gc := C.XCreateGC(display.xDisplay, C.Drawable(pixmap), 0, nil)
//	C.XSetForeground(display.xDisplay, gc, 0xFFFF00FF)
//	C.XFillRectangle(display.xDisplay, C.Drawable(pixmap), gc, 0, 0, C.uint(width), C.uint(height))
*/

	var eglAttribs *C.EGLint
	if attribList != nil {
		eglAttribs = (*C.EGLint)(&(attribList[0]))
	}
	eglSurface := C.eglCreatePixmapSurface(display.eglDisplay, C.EGLConfig(config), C.EGLNativePixmapType(pixmap), eglAttribs)
	if eglSurface == noSurface {
		C.XFreePixmap(display.xDisplay, pixmap)
		return nil, getError()
	}

	surface := new(Surface)
	//runtime.SetFinalizer(surface, destroySurface)
	surface.Display = display
	surface.eglSurface = eglSurface
	surface.xPixmap = pixmap
	return surface, nil
}

func (surface *Surface) Destroy() error {
	var result error

	success := C.eglDestroySurface(surface.Display.eglDisplay, surface.eglSurface)
	if success == C.EGL_FALSE {
		result = getError()
	}

	if surface.xPixmap != 0 {
		C.XFreePixmap(surface.Display.xDisplay, surface.xPixmap)
	}

	return result
}

func (surface *Surface) CopyBuffers() (*image.NRGBA, error) {
	width, widthErr := surface.Query(Width)
	if widthErr != nil {
		return nil, widthErr
	}

	height, heightErr := surface.Query(Height)
	if heightErr != nil {
		return nil, heightErr
	}

	display := surface.Display
	xDisplay := display.xDisplay
	pixmap := surface.xPixmap
	if pixmap == 0 {
		pixmap = C.XCreatePixmap(xDisplay, C.Drawable(C.XDefaultRootWindow(xDisplay)), C.uint(width), C.uint(height), pixelBitCount)
//		pixmap = C.XCreatePixmap(xDisplay, C.Drawable(C.XDefaultScreen(xDisplay)), C.uint(width), C.uint(height), pixelBitCount)
		fmt.Printf("surface is not a pixmap surface, created temporary pixmap == %d\n", pixmap)
		defer C.XFreePixmap(xDisplay, pixmap)

		C.eglCopyBuffers(display.eglDisplay, surface.eglSurface, C.EGLNativePixmapType(pixmap))
	}

	xImage := C.XGetImage(xDisplay, C.Drawable(pixmap), 0, 0, C.uint(width), C.uint(height), pixelMask, C.ZPixmap)
	if xImage == nil {
		return nil, errors.New("XGetImage returned nil")
	}

/*
	fmt.Printf("xImage 0x%X contains:\n", unsafe.Pointer(xImage))
	fmt.Printf("\twidth == %d\n", xImage.width)
	fmt.Printf("\theight == %d\n", xImage.height)
	fmt.Printf("\tformat == %d\n", xImage.format)
	fmt.Printf("\tdata == 0x%X\n", xImage.data)
	fmt.Printf("\t*data == 0x%X\n", *(*uint)(unsafe.Pointer(xImage.data)))
	fmt.Printf("\tbyte_order == %d\n", xImage.byte_order)
	fmt.Printf("\tbitmap_unit == %d\n", xImage.bitmap_unit)
	fmt.Printf("\tbitmap_bit_order == %d\n", xImage.bitmap_bit_order)
	fmt.Printf("\tbitmap_pad == %d\n", xImage.bitmap_pad)
	fmt.Printf("\tdepth == %d\n", xImage.depth)
	fmt.Printf("\tbytes_per_line == %d\n", xImage.bytes_per_line)
	fmt.Printf("\tbits_per_pixel == %d\n", xImage.bits_per_pixel)
	fmt.Printf("\tred_mask == 0x%X\n", xImage.red_mask)
	fmt.Printf("\tgreen_mask == 0x%X\n", xImage.green_mask)
	fmt.Printf("\tblue_mask == 0x%X\n", xImage.blue_mask)
*/

	bounds := image.Rect(0, 0, int(xImage.width), int(xImage.height))
	goImage := image.NewNRGBA(bounds)
	rowBytes := int(xImage.width) * pixelByteCount
	xLength := xImage.height * xImage.bytes_per_line
	xSlice := C.GoBytes(unsafe.Pointer(xImage.data), xLength)
	for y := 0; y < int(xImage.height); y++ {
		xOffset := y * int(xImage.bytes_per_line)
		goOffset := y * goImage.Stride

		copy(goImage.Pix[goOffset:goOffset + rowBytes], xSlice[xOffset:xOffset + rowBytes])
	}

	return goImage, nil
}
