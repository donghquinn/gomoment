# Go Moment

* It's highly inspired by moment.js

## Installation

```zsh
go get -u github.com/donghquinn/gomement
```

## Usage

* Same as moment.js

```go

import (
    	"time"
        "github.com/donghquinn/gomoment"
)

func main() {
    now, formatErr := gomoment.NewMoment(time.now()).Format("YYYY-MM-DD HH:mm:ss")

    // 2025-04-10 12:54:23
}
```