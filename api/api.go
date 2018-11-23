package wordnet

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/wolfgangshilei/gown"
)

var wordnet *gown.WN

type errorResponse struct {
	Error string `json:"error"`
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

var ErrWordnetNotLoaded = errors.New("Wordnet is not loaded.")

// Load loads the wordnet database into memory.
func Load(dirName string) string {
	if wordnet != nil {
		return makeJSONError(nil)
	}

	var (
		err error
		wn  *gown.WN
	)

	if wn, err = load(dirName); err == nil {
		wordnet = wn
	}
	return makeJSONError(err)
}

func load(dirName string) (*gown.WN, error) {
	// fmt.Printf("Loading wordnet %s\n", dirName)
	wn, err := gown.LoadWordNet(dirName)
	if err != nil {
		return nil, err
	}
	wn.InitMorphData(dirName)
	return wn, nil
}

// LookupWithPartOfSpeech returns a JSON-formatted string of the data index
// based on the input lemma and part of speech.
func LookupWithPartOfSpeech(lemma string, pos int) string {
	if wordnet == nil {
		return makeJSONError(ErrWordnetNotLoaded)
	}
	return makeJSONResponse(wordnet.LookupWithPartOfSpeech(lemma, pos))
}

// LookupSensesWithPartOfSpeech returns a JSON-formatted string of the sense
// index list based on the input lemma and part of speech.
func LookupSensesWithPartOfSpeech(lemma string, pos int) string {
	if wordnet == nil {
		return makeJSONError(ErrWordnetNotLoaded)
	}
	return makeJSONResponse(wordnet.LookupSensesWithPartOfSpeech(lemma, pos))
}

// LookupWithPartOfSpeechAndSense returns a JSON-formmatted string of the
// sense index entry based on the input lemma, part of speech and sense id.
func LookupWithPartOfSpeechAndSense(lemma string, pos int, senseID int) string {
	if wordnet == nil {
		return makeJSONError(ErrWordnetNotLoaded)
	}
	return makeJSONResponse(wordnet.LookupWithPartOfSpeechAndSense(lemma, pos, senseID))
}

// Lookup returns a JSON-formmatted string of the list of sense index entries of
// the input lemma.
func Lookup(lemma string) string {
	if wordnet == nil {
		return makeJSONError(ErrWordnetNotLoaded)
	}
	return makeJSONResponse(wordnet.Lookup(lemma))
}

// GetSynset returns a JSON-formatted string of the synset(data) in database.
func GetSynset(pos int, synsetOffset int) string {
	if wordnet == nil {
		return makeJSONError(ErrWordnetNotLoaded)
	}
	return makeJSONResponse(wordnet.GetSynset(pos, synsetOffset))
}

// GetSynsetsWithLemma returns a JSON-formatted string which contains
// a list of synsets of given lemma.
func GetSynsetsWithLemma(lemma string) string {
	if synsets, err := getSynsetsWithLemma(wordnet, lemma); err != nil {
		return makeJSONError(err)
	} else {
		return makeJSONResponse(synsets)
	}
}

func getSynsetsWithLemma(wn *gown.WN, lemma string) (synsets []*gown.Synset, err error) {
	if wn == nil {
		err = ErrWordnetNotLoaded
		return
	}

	for _, pos := range allPos() {
		senses := wn.LookupSensesWithPartOfSpeech(strings.ToLower(lemma), pos)
		for _, s := range senses {
			synsets = append(synsets, s.GetSynsetPtr())
		}
	}
	return
}

// GetSynsetsWithLemmaAndPos returns a JSON-formatted string which contains
// a list of synsets of given lemma and its part of speech.
func GetSynsetsWithLemmaAndPos(lemma string, pos int) string {
	if synsets, err := getSynsetsWithLemmaAndPos(wordnet, lemma, pos); err != nil {
		return makeJSONError(err)
	} else {
		return makeJSONResponse(synsets)
	}
}

func getSynsetsWithLemmaAndPos(wn *gown.WN, lemma string, pos int) (synsets []*gown.Synset, err error) {
	if wn == nil {
		err = ErrWordnetNotLoaded
		return
	}

	senses := wn.LookupSensesWithPartOfSpeech(strings.ToLower(lemma), pos)
	for _, s := range senses {
		synsets = append(synsets, s.GetSynsetPtr())
	}
	return
}

// Morph returns a JSON-formatted string of a list of possible words which
// are a result of the original word being morphologically processed.
func Morph(origWord string) string {
	if morphedPosMap, err := morph(wordnet, origWord); err != nil {
		return makeJSONError(err)
	} else {
		return makeJSONResponse(morphedPosMap)
	}
}

func morph(wn *gown.WN, origWord string) (morphedPosMap map[string][]int, err error) {
	if wn == nil {
		err = ErrWordnetNotLoaded
		return
	}

	morphedPosMap = map[string][]int{}
	for _, pos := range allPos() {
		if morphed := wn.Morph(origWord, pos); morphed != "" {
			poses := morphedPosMap[morphed]
			morphedPosMap[morphed] = append(poses, pos)
		}
	}
	return
}

func makeJSONError(err error) string {
	errString := ""

	if err != nil {
		errString = err.Error()
	}

	res := errorResponse{
		Error: errString,
	}

	jres, _ := json.Marshal(res)
	return string(jres)
}

func makeJSONResponse(response interface{}) string {
	res := dataResponse{
		Data: response,
	}

	jres, _ := json.Marshal(res)
	return string(jres)
}

func allPos() []int {
	return []int{
		gown.POS_NOUN,
		gown.POS_VERB,
		gown.POS_ADJECTIVE,
		gown.POS_ADVERB,
		gown.POS_ADJECTIVE_SATELLITE,
	}
}
