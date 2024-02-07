# This demo will create a test card in exr format

include "testcard/uhd.c"

main() {
    println( "Demo UHD Test Card in EXR format" )

    // create a context for 4K resolution
//    ctx:= animGraphic.New4k()
    ctx:= animGraphic.New(image.New4K16())
//    ctx:= animGraphic.New(image.New720p())
//    ctx:= animGraphic.New(image.New720p16())

    testCardUHD( ctx )

    // exr as a single file
    try( f := os.Create("test.exr") ) {
        render.Exr().Encoder().Compress(true).Encode(f,ctx.Image())
    }

    try( f := os.Create("test.png") ) {
        render.Png().Encode(f,ctx.Image())
    }

    try( f := os.Create("test.jpg") ) {
        e:=render.Jpeg().Encode(f,ctx.Image())
    }

    try( f := os.Create("test.tif") ) {
        render.Tiff().Encode(f,ctx.Image())
    }

    // exr vua a renderer
    try( encoder := render.New( "test-%05d.exr", 1 ) ) {
        encoder.WriteImage(ctx.Image())
    }
}
