# GoMoment

[English](./README.md) | **한국어**

[![Test](https://github.com/donghquinn/gomoment/actions/workflows/test.yml/badge.svg)](https://github.com/donghquinn/gomoment/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/go-1.12+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/donghquinn/gomoment.svg)](https://pkg.go.dev/github.com/donghquinn/gomoment)

[moment.js](https://momentjs.com/) 스타일의 포맷 토큰으로 날짜를 파싱하고, 포맷하고, 타임존을 변환하는 Go 라이브러리입니다. Go의 기준 시각 `2006-01-02` 대신 `YYYY-MM-DD`처럼 작성할 수 있습니다.

```go
moment, _ := gomoment.NewMoment("2023-05-15T14:30:45Z")
kst, _ := moment.TZ("KST")
fmt.Println(kst.Must("YYYY-MM-DD HH:mm Z")) // 2023-05-15 23:30 +09:00
```

## 기능

- **moment.js 스타일 토큰** — `YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`, `SSS`, `A`, `Z`, ...
- **타임존 변환** — `TZ("KST")`, `TZ("America/New_York")`, `UTC()`, `Local()`
- **유연한 파싱** — 자주 쓰는 날짜 문자열 형식을 자동으로 인식합니다
- **엄격한 검증** — 지원하지 않는 토큰은 잘못된 결과를 내지 않고 에러를 반환합니다
- **의존성 없음** — 표준 라이브러리만 사용합니다

## 설치

```bash
go get -u github.com/donghquinn/gomoment
```

## 빠른 시작

```go
package main

import (
    "fmt"

    "github.com/donghquinn/gomoment"
)

func main() {
    now, err := gomoment.NewMoment()
    if err != nil {
        panic(err)
    }

    formatted, err := now.Format("YYYY-MM-DD HH:mm:ss")
    if err != nil {
        panic(err)
    }

    fmt.Println(formatted) // 예: 2025-08-06 14:30:45
}
```

## 사용법

### Moment 생성

#### 현재 시각
```go
moment, err := gomoment.NewMoment()
// 또는
moment := gomoment.Now()
```

#### 문자열에서 생성
```go
moment, err := gomoment.NewMoment("2025-08-06 14:30:45")
// 또는
moment, err := gomoment.Parse("2025-08-06 14:30:45")
```

#### time.Time에서 생성
```go
t := time.Now()
moment, err := gomoment.NewMoment(t)
```

그 밖의 타입을 인자로 넘기면 에러를 반환합니다.

### 포맷

| 토큰 | 설명 | 예시 |
|-------|-------------|---------|
| `YYYY` | 4자리 연도 | 2025 |
| `YY` | 2자리 연도 | 25 |
| `MM` | 월 (01-12) | 08 |
| `M` | 월 (1-12) | 8 |
| `DD` | 일 (01-31) | 06 |
| `D` | 일 (1-31) | 6 |
| `HH` | 시, 24시간제 (00-23) | 09 |
| `H` | 시, 24시간제 (0-23) | 9 |
| `hh` | 시, 12시간제 (01-12) | 02 |
| `h` | 시, 12시간제 (1-12) | 2 |
| `mm` | 분 (00-59) | 05 |
| `m` | 분 (0-59) | 5 |
| `ss` | 초 (00-59) | 07 |
| `s` | 초 (0-59) | 7 |
| `SSS` | 밀리초 (000-999) | 123 |
| `A` | AM/PM | PM |
| `a` | am/pm | pm |
| `ZZ` | UTC 오프셋 | +0900 |
| `Z` | 콜론을 포함한 UTC 오프셋 | +09:00 |

포맷 규칙:

- 토큰은 같은 문자가 연속된 묶음입니다. 표에 없는 묶음(예: `MMM`, `DDDD`)은 에러입니다.
- 대괄호 안의 텍스트는 그대로 출력합니다: `[Today is] YYYY-MM-DD` → `Today is 2025-08-06`.
- ASCII 영문자가 아닌 문자(숫자, 구두점, 공백, 한글 등)는 그대로 출력합니다.

#### 예시

```go
moment, _ := gomoment.NewMoment("2025-08-06 14:30:45")

fmt.Println(moment.Must("YYYY-MM-DD"))          // 2025-08-06
fmt.Println(moment.Must("MM/DD/YYYY"))          // 08/06/2025
fmt.Println(moment.Must("YYYY-MM-DD HH:mm:ss")) // 2025-08-06 14:30:45
fmt.Println(moment.Must("h:mm A"))              // 2:30 PM
fmt.Println(moment.Must("[Date:] YYYY.MM.DD"))  // Date: 2025.08.06
fmt.Println(moment.Must("YYYY년 M월 D일"))        // 2025년 8월 6일
```

### Must

`Must`는 `Format`과 같지만 포맷 문자열이 잘못되면 panic이 발생합니다. 상수처럼 유효함이 확실한 포맷 문자열에만 사용하십시오.

```go
moment := gomoment.Now()
formatted := moment.Must("YYYY-MM-DD HH:mm:ss")
```

### 지원하는 입력 형식

`NewMoment(string)`과 `Parse`는 아래 레이아웃을 순서대로 시도합니다(Go 기준 시각 표기).

- `2006-01-02`
- `2006/01/02`
- `2006-01-02 15:04:05`
- `2006/01/02 15:04:05`
- `2006-01-02T15:04:05`
- `2006-01-02T15:04:05Z`
- `2006-01-02T15:04:05-07:00`
- `01/02/2006`
- `01-02-2006`
- `01/02/2006 15:04:05`
- `01-02-2006 15:04:05`
- `15:04:05`

오프셋이 없는 문자열은 UTC로 해석합니다. 다른 타임존이 필요하면 결과에 `TZ`를 사용하십시오.

### 내부 time.Time 가져오기

```go
moment, _ := gomoment.NewMoment()
t := moment.Time() // time.Time
```

### 타임존

모든 변환은 새 `Moment`를 반환하며, 원래 값은 바뀌지 않습니다.

#### 특정 타임존으로 변환

```go
moment, _ := gomoment.NewMoment("2023-05-15T14:30:45Z")

kstMoment, _ := moment.TZ("KST")
fmt.Println(kstMoment.Must("YYYY-MM-DD HH:mm")) // 2023-05-15 23:30

estMoment, _ := moment.TZ("EST")
fmt.Println(estMoment.Must("YYYY-MM-DD HH:mm")) // 2023-05-15 10:30 (5월은 EDT)

// IANA 타임존 이름도 사용할 수 있습니다
nyMoment, _ := moment.TZ("America/New_York")
tokyoMoment, _ := moment.TZ("Asia/Tokyo")
```

#### 지원하는 타임존 약어

약어는 대소문자를 구분하지 않습니다. 각 약어는 IANA 지역에 대응하므로 해당 지역이 서머타임을 적용하면 서머타임도 반영됩니다(예: `EST`는 뉴욕을 따르므로 여름에는 EDT가 됩니다).

| 약어 | IANA 이름 | 설명 |
|--------------|-----------|-------------|
| `UTC` | UTC | 협정 세계시 |
| `GMT` | GMT | 그리니치 평균시 |
| `KST` | Asia/Seoul | 한국 표준시 |
| `JST` | Asia/Tokyo | 일본 표준시 |
| `EST` | America/New_York | 미국 동부 시간 |
| `PST` | America/Los_Angeles | 미국 태평양 시간 |
| `CST` | America/Chicago | 미국 중부 시간 |
| `MST` | America/Denver | 미국 산악 시간 |
| `CET` | Europe/Paris | 중앙유럽 시간 |
| `IST` | Asia/Kolkata | 인도 표준시 |

#### 기타 타임존 메서드

```go
moment, _ := gomoment.NewMoment()

utcMoment := moment.UTC()     // UTC로 변환
localMoment := moment.Local() // 로컬 타임존으로 변환

name, offset := moment.Zone() // 예: "KST", 32400
offset = moment.Offset()      // UTC 기준 동쪽 방향 초 단위 오프셋
```

## 에러 처리

```go
// 잘못된 날짜 문자열
_, err := gomoment.NewMoment("invalid-date")
fmt.Println(err)
// unable to parse date string "invalid-date" (supported formats: YYYY-MM-DD, YYYY/MM/DD, RFC3339, etc.)

// 지원하지 않는 포맷 토큰
moment := gomoment.Now()
_, err = moment.Format("YYYY-QQ-MMM")
fmt.Println(err)
// unsupported format tokens: [QQ MMM] (supported tokens: YYYY, YY, MM, M, DD, D, HH, H, hh, h, mm, m, ss, s, SSS, A, a, ZZ, Z)

// 알 수 없는 타임존
_, err = moment.TZ("Mars/Olympus")
fmt.Println(err) // invalid timezone: Mars/Olympus
```

## 성능

`Format`은 포맷 문자열을 한 번만 순회하며 결과를 바로 씁니다. 중간에 Go 레이아웃을 만들지 않으므로 호출당 할당은 보통 1~2회입니다.

```
goos: darwin
goarch: arm64
cpu: Apple M5
```

| 벤치마크 | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `Format("YYYY-MM-DD")` | 54.6 | 16 | 1 |
| `Format("YYYY-MM-DD HH:mm:ss Z")` | 130.7 | 80 | 2 |
| 표준 라이브러리 `time.Format("2006-01-02 15:04:05 -07:00")` | 66.0 | 32 | 1 |
| `NewMoment()` | 34.6 | 24 | 1 |
| `NewMoment("2023-05-15 14:30:45")` | 169.0 | 312 | 9 |
| `Now()` | 27.5 | 0 | 0 |

`Format`은 같은 레이아웃의 표준 라이브러리 `time.Format`보다 약 2배 느립니다. 대신 Go 기준 시각이 아닌 moment.js 토큰으로 작성할 수 있습니다.

직접 벤치마크를 실행하려면:

```bash
go test -run '^$' -bench . -benchmem
```

## 기여

기여를 환영합니다. 자유롭게 Pull Request를 보내 주십시오.

PR을 보내기 전에 CI와 같은 검사를 실행하십시오.

```bash
go vet ./...
golangci-lint run
go test -race ./...
```

GitHub Actions는 `main`/`dev` 브랜치 push와 모든 Pull Request마다 `.github/workflows/test.yml`을 실행합니다. Go 1.12(`go.mod`의 최소 버전), 직전 안정 버전, 최신 안정 버전에서 테스트하고 `golangci-lint`를 실행합니다.

## 릴리스

`.github/workflows/release.yml`은 전체 테스트 워크플로를 실행한 뒤, 자동 생성된 릴리스 노트로 GitHub Release를 만들고 Go 모듈 프록시에 버전을 등록합니다. 시작하는 방법은 두 가지입니다.

- **GitHub에서 실행(권장)** — Actions → Release → `main`에서 Run workflow를 누르고 `patch`, `minor`, `major` 중 하나를 고릅니다. 워크플로가 최신 `vX.Y.Z` 태그를 올리고 새 태그를 직접 push합니다.
- **태그로 실행** — semver 태그를 push합니다. `v1.3.0-rc.1`처럼 접미사가 붙은 태그는 pre-release로 게시합니다.

  ```bash
  git tag v1.2.0
  git push origin v1.2.0
  ```

v2 이상으로 올리는 `major` 릴리스는 Go 모듈 규칙에 따라 `go.mod`의 모듈 경로가 `/v2`(또는 해당 메이저 버전)로 끝나기 전까지 의도적으로 실패합니다.

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참고하십시오.

## 영감

이 라이브러리는 JavaScript 라이브러리 [moment.js](https://momentjs.com/)에서 영감을 받았습니다.
