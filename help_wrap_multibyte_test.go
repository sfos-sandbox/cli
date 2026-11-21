package cli

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapMultiByte(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		offset int
		wrapAt int
		want   string
	}{
		{
			name:   "umlauts break at character width",
			input:  "app - ändert die Größe für Überprüfung",
			wrapAt: 30,
			want:   "app - ändert die Größe für\nÜberprüfung",
		},
		{
			name:   "line fitting in characters is not wrapped",
			input:  "ab  éé",
			wrapAt: 6,
			want:   "ab  éé",
		},
		{
			name:   "padding counted in characters",
			input:  "größe größe größe",
			offset: 2,
			wrapAt: 13,
			want:   "größe größe\n  größe",
		},
		{
			name:   "ascii unchanged",
			input:  "one two three four",
			wrapAt: 9,
			want:   "one two\nthree\nfour",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, wrap(tt.input, tt.offset, tt.wrapAt))
		})
	}
}

func TestWrappedHelpMultiByte(t *testing.T) {
	defer func(old HelpPrinterFunc) {
		HelpPrinter = old
	}(HelpPrinter)

	HelpPrinter = func(w io.Writer, templ string, data any) {
		HelpPrinterCustom(w, templ, data, map[string]any{
			"wrapAt": func() int { return 30 },
		})
	}

	output := new(bytes.Buffer)
	cmd := &Command{
		Writer:    output,
		Name:      "app",
		Usage:     "ändert die Größe für Überprüfung",
		UsageText: "ab  éé",
		HideHelp:  true,
	}

	_ = ShowRootCommandHelp(cmd)

	expected := `NAME:
   app - ändert die Größe für
         Überprüfung

USAGE:
   ab  éé
`

	assert.Equal(t, expected, output.String(), "Unexpected wrapping")
}
