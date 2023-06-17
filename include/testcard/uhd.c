# UHD Testcard

// Creates a common state object
testCardUHD_Init(ctx) {
    // Common config
    s := map(
        "w": math.Float(ctx.Width()),
        "h": math.Float(ctx.Height())
    )
    s.w2 = s.w / 2.0
    s.dW = s.w / 30.0
    s.dH = s.h / 27.0
    s.dW2 = s.dW * 2.0
    s.dWd2 = s.dW / 2.0
    s.dHd2 = s.dH / 2.0
    s.dH2 = s.dH * 2.0
    s.h2 = s.h / 2.0
    s.dH23 = s.dH * 2.0 / 3.0
    s.dx = (s.w - s.dW2) / 15.0
    s.dd = s.dW + (2.0 * s.dx)
    s.dy = (s.h - s.dH2) / 8.0

    s.black = colour.Colour("black")
    s.white = colour.Colour("white")
    s.grey = colour.Colour("grey")
    s.lightblue = colour.Colour("Lightblue")
    s.blue = colour.Colour("Blue")
    s.cyan = colour.Colour("Cyan")
    s.lightgreen = colour.Colour("Lightgreen")
    s.green = colour.Colour("Green")
    s.darkgreen = colour.Colour("Darkgreen")
    s.magenta = colour.Colour("Magenta")
    s.red = colour.Colour("Red")
    s.yellow = colour.Colour("Yellow")

    return s
}

testCardUHD_Base(ctx, s) {
    image.Fill(ctx,s.grey)
}

testCardUHD_Lower(ctx, s) {
    try( ctx ) {
        gc := ctx.Gc()

        // Grid lines
        gc.SetStrokeColor(s.white)
        gc.SetLineWidth(10)
        gc.BeginPath()
        animGraphic.Rectangle(gc, s.dW, s.dH, s.w-(s.dW*2), s.h-(s.dH*2))
        gc.Stroke()

        gc.BeginPath()

        y := s.dH
        for i := 0; i < 8; i=i+1 {
            if i == 4 {
                gc.MoveTo(s.dd, y)
                gc.LineTo(s.w-s.dd, y)
            } else {
                gc.MoveTo(0, y)
                gc.LineTo(s.w, y)
            }
            y = y + s.dy
        }

        x := s.dW
        for i := 0; i < 15; i=i+1 {
            y1 := s.h
            if i > 3 && i < 12 {
                y1 = y - s.dy
            }
            gc.MoveTo(x, 0)
            gc.LineTo(x, y1)
            x = x + s.dx
        }

        gc.Stroke()
    }
}

testCardUHD_Upper1(ctx,gc,s,dd, l, r ) {
    animGraphic.RelLine(gc,
        s.dW+(l*s.dx)+dd, s.dH,
        0, (2.0*s.dy)-dd,
        (r*s.dx)-(2*dd), 0,
        0, dd-(2.0*s.dy))
}

