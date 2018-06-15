package kaldigst

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ziutek/glib"
	"github.com/ziutek/gst"
)

type Decoder struct {
	src                                  *AppSrc
	decode, convert, sink, resample, asr *gst.Element
	bus                                  *gst.Bus
	pipe                                 *gst.Pipeline

	OnEOS             func()
	OnError           func(err *glib.Error, debug string)
	OnFullFinalResult func(res FullFinalResult)
}

func NewDecoder(props Props) (*Decoder, error) {
	d := &Decoder{}

	d.OnEOS = func() {}
	d.OnError = func(err *glib.Error, debug string) {}
	d.OnFullFinalResult = func(res FullFinalResult) {}

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
		d.Stop()
		d.OnError(nil, "")
	}, nil)
	d.bus.Connect("message::eos", func(bus *gst.Bus, msg *gst.Message) {
		d.Stop()
		d.OnEOS()
	}, nil)

	d.asr.Connect("full-final-result", func(bus *gst.Bus, data string) {
		var res FullFinalResult
		if err := json.Unmarshal([]byte(data), &res); err != nil {
			d.OnError(nil, err.Error())
		}

		d.OnFullFinalResult(res)
	}, nil)

	return d, nil
}

func (d *Decoder) Start(caps *gst.Caps, adaptationState string) error {
	state, _, ret := d.pipe.GetState(int64(time.Second))
	if ret != gst.STATE_CHANGE_SUCCESS {
		return fmt.Errorf("start: check pipeline state failed")
	}
	if state != gst.STATE_NULL {
		return fmt.Errorf("start: pipeline not ready (state=%v)", state)
	}

	d.src.SetCaps(caps)

	ret = d.pipe.SetState(gst.STATE_PLAYING)
	if ret == gst.STATE_CHANGE_FAILURE {
		return fmt.Errorf("start: set playing failed")
	}

	d.asr.SetProperty("adaptation-state", adaptationState)

	return nil
}

func (d *Decoder) Stop() error {
	ret := d.pipe.SetState(gst.STATE_NULL)
	if ret != gst.STATE_CHANGE_SUCCESS {
		return fmt.Errorf("decoder: set null failed")
	}

	return nil
}

func (d *Decoder) CloseWrite() error {
	return d.src.EOS()
}

func (d *Decoder) Write(b []byte) (int, error) {
	return d.src.Write(b)
}
