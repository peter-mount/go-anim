# This demo will create a test card in exr format

include "testcard/uhd.c"

main() {
    println( "Demo UHD Test Card in EXR format" )

    // create a context for 4K resolution
    ctx:= animGraphic.New4k()

    // exr as a single file
    try( f := os.Create("test.exr") ) {
        testCardUHD( ctx )
        render.Exr().Encode(f,ctx.Image())
    }
    try( f := os.Create("test.png") ) {
        testCardUHD( ctx )
        render.Png().Encode(f,ctx.Image())
    }
    try( f := os.Create("test.jpg") ) {
        testCardUHD( ctx )
        e:=render.Jpeg().Encode(f,ctx.Image())
    }
    try( f := os.Create("test.tif") ) {
        testCardUHD( ctx )
        render.Tiff().Encode(f,ctx.Image())
    }

    // exr vua a renderer
    try( encoder := render.New( "test-%05d.exr", 1 ) ) {
        testCardUHD( ctx )
        encoder.WriteImage(ctx.Image())
    }

}
