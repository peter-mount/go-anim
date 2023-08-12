# example of decoding exif data from an image

main() {
    try( f := os.Open("test.jpg") ) {
        t := exif.Decode(f)

        // It's a map so this will print all tags found
        for k,v := range t {
            fmt.Printf("%32s %s\n",k,v)
        }

        // Get a single tag. This will error if the tag does not exist
        fmt.Println(t.Make)

        // Does a tag exist?
        if t.Contains("Make") {
            fmt.Println("Make tag exists")
        }

        // Get the string value of tags, or use the defaults if they
        // do not exist
        fmt.Println(t.String("Make",""),t.String("Model",""))

        // Getting a string with a default value.
        // As this is not a standard tag this will always print "N/A"
        fmt.Println(t.String("Wibble","N/A"))

        // Example of getting the number of values in a tag.
        //  0 = tag does not exist
        //  1 = single value
        // >1 = Array
        //
        // As this is not a standard tag it will never run
        if t.Count("Zaphod")>0 {
            fmt.Printf("Field Zaphod has %d entries\n", t.Zaphod.Count)
        }

        // Get the first integer in a tag
        //
        // This will return -1 if the tag doesn't exist as we pass that as the
        // default value.
        //
        // Unlike the underlying exif library, this will return an int value
        // for a float tag
        fmt.Printf("int value %d\n", t.Int(0,"ISOSpeedRatings",-1))

        // Get the first float in a tag
        //
        // This will return -1.1 if the tag doesn't exist as we pass that as the
        // default value
        //
        // Unlike the underlying exif library, this will return a float value
        // for an int tag
        fmt.Printf("float value %f\n", t.Float(0,"ISOSpeedRatings",-1.1))
    }
}