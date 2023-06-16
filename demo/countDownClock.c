# A simple demo to render a countdown clock counting down from 10 to 0
#
include "clock/dial.c"
include "clock/clockFace.c"
include "clock/hands.c"

main() {
    println( "Demo Countdown Clock - Smooth")

    startTime := 10  // Start at 10 seconds
    frameRate := 30.0  // Frame rate of video

    // The end frame number - here we want startTime seconds + a buffer at the end
    end := (startTime+1)*frameRate

    // create a context with start, end frame numbers, the frame rate and the duration
    ctx:= animGraphic.NewContext()

    try( encoder := ffmpeg.New( "test.mp4", frameRate ) ) {
        for frame:=0; frame<=end; frame=frame+1 {
            // the time since the startTime for this frame
            // we use Max so that when it hits 0 the clock hand stays there
            sec := math.Max( 0, startTime - (frame/frameRate))

            // Clear the frame to all black
            image.Fill(ctx, colour.Colour("black"))

            // Draw the clock
            demoCountdown(ctx, sec)

            // Render the image
            encoder.WriteImage(ctx.Image())
        }
    }

}

demoCountdown(ctx, sec) {
    clockDialBackground(ctx, 1920, 1080, 540, colour.Colour("#000000"))
    clockDialForeground(ctx, 1920, 1080, 540, colour.Colour("#ffffff"))
    countdownClockFace(ctx, 1920, 1080, 540)
    clockSecondHand(ctx, 1920, 1080, 540, false, sec)
}
