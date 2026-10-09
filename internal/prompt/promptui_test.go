package prompt

import (
	"slices"
	"testing"
)

func TestSplitList(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: "k8s,cilium", want: []string{"k8s", "cilium"}},
		{input: " k8s , cilium ,", want: []string{"k8s", "cilium"}},
		{input: " , ,", want: nil},
		{input: "", want: nil},
	}

	for _, tt := range tests {
		if got := splitList(tt.input); !slices.Equal(got, tt.want) {
			t.Errorf("splitList(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRequireList(t *testing.T) {
	if err := requireList(" , "); err == nil {
		t.Error("requireList accepted a list with no values")
	}
	if err := requireList("k8s"); err != nil {
		t.Errorf("requireList(%q): %v", "k8s", err)
	}
}
