# This demo will create a test card in exr format

include "testcard/uhd.c"

main() {
    println( "Demo UHD Test Card in EXR format" )

    // create a context for 4K resolution
    ctx:= animGraphic.New4k()

    try( encoder := render.New( "test-%05d.exr", 1 ) ) {
        renderImage(ctx,encoder)
    }

//    try( encoder := render.New( "test-%05d.png", 1 ) ) {
//        renderImage(ctx,encoder)
//    }

}

renderImage(ctx,encoder) {
    testCardUHD( ctx )
    encoder.WriteImage(ctx.Image())
}