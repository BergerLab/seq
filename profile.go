package seq

import (
	"bytes"
	"fmt"
	"math"
	"text/tabwriter"
)

// Profile represents a sequence profile in terms of log-odds scores.
type Profile struct {
	// The columns of a profile.
	Emissions []EProbs

	// The alphabet of the profile. The length of the alphabet should be
	// equal to the number of rows in the profile.
	// There are no restrictions on the alphabet. (i.e., Gap characters are
	// allowed but they are not treated specially.)
	Alphabet Alphabet
}

// NewProfile initializes a profile with a default
// alphabet that is compatible with this package's BLOSUM62 matrix.
// Emission probabilities are set to the minimum log-odds probability.
func NewProfile(columns int) *Profile {
	return NewProfileAlphabet(columns, AlphaBlosum62)
}

// NewProfileAlphabet initializes a profile with the given alphabet.
// Emission probabilities are set to the minimum log-odds probability.
func NewProfileAlphabet(columns int, alphabet Alphabet) *Profile {
	emits := make([]EProbs, columns)
	for i := 0; i < columns; i++ {
		emits[i] = NewEProbs(alphabet)
	}
	return &Profile{emits, alphabet}
}

func (p *Profile) Len() int {
	return len(p.Emissions)
}

func (p *Profile) String() string {
	buf := new(bytes.Buffer)
	tabw := tabwriter.NewWriter(buf, 4, 0, 3, ' ', 0)
	pf := func(ft string, v ...interface{}) { fmt.Fprintf(tabw, ft, v...) }
	for _, r := range p.Alphabet {
		pf("%c", rune(r))
		for _, eprobs := range p.Emissions {
			pf("\t%0.4f", eprobs[r])
		}
		pf("\n")
	}
	tabw.Flush()
	return buf.String()
}

// FrequencyProfile represents a sequence profile in terms of raw frequencies.
// A FrequencyProfile is useful as an intermediate representation. It can be
// used to incrementally build a Profile.
type FrequencyProfile struct {
	// The columns of a frequency profile.
	Freqs []map[Residue]int

	// The alphabet of the profile. The length of the alphabet should be
	// equal to the number of rows in the frequency profile.
	// There are no restrictions on the alphabet. (i.e., Gap characters are
	// allowed but they are not treated specially.)
	Alphabet Alphabet
}

func (fp *FrequencyProfile) String() string {
	buf := new(bytes.Buffer)
	tabw := tabwriter.NewWriter(buf, 4, 0, 3, ' ', 0)
	pf := func(ft string, v ...interface{}) { fmt.Fprintf(tabw, ft, v...) }
	for _, r := range fp.Alphabet {
		pf("%c", rune(r))
		for _, column := range fp.Freqs {
			pf("\t%d", column[r])
		}
		pf("\n")
	}
	tabw.Flush()
	return buf.String()
}

// NewNullProfile initializes a frequency profile that can be used to tabulate
// a null model. This is equivalent to calling NewFrequencyProfile with the
// number of columns set to 1.
func NewNullProfile() *FrequencyProfile {
	return NewFrequencyProfile(1)
}

// NewFrequencyProfile initializes a frequency profile with a default
// alphabet that is compatible with this package's BLOSUM62 matrix.
func NewFrequencyProfile(columns int) *FrequencyProfile {
	return NewFrequencyProfileAlphabet(columns, AlphaBlosum62)
}

// NewFrequencyProfileAlphabet initializes a frequency profile with the
// given alphabet.
func NewFrequencyProfileAlphabet(
	columns int,
	alphabet Alphabet,
) *FrequencyProfile {
	freqs := make([]map[Residue]int, columns)
	for i := 0; i < columns; i++ {
		freqs[i] = make(map[Residue]int, len(alphabet))
		for _, residue := range alphabet {
			freqs[i][residue] = 0
		}
	}
	return &FrequencyProfile{freqs, alphabet}
}

// Len returns the number of columns in the frequency profile.
func (fp *FrequencyProfile) Len() int {
	return len(fp.Freqs)
}

// Add adds the sequence to the given profile. The sequence must have length
// equivalent to the number of columns in the profile. The sequence must also
// only contain residues that are in the alphabet for the profile.
//
// As a special case, if the alphabet contains the 'X' residue, then any
// unrecognized residues in the sequence with respect to the profile's alphabet
// will be considered as an 'X' residue.
func (fp *FrequencyProfile) Add(s Sequence) {
	if fp.Len() != s.Len() {
		panic(fmt.Sprintf("Profile has length %d but sequence has length %d",
			fp.Len(), s.Len()))
	}
	for column := 0; column < fp.Len(); column++ {
		r := s.Residues[column]
		if _, ok := fp.Freqs[column][r]; ok {
			fp.Freqs[column][r] += 1
		} else if _, ok := fp.Freqs[column]['X']; ok {
			fp.Freqs[column]['X'] += 1
		} else {
			panic(fmt.Sprintf("Unrecognized residue %c while using an "+
				"alphabet without a wildcard: '%s'.", r, fp.Alphabet))
		}
	}
}

// Profile converts a raw frequency profile to a profile that uses a log-odds
// representation. The log-odds scores are computed with the given null model,
// which is itself just a raw frequency profile with a single column.
func (fp *FrequencyProfile) Profile(null *FrequencyProfile) *Profile {
	if null.Len() != 1 {
		panic(fmt.Sprintf("null model has %d columns; should have 1",
			null.Len()))
	}
	if !fp.Alphabet.Equals(null.Alphabet) {
		panic(fmt.Sprintf("freq profile alphabet '%s' is not equal to "+
			"null profile alphabet '%s'.", fp.Alphabet, null.Alphabet))
	}
	p := NewProfileAlphabet(fp.Len(), fp.Alphabet)

	// Compute the background emission probabilities.
	nulltot := freqTotal(null.Freqs[0])
	nullemit := make(map[Residue]float64, fp.Alphabet.Len())
	for _, residue := range null.Alphabet {
		nullemit[residue] = float64(null.Freqs[0][residue]) / float64(nulltot)
	}

	// Now compute the emission probabilities and convert to log-odds.
	for column := 0; column < fp.Len(); column++ {
		tot := freqTotal(fp.Freqs[column])
		for _, residue := range fp.Alphabet {
			if null.Freqs[0][residue] == 0 || fp.Freqs[column][residue] == 0 {
				p.Emissions[column][residue] = MinProb
			} else {
				prob := float64(fp.Freqs[column][residue]) / float64(tot)
				logOdds := -Prob(math.Log(prob / nullemit[residue]))
				p.Emissions[column][residue] = logOdds
			}
		}
	}
	return p
}

// freqTotal computes the total frequency in a single column.
func freqTotal(column map[Residue]int) int {
	tot := 0
	for _, freq := range column {
		tot += freq
	}
	return tot
}
