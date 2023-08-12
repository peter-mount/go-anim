# example of decoding exif data from an image

main() {
    try( f := os.Open("test.jpg") ) {
        t := exif.Decode(f)

        fmt.Println(t)
        fmt.Println(t.String("Make",""),t.String("Model",""))
        fmt.Println(t.String("Wibble","N/A"))
    }
}