testCardUHD_Upper(ctx, s) {
    try( ctx ) {
        gc := ctx.Gc()
        
        dd := 7.0
        gc.SetStrokeColor(s.black)
        gc.SetLineWidth(dd)
        gc.BeginPath()
    
        // Black lines top row
        animGraphic.RelLine(gc,
            s.dW, s.dH+(2*s.dy)-dd,
            (3*s.dx)-dd, 0,
            0, dd-(2*s.dy))
    
        testCardUHD_Upper1(ctx,gc,s,dd, 3.0, 4.0)
        testCardUHD_Upper1(ctx,gc,s,dd, 7.0, 1.0)
        testCardUHD_Upper1(ctx,gc,s,dd, 8.0, 4.0)
    
        animGraphic.RelLine(gc,
            s.dW+(12*s.dx)+dd, s.dH,
            0, (2*s.dy)-dd,
            3*s.dx+dd, 0)
    
        // Bottom row
        animGraphic.RelLine(gc,
            s.dW, s.h-s.dH-(2*s.dy)+dd,
            (3*s.dx)-dd, 0,
            0, 2*s.dy)
    
        animGraphic.RelLine(gc,
            s.dW+(3*s.dx)+dd, (s.h-s.dH)+dd+dd,
            0, (-2*s.dy)-dd,
            (4*s.dx)-dd+dd, 0,
            0, s.dy-dd)
    
        animGraphic.RelLine(gc,
            s.dW+(7*s.dx)+dd, (s.h-s.dH)-s.dy+dd,
            0, (-s.dy)+dd+dd,
            (s.dx)-dd+dd, 0,
            0, s.dy-dd)
    
        animGraphic.RelLine(gc,
            s.dW+(8*s.dx)+dd, (s.h-s.dH)-s.dy+dd,
            0, (-s.dy)+dd+dd,
            (4*s.dx)-dd+dd, 0,
            0, (2*s.dy)-dd)
    
        animGraphic.RelLine(gc,
            s.dW+(12*s.dx)+dd, (s.h-s.dH)-dd,
            0, (-2*s.dy)+dd+dd,
            (3*s.dx)-dd-dd, 0)
    
        // side black lines
        animGraphic.RelLine(gc,
            s.dW, s.dH+(2*s.dy)+dd,
            (3*s.dx)-dd, 0,
            0, (2*s.dy)-dd-dd,
            -s.dx+dd+dd, 0,
            0, -s.dy/2.0)
    
        animGraphic.RelLine(gc,
            s.dW, s.h-(s.dH+(2*s.dy)+dd),
            (3*s.dx)-dd, 0,
            0, -((2 * s.dy) - dd - dd),
            -s.dx+dd+dd, 0,
            0, s.dy/2.0)
    
        animGraphic.RelLine(gc,
            s.w-s.dW, s.dH+(2*s.dy)+dd,
            -((3 * s.dx) - dd), 0,
            0, (2*s.dy)-dd-dd,
            -(-s.dx + dd + dd), 0,
            0, -s.dy/2.0)
    
        animGraphic.RelLine(gc,
            s.w-s.dW, s.h-(s.dH+(2*s.dy)+dd),
            -((3 * s.dx) - dd), 0,
            0, -((2 * s.dy) - dd - dd),
            -(-s.dx + dd + dd), 0,
            0, s.dy/2.0)
    
        // Odd 2 lines on either side of center
        animGraphic.RelLine(gc, s.dW+(2*s.dx)-dd, s.dH+(3.5*s.dy)-dd, 0, s.dy+dd+dd)
        animGraphic.RelLine(gc, (s.w-s.dW)-(2*s.dx)-dd, s.dH+(3.5*s.dy)-dd, 0, s.dy+dd+dd)
    
        // Draw black lines
        gc.Stroke()
    }
}

testCardUHD_Top(ctx, s) {
    try( ctx ) {
        gc := ctx.Gc()

        testCardUHD_cornerCalibrationMarks(ctx,s)
        testCardUHD_leftGradient(ctx,s)
        testCardUHD_rightGradient(ctx,s)

        // Clear bottom border
        animGraphic.FillRectangle(gc, 0, s.h-s.dH, s.w, s.dH, s.black)

        // cornerWhiteSquares
        animGraphic.FillRectangle(gc, 0, 0, s.dW2, s.dH, s.white)
        animGraphic.FillRectangle(gc, s.w-s.dW2, 0, s.dW2, s.dH, s.white)
        animGraphic.FillRectangle(gc, s.w-s.dW, s.h-s.dH, s.dW, s.dH, s.white)

        // Top border
        animUtil.DrawColourBars(gc,
            animUtil.Rect(s.dW2,0,s.w-s.dW2,s.dH),
            s.yellow, s.cyan, s.green, s.magenta, s.red, s.blue, s.black
        )

        // Side borders
        animUtil.DrawColourBarsVertical(gc,
            animUtil.Rect(0,s.dH,s.dW,s.h-s.dH),
            s.lightgreen, s.red, s.black, s.lightblue, s.darkgreen
        )
        animUtil.DrawColourBarsVertical(gc,
            animUtil.Rect(s.w-s.dW,s.dH,s.w,s.h-s.dH),
            s.yellow, s.cyan, s.green, s.magenta, s.red, s.blue, s.black
        )

        testCardUHD_bottomBorders(ctx,s)

        // inner border
        try (ctx) {
            gc.SetStrokeColor(s.white)
        	gc.SetLineWidth(10)
        	gc.BeginPath()
        	animGraphic.Rectangle(gc, s.dW, s.dH, s.w-(s.dW*2), s.h-(s.dH*2))
        	gc.Stroke()
        }

        // triangles
        try (ctx) {
            // Diamonds left/right center
            diaDw := s.dW - 12
            //for _, x := range append(newArray(), 0.0, s.dW2 - 15.0, s.w - s.dW2, s.w - s.dW2 - s.dW2 + 17.0) {
            animGraphic.FillPolyRel(gc, s.white, 0, s.h2, s.dW, -s.dH23, diaDw, s.dH23, -diaDw, s.dH23)
            animGraphic.FillPolyRel(gc, s.white, s.dW2-15, s.h2, s.dW, -s.dH23, diaDw, s.dH23, -diaDw, s.dH23)
            animGraphic.FillPolyRel(gc, s.white, s.w-s.dW2, s.h2, s.dW, -s.dH23, diaDw, s.dH23, -diaDw, s.dH23)
            animGraphic.FillPolyRel(gc, s.white, ((s.w-s.dW2)-s.dW2)+17.0, s.h2, s.dW, -s.dH23, diaDw, s.dH23, -diaDw, s.dH23)

            // left center triangle
            animGraphic.FillPolyRel(gc, s.white, (s.dW2*2)-30, s.h2, s.dW-2, -s.dH23, 0, s.dH23*2)

            // right center triangle
            animGraphic.FillPolyRel(gc, s.white, ((s.w-(s.dW2*2))+30)-12, s.h2, -diaDw+2, -s.dH23, 0, s.dH23*2)

            // Top/Bottom triangles
            animGraphic.FillPolyRel(gc, s.white, s.w2, s.h-1, -s.dH/2, -s.dH, s.dH, 0)
            animGraphic.FillPolyRel(gc, s.white, s.w2, 0, -s.dH/2, s.dH, s.dH, 0)
        }
    }
}

