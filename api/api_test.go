package wordnet

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wolfgangshilei/gown"
)

const dictDirName string = "../wn-dict"

func sameError(err1 error, err2 error) bool {
	if err1 == nil && err2 != nil {
		return false
	}

	if err1 != nil && err2 == nil {
		return false
	}

	if err1 == nil && err2 == nil {
		return true
	}

	if err1.Error() != err2.Error() {
		return false
	}

	return true
}

func ensureWordnetLoaded(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		dictDirName string
		err         error
	}{
		{
			dictDirName: dictDirName,
			err:         nil,
		},
		{
			dictDirName: "unknown",
			err:         fmt.Errorf("can't open unknown/index.noun: open unknown/index.noun: no such file or directory"),
		},
	}

	for _, tc := range testCases {
		_, err := load(tc.dictDirName)
		if !sameError(err, tc.err) {
			t.Fatalf("TestLoad failed with error.\n expected: %s\n actual: %s\n", tc.err, err)
		}
	}
}

func TestGetSynsetsWithLemma(t *testing.T) {
	wn, err := load(dictDirName)
	ensureWordnetLoaded(t, err)

	testCases := []struct {
		description string
		wn          *gown.WN
		lemma       string
		err         error
		ssNum       int
	}{
		{
			description: "Test wordnet is not loaded.",
			wn:          nil,
			lemma:       "test",
			err:         ErrWordnetNotLoaded,
		},
		{
			description: "Test lookup an existing word.",
			wn:          wn,
			lemma:       "test",
			ssNum:       13,
		},
		{
			description: "Test lookup a capitalized word.",
			wn:          wn,
			lemma:       "Test",
			ssNum:       13,
		},
	}

	for _, tc := range testCases {
		synsets, err := getSynsetsWithLemma(tc.wn, tc.lemma)
		if !sameError(err, tc.err) {
			t.Fatal("\nGot different error.\n", "Actual: ", err, "\nExpected: ", tc.err, "\n\n")
		}

		Convey(tc.description, t, func() {
			So(len(synsets), ShouldEqual, tc.ssNum)
		})
	}
}

func TestGetSynsetsWithLemmaAndPos(t *testing.T) {
	wn, err := load(dictDirName)
	ensureWordnetLoaded(t, err)

	testCases := []struct {
		description string
		wn          *gown.WN
		lemma       string
		pos         int
		err         error
		ssNum       int
	}{
		{
			description: "Test wordnet is not loaded.",
			wn:          nil,
			lemma:       "test",
			pos:         1,
			err:         ErrWordnetNotLoaded,
		},
		{
			description: "Test lookup an existing word for its synsets as a noun.",
			wn:          wn,
			lemma:       "test",
			pos:         1,
			ssNum:       6,
		},
		{
			description: "Test lookup a capitalized word.",
			wn:          wn,
			lemma:       "Test",
			pos:         1,
			ssNum:       6,
		},
	}

	for _, tc := range testCases {
		synsets, err := getSynsetsWithLemmaAndPos(tc.wn, tc.lemma, tc.pos)
		if !sameError(err, tc.err) {
			t.Fatal("\nGot different error.\n", "Actual: ", err, "\nExpected: ", tc.err, "\n\n")
		}

		Convey(tc.description, t, func() {
			So(len(synsets), ShouldEqual, tc.ssNum)
		})
	}
}

func TestMorph(t *testing.T) {
	wn, err := load(dictDirName)
	ensureWordnetLoaded(t, err)

	testCases := []struct {
		description           string
		wn                    *gown.WN
		origWord              string
		err                   error
		expectedMorphedPosMap map[string][]int
	}{
		{
			description: "Test wordnet is not loaded.",
			wn:          nil,
			err:         ErrWordnetNotLoaded,
		},
		{
			description: "Simply morph for word 'shone'.",
			wn:          wn,
			origWord:    "shone",
			expectedMorphedPosMap: map[string][]int{
				"shine": []int{2},
			},
		},
		{
			description: "Complex morph for word 'lives'.",
			wn:          wn,
			origWord:    "lives",
			expectedMorphedPosMap: map[string][]int{
				"life": []int{1},
				"live": []int{2},
			},
		},
		{
			description: "morph is case-sensitive, e.g. 'Lives' is different from 'lives'.",
			wn:          wn,
			origWord:    "Lives",
			expectedMorphedPosMap: map[string][]int{
				"Live": []int{2},
			},
		},
	}

	for _, tc := range testCases {
		morphedPosMap, err := morph(tc.wn, tc.origWord)
		if !sameError(err, tc.err) {
			t.Fatal("\nGot different error.\n", "Actual: ", err, "\nExpected: ", tc.err, "\n\n")
		}

		Convey(tc.description, t, func() {
			for k, v := range tc.expectedMorphedPosMap {
				poses, exist := morphedPosMap[k]
				So(exist, ShouldBeTrue)
				So(len(poses), ShouldEqual, len(v))
				for i, pos := range v {
					So(poses[i], ShouldEqual, pos)
				}
			}
		})
	}
}
