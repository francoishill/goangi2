package httpUtils

import (
	"github.com/astaxie/beego/context"
)

const (
	CSS_CACHE_CONTROL_HEADER   = "public, max-age=31536000"
	JS_CACHE_CONTROL_HEADER    = "public, max-age=31536000"
	JSON_CACHE_CONTROL_HEADER  = "public, max-age=31536000"
	IMAGE_CACHE_CONTROL_HEADER = "public, max-age=31536000"
	FONT_CACHE_CONTROL_HEADER  = "public, max-age=31536000"
	HTML_CACHE_CONTROL_HEADER  = "public, max-age=31536000"
)

var fileExtensionWithCacheControl map[string]string = map[string]string{
	".css": CSS_CACHE_CONTROL_HEADER,

	".js": JS_CACHE_CONTROL_HEADER,

	".json": JSON_CACHE_CONTROL_HEADER,

	".ico":  IMAGE_CACHE_CONTROL_HEADER,
	".gif":  IMAGE_CACHE_CONTROL_HEADER,
	".png":  IMAGE_CACHE_CONTROL_HEADER,
	".jpg":  IMAGE_CACHE_CONTROL_HEADER,
	".jpeg": IMAGE_CACHE_CONTROL_HEADER,

	".woff": FONT_CACHE_CONTROL_HEADER,
	".eot":  FONT_CACHE_CONTROL_HEADER,
	".ttf":  FONT_CACHE_CONTROL_HEADER,
	".svg":  FONT_CACHE_CONTROL_HEADER,

	".html": HTML_CACHE_CONTROL_HEADER,
}

var fileExtensionsWithVaryAcceptEncodingHeader = []string{".js", ".json", ".css", ".xml", ".gz", ".html", ".woff", ".eot", ".ttf", ".svg"}

var fileExtensionsWithContentType = map[string]string{
	".json": "application/x-javascript",
}

func SetCacheControlHeader(ctx *context.Context, headerValue string) {
	ctx.Output.Header("Cache-Control", headerValue)
}

func SetVaryHeaderToAcceptEncoding(ctx *context.Context) {
	ctx.Output.Header("Vary", "Accept-Encoding")
}

func SetContentType(ctx *context.Context, headerValue string) {
	ctx.Output.Header("Content-Type", headerValue)
}
