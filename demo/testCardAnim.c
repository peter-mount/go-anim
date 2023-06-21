
include "testcard/uhd.c"
include "testcard/uhd-countdown-slider.c"

main() {
    println( "Demo Test Card UHD with slider" )

    frameRate   := 30 // Frames per second
    duration    := 10 // Duration of animation

    ctx:= animGraphic.NewContext()

    totalFrames := frameRate * duration

    try( encoder := ffmpeg.Raw( "test.mp4", frameRate, ctx.Image() ) ) {
        for frame:=0; frame< totalFrames; frame=frame+1 {
            // The shared state
            state := testCardUHD_Init(ctx)
            state.frameRate = frameRate
            state.frameNumber = math.Float(frame)

            testCardUHD_Base(ctx, state)
            testCardUHD_Lower(ctx, state)
            testCardUHD_Slicer_Lower(ctx, state)
            testCardUHD_Upper(ctx, state)
            testCardUHD_Slicer_Upper(ctx, state)
            testCardUHD_Top(ctx, state)

            encoder.WriteImage(ctx.Image())
        }
    }

}
