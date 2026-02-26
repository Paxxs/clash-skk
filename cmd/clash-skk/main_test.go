package main

import (
	"bytes"
	"testing"
)

func TestConvertDomainRuleToClassic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
		ok    bool
	}{
		{
			name:  "plus dot to domain suffix",
			input: "+.example.com",
			want:  "DOMAIN-SUFFIX,example.com",
			ok:    true,
		},
		{
			name:  "leading dot to domain wildcard",
			input: ".example.com",
			want:  "DOMAIN-WILDCARD,*.example.com",
			ok:    true,
		},
		{
			name:  "plain domain to domain",
			input: "example.com",
			want:  "DOMAIN,example.com",
			ok:    true,
		},
		{
			name:  "prefix wildcard dash to domain keyword",
			input: "*-events.adjust.com",
			want:  "DOMAIN-KEYWORD,-events.adjust.com",
			ok:    true,
		},
		{
			name:  "leading wildcard label to domain regex",
			input: "*.baidu.com",
			want:  "DOMAIN-REGEX,^[^.]+\\.baidu\\.com$",
			ok:    true,
		},
		{
			name:  "middle wildcard label to domain regex",
			input: "xbox.*.microsoft.com",
			want:  "DOMAIN-REGEX,^xbox\\.[^.]+\\.microsoft\\.com$",
			ok:    true,
		},
		{
			name:  "unsupported rule",
			input: "http://bad.example.com",
			want:  "",
			ok:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := convertDomainRuleToClassic(tt.input)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("converted rule = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseDomainAsClassicPreservesOrderAndCommentsUnsupported(t *testing.T) {
	t.Parallel()

	lines := []string{
		"# Header A",
		"# Header B",
		"",
		"+.example.com",
		"http://bad.example.com",
		"xbox.*.microsoft.com",
	}

	var warnBuf bytes.Buffer
	header, payload := parseDomainAsClassic(lines, &warnBuf)

	if len(header) != 3 {
		t.Fatalf("header length = %d, want 3", len(header))
	}
	if header[0] != "# Header A" || header[1] != "# Header B" {
		t.Fatalf("unexpected header: %#v", header)
	}
	if header[2] != "" {
		t.Fatalf("header blank separator missing: %#v", header)
	}

	wantPayload := []payloadItem{
		{Value: "DOMAIN-SUFFIX,example.com"},
		{Value: "# unsupported domain rule: http://bad.example.com", IsComment: true},
		{Value: "DOMAIN-REGEX,^xbox\\.[^.]+\\.microsoft\\.com$"},
	}

	if len(payload) != len(wantPayload) {
		t.Fatalf("payload length = %d, want %d", len(payload), len(wantPayload))
	}
	for i := range payload {
		if payload[i] != wantPayload[i] {
			t.Fatalf("payload[%d] = %#v, want %#v", i, payload[i], wantPayload[i])
		}
	}

	warn := warnBuf.String()
	if want := "warn: unsupported domain rule \"http://bad.example.com\""; !bytes.Contains([]byte(warn), []byte(want)) {
		t.Fatalf("warn log %q does not contain %q", warn, want)
	}

	gotYAML := buildYAMLItems(header, payload, false)
	wantYAML := "# Header A\n# Header B\n\npayload:\n- DOMAIN-SUFFIX,example.com\n# unsupported domain rule: http://bad.example.com\n- DOMAIN-REGEX,^xbox\\.[^.]+\\.microsoft\\.com$\n"
	if gotYAML != wantYAML {
		t.Fatalf("yaml output mismatch\n--- got ---\n%s\n--- want ---\n%s", gotYAML, wantYAML)
	}
}
