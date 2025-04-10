# Go Moment

* It's highly inspired by moment.js

## Installation

```zsh
go get -u github.com/donghquinn/gomement
```

## Usage

* Same as moment.js
* Default Format is YYYY-MM-DD HH:mm:ss

```go

import (
    	"time"
        "github.com/donghquinn/gomoment"
)

func main() {
    now, formatErr := gomoment.NewMoment(time.now()).Format("YYYY-MM-DD HH:mm:ss")

    // now: 2025-04-10 12:54:23
}
```

### Must
* Must() is the strong verifying method
    * It will return formatted time string
    * But It will throw panic if error found

```go
import (
    	"time"
        "github.com/donghquinn/gomoment"
)

func main() {
    now := gomoment.NewMoment(time.now()).Must("YYYY-MM-DD HH:mm:ss")

    // now: 2025-04-10 12:54:23
}
```