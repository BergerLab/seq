package seq

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// HMM states in the Plan7 architecture.
const (
	Match HMMState = iota
	Deletion
	Insertion
	Begin
	End
)

type HMMState int

type HMM struct {
	// An ordered list of HMM nodes.
	Nodes []HMMNode

	// The alphabet as defined by an ordering of residues.
	// Indices in this slice correspond to indices in match/insertion emissions.
	Alphabet Alphabet

	// NULL model. (Amino acid background frequencies.)
	// HMMER hmm files don't have this, but HHsuite hhm files do.
	// In the case of HHsuite, the NULL model is used for insertion emissions
	// in every node.
	Null EProbs
}

type HMMNode struct {
	HMM                 *HMM
	Residue             Residue
	NodeNum             int
	InsEmit             EProbs
	MatEmit             EProbs
	Transitions         TProbs
	NeffM, NeffI, NeffD Prob
}

// EProbs represents emission probabilities, as log-odds scores.
type EProbs map[Residue]Prob

// NewEProbs creates a new EProbs map from the given alphabet. Keys of the map
// are residues defined in the alphabet, and values are defaulted to the
// minimum probability.
func NewEProbs(alphabet Alphabet) EProbs {
	ep := make(EProbs, len(alphabet))
	for _, residue := range alphabet {
		ep[residue] = MinProb
	}
	return ep
}

// Returns the emission probability for a particular residue.
func (ep EProbs) EmitProb(r Residue) Prob {
	return ep[r]
}

func (ep *EProbs) MarshalJSON() ([]byte, error) {
	strmap := make(map[string]Prob, len(*ep))
	for k, v := range *ep {
		strmap[string(k)] = v
	}
	return json.Marshal(strmap)
}

func (ep *EProbs) UnmarshalJSON(bs []byte) error {
	var strmap map[string]Prob
	if err := json.Unmarshal(bs, &strmap); err != nil {
		return err
	}
	if *ep == nil {
		*ep = make(EProbs, len(strmap))
	}
	for k, v := range strmap {
		(*ep)[Residue(k[0])] = v
	}
	return nil
}

// TProbs represents transition probabilities, as log-odds scores.
// Note that ID and DI are omitted (Plan7).
type TProbs struct {
	MM, MI, MD, IM, II, DM, DD Prob
}

// Prob represents a transition or emission probability.
type Prob float64

var invalidProb = Prob(math.NaN())

// The value representing a minimum emission/transition probability.
// Remember, max in negative log space is minimum probability.
var MinProb = Prob(math.MaxFloat64)

// NewProb creates a new probability value from a string (usually read from
// an hmm or hhm file). If the string is equivalent to the special value "*",
// then the probability returned is guaranteed to be minimal. Otherwise, the
// string is parsed as a float, and an error returned if parsing fails.
func NewProb(fstr string) (Prob, error) {
	if fstr == "*" {
		return MinProb, nil
	}

	f, err := strconv.ParseFloat(fstr, 64)
	if err != nil {
		return invalidProb,
			fmt.Errorf("Could not convert '%s' to a log probability: %s",
				fstr, err)
	}
	return Prob(f), nil
}

// Less returns true if `p1` represents a smaller probability than `p2`.
func (p1 Prob) Less(p2 Prob) bool {
	return p1 > p2
}

// IsMin returns true if the probability is minimal.
func (p Prob) IsMin() bool {
	return p == MinProb
}

// Ratio returns the log probability as a ratio in the range [0, 1].
// (This assumes that `p` is a log-odds score.)
func (p Prob) Ratio() float64 {
	if p.IsMin() {
		return 0.0
	}
	return math.Exp(-(float64(p)))
}

// Distance returns the distance between `p1` and `p2`.
func (p1 Prob) Distance(p2 Prob) float64 {
	return math.Abs(float64(p1) - float64(p2))
}

// String returns a string representation of the probability.
// When `p` is the minimum probability, then "*" is used.
// Otherwise, the full number is written.
func (p Prob) String() string {
	if p.IsMin() {
		return "*"
	}
	return fmt.Sprintf("%v", float64(p))
}

func (p Prob) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Prob) UnmarshalJSON(bs []byte) error {
	var str string
	var err error
	if err := json.Unmarshal(bs, &str); err != nil {
		return err
	}
	if *p, err = NewProb(str); err != nil {
		return err
	}
	return nil
}

// NewHMM creates a new HMM from a list of nodes, an ordered alphabet and a
// set of null probabilities (which may be nil).
func NewHMM(nodes []HMMNode, alphabet []Residue, null EProbs) *HMM {
	return &HMM{
		Nodes:    nodes,
		Alphabet: alphabet,
		Null:     null,
	}
}

// Slice returns a slice of the HMM given. A slice of an HMM returns only the
// HMM nodes (i.e., columns or match/delete states) in the region specified
// by the slice. Also, the transition probabilities of the last state are
// specially set: M->M = 0, M->I = *, M->D = *, I->M = 0, I->I = *, D->M = 0,
// and D->D = *.
// No other modifications are made.
func (hmm *HMM) Slice(start, end int) *HMM {
	nodes := make([]HMMNode, end-start)
	copy(nodes, hmm.Nodes[start:end])
	last := len(nodes) - 1
	nodes[last].Transitions.MM = 0
	nodes[last].Transitions.MI = MinProb
	nodes[last].Transitions.MD = MinProb
	nodes[last].Transitions.IM = 0
	nodes[last].Transitions.II = MinProb
	nodes[last].Transitions.DM = 0
	nodes[last].Transitions.DD = MinProb

	return &HMM{
		Nodes:    nodes,
		Alphabet: hmm.Alphabet,
		Null:     hmm.Null,
	}
}