testCardUHD_cornerCalibrationMarks(ctx,s) {
	for sy := 0; sy < 2; sy=sy+1 {
		for sx := 0; sx < 2; sx=sx+1 {
			testCardUHD_cornerCalibrationMark(ctx,s, sx, sy)
		}
	}
}

testCardUHD_cornerCalibrationMark(ctx,s, sx, sy) {
    try( ctx ) {
        gc := ctx.Gc()

        tx:=0
        ty:=0
        rot:=0
        //sw:= sx + (sy << 1)
        sw := sx + (sy * 2)
        if sw==0 {
            //tx, ty, rot = s.dW, -s.dH, 45.0
            tx = s.dW
            ty =  -s.dH
            rot = 45.0
        } else if sw==1 {
            //tx, ty, rot = s.w+s.dW*.25+20, s.dH+10, 135
            tx = s.w+s.dW*0.25+20
            ty =  s.dH+10
            rot = 135.0
        } else if sw==2 {
            //tx, ty, rot = s.dW-s.dx*.75, s.h-s.dH-15, -45
            tx = s.dW-s.dx*0.75
            ty =  s.h-s.dH-15
            rot = -45.0
        } else if sw==3 {
            //tx, ty, rot = s.w+s.dW*.25-s.dx*.7+10, s.h+s.dH, -135
            tx = s.w+s.dW*0.25-s.dx*0.7+10
            ty =  s.h+s.dH
            rot = -135.0
        }

        gc.Translate(tx, ty)
        gc.Rotate(math.Rad(rot))

        //bw, bh := 2.45*s.dx, s.dy
        bw := 2.45*s.dx
        bh := s.dy
        animGraphic.FillRectangle(gc, 0, 0, bw, bh, s.white)

        gc.BeginPath()
        gc.SetStrokeColor(s.black)
        gc.SetLineWidth(10)
        for i := 0.1; i <= 0.9; i = i + 0.1 {
            animGraphic.RelLine(gc, 0, s.dy*i, bw-(s.dy/10.0), 0)
        }
        gc.Stroke()
	}
}

testCardUHD_leftGradient(ctx,s) {
    try( ctx ) {
        gc := ctx.Gc()

        x := s.dW+(3*s.dx)+5
        y := s.dH+(2*s.dy)+5

        rows := 7

        // c colour index, dc change between each colour, cdy height of each entry
        c := 255.0
        dc := 256.0/math.Float(rows-1)
        cdy := ((4.0*s.dy)-10)/math.Float(rows)
        for i := 0; i < rows; i=i+1 {
            animGraphic.FillRectangle(gc, x, y, s.dx-10, cdy, colour.GreyScale(c) )
            y = y + cdy
            c = math.Max(0.0, c-dc)
        }
    }
}

