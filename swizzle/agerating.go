package swizzle

const ESRBERating ESRBString = "E"
const ESRBEPlusRating ESRBString = "E10+"
const ESRBTeenRating ESRBString = "T"
const ESRBMatureRating ESRBString = "M"
const ESRBAdultRating ESRBString = "AO"

type AgeRating struct {
	ESRB ESRBString `json:"esrb,omitempty" yaml:"esrb,omitempty"`
}

type ESRBString string
