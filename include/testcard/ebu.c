# EBU test card. Colours based on
# https://en.wikipedia.org/wiki/Test_card#/media/File:EBU_Colorbars.svg

testCardEBU(ctx) {
    gc := ctx.Gc()

    animUtil.DrawColourBars(gc, ctx.Bounds(),
        colour.Colour("#ffffff"),
        colour.Colour("#c0c000"),
        colour.Colour("#00c0c0"),
        colour.Colour("#00c000"),
        colour.Colour("#c000c0"),
        colour.Colour("#c00000"),
        colour.Colour("#0000c0"),
        colour.Colour("#000000")
    )
}

main() {
    println( "Test Card: EBU" )

    ctx := animGraphic.NewContext()

    testCardEBU(ctx)

    try( f:=os.Create("/home/peter/test.png") ) {
        png.Encode(f,ctx.Image())
    }
}
