# Clash-SKK

将 Sukka Ruleset 的 txt 规则自动转换为 Mihomo (Clash.Meta) 可直接引用的 YAML rule-provider。

## 背景

Sukka Ruleset 只提供 txt 版本规则。subconverter 转换时不会为 rule provider 自动补上 `format: text`，导致 Clash Verge 将其误判为 YAML 并报错。本项目会把 txt 规则转换为 YAML（带 `payload`），因此可以直接引用，无需再写 `type=txt` 之类的字段。

## 自动更新

- 数据来源：`https://ruleset.skk.moe`
- 输出路径：`Clash/<category>/<name>.yaml`
- 更新频率：GitHub Actions 每日两次
- CDN 加速：`https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main`

## Mihomo 规则表（本项目对应路径）

> 本项目输出均为 YAML；如需显式声明，可使用 `format: yaml`。

| 目录               | Mihomo behavior         | 说明                            |
| ------------------ | ----------------------- | ------------------------------- |
| `Clash/domainset/` | `domain`                | 仅域名的规则组，不触发 DNS 解析 |
| `Clash/non_ip/`    | `classical`             | 非 IP 的规则组，不触发 DNS 解析 |
| `Clash/ip/`        | `classical` 或 `ipcidr` | IP 相关规则组，会触发 DNS 解析  |

> 按 `domainset`、`non_ip`、`ip` 的顺序引入规则组，避免不必要的 DNS 解析。

## 命令行使用

```bash
clash-skk -t classic -u https://ruleset.skk.moe/Clash/non_ip/reject-no-drop.txt -o Clash/non_ip/reject-no-drop.yaml
clash-skk -t domain  -u https://ruleset.skk.moe/Clash/domainset/reject.txt -o Clash/domainset/reject.yaml
clash-skk -t ipcidr  -u https://ruleset.skk.moe/Clash/ip/china_ip.txt -o Clash/ip/china_ip.yaml
```

## 配置示例（Mihomo / Clash.Meta）

> 以下示例按 Sukka Ruleset README 的条目整理，并将 URL 替换为本项目的 jsDelivr 地址。

### 广告拦截 / 隐私保护 / Malware 拦截 / Phishing 拦截

```yaml
rule-providers:
  reject_non_ip_no_drop:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/reject-no-drop.yaml
    path: ./sukkaw_ruleset/reject_non_ip_no_drop.yaml
  reject_non_ip_drop:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/reject-drop.yaml
    path: ./sukkaw_ruleset/reject_non_ip_drop.yaml
  reject_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/reject.yaml
    path: ./sukkaw_ruleset/reject_non_ip.yaml
  reject_domainset:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/reject.yaml
    path: ./sukkaw_ruleset/reject_domainset.yaml
  reject_extra_domainset:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/reject_extra.yaml
    path: ./sukkaw_ruleset/reject_extra_domainset.yaml
  reject_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/reject.yaml
    path: ./sukkaw_ruleset/reject_ip.yaml

rules:
  - RULE-SET,reject_non_ip_drop,REJECT-DROP
  - RULE-SET,reject_domainset,REJECT
  - RULE-SET,reject_extra_domainset,REJECT
  - RULE-SET,reject_non_ip,REJECT
  - RULE-SET,reject_non_ip_no_drop,REJECT
  - RULE-SET,reject_ip,REJECT
```

### 搜狗输入法

```yaml
rule-providers:
  sogouinput:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/sogouinput.yaml
    path: ./sukkaw_ruleset/sogouinput.yaml

rules:
  - RULE-SET,sogouinput,REJECT
```

### Speedtest 测速域名

```yaml
rule-providers:
  speedtest:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/speedtest.yaml
    path: ./sukkaw_ruleset/speedtest.yaml

rules:
  - RULE-SET,speedtest,[Replace with your policy]
```

### 常见静态 CDN

