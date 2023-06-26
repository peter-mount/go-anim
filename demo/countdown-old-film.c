# countdown-old-film.c
#
# This shows two circles centered in the frame, an number showing
# the number of seconds remaining and a line circling the center once
# per second, with the background becoming darker as it progresses

main() {
    println( "Demo CountDown Old Film")

    frameRate   := 30 // Frames per second
    duration    := 3 // Duration of animation in seconds

    // Colours of the circles and number
    circleColour    := colour.Colour("black")
    numberColour    := colour.Colour("black")

    ctx:= animGraphic.NewContext()

    try( encoder := ffmpeg.New( "test.mp4", frameRate ) ) {
        for second:=duration; second>=0; second=second-1 {
            drawBackground(ctx)

            for frame:=0; frame<frameRate; frame=frame+1 {
            drawSwipe(ctx,frame,frameRate)
            drawCircles(ctx,circleColour)
            drawCounter(ctx,second,numberColour)
                encoder.WriteImage(ctx.Image())
            }
        }
    }
}

drawBackground(ctx) {
    image.Fill(ctx,colour.Colour("grey"))
    // Hint, you could place a static image here instead of a fill
}

drawSwipe(ctx,frame,frameRate) {
    bounds := ctx.Bounds()
    center := ctx.Center()
    radius := math.Sqrt( (center[0]*center[0])+(center[1]*center[1]))

    // The angle involved
    theta:= 2.0*math.Pi*math.Float(frame)/math.Float(frameRate)
    deg := math.Deg(theta)
    if theta > 0.00001 {
        try( ctx ) {
            gc := ctx.Gc()

            // Save the underlying image then create a mask to apply
            image := ctx.Image()
            ctx.NewImage()

            col := colour.Colour("#40404040")
            println(col)
            gc.SetFillColor(col)
            gc.BeginPath()
            gc.MoveTo(center[0],center[1])
            gc.LineTo(center[0],0)

            if deg >= 45.0 { gc.LineTo(bounds.X2,0)}
            if deg >= 135.0 { gc.LineTo(bounds.X2, bounds.Y2)}
            if deg >= 225.0 { gc.LineTo(bounds.X1,bounds.Y2)}
            if deg >= 315.0 { gc.LineTo(bounds.X1,bounds.Y1)}

            sc := math.Sincos(theta)
            x := radius * sc[0]
            y := radius * sc[1]
            gc.LineTo(center[0]+x,center[1]-y)

            gc.LineTo(center[0],center[1])
            gc.Fill()

            // Now overlay the two images

        }
    }
}

drawCircles(ctx,circleColour) {
    try( ctx ) {
        gc := ctx.Gc()

        center := ctx.Center()
        radius := math.Min( center[0], center[1] ) * 0.9

        gc.SetLineWidth(20)

        gc.SetStrokeColor(circleColour)
        gc.SetFillColor(circleColour)

        gc.BeginPath()
        gc.ArcTo(center[0],center[1],radius,radius,0,2*math.Pi)
        gc.Stroke()

        radius = radius -40
        gc.BeginPath()
        gc.ArcTo(center[0],center[1],radius-20,radius-20,0,2*math.Pi)
        gc.Stroke()
    }
}

drawCounter(ctx,second,numberColour) {
    try( ctx ) {
        gc := ctx.Gc()

        center := ctx.Center()
        radius := math.Min( center[0], center[1] ) * 0.9

        animGraphic.SetFont(gc,"luxi 560 mono bold")
        gc.SetStrokeColor(numberColour)
        gc.SetFillColor(numberColour)

        animUtil.DrawString(gc, center[0], center[1], "%d", second )
    }
}