package egl

/*
#cgo LDFLAGS: -lEGL

#include <EGL/egl.h>
*/
import "C"

import (
	"fmt"
)

// Out-of-band attribute value
const (
	DontCare = C.EGL_DONT_CARE
)

// Config attributes
const (
	BufferSize = C.EGL_BUFFER_SIZE
	AlphaSize = C.EGL_ALPHA_SIZE
	BlueSize = C.EGL_BLUE_SIZE
	GreenSize = C.EGL_GREEN_SIZE
	RedSize = C.EGL_RED_SIZE
	DepthSize = C.EGL_DEPTH_SIZE
	StencilSize = C.EGL_STENCIL_SIZE
	ConfigCaveat = C.EGL_CONFIG_CAVEAT
	ConfigId = C.EGL_CONFIG_ID
	Level = C.EGL_LEVEL
	MaxPbufferHeight = C.EGL_MAX_PBUFFER_HEIGHT
	MaxPbufferPixels = C.EGL_MAX_PBUFFER_PIXELS
	MaxPbufferWidth = C.EGL_MAX_PBUFFER_WIDTH
	NativeRenderable = C.EGL_NATIVE_RENDERABLE
	NativeVisualId = C.EGL_NATIVE_VISUAL_ID
	NativeVisualType = C.EGL_NATIVE_VISUAL_TYPE
	Samples = C.EGL_SAMPLES
	SampleBuffers = C.EGL_SAMPLE_BUFFERS
	SurfaceType = C.EGL_SURFACE_TYPE
	TransparentType = C.EGL_TRANSPARENT_TYPE
	TransparentBlueValue = C.EGL_TRANSPARENT_BLUE_VALUE
	TransparentGreenValue = C.EGL_TRANSPARENT_GREEN_VALUE
	TransparentRedValue = C.EGL_TRANSPARENT_RED_VALUE
	None = C.EGL_NONE // Attrib list terminator
	BindToTextureRgb = C.EGL_BIND_TO_TEXTURE_RGB
	BindToTextureRgba = C.EGL_BIND_TO_TEXTURE_RGBA
	MinSwapInterval = C.EGL_MIN_SWAP_INTERVAL
	MaxSwapInterval = C.EGL_MAX_SWAP_INTERVAL
	LuminanceSize = C.EGL_LUMINANCE_SIZE
	AlphaMaskSize = C.EGL_ALPHA_MASK_SIZE
	ColorBufferType = C.EGL_COLOR_BUFFER_TYPE
	RenderableType = C.EGL_RENDERABLE_TYPE
	MatchNativePixmap = C.EGL_MATCH_NATIVE_PIXMAP // Pseudo-attribute (not queryable)
	Conformant = C.EGL_CONFORMANT
)

var AllConfigAttribNames [34]Attrib = [34]Attrib{
	BufferSize,
	AlphaSize,
	BlueSize,
	GreenSize,
	RedSize,
	DepthSize,
	StencilSize,
	ConfigCaveat,
	ConfigId,
	Level,
	MaxPbufferHeight,
	MaxPbufferPixels,
	MaxPbufferWidth,
	NativeRenderable,
	NativeVisualId,
	NativeVisualType,
	Samples,
	SampleBuffers,
	SurfaceType,
	TransparentType,
	TransparentBlueValue,
	TransparentGreenValue,
	TransparentRedValue,
	None,
	BindToTextureRgb,
	BindToTextureRgba,
	MinSwapInterval,
	MaxSwapInterval,
	LuminanceSize,
	AlphaMaskSize,
	ColorBufferType,
	RenderableType,
	MatchNativePixmap,
	Conformant,
}

// Config attribute mask bits
const (
	PbufferBit = C.EGL_PBUFFER_BIT
	PixmapBit = C.EGL_PIXMAP_BIT
	WindowBit = C.EGL_WINDOW_BIT
	VgColorspaceLinearBit = C.EGL_VG_COLORSPACE_LINEAR_BIT
	VgAlphaFormatPreBit = C.EGL_VG_ALPHA_FORMAT_PRE_BIT
	MultisampleResolveBoxBit = C.EGL_MULTISAMPLE_RESOLVE_BOX_BIT
	SwapBehaviorPreservedBit = C.EGL_SWAP_BEHAVIOR_PRESERVED_BIT
	OpenglEsBit = C.EGL_OPENGL_ES_BIT
	OpenvgBit = C.EGL_OPENVG_BIT
	OpenglEs2Bit = C.EGL_OPENGL_ES2_BIT
	OpenglBit = C.EGL_OPENGL_BIT
)