```yaml
rule-providers:
  cdn_domainset:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/cdn.yaml
    path: ./sukkaw_ruleset/cdn_domainset.yaml
  cdn_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/cdn.yaml
    path: ./sukkaw_ruleset/cdn_non_ip.yaml

rules:
  - RULE-SET,cdn_domainset,[Replace with your policy]
  - RULE-SET,cdn_non_ip,[Replace with your policy]
```

### 流媒体

```yaml
rule-providers:
  stream_us_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_us.yaml
    path: ./sukkaw_ruleset/stream_us_non_ip.yaml
  stream_us_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_us.yaml
    path: ./sukkaw_ruleset/stream_us_ip.yaml
  stream_eu_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_eu.yaml
    path: ./sukkaw_ruleset/stream_eu_non_ip.yaml
  stream_eu_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_eu.yaml
    path: ./sukkaw_ruleset/stream_eu_ip.yaml
  stream_jp_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_jp.yaml
    path: ./sukkaw_ruleset/stream_jp_non_ip.yaml
  stream_jp_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_jp.yaml
    path: ./sukkaw_ruleset/stream_jp_ip.yaml
  stream_kr_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_kr.yaml
    path: ./sukkaw_ruleset/stream_kr_non_ip.yaml
  stream_kr_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_kr.yaml
    path: ./sukkaw_ruleset/stream_kr_ip.yaml
  stream_hk_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_hk.yaml
    path: ./sukkaw_ruleset/stream_hk_non_ip.yaml
  stream_hk_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_hk.yaml
    path: ./sukkaw_ruleset/stream_hk_ip.yaml
  stream_tw_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream_tw.yaml
    path: ./sukkaw_ruleset/stream_tw_non_ip.yaml
  stream_tw_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream_tw.yaml
    path: ./sukkaw_ruleset/stream_tw_ip.yaml
  stream_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/stream.yaml
    path: ./sukkaw_ruleset/stream_non_ip.yaml
  stream_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/stream.yaml
    path: ./sukkaw_ruleset/stream_ip.yaml

rules:
  - RULE-SET,stream_us_non_ip,[Replace with your policy]
  - RULE-SET,stream_eu_non_ip,[Replace with your policy]
  - RULE-SET,stream_jp_non_ip,[Replace with your policy]
  - RULE-SET,stream_kr_non_ip,[Replace with your policy]
  - RULE-SET,stream_hk_non_ip,[Replace with your policy]
  - RULE-SET,stream_tw_non_ip,[Replace with your policy]
  - RULE-SET,stream_non_ip,[Replace with your policy]
  - RULE-SET,stream_us_ip,[Replace with your policy]
  - RULE-SET,stream_eu_ip,[Replace with your policy]
  - RULE-SET,stream_jp_ip,[Replace with your policy]
  - RULE-SET,stream_kr_ip,[Replace with your policy]
  - RULE-SET,stream_hk_ip,[Replace with your policy]
  - RULE-SET,stream_tw_ip,[Replace with your policy]
  - RULE-SET,stream_ip,[Replace with your policy]
```

### AI

```yaml
rule-providers:
  ai_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/ai.yaml
    path: ./sukkaw_ruleset/ai_non_ip.yaml
  apple_intelligence_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/apple_intelligence.yaml
    path: ./sukkaw_ruleset/apple_intelligence_non_ip.yaml

rules:
  - RULE-SET,ai_non_ip,[Replace with your policy]
  - RULE-SET,apple_intelligence_non_ip,[Replace with your policy]
```

### Telegram

```yaml
rule-providers:
  telegram_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/telegram.yaml
    path: ./sukkaw_ruleset/telegram_non_ip.yaml
  telegram_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/telegram.yaml
    path: ./sukkaw_ruleset/telegram_ip.yaml

rules:
  - RULE-SET,telegram_non_ip,[Replace with your policy]
  - RULE-SET,telegram_ip,[Replace with your policy]
```

### Apple CDN

```yaml
rule-providers:
  apple_cdn:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/apple_cdn.yaml
    path: ./sukkaw_ruleset/apple_cdn.yaml

rules:
  - RULE-SET,apple_cdn,[Replace with your policy]
```

