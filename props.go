package kaldigst

import (
	"github.com/ziutek/gst"
)

type Props struct {
	// nnet mode: 2 for nnet2, 3 for nnet3
	NNetMode int `json:"nnetMode"`
	// Silence the decoder
	Silent bool `json:"silent"`
	// Filename of the acoustic model
	Model string `json:"model"`
	// Filename of the HCLG FST
	FST string `json:"fst"`
	// Name of word symbols file (typically words.txt)
	WordSyms string `json:"wordSyms"`
	// Name of phoneme symbols file (typically phones.txt)
	PhoneSyms string `json:"phoneSyms"`
	// If true, output phoneme-level alignment
	DoPhoneAlignment bool `json:"doPhoneAlignment"`
	// If true, apply endpoint detection, and split the audio at endpoints
	DoEndpointing bool `json:"doEndpointing"`
	// Current adaptation state, in stringified form, set to empty string to reset
	AdaptationState string `json:"adaptationState"`
	// If true, inverse the acoustic scaling of the output lattice
	InverseScale bool `json:"inverseScale"`
	// LM scaling for the output lattice, usually in conjunction with inverse-scaling=true
	LMWTScale float64 `json:"lmwtScale"`
	// Smaller values decrease latency, bigger values (e.g. 0.2) improve speed if multithreaded BLAS/MKL is used
	ChunkLengthSecs float64 `json:"chunkLengthsSecs"`
	// Time period after which new interim recognition result is sent
	TracebackPeriodSecs float64 `json:"tracebackPeriodSecs"`
	// Language language model FST (G.fst), only needed when rescoring with the constant ARPA LM
	LMFST string `json:"lmFST"`
	// Big language model in constant ARPA format (typically G.carpa), to be used for rescoring final lattices. Also requires 'lm-fst' property
	BigLMConstARPA string `json:"bigLMConstARPA"`
	// Use a decoder that does feature calculation and decoding in separate threads (NB! must be set before other properties)
	UseThreadedDecoder bool `json:"useThreadedDecoder"`
	// number of hypotheses in the full final results
	NumNBest int `json:"numNBest"`
	// number of hypotheses where alignment should be done
	NumPhoneAlignment int `json:"numPhoneAlignment"`
	// Word-boundary file. Setting this property triggers generating word alignments in full results
	//
	// Word-boundary file has format (on each line): <integer-phone-id> [begin|end|singleton|internal|nonword]
	WordBoundaryFile string `json:"wordBoundaryFile"`
	// Minimal number of words in the first transcription for triggering update of the adaptation state
	MinWordsForIVector int `json:"minWordsForIVector"`

	// For NNET3

	MFCCConfig              string
	IVectorExtractionConfig string

	AcousticScale          float64
	Beam                   float64
	FrameSubsamplingFactor float64
	LatticeBeam            float64
	MaxActive              float64
	MaxMem                 float64
	EndpointSilencePhones  string
}

func (p Props) set(asr *gst.Element) {
	// These must be first
	asr.SetProperty("use-threaded-decoder", p.UseThreadedDecoder)
	if p.NNetMode > 0 {
		asr.SetProperty("nnet-mode", p.NNetMode)
	}

	if p.ChunkLengthSecs > 0 {
		asr.SetProperty("chunk-length-in-secs", p.ChunkLengthSecs)
	}
	if p.NumNBest > 0 {
		asr.SetProperty("num-nbest", p.NumNBest)
	}
	if p.NumPhoneAlignment > 0 {
		asr.SetProperty("num-phone-alignment", p.NumPhoneAlignment)
	}
	if p.TracebackPeriodSecs > 0 {
		asr.SetProperty("traceback-period-in-secs", p.TracebackPeriodSecs)
	}
	if p.WordBoundaryFile != "" {
		asr.SetProperty("word-boundary-file", p.WordBoundaryFile)
	}
	if p.WordSyms != "" {
		asr.SetProperty("word-syms", p.WordSyms)
	}
	if p.PhoneSyms != "" {
		asr.SetProperty("phone-syms", p.PhoneSyms)
	}
	if p.LMWTScale != 0 {
		asr.SetProperty("lmwt-scale", p.LMWTScale)
	}
	if p.LMFST != "" {
		asr.SetProperty("lm-fst", p.LMFST)
	}
	if p.BigLMConstARPA != "" {
		asr.SetProperty("big-lm-const-arpa", p.BigLMConstARPA)
	}
	if p.MinWordsForIVector > 0 {
		asr.SetProperty("min-words-for-ivector", p.MinWordsForIVector)
	}
	asr.SetProperty("do-phone-alignment", p.DoPhoneAlignment)
	asr.SetProperty("do-endpointing", p.DoEndpointing)
	asr.SetProperty("silent", p.Silent)
	asr.SetProperty("inverse-scale", p.InverseScale)

	if p.NNetMode == 3 {
		if p.AcousticScale != 0 {
			asr.SetProperty("acoustic-scale", p.AcousticScale)
		}
		if p.Beam != 0 {
			asr.SetProperty("beam", p.Beam)
		}
		if p.FrameSubsamplingFactor != 0 {
			asr.SetProperty("frame-subsampling-factor", p.FrameSubsamplingFactor)
		}
		if p.LatticeBeam != 0 {
			asr.SetProperty("lattice-beam", p.LatticeBeam)
		}
		if p.MaxActive != 0 {
			asr.SetProperty("max-active", p.MaxActive)
		}
		if p.MaxMem != 0 {
			asr.SetProperty("max-mem", p.MaxMem)
		}
		if p.MFCCConfig != "" {
			asr.SetProperty("mfcc-config", p.MFCCConfig)
		}
		if p.IVectorExtractionConfig != "" {
			asr.SetProperty("ivector-extraction-config", p.IVectorExtractionConfig)
		}
		if p.EndpointSilencePhones != "" {
			asr.SetProperty("endpoint-silence-phones", p.EndpointSilencePhones)
		}
	}

	// These must be last
	if p.FST != "" {
		asr.SetProperty("fst", p.FST)
	}
	if p.Model != "" {
		asr.SetProperty("model", p.Model)
	}
}
