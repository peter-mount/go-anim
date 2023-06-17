# This demo will display a set of test cards in sequence.
#
# Each test card will be displayed for 5 seconds each.
#
# Similar to countDownClockFast, this will render each card just once and then
# pass those frames to ffmpeg

include "testcard/bbc_unknown.c"
include "testcard/ebu.c"
include "testcard/smpte.c"
include "testcard/uhd.c"

main() {
    println( "Demo Test Card Sequence" )

    frameRate       := 30 // Frames per second
    durationPerCard :=  5 // Number of seconds between each card

    ctx:= animGraphic.NewContext()

    try( encoder := ffmpeg.New( "test.mp4", frameRate ) ) {

        testCardSMPTE( ctx )
        renderFrame(ctx, encoder, frameRate, durationPerCard)

        testCardEBU( ctx )
        renderFrame(ctx, encoder, frameRate, durationPerCard)

        testCardBBCUnknown( ctx )
        renderFrame(ctx, encoder, frameRate, durationPerCard)

        testCardUHD( ctx )
        renderFrame(ctx, encoder, frameRate, durationPerCard*3)
    }

}

renderFrame(ctx, encoder, frameRate, duration) {
    // Render the image to a buffer
    frame := png.EncodeBytes( ctx.Image() )

    frameCount := frameRate * duration

    // Now render it frameRate times.
    // this time we use Write() instead of WriteImage()
    // as frame has already been encoded
    for i:=0; i<frameCount; i=i+1 {
        encoder.Write(frame)
    }
}