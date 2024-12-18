module github.com/peter-mount/go-anim

go 1.19

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/llgcode/draw2d v0.0.0-20240627062922-0ed1ff131195
	github.com/peter-mount/go-build v0.0.0-20240804094359-01252fe8316a
	github.com/peter-mount/go-kernel/v2 v2.0.3-0.20240514072728-897c39470117
	github.com/peter-mount/go-script v0.0.0-20241218090358-129a6c764bf4
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/x448/float16 v0.8.4
	golang.org/x/image v0.23.0
)

//replace github.com/peter-mount/go-script v0.0.0-20241218090358-129a6c764bf4 => ../script

require (
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
