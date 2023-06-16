
clockSecondHand(ctx, cx, cy, radius, smooth, sec ) {
    radius = radius - 60
    radiusInner := radius / 2.0

    // If not in smooth mode then keep hands snapped at the exact second
    if !smooth { sec = math.Floor(sec) }

    try( ctx ) {
        sc := math.Sincos( math.Rad(360.0*(45-sec)/60.0) )

        gc := ctx.Gc()
        gc.BeginPath()
        gc.MoveTo(cx+((radiusInner+10)*sc[1]), cy+((radiusInner+10)*sc[0]))
        gc.LineTo(cx+((radius-20)*sc[1]), cy+((radius-20)*sc[0]))
        gc.SetStrokeColor(colour.Colour("red"))
        gc.SetLineWidth(5)
        gc.FillStroke()
    }
}