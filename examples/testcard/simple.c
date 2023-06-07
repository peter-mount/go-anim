#! ./builds/linux/amd64/bin/goanim

main() {
    println( "Simple Test Card Example" )

    min = 5
    max = 10
    for i=0;i<12;i=i+.5 {
        println( i, min, max, between(i,min,max) )
    }
}