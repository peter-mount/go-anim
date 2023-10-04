# A clock which counts down from a set time to 0:00
# The border is animated so you see the progress of the current second.
# The border also cycles through a set of colours, so for 2 seconds it shows
# one color (first drawing the second, then removing the second from the display)
#
# In the center is the countdown timer.
#
# Based on the 5 minute countdown clock at the start of this stream:
# https://www.youtube.com/watch?v=SJo5hVNk0BQ

main() {
    // Start at 5 minutes
    startTime := 5
    // Show frame counter
    showFrames := true

    // Frame rate of video
    frameRate := 30

    // Colours we will rotate though, must be at least one but can be more than two
    colours := append(newArray(), colour.Colour("yellow"), colour.Colour("white") )

    // colour beneath the animated circle
    baseColour := colour.Colour("grey")

    // create a context for 4K resolution
    ctx:= animGraphic.New4k()
    cx := ctx.Width()/2.0
    cy := ctx.Height()/2.0

    // To speed up development, use 1080p but keep the 4K coordinates by scaling by 0.5
    // This will speed up rendering rates. Only switch back to 4K for the final render.
    //ctx:= animGraphic.New1080p().Scale(0.5,0.5)
    //cx := ctx.Width()*2/2.0
    //cy := ctx.Height()*2/2.0

    // Alternatively do it at 720p, which is 2/3 the size of 1080p or 1/3 of 4K
    //ctx:= animGraphic.New720p().Scale(1/3.0,1/3.0)
    //cx := ctx.Width()*3.0/2.0
    //cy := ctx.Height()*3.0/2.0


    radius := math.Min(cx,cy)*0.6
    lineWidth := math.Min(cx,cy)*0.05

    pi2=2.0*math.Pi

    try( encoder := render.New( "test.mp4", frameRate ) ) {
        // Cycle through each minute
        min:=startTime
        while min>0 {
            min--

            for sec:=59; sec>=0; sec-- {
                // Select the colour to use
                activeColour := colours[(sec/2) % len(colours)]

                for frame:=0; frame<frameRate; frame++ {
                    image.Fill(ctx, colour.Colour("black"))

                    try( ctx ) {
                        gc:= ctx.Gc()

                        gc.BeginPath()
                        gc.ArcTo(cx,cy,radius,radius,0,pi2)
                        gc.SetStrokeColor( baseColour )
                        gc.SetLineWidth(lineWidth)
                        gc.Stroke()
                    }

                    try( ctx ) {
                        gc:= ctx.Gc()

                        gc.Translate(cx,cy)
                        gc.Rotate(-math.Pi/2.0)
                        gc.SetStrokeColor( activeColour )
                        gc.SetLineWidth(lineWidth)
                        gc.SetLineCap(1)
                        gc.BeginPath()

                        p := frame * (pi2/frameRate)
                        if sec%2 == 0 p-=pi2
                        gc.ArcTo(0,0,radius,radius,0,p)

                        gc.Stroke()
                    }

                    try( ctx ) {
                        gc:= ctx.Gc()

                        animGraphic.SetFont(gc,"freemono 120 mono bold")
                        gc.SetStrokeColor(colour.Colour("blue"))
                        gc.SetFillColor(colour.Colour("white"))
                        if showFrames {
                            animUtil.DrawString(gc, cx, cy, "%02d:%02d:%02d", min, sec, frameRate-frame-1 )
                        } else {
                            animUtil.DrawString(gc, cx, cy, "%02d:%02d", min, sec )
                        }
                    }

                    encoder.WriteImage(ctx.Image())
                }

            }

        }
    }
}