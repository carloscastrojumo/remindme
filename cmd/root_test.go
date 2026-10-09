package cmd

import (
	"strings"
	"testing"
)

func TestNeedsStorage(t *testing.T) {
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()

	tests := []struct {
		args []string
		want bool
	}{
		{args: []string{"list"}, want: true},
		{args: []string{"list", "tags"}, want: true},
		{args: []string{"help"}, want: false},
		{args: []string{"completion", "bash"}, want: false},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			cmd, _, err := rootCmd.Find(tt.args)
			if err != nil {
				t.Fatalf("Find(%v): %v", tt.args, err)
			}
			if got := needsStorage(cmd); got != tt.want {
				t.Errorf("needsStorage(%s) = %v, want %v", cmd.CommandPath(), got, tt.want)
			}
		})
	}
}
