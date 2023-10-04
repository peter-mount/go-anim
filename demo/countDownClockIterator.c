# A simple demo to render a countdown clock counting down from 10 to 0
#
# This one is similar to countDownClock but where you could just set the
# smooth flag in the call to clockSecondHand to false, as that demo only shows
# the clock, that meant that for each second we were rendering effectively the
# same frame.
#
# So this demo shows how we can speed that up by rendering each seconds frame
# once and passing that to ffmpeg 30 times so we get the video rendered faster.
#
# This method only works if a frame does not change often but it can be handy
# if you want to speed up rendering times when nothing is happening between
# blocks of frames.
#
include "clock/dial.c"
include "clock/clockFace.c"
include "clock/hands.c"

main() {
    println( "Demo Countdown Clock - Iterator")

    startTime := 50  // Start at 10 seconds
    frameRate := 25.0  // Frame rate of video

    // create a context for 4K resolution
    //ctx:= animGraphic.New4k()

    // To speed up development, use 1080p but keep the 4K coordinates by scaling by 0.5
    // This will speed up rendering rates. Only switch back to 4K for the final render.
//    ctx:= animGraphic.New1080p().Scale(0.5,0.5)

// Alternatively do it at 720p, which is 2/3 the size of 1080p or 1/3 of 4K
    ctx:= animGraphic.New720p().Scale(1/3.0,1/3.0)

    try( encoder := render.New( "test.mp4", frameRate ) ) {
        //encoder.TimeCode().Set("09:25:30")

        // The end frame number - here we want startTime seconds + a buffer at the end
        end := (startTime+1)*frameRate

        // here we use ForFrames providing it with the number of frames
        // for this loop.
        //
        // Note: frameNum will be the frame number within the loop & not
        // for the entire animation.
        for frameNum ,_ = range encoder.ForFrames( end ) {
            image.Fill(ctx, colour.Colour("black"))

            // Max 0 here means it stops at 0 and doesn't continue past it as we
            // have an extra second on the end showing the clock stopped
            demoCountdown(ctx, math.Max(0,startTime - (math.Float(frameNum)/frameRate)) )

            encoder.WriteImage(ctx.Image())
        }
    }

}

demoCountdown(ctx, sec) {
    clockDialBackground(ctx, 1920, 1080, 540, colour.Colour("#000000"))
    clockDialForeground(ctx, 1920, 1080, 540, colour.Colour("#ffffff"))
    countdownClockFace(ctx, 1920, 1080, 540)
    // Change true here to false to have the hand's move on each second
    // rather than smooth.
    //
    // However look at the other demo on a way to do that even faster
    clockSecondHand(ctx, 1920, 1080, 540, true, sec)
}
