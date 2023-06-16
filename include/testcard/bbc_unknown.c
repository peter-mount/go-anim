#
# BBC Unknown testcard
# Unknown based on the BBC Unknown test card
# https://en.wikipedia.org/wiki/List_of_BBC_test_cards

testCardBBCUnknown(ctx) {
    White := colour.Colour("white")
    Yellow:= colour.Colour("yellow")
    Cyan:= colour.Colour("cyan")
    Green:= colour.Colour("Green")
    Magenta:= colour.Colour("magenta")
    Red:= colour.Colour("red")
    Blue:= colour.Colour("blue")
    Black:= colour.Colour("black")
    Lightgrey:=colour.Colour("Lightgrey")

    image.Fill(ctx,Black)

    try( ctx ) {
        gc := ctx.Gc()

        // Grid pattern 16x12 but offset by half
        cw := 16.0
        ch := 12.0

        w := ctx.Width()
        h := ctx.Height()
        dw := w/cw
        dh := h/ch
        dx := dw/2.0
        dy := dh/2.0

        gc.SetLineWidth(2)

        gc.SetStrokeColor(White)
        gc.BeginPath()

        // Vertical lines
        for i := 0.0; i < cw; i=i+1 {
            animGraphic.RelLine(gc, dx+(dw*i), 0, 0, h)
        }

        // Horizontal lines
        for i := 0.0; i < ch; i=i+1 {
            animGraphic.RelLine(gc, 0, dy+(dh*i), w, 0)
        }

        gc.Stroke()

        // Vertical line spacers
        gc.SetLineWidth(20)
        gc.SetStrokeColor(Black)
        for i := 1.0; i < cw; i=i+1 {
            animGraphic.RelLine(gc, dw*i, 0, 0, h)
        }
        gc.Stroke()

        // FIXME need array to args, eg gradient()...
        //animUtil.DrawColourBars(gc,
        //    animUtil.Rect(10, 2*dh, w-10, 5*dh),
        //    colour.Gradient(7, Black, White))

        bdw := animUtil.DrawColourBars(gc,
            animUtil.Rect(10, 7*dh, w-10, 9*dh),
            White,
            Yellow,
            Cyan,
            Green,
            Magenta,
            Red,
            Blue,
            Black
        )

        animUtil.DrawColourBars(gc,
            animUtil.Rect(0, 9*dh, w-50+bdw[0], 10*dh),
            Blue,
            Black,
            Magenta,
            Black,
            Cyan,
            Black,
            Lightgrey,
            Black,
            White
        )

    }
}

main() {
    println( "Test Card: BBC Unknown" )

    ctx := animGraphic.NewContext(0,5*30,30,5*30)

    testCardBBCUnknown(ctx)

    try( f:=os.Create("/home/peter/test.png") ) {
        png.Encode(f,ctx.Image())
    }
}