testCardUHD_rightGradient(ctx,s) {
    try( ctx ) {
        gc := ctx.Gc()

        x := s.dW+(11*s.dx)+5
        y := s.dH+(2*s.dy)+8
        y1 := y
    
        rows := 6
        cdy := ((4.0 * s.dy) - 10) / math.Float(rows)
    
        // Fill in block
        animGraphic.FillRectangle(gc, x, y, s.dx-10, cdy*math.Float(rows), s.white)
    
        gc.SetStrokeColor(s.black)
        cols := 7
        for i := 0; i < rows; i=i+1 {
            cdx := (s.dx - 15) / math.Float(cols)
    
            cx := 0.0
            if (cols % 2) == 0 {
                cx = cx + cdx/2.0
            }
    
            gc.BeginPath()
            gc.SetLineWidth(cdx / 2.0)
            for j := 0; j <= cols; j=j+1 {
                animGraphic.RelLine(gc, x+cx+5, y, 0, cdy)
                cx = cx + cdx
            }
            gc.Stroke()
    
            //log.Printf("lg %d %f %f (%f %f %f %f)", i, c, dc, x, y, dx, dy)
            //draw2d.FillRectangle(gc, x, y, dx-10, cdy, color.White)
            y = y + cdy
            cols = cols + (cols - 1)
        }
    
        // Borders
        gc.SetLineWidth(10)
        gc.BeginPath()
        gc.SetStrokeColor(s.white)
        animGraphic.Rectangle(gc, x-5, y1-5, s.dx, cdy*math.Float(rows)+5)
        gc.Stroke()
    }
}

testCardUHD_bottomBorders(ctx,s) {
    try( ctx ) {
        gc := ctx.Gc()

        // Bottom left border
        x := 0.0
        dx := (s.w2 - s.dW) / 256.0
        for i := 255; i > 0; i=i-1 {
            animGraphic.FillRectangle(gc, x, s.h-s.dH, dx, s.dHd2, colour.Grey(32 + ((255 - i)/2)))
            xy := animGraphic.FillRectangle(gc, x, s.h-s.dHd2, dx, s.dHd2, colour.Grey(i))
            x = xy[0]
        }
    
        // Bottom right border
        grad1 :=colour.Gradient(256, colour.Colour("#3cad7c"), colour.Colour("#808080"))
        grad2 :=colour.Gradient(256, colour.Colour("#808080"), colour.Colour("#d25580"))
        grad3 :=colour.Gradient(256, colour.Colour("#829600"), colour.Colour("#808080"))
        grad4 :=colour.Gradient(256, colour.Colour("#808080"), colour.Colour("#7a6bf1"))
        x = math.Float(s.w2+s.dW)
        dw := (s.w2+s.dW2)/2.0
        dx := dw / 256.0
        dw = dw - s.dW2
        for i := 0; i < len(grad1); i=i+1 {
            animGraphic.FillRectangle(gc, x, s.h-s.dH, dx, s.dHd2, grad1[i])
            animGraphic.FillRectangle(gc, x+dw, s.h-s.dH, dx, s.dHd2, grad2[i])
            animGraphic.FillRectangle(gc, x, s.h-s.dHd2, dx, s.dHd2, grad3[i])
            animGraphic.FillRectangle(gc, x+dw, s.h-s.dHd2, dx, s.dHd2, grad4[i])
            x = (x + dx) -1
        }

    }
}

xxtestCardUHD_bottomBorders(ctx,s) {
    try( ctx ) {
        gc := ctx.Gc()


    }
}

// Render the full default UHD test card - see main() below on
// how to modify it
testCardUHD(ctx) {
    try( ctx ) {
        state := testCardUHD_Init(ctx)

        // Render each layer in sequence.
        //
        // This is split like this so that you can insert additional features
        // between each layer
        testCardUHD_Base(ctx, state)
        testCardUHD_Lower(ctx, state)
        testCardUHD_Upper(ctx, state)
        testCardUHD_Top(ctx, state)
    }
}

main() {
    println( "Test Card: UHD TestCard" )

    ctx := animGraphic.NewContext()

    // The shared state
    state := testCardUHD_Init(ctx)

    // Render each layer in sequence.
    //
    // This is split like this so that you can insert additional features
    // between each layer
    testCardUHD_Base(ctx, state)
    testCardUHD_Lower(ctx, state)
    testCardUHD_Upper(ctx, state)
    testCardUHD_Top(ctx, state)

    try( f:=os.Create("/home/peter/test.png") ) {
        png.Encode(f,ctx.Image())
    }
}
