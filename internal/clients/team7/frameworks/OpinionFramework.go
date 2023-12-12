package frameworks


type IOpinionFramework[I, SN, O any] interface {
	GetOpinion(inputs I, socialNetwork SN) O
}

type OpinionFrameworkInputs struct {
}

type OpinionFrameworkOutputs struct {
}

type OpinionFramework struct {
	IOpinionFramework[OpinionFrameworkInputs, SocialNetwork, OpinionFrameworkOutputs]
	inputs *OpinionFrameworkInputs
}

// Constructor for OpinionFramework
func NewOpinionFramework(of OpinionFrameworkInputs) *OpinionFramework {
	return &OpinionFramework{inputs: &of}
}

func (of *OpinionFramework) GetOpinion(inputs OpinionFrameworkInputs, socialNetwork SocialNetwork) OpinionFrameworkOutputs {
	of.inputs = &inputs
	// TODO: Formulate opinion
	return OpinionFrameworkOutputs{}
}
