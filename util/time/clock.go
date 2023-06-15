package time

import (
	"git.area51.dev/peter/videoident/renderer"
	"math"
)

// GetCountdownTime returns the number of seconds remaining on the countdown clock.
// This is based on the duration of the clip and the frame being rendered.
// Returned is the time in seconds and as an integer second/frame in second values.
func GetCountdownTime(ctx renderer.Context) (float64, float64, float64) {
	sec := ctx.Duration() - (ctx.FrameF()+1)/ctx.FrameRate()
	secI, secF := math.Modf(sec)

	return sec, secI, math.Mod(secF*ctx.FrameRate(), ctx.FrameRate())
}