### Apple Service

```yaml
rule-providers:
  apple_services:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/apple_services.yaml
    path: ./sukkaw_ruleset/apple_services.yaml

rules:
  - RULE-SET,apple_services,[Replace with your policy]
```

### Apple CN

```yaml
rule-providers:
  apple_cn_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/apple_cn.yaml
    path: ./sukkaw_ruleset/apple_cn_non_ip.yaml

rules:
  - RULE-SET,apple_cn_non_ip,[Replace with your policy]
```

### Microsoft CDN

```yaml
rule-providers:
  microsoft_cdn_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/microsoft_cdn.yaml
    path: ./sukkaw_ruleset/microsoft_cdn_non_ip.yaml

rules:
  - RULE-SET,microsoft_cdn_non_ip,[Replace with your policy]
```

### Microsoft

```yaml
rule-providers:
  microsoft_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/microsoft.yaml
    path: ./sukkaw_ruleset/microsoft_non_ip.yaml

rules:
  - RULE-SET,microsoft_non_ip,[Replace with your policy]
```

### 网易云音乐

```yaml
rule-providers:
  neteasemusic_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/neteasemusic.yaml
    path: ./sukkaw_ruleset/neteasemusic_non_ip.yaml
  neteasemusic_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/neteasemusic.yaml
    path: ./sukkaw_ruleset/neteasemusic_ip.yaml

rules:
  - RULE-SET,neteasemusic_non_ip,[Replace with your policy]
  - RULE-SET,neteasemusic_ip,[Replace with your policy]
```

### 软件更新、操作系统等大文件下载

```yaml
rule-providers:
  download_domainset:
    type: http
    behavior: domain
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/domainset/download.yaml
    path: ./sukkaw_ruleset/download_domainset.yaml
  download_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/download.yaml
    path: ./sukkaw_ruleset/download_non_ip.yaml

rules:
  - RULE-SET,download_domainset,[Replace with your policy]
  - RULE-SET,download_non_ip,[Replace with your policy]
```

### 内网域名和局域网 IP

```yaml
rule-providers:
  lan_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/lan.yaml
    path: ./sukkaw_ruleset/lan_non_ip.yaml
  lan_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/lan.yaml
    path: ./sukkaw_ruleset/lan_ip.yaml

rules:
  - RULE-SET,lan_non_ip,DIRECT
  - RULE-SET,lan_ip,DIRECT
```

### Misc

```yaml
rule-providers:
  domestic_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/domestic.yaml
    path: ./sukkaw_ruleset/domestic_non_ip.yaml
  direct_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/direct.yaml
    path: ./sukkaw_ruleset/direct_non_ip.yaml
  global_non_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/non_ip/global.yaml
    path: ./sukkaw_ruleset/global_non_ip.yaml
  domestic_ip:
    type: http
    behavior: classical
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/domestic.yaml
    path: ./sukkaw_ruleset/domestic_ip.yaml

rules:
  - RULE-SET,domestic_non_ip,[Replace with your policy]
  - RULE-SET,direct_non_ip,[Replace with your policy]
  - RULE-SET,global_non_ip,[Replace with your policy]
  - RULE-SET,domestic_ip,[Replace with your policy]
```

### chnroute CIDR

```yaml
rule-providers:
  china_ip:
    type: http
    behavior: ipcidr
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/china_ip.yaml
    path: ./sukkaw_ruleset/china_ip.yaml
  china_ip_ipv6:
    type: http
    behavior: ipcidr
    interval: 43200
    url: https://cdn.jsdelivr.net/gh/Paxxs/clash-skk@main/Clash/ip/china_ip_ipv6.yaml
    path: ./sukkaw_ruleset/china_ip_ipv6.yaml

rules:
  - RULE-SET,china_ip,[Replace with your policy]
  # Only use it if you are using IPv6
  # - RULE-SET,china_ip_ipv6,[Replace with your policy]
```