// QueryString targets
const (
	Vendor = C.EGL_VENDOR
	Version = C.EGL_VERSION
	Extensions = C.EGL_EXTENSIONS
	ClientAPIs = C.EGL_CLIENT_APIS
)

// QuerySurface / SurfaceAttrib / CreatePbufferSurface targets
const (
	Height = C.EGL_HEIGHT
	Width = C.EGL_WIDTH
	LargestPbuffer = C.EGL_LARGEST_PBUFFER
	TextureFormat = C.EGL_TEXTURE_FORMAT
	TextureTarget = C.EGL_TEXTURE_TARGET
	MipmapTexture = C.EGL_MIPMAP_TEXTURE
	MipmapLevel = C.EGL_MIPMAP_LEVEL
	RenderBuffer = C.EGL_RENDER_BUFFER
	VgColorspace = C.EGL_VG_COLORSPACE
	VgAlphaFormat = C.EGL_VG_ALPHA_FORMAT
	HorizontalResolution = C.EGL_HORIZONTAL_RESOLUTION
	VerticalResolution = C.EGL_VERTICAL_RESOLUTION
	PixelAspectRatio = C.EGL_PIXEL_ASPECT_RATIO
	SwapBehavior = C.EGL_SWAP_BEHAVIOR
	MultisampleResolve = C.EGL_MULTISAMPLE_RESOLVE
)

// CreateContext attributes
const (
	ContextClientVersion = C.EGL_CONTEXT_CLIENT_VERSION
)

// BindAPI/QueryAPI targets
const (
	OpenglEsAPI = C.EGL_OPENGL_ES_API
	OpenvgAPI = C.EGL_OPENVG_API
	OpenglAPI = C.EGL_OPENGL_API
)

func (name Attrib) String() string {
	switch name {
		// Config attributes
		case BufferSize:
			return "total color component bits in the color buffer"
		case RedSize:
			return "bits of Red in the color buffer"
		case GreenSize:
			return "bits of Green in the color buffer"
		case BlueSize:
			return "bits of Blue in the color buffer"
		case LuminanceSize:
			return "bits of Luminance in the color buffer"
		case AlphaSize:
			return "bits of Alpha in the color buffer"
		case BindToTextureRgb:
			return "True if bindable to RGB textures"
		case BindToTextureRgba:
			return "True if bindable to RGBA textures"
		case ColorBufferType:
			return "color buffer type"
		case ConfigCaveat:
			return "any caveats for the configuration"
		case ConfigId:
			return "unique EGLConfig identifier"
		case Conformant:
			return "whether contexts created with this config are conformant"
		case DepthSize:
			return "bits of Z in the depth buffer"
		case Level:
			return "frame buffer level"
		case MaxPbufferWidth:
			return "maximum width of pbuffer"
		case MaxPbufferHeight:
			return "maximum height of pbuffer"
		case MaxPbufferPixels:
			return "maximum size of pbuffer"
		case MaxSwapInterval:
			return "maximum swap interval"
		case MinSwapInterval:
			return "minimum swap interval"
		case NativeRenderable:
			return "EGL_TRUE if native rendering APIs can render to surface"
		case NativeVisualId:
			return "handle of corresponding native visual"
		case NativeVisualType:
			return "native visual type of the associated visual"
		case RenderableType:
			return "which client APIs are supported"
		case SampleBuffers:
			return "number of multisample buffers"
		case Samples:
			return "number of samples per pixel"
		case StencilSize:
			return "bits of Stencil in the stencil buffer"
		case SurfaceType:
			return "which types of EGL surfaces are supported"
		case TransparentType:
			return "type of transparency supported"
		case TransparentRedValue:
			return "transparent red value"
		case TransparentGreenValue:
			return "transparent green value"
		case TransparentBlueValue:
			return "transparent blue value"

		// Surface attributes
		case VgAlphaFormat:
			return "Alpha format for OpenVG"
		case VgColorspace:
			return "Color space for OpenVG"
		case Height:
			return "Height of surface"
		case HorizontalResolution:
			return "Horizontal dot pitch"
		case LargestPbuffer:
			return "If true, create largest pbuffer possible"
		case MipmapTexture:
			return "True if texture has mipmaps"
		case MipmapLevel:
			return "Mipmap level to render to"
		case MultisampleResolve:
			return "Multisample resolve behavior"
		case PixelAspectRatio:
			return "Display aspect ratio"
		case RenderBuffer:
			return "Render buffer"
		case SwapBehavior:
			return "Buffer swap behavior"
		case TextureFormat:
			return "Format of texture"
		case TextureTarget:
			return "Type of texture"
		case VerticalResolution:
			return "Vertical dot pitch"
		case Width:
			return "Width of surface"
	}
	return fmt.Sprintf("EGL attribute name %d", name)
}
