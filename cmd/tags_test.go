package cmd

import (
	"slices"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestTagsFlagSplitsCommas(t *testing.T) {
	for _, c := range []*cobra.Command{addCmd, listCmd, removeCmd} {
		t.Run(c.Name(), func(t *testing.T) {
			flags := c.Flags()
			t.Cleanup(func() {
				flags.Lookup("tags").Value.(pflag.SliceValue).Replace(nil)
				flags.Lookup("tags").Changed = false
			})

			for _, value := range []string{"k8s,cilium", "net"} {
				if err := flags.Set("tags", value); err != nil {
					t.Fatalf("Set(tags, %q): %v", value, err)
				}
			}

			got, err := flags.GetStringSlice("tags")
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"k8s", "cilium", "net"}; !slices.Equal(got, want) {
				t.Errorf("tags = %q, want %q", got, want)
			}
		})
	}
}
