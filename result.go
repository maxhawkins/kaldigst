package kaldigst

type Status int

const (
	StatusOK Status = 0
)

type FullFinalResult struct {
	Status        Status  `json:"status"`
	Result        Result  `json:"result"`
	SegmentStart  float64 `json:"segment-start"`
	SegmentLength float64 `json:"segment-length,omitempty"`
	TotalLength   float64 `json:"total-length,omitempty"`
}

type Result struct {
	Final      bool          `json:"final"`
	Hypotheses []NBestResult `json:"hypotheses"`
}

type NBestResult struct {
	Transcript     string               `json:"transcript"`
	Likelihood     float64              `json:"likelihood"`
	PhoneAlignment []PhoneAlignmentInfo `json:"phone-alignment"`
	WordAlignment  []WordAlignmentInfo  `json:"word-alignment"`
}

type WordAlignmentInfo struct {
	Start      float64 `json:"start"`
	Length     float64 `json:"length"`
	Word       string  `json:"word"`
	Confidence float32 `json:"confidence"`
}

type PhoneAlignmentInfo struct {
	Start      float64 `json:"start"`
	Length     float64 `json:"length"`
	Phone      string  `json:"phone"`
	Confidence float32 `json:"confidence"`
}
