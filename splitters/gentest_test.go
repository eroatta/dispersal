package splitters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit_OnGenTest_ShouldReturnValidSplits(t *testing.T) {
	cases := []struct {
		ID       int
		token    string
		expected []string
	}{
		{0, "car", []string{"car"}},
		{1, "getString", []string{"get", "String"}},
		{2, "GPSstate", []string{"GPS", "state"}},
		{3, "ASTVisitor", []string{"AST", "Visitor"}},
		{4, "notype", []string{"no", "type"}},
	}

	gentest := NewGenTest()
	for _, c := range cases {
		got, err := gentest.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func TestGeneratePotentialSplits_ShouldReturnEveryPossibleCombination(t *testing.T) {
	cases := []struct {
		token    string
		expected []string
	}{
		{"car", []string{"car", "c_ar", "c_a_r", "ca_r"}},
		{"bond", []string{"bond", "b_ond", "b_o_nd", "b_on_d", "bo_nd", "bo_n_d", "bon_d"}},
	}

	for _, c := range cases {
		var got []string
		for _, potentialSplit := range generatePotentialSplits(c.token) {
			got = append(got, potentialSplit.split)
		}

		assert.ElementsMatch(t, c.expected, got, "elements should match")
	}
}

func TestFindExpansions_OnGenTestWithCustomList_ShouldReturnAllMatches(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expansions []string
	}{
		{"input_st", "st", []string{"string", "steer", "set"}},
		{"input_rlen", "rlen", []string{"riflemen"}},
		{"input_str", "str", []string{"steer", "string"}},
		{"input_len", "len", []string{"lender", "length"}},
		{"empty_input", "", []string{}},
		{"blankspace_input", " ", []string{}},
	}

	gentest := NewGenTest()
	gentest.list = "car string steer set riflemen lender bar length kamikaze"
	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			got := gentest.findExpansions(fixture.input)

			assert.ElementsMatch(t, fixture.expansions, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestNewPotentialSplit_OnMarkedHardword_ShouldReturnPotentialSplit(t *testing.T) {
	got := NewPotentialSplit("foo_bar")

	assert.Equal(t, "foo_bar", got.split, "split should match de input")
	assert.ElementsMatch(t, []string{"foo", "bar"}, got.softwords, "elements should match")
	assert.Equal(t, 0, len(got.expansions), "expansions map should be empty")
}

func TestNewPotentialSplit_OnEmptyHardword_ShouldReturnEmptyPotentialSplit(t *testing.T) {
	got := NewPotentialSplit("")

	assert.Equal(t, "", got.split, "split should be empty")
	assert.ElementsMatch(t, []string{}, got.softwords, "there should be no softwords")
	assert.Equal(t, 0, len(got.expansions), "there should be no elements")
}

func BenchmarkGenerate(b *testing.B) {
	benchmarks := []struct {
		name  string
		token string
	}{
		{"Short token", "car"},
		{"Medium token", "numsize"},
		{"Long token", "allocatedsize"},
		{"Longest token", "veryverylongtokennameforsplitting"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				generatePotentialSplits(bm.token)
			}
		})
	}
}
