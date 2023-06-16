# clockFace.c - Display the clock face

// clockFace draws the ticks at every second & large every 5 seconds
clockFace(ctx, cx, cy, radius) {
    _clockFace(ctx,cx,cy,radius, 5, 60)
}

// countdownClockFace draws large ticks every 10 seconds and every second < 10
countdownClockFace(ctx, cx, cy, radius) {
    _clockFace(ctx,cx,cy,radius, 10, 10)
}

_clockFace(ctx, cx, cy, radius, majorStep, minorMax) {
    gc := ctx.Gc()

    radius = radius - 60
    radiusInner := radius / 2.0

    // Clock dial
    try( ctx ) {
        gc.SetStrokeColor(colour.Colour("white"))
        gc.SetLineWidth(5)

        // Outer dial
        gc.BeginPath()
        gc.ArcTo(cx,cy,radius,radius,0, 2*math.Pi)
        gc.Stroke()

        // Inner dial
        gc.BeginPath()
        gc.ArcTo(cx,cy,radiusInner,radiusInner,0, 2*math.Pi)
        gc.Stroke()

        // Second ticks
        for sec :=0; sec < 60; sec=sec+1 {
            if (sec%majorStep)==0 || sec <= minorMax {
                deg := 360.0 * (45-sec) / 60.0

                // Small or Large ticks
                r1 := radius
                r2 := radius-(2*5)
                if sec == 5 || (sec%majorStep) == 0 {
                    r1 = radius+(3*5)
                    r2 = radius-(3*5)
                }

                sc1 := math.Sincos(math.Rad(deg - 0.25))
                sc2 := math.Sincos(math.Rad(deg + 0.25))

                gc.MoveTo(cx+(r1*sc1[1]), cy+(r1*sc1[0]))
                gc.LineTo(cx+(r2*sc1[1]), cy+(r2*sc1[0]))
                gc.LineTo(cx+(r2*sc2[1]), cy+(r2*sc2[0]))
                gc.LineTo(cx+(r1*sc2[1]), cy+(r1*sc2[0]))
                gc.Close()
            }
        }
        gc.FillStroke()
    }

    // Legend
    try( ctx ) {
        animGraphic.SetFont(gc,"luxi 30 mono bold")
        gc.SetStrokeColor(colour.Colour("blue"))
        gc.SetFillColor(colour.Colour("white"))

        r := radius + 38

        for sec:=0; sec < 60; sec=sec+1 {
            if (sec%majorStep)==0 || sec==5 {
                try( ctx ) {
                    rad := math.Rad(360.0*(45.0-sec) / 60.0 )
                    sc := math.Sincos(rad)

                    gc.Translate( cx+(r*sc[1]), cy+(r*sc[0]) )
                    gc.Rotate( rad + (math.Pi/2.0))
                    animUtil.DrawString(gc, 0, 0, "%d", sec )
                }
            }
        }
    }
}
