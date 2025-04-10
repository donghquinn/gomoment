# Go Moment

* It's highly inspired from moment.js
    * It was weird to me formatting time value with such strings; '2006-~~~'
    * or RFC format either. I'm familiar with moment.js
* gomoment creates Formatted time strings with formats

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