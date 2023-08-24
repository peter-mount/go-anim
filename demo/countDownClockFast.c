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
    println( "Demo Countdown Clock - Fast")

    startTime := 50  // Start at 10 seconds
    frameRate := 25.0  // Frame rate of video

    // The end frame number - here we want startTime seconds + a buffer at the end
    end := (startTime+1)*frameRate

    // create a context for 4K resolution
    ctx:= animGraphic.New4k()

    // To speed up development, use 1080p but keep the 4K coordinates by scaling by 0.5
    // This will speed up rendering rates. Only switch back to 4K for the final render.
//    ctx:= animGraphic.New1080p().Scale(0.5,0.5)

// Alternatively do it at 720p, which is 2/3 the size of 1080p or 1/3 of 4K
//    ctx:= animGraphic.New720p().Scale(1/3.0,1/3.0)

    try( encoder := render.New( "test.mp4", frameRate ) ) {
        encoder.TimeCode().Set("09:25:30")
        // here we use sec as the main loop
        for sec:=startTime; sec>=0; sec=sec-1 {
            // Clear the frame to all black
            image.Fill(ctx, colour.Colour("black"))

            // Draw the clock
            demoCountdown(ctx, sec)

            // Now render it frameRate times.
            // This works because we only change the scene once every second.
            // Using WriteImageMulti is better than a for loop as it's faster
            encoder.WriteImageN(ctx.Image(),frameRate)
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
