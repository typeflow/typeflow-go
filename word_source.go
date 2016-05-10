package typeflow

import (
    "github.com/alediaferia/prefixmap"
    "strings"
    "regexp"
)

type WordSource struct {
    prefixMap *prefixmap.PrefixMap
    wc   int
}

const allowedKeyChars = `[^a-zA-Z0-9]+`
var allowedKeyCharsRegexp *regexp.Regexp

func init() {
    allowedKeyCharsRegexp = regexp.MustCompile(allowedKeyChars)
}

// Represents a match
// found
type Match struct {
    // The actual string which matched
    // the requested string
    String string `json:"string"`

    // The similarity value for the
    // current word.
    // Similarity is 1 when
    // the two words are equal.
    // See github.com/typeflow/typeflow-go for
    // more details on how the similarity is computed.
    Similarity float32 `json:"similarity"`
}

func computeSimilarity(lenw1, lenw2, ld int) float32 {
    den := maximum(lenw1, lenw2)

    return 1.0 - float32(ld)/float32(den)
}

// NewWordSource Initializes a new empty WordSource
func NewWordSource() (ws *WordSource) {
    ws = new(WordSource)
    ws.prefixMap = prefixmap.New()
    ws.wc = 0
    return
}

// SetSource sets the given strings slice as the
// current source after applying the given
// filters in order
func (ws *WordSource) SetSource(strs []string) {    
    for _, s := range strs {
        keys := allowedKeyCharsRegexp.Split(strings.TrimSpace(s), -1)
        for _, k := range keys {
            ws.prefixMap.Insert(k, s)
        }
    }
}

type dirty_range struct {
    low    int
    length int
}

type Matches []Match

func (matches Matches) Len() int {
    return len(matches)
}

func (matches Matches) Less(i, j int) bool {
    return matches[i].Similarity > matches[j].Similarity
}

func (matches Matches) Swap(i, j int) {
    mi := matches[i]
    mj := matches[j]
    matches[i] = mj
    matches[j] = mi
}

// FindMatch Finds a match among the current words
// in the source.
// minSimilarity is the minimum accepted similarity
// to use when filling the matches slice.
// Param substr is the string to match against.
func (ws *WordSource) FindMatch(substr string, minSimilarity float32) (matches []Match, err error) {
    matches = make([]Match, 0, ws.wc)
    err = nil
    
    values := ws.prefixMap.GetByPrefix(substr)
    for _, v := range values {
        value := v.(string)
        similarity := computeSimilarity(len(value), len(substr), LevenshteinDistance(value, substr))
        if similarity >= minSimilarity {
            matches = append(matches, Match{value, similarity})
        }
    }
    
    return matches, err
}
