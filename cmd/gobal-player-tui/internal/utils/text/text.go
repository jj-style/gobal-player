package text

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

// formats a dictionary of key shortcuts as a help message to display to the screen
func FormatHelp(help map[string]string) string {
	elems := lo.MapToSlice(help, func(key string, desc string) string {
		return fmt.Sprintf("%s: %s", key, desc)
	})
	slices.Sort(elems)
	return strings.Join(elems, ", ")
}
