# GoFritzBox

GoFritzBox is a Golang utility to read or modify Fritz!Box parameters.

### Why it is not complete?
Because AVMDE didn't release a documentation for its APIs, building this library is very hard and I won't waste my time building functions that nobody will use.

### Dude! There's no good documentation!
Nice point, I really appreciate your attention.
I don't really know what some things are due to the very confusing APIs. I've tried to do my best, you'll have some happy time trying to figure out what's going on.

## Supported features:
* Login
* LoadInfo 
* GetStats

If you need any other feature you can open an issue and I will try to add it.  
Soon, ways to edit the Fritz!Box configuration will be added with an auto CRC32 signature.

## Example
```go
package main

import (
	"fmt"

	"github.com/LucaTheHacker/GoFritzBox"
)

func main() {
	session, err := GoFritzBox.Login("http://IP", "USERNAME", "PASSWORD")
	if err != nil {
		panic(err)
	}

	data, err := session.LoadInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println("Fritz!OS version: ", data.FritzOS.Version)

	stats, err := session.GetStats()
	if err != nil {
		panic(err)
	}

	stats.Load()
	fmt.Println(stats.DownstreamTotal)
	fmt.Println(stats.UpstreamTotal)
}

```
