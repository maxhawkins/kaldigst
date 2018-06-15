package kaldigst

import (
	"fmt"
	"unsafe"

	"github.com/ziutek/gst"
)

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo LDFLAGS: -lgstapp-1.0
#include <stdlib.h>
#include <string.h>
#include <gst/gst.h>
#include <gst/app/gstappsrc.h>
*/
import "C"

type FlowReturn C.GstFlowReturn

const (
	GST_FLOW_OK             = FlowReturn(C.GST_FLOW_OK)
	GST_FLOW_FLUSHING       = FlowReturn(C.GST_FLOW_FLUSHING)
	GST_FLOW_NOT_LINKED     = FlowReturn(C.GST_FLOW_NOT_LINKED)
	GST_FLOW_NOT_NEGOTIATED = FlowReturn(C.GST_FLOW_NOT_NEGOTIATED)
	GST_FLOW_ERROR          = FlowReturn(C.GST_FLOW_ERROR)
	GST_FLOW_NOT_SUPPORTED  = FlowReturn(C.GST_FLOW_NOT_SUPPORTED)
)

func (f FlowReturn) String() string {
	switch f {
	case GST_FLOW_OK:
		return "GST_FLOW_OK"
	case GST_FLOW_FLUSHING:
		return "GST_FLOW_FLUSHING"
	case GST_FLOW_NOT_LINKED:
		return "GST_FLOW_NOT_LINKED"
	case GST_FLOW_NOT_NEGOTIATED:
		return "GST_FLOW_NOT_NEGOTIATED"
	case GST_FLOW_ERROR:
		return "GST_FLOW_ERROR"
	case GST_FLOW_NOT_SUPPORTED:
		return "GST_FLOW_NOT_SUPPORTED"
	default:
		return fmt.Sprintf("flow error: %d", f)
	}
}

type AppSrc struct {
	*gst.Element
}

func NewAppSrc(name string) *AppSrc {
	element := gst.ElementFactoryMake("appsrc", name)

	element.SetProperty("is-live", true)
	element.SetProperty("block", true)

	return &AppSrc{element}
}

func (a *AppSrc) g() *C.GstAppSrc {
	return (*C.GstAppSrc)(a.GetPtr())
}

func (a *AppSrc) SetCaps(caps *gst.Caps) {
	p := unsafe.Pointer(caps) // HACK
	C.gst_app_src_set_caps(a.g(), (*C.GstCaps)(p))
}

func (a *AppSrc) Close() error {
	ret := FlowReturn(C.gst_app_src_end_of_stream(a.g()))
	if FlowReturn(ret) != GST_FLOW_OK {
		return fmt.Errorf("close appsrc: %v", ret)
	}
	return nil
}

func (a *AppSrc) Write(d []byte) (int, error) {
	buf := C.gst_buffer_new_allocate(nil, C.gsize(len(d)), nil)
	n := C.gst_buffer_fill(buf, C.gsize(0), (C.gconstpointer)(C.CBytes(d)), C.gsize(len(d)))

	ret := FlowReturn(C.gst_app_src_push_buffer((*C.GstAppSrc)(a.GetPtr()), buf))
	if ret != GST_FLOW_OK {
		return 0, fmt.Errorf("appsrc push buffer failed: %v", ret)
	}

	return int(n), nil
}
