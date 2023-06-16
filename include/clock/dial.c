#
# Clock dial generates a bordered circular area where
#
include "testcard/smpte.c"
include "clock/clockFace.c"
include "clock/hands.c"

clockDialBackground(ctx, cx, cy, radius, fill) {
    try( ctx ) {
        gc := ctx.Gc()
        gc.SetFillColor( fill )
        gc.BeginPath()
        gc.ArcTo(cx, cy, radius, radius, 0, 2*math.Pi)
        gc.Fill()
    }
}

clockDialForeground(ctx, cx, cy, radius, stroke) {
    try( ctx ) {
        gc := ctx.Gc()
        gc.SetLineWidth(10)
        gc.SetStrokeColor(stroke)
        gc.SetFillColor(stroke)
        gc.BeginPath()
        gc.ArcTo(cx,cy,radius,radius,0,2*math.Pi)
        gc.Stroke()
    }
}

main() {
    println( "Test Card: SMPTE + Countdown Clock" )

    ctx := animGraphic.NewContext(0,5*30,30,5*30)

    testCardSMPTE(ctx)

    clockDialBackground(ctx, 1920, 1080, 540, colour.Colour("#000000"))
    clockDialForeground(ctx, 1920, 1080, 540, colour.Colour("#ffffff"))
    //clockFace(ctx, 1920, 1080, 540)
    countdownClockFace(ctx, 1920, 1080, 540)
    clockSecondHand(ctx, 1920, 1080, 540, false, 19)

    try( f:=os.Create("/home/peter/test.png") ) {
        image.WritePNG(f,ctx.Image())
    }
}