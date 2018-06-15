package kaldigst

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ziutek/glib"
	"github.com/ziutek/gst"
)

type Decoder struct {
	src                                  *AppSrc
	decode, convert, sink, resample, asr *gst.Element
	bus                                  *gst.Bus
	pipe                                 *gst.Pipeline
}

func NewDecoder(props Props) (*Decoder, error) {
	d := &Decoder{}

	d.src = NewAppSrc("appsrc")
	if d.src == nil {
		return nil, fmt.Errorf("failed to create appsrc")
	}

	d.decode = gst.ElementFactoryMake("decodebin", "decodebin")
	if d.decode == nil {
		return nil, fmt.Errorf("failed to create decodebin")
	}
	d.convert = gst.ElementFactoryMake("audioconvert", "audioconvert")
	if d.convert == nil {
		return nil, fmt.Errorf("failed to create audioconvert")
	}
	d.sink = gst.ElementFactoryMake("fakesink", "fakesink")
	if d.sink == nil {
		return nil, fmt.Errorf("failed to create fakesink")
	}
	d.resample = gst.ElementFactoryMake("audioresample", "audioresample")
	if d.resample == nil {
		return nil, fmt.Errorf("failed to create audioresample")
	}
	d.asr = gst.ElementFactoryMake("kaldinnet2onlinedecoder", "asr")
	if d.asr == nil {
		return nil, fmt.Errorf("failed to create kaldinnet2onlinedecoder")
	}

	props.set(d.asr)

	d.pipe = gst.NewPipeline("pipeline")
	d.pipe.Add(d.src.Element, d.decode, d.convert, d.resample, d.asr, d.sink)

	d.src.Link(d.decode)
	d.decode.ConnectNoi("pad-added", func(convertSink, pad *gst.Pad) {
		ret := pad.Link(convertSink)
		if ret != gst.PAD_LINK_OK {
			log.Fatal("audioconvert link error")
		}
	}, d.convert.GetStaticPad("sink"))
	d.convert.Link(d.resample)
	d.resample.Link(d.asr)
	d.asr.Link(d.sink)

	d.bus = d.pipe.GetBus()
	d.bus.AddSignalWatch()

	d.bus.Connect("message::error", func(bus *gst.Bus, msg *gst.Message) {
		d.pipe.SetState(gst.STATE_NULL)
	}, nil)
	d.bus.Connect("message::eos", func(bus *gst.Bus, msg *gst.Message) {
		d.pipe.SetState(gst.STATE_NULL)
	}, nil)

	return d, nil
}

func (d *Decoder) OnError(f func(err *glib.Error, debug string)) {
	d.bus.Connect("message::error", func(bus *gst.Bus, msg *gst.Message) {
		// msg.ParseError() // TODO(maxhawkins): why does this crash?
		f(nil, "")
	}, nil)
}

func (d *Decoder) OnEOS(f func()) {
	d.bus.Connect("message::eos", func(bus *gst.Bus, msg *gst.Message) {
		f()
	}, nil)
}

func (d *Decoder) Start(caps *gst.Caps) {
	d.src.SetCaps(caps)
	d.pipe.SetState(gst.STATE_PLAYING)
}

func (d *Decoder) OnFullFinalResult(f func(FullFinalResult)) {
	d.asr.Connect("full-final-result", func(bus *gst.Bus, data string) {
		var res FullFinalResult
		if err := json.Unmarshal([]byte(data), &res); err != nil {
			log.Fatal(err)
		}

		f(res)
	}, nil)
}

func (d *Decoder) CloseWrite() error {
	return d.src.Close()
}

func (d *Decoder) Write(b []byte) (int, error) {
	return d.src.Write(b)
}
