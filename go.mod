module github.com/peter-mount/go-anim

go 1.24

toolchain go1.24.0

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/llgcode/draw2d v0.0.0-20240627062922-0ed1ff131195
	github.com/peter-mount/go-build v0.0.0-20250218200125-f187f75a6a5d
	github.com/peter-mount/go-kernel/v2 v2.0.3-0.20250218195942-5604474bedd7
	github.com/peter-mount/go-script v0.0.0-20250218200359-943ffa62e818
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/x448/float16 v0.8.4
	golang.org/x/image v0.24.0
)

//replace github.com/peter-mount/go-script v0.0.0-20241218090358-129a6c764bf4 => ../script

require (
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
