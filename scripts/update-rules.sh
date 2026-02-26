#!/usr/bin/env bash
set -euo pipefail

ROOT_URL="${ROOT_URL:-https://ruleset.skk.moe}"
OUT_DIR="${OUT_DIR:-.}"
BIN_PATH="${BIN_PATH:-./bin/clash-skk}"

links=$(python3 - "$ROOT_URL" <<'PY'
import sys
import urllib.request
from html.parser import HTMLParser

if len(sys.argv) < 2:
    raise SystemExit("missing base url")
base = sys.argv[1]
print(f"Fetching links from {base}...", file=sys.stderr)
req = urllib.request.Request(
    base,
    headers={
        "User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
        "Accept": "text/html,application/xhtml+xml",
    },
)
html = urllib.request.urlopen(req, timeout=30).read().decode("utf-8", errors="ignore")
links = []

class Parser(HTMLParser):
    def handle_starttag(self, tag, attrs):
        if tag != "a":
            return
        href = dict(attrs).get("href", "")
        if href.startswith("/Clash/") and href.endswith(".txt"):
            links.append(href)

Parser().feed(html)
for href in sorted(set(links)):
    print(href)
PY
)

while IFS= read -r path; do
  make_domain_classical_copy=0
  case "$path" in
    /Clash/domainset/*)
      rule_type="domain"
      make_domain_classical_copy=1
      ;;
    /Clash/non_ip/*)
      rule_type="classic"
      ;;
    /Clash/ip/*)
      rule_type="ipcidr"
      ;;
    *)
      continue
      ;;
  esac
  
  echo "${ROOT_URL}${path}"
  output_path="${OUT_DIR%/}${path%.txt}.yaml"
  echo $output_path
  mkdir -p "$(dirname "$output_path")"
  "$BIN_PATH" -t "$rule_type" -u "${ROOT_URL}${path}" -o "$output_path"

  # For Shadowrocket compatibility, domainset rules need an extra classical copy.
  if [[ "$make_domain_classical_copy" == "1" ]]; then
    name="$(basename "${path%.txt}")"
    classical_output_path="${OUT_DIR%/}/Clash/non_ip/${name}_classical.yaml"
    echo "$classical_output_path"
    mkdir -p "$(dirname "$classical_output_path")"
    "$BIN_PATH" -t "domain-classical" -u "${ROOT_URL}${path}" -o "$classical_output_path"
  fi
done <<< "$links"
