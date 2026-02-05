package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ruleTypeClassic = "classic"
	ruleTypeDomain  = "domain"
	ruleTypeIPCIDR  = "ipcidr"
)

func main() {
	var ruleType string
	var sourceURL string
	var outputPath string

	flag.StringVar(&ruleType, "t", "", "rule type: classic, domain, ipcidr")
	flag.StringVar(&sourceURL, "u", "", "source url")
	flag.StringVar(&outputPath, "o", "", "output file path")
	flag.Parse()

	if err := run(ruleType, sourceURL, outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(ruleType, sourceURL, outputPath string) error {
	if ruleType == "" || sourceURL == "" || outputPath == "" {
		usage := "usage: clash-skk -t classic|domain|ipcidr -u <url> -o <output>"
		return errors.New(usage)
	}

	if ruleType != ruleTypeClassic && ruleType != ruleTypeDomain && ruleType != ruleTypeIPCIDR {
		return fmt.Errorf("unknown rule type: %s", ruleType)
	}

	body, err := fetchURL(sourceURL)
	if err != nil {
		return err
	}

	lines, err := readLines(bytes.NewReader(body))
	if err != nil {
		return err
	}

	var output string
	switch ruleType {
	case ruleTypeClassic:
		header, payload := parseClassic(lines)
		output = buildYAML(header, payload, false)
	case ruleTypeDomain:
		header, payload := parseDomain(lines)
		output = buildYAML(header, payload, true)
	case ruleTypeIPCIDR:
		payload := parseIPCIDR(lines)
		output = buildYAML(nil, payload, true)
	}

	if err := ensureDir(outputPath); err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, []byte(output), 0o644); err != nil {
		return err
	}

	return nil
}

func fetchURL(sourceURL string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "clash-skk/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch url: unexpected status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return data, nil
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func parseClassic(lines []string) ([]string, []string) {
	return parseWithHeader(lines, func(line string) (string, bool) {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			return "", false
		}
		return trimmed, true
	})
}

func parseDomain(lines []string) ([]string, []string) {
	return parseWithHeader(lines, func(line string) (string, bool) {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			return "", false
		}
		if strings.HasPrefix(trimmed, "+.") {
			trimmed = "." + strings.TrimPrefix(trimmed, "+.")
		}
		return trimmed, true
	})
}

func parseIPCIDR(lines []string) []string {
	var payload []string
	for _, raw := range lines {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "ruleset.skk.moe") {
			continue
		}
		payload = append(payload, line)
	}
	return payload
}

func parseWithHeader(lines []string, transform func(string) (string, bool)) ([]string, []string) {
	var header []string
	var payload []string
	payloadStarted := false
	skippedTopSeparator := false

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !payloadStarted && len(header) > 0 {
				header = append(header, "")
			}
			continue
		}

		if strings.Contains(line, "ruleset.skk.moe") {
			continue
		}

		isComment := strings.HasPrefix(strings.TrimSpace(line), "#")
		if !payloadStarted && isComment {
			if !skippedTopSeparator && isSeparatorLine(line) && len(header) == 0 {
				skippedTopSeparator = true
				continue
			}
			if shouldSkipHeaderLine(line) {
				continue
			}
			header = append(header, line)
			continue
		}

		if !payloadStarted {
			payloadStarted = true
		}
		if isComment {
			continue
		}

		if out, ok := transform(line); ok {
			if strings.Contains(out, "ruleset.skk.moe") {
				continue
			}
			payload = append(payload, out)
		}
	}

	return header, payload
}

func buildYAML(header []string, payload []string, quote bool) string {
	var buf bytes.Buffer
	for _, line := range header {
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	buf.WriteString("payload:\n")
	for _, item := range payload {
		if quote {
			buf.WriteString("- '")
			buf.WriteString(escapeSingleQuotes(item))
			buf.WriteString("'\n")
			continue
		}
		buf.WriteString("- ")
		buf.WriteString(item)
		buf.WriteByte('\n')
	}
	return buf.String()
}

func escapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func shouldSkipHeaderLine(line string) bool {
	return strings.Contains(line, "License:") ||
		strings.Contains(line, "Homepage:") ||
		strings.Contains(line, "GitHub:")
}

func isSeparatorLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	for _, r := range trimmed {
		if r != '#' {
			return false
		}
	}
	return true
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	return nil
}
