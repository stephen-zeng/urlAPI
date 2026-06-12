#!/usr/bin/env sh

set -u

BASE_URL="${BASE_URL:-https://api.0x535a.cn}"
ALLOWED_REFERER="${ALLOWED_REFERER:-https://0x535a.cn/}"
DISALLOWED_REFERER="${DISALLOWED_REFERER:-https://evil.example/}"

PASS_COUNT=0
FAIL_COUNT=0

print_header() {
  printf '\n== %s ==\n' "$1"
}

record_result() {
  name="$1"
  expected="$2"
  actual="$3"
  body="$4"

  if [ "$actual" = "$expected" ]; then
    PASS_COUNT=$((PASS_COUNT + 1))
    printf '[通过] %s：HTTP %s\n' "$name" "$actual"
  else
    FAIL_COUNT=$((FAIL_COUNT + 1))
    printf '[失败] %s：期望 HTTP %s，实际 HTTP %s\n' "$name" "$expected" "$actual"
    printf '       响应内容：%s\n' "$body"
  fi
}

request() {
  method="$1"
  url="$2"
  referer="$3"

  tmp="$(mktemp)"
  if [ -n "$referer" ]; then
    status="$(curl -sS -L -o "$tmp" -w '%{http_code}' -X "$method" -H "Referer: $referer" "$url" 2>/dev/null || printf '000')"
  else
    status="$(curl -sS -L -o "$tmp" -w '%{http_code}' -X "$method" "$url" 2>/dev/null || printf '000')"
  fi
  body="$(tr '\n' ' ' < "$tmp")"
  rm -f "$tmp"
  printf '%s\n%s\n' "$status" "$body"
}

expect_status() {
  name="$1"
  expected="$2"
  method="$3"
  url="$4"
  referer="$5"

  result="$(request "$method" "$url" "$referer")"
  status="$(printf '%s' "$result" | sed -n '1p')"
  body="$(printf '%s' "$result" | sed -n '2p')"
  record_result "$name" "$expected" "$status" "$body"
}

expect_forbidden_like() {
  name="$1"
  method="$2"
  url="$3"
  referer="$4"

  result="$(request "$method" "$url" "$referer")"
  status="$(printf '%s' "$result" | sed -n '1p')"
  body="$(printf '%s' "$result" | sed -n '2p')"

  case "$status" in
    400|401|403|429)
      PASS_COUNT=$((PASS_COUNT + 1))
      printf '[通过] %s：已被拒绝，HTTP %s\n' "$name" "$status"
      ;;
    *)
      FAIL_COUNT=$((FAIL_COUNT + 1))
      printf '[失败] %s：期望被拒绝 HTTP 400/401/403/429，实际 HTTP %s\n' "$name" "$status"
      printf '       响应内容：%s\n' "$body"
      ;;
  esac
}

print_header "配置"
printf 'BASE_URL=%s\n' "$BASE_URL"
printf 'ALLOWED_REFERER=%s\n' "$ALLOWED_REFERER"
printf 'DISALLOWED_REFERER=%s\n' "$DISALLOWED_REFERER"

print_header "风控验证：Referer 白名单"
expect_status \
  "拒绝空 Referer" \
  "403" \
  "GET" \
  "$BASE_URL/txt?prompt=poem&api=openai&format=json" \
  ""

expect_status \
  "拒绝非白名单 Referer" \
  "403" \
  "GET" \
  "$BASE_URL/txt?prompt=poem&api=openai&format=json" \
  "$DISALLOWED_REFERER"

print_header "风控验证：Prompt 与 API 参数"
expect_forbidden_like \
  "拒绝非法文字生成供应商" \
  "GET" \
  "$BASE_URL/txt?prompt=poem&api=invalid-provider&format=json" \
  "$ALLOWED_REFERER"

expect_forbidden_like \
  "拒绝非法图片生成供应商" \
  "GET" \
  "$BASE_URL/img?prompt=cat&api=invalid-provider&size=1024x1024&format=json" \
  "$ALLOWED_REFERER"

expect_forbidden_like \
  "拒绝非法随机图片供应商" \
  "GET" \
  "$BASE_URL/rand?api=invalid-provider&user=test&repo=test&format=json" \
  "$ALLOWED_REFERER"

expect_forbidden_like \
  "拒绝非白名单网页目标站点" \
  "GET" \
  "$BASE_URL/web?img=https://not-allowed.invalid/page&format=json" \
  "$ALLOWED_REFERER"

print_header "流控验证：并发突发请求"
RATE_LIMIT_HIT=0
RATE_LIMIT_LAST_STATUS=""
RATE_LIMIT_LAST_BODY=""

burst_dir="$(mktemp -d)"
i=1
while [ "$i" -le 30 ]; do
  (
    result="$(request "GET" "$BASE_URL/txt?prompt=poem&api=invalid-provider&format=json&burst=$i" "$ALLOWED_REFERER")"
    printf '%s\n' "$result" > "$burst_dir/$i.out"
  ) &
  i=$((i + 1))
done
wait

i=1
while [ "$i" -le 30 ]; do
  if [ -f "$burst_dir/$i.out" ]; then
    status="$(sed -n '1p' "$burst_dir/$i.out")"
    body="$(sed -n '2p' "$burst_dir/$i.out")"
    RATE_LIMIT_LAST_STATUS="$status"
    RATE_LIMIT_LAST_BODY="$body"
    case "$body" in
      *"accessed too frequently"*)
        RATE_LIMIT_HIT=1
        break
        ;;
    esac
  fi
  i=$((i + 1))
done
rm -rf "$burst_dir"

if [ "$RATE_LIMIT_HIT" = "1" ]; then
  PASS_COUNT=$((PASS_COUNT + 1))
  printf '[通过] 并发突发请求已触发流控，HTTP %s\n' "$RATE_LIMIT_LAST_STATUS"
else
  FAIL_COUNT=$((FAIL_COUNT + 1))
  printf '[失败] 30 个并发请求后仍未触发流控；最后一次 HTTP %s\n' "$RATE_LIMIT_LAST_STATUS"
  printf '       最后响应内容：%s\n' "$RATE_LIMIT_LAST_BODY"
fi

print_header "汇总"
printf '通过=%s 失败=%s\n' "$PASS_COUNT" "$FAIL_COUNT"

if [ "$FAIL_COUNT" -gt 0 ]; then
  exit 1
fi

exit 0
