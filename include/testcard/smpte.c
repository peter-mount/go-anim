#
# SMPTE Test Card
#

main() {
    println( "Test Card: SMPTE" )

    ctx := animGraphic.NewContext(0,5*30,30,5*30)

    black := colour.Colour("black")
    white := colour.Colour("white")

    w := ctx.Width()
    h := ctx.Height()
    ub := animUtil.Rect(0, 0, w, h*336.0/504.0)
    mb := animUtil.Rect(0, ub.Y2, w, ub.Y2+(h*96.0/504.0))
    lb := animUtil.Rect(0, mb.Y2, w, h) // mb.Y2+(h*126/504))

    image.Fill(ctx,black)

    try( ctx ) {
        gc := ctx.Gc()

        animUtil.DrawColourBars(gc, ub,
            white,
            colour.Colour("#c0c000"),
            colour.Colour("#00c0c0"),
            colour.Colour("#00c000"),
            colour.Colour("#c000c0"),
            colour.Colour("#c00000"),
            colour.Colour("#0000c0")
        )

        animUtil.DrawColourBars(gc, mb,
            colour.Colour("#0000c0"),
            black,
            colour.Colour("#c000c0"),
            black,
            colour.Colour("#00c0c0"),
            black,
            colour.Colour("silver")
        )

        // Bottom row is 5 columns but the last one is itself split into 4
        dw := animUtil.DrawColourBars(gc, lb,
            colour.Colour("#00214c"),
            white,
            colour.Colour("#32006a"),
            black,
            // Placeholder for next block
            black
        )
        animUtil.DrawColourBars(gc,
            animUtil.Rect(dw[0]*(dw[1]-1), lb.Y1, lb.X2, lb.Y2),
            colour.Colour("#090909"),
            black,
            colour.Colour("#1d1d1d"),
            black
        )
    }

    try( f:=os.Create("/home/peter/test.png") ) {
        image.WritePNG(f,ctx.Image())
    }

}