package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSanitizerFlags(t *testing.T) {
	tests := []struct {
		name            string
		sanitizer       string
		wantCxxFlags    string
		wantLinkerFlags string
	}{
		{
			name:            "Empty sanitizer",
			sanitizer:       "",
			wantCxxFlags:    "",
			wantLinkerFlags: "",
		},
		{
			name:            "ASan",
			sanitizer:       "asan",
			wantCxxFlags:    " -fsanitize=address -fno-omit-frame-pointer",
			wantLinkerFlags: "-fsanitize=address",
		},
		{
			name:            "TSan",
			sanitizer:       "tsan",
			wantCxxFlags:    " -fsanitize=thread",
			wantLinkerFlags: "-fsanitize=thread",
		},
		{
			name:            "MSan",
			sanitizer:       "msan",
			wantCxxFlags:    " -fsanitize=memory -fno-omit-frame-pointer",
			wantLinkerFlags: "-fsanitize=memory",
		},
		{
			name:            "UBSan",
			sanitizer:       "ubsan",
			wantCxxFlags:    " -fsanitize=undefined",
			wantLinkerFlags: "-fsanitize=undefined",
		},
		{
			name:            "Unknown sanitizer",
			sanitizer:       "unknown",
			wantCxxFlags:    "",
			wantLinkerFlags: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCxx, gotLinker := GetSanitizerFlags(tt.sanitizer)
			assert.Equal(t, tt.wantCxxFlags, gotCxx)
			assert.Equal(t, tt.wantLinkerFlags, gotLinker)
		})
	}
}
