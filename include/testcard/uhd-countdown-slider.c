# uhd-countdown-slider.c
#
# This displays a simple slider which expands for each frame with one cycle
# representing a single second.
#
# With uhd.c this is in two parts, the lower one is rendered after testCardUHD_Lower and
# the second after testCardUHD_Upper.
#
# It also requires the following entries in state:
#
# frameRate     number of frames per second.
# frameNumber   the frame number being rendered.
#

testCardUHD_Slicer_Lower(ctx, s) {
    try( ctx ) {
        gc := ctx.Gc()

        animGraphic.FillRectangle(gc,
            s.dW+(5.8*s.dx), s.dH+(0.6*s.dy),
            3.4*s.dx, s.dy*0.8,
            s.white)
    }
}

testCardUHD_Slicer_Upper(ctx, s) {
    try( ctx ) {
        gc := ctx.Gc()

        // convert into seconds [0] and frame in second [1]
        frameRate:= math.Float(s.frameRate)
        second := math.Modf( math.Float(s.frameNumber) / frameRate )
		sz := second[1]

        // Because 1 frame is blank adjust sz so that frame (frameRate-1) is 100% rather than just under.
        // This means that on frame (frameRate-1) the slicer is full and on frame (frameRate) it is then blank.
        // Without this you would never see a full slicer
		sz = ( sz * frameRate) / (frameRate-1)

        animGraphic.FillRectangle(gc,
			s.dW+(6.4*s.dx), s.dH+(0.8*s.dy),
			2.2*s.dx, s.dy*0.4,
			s.black)

		if sz > 0.0001 {
			//sz = sz + (1 / frameRate)
			animGraphic.FillRectangle(gc,
				s.dW+(6.4*s.dx)+10, s.dH+(0.8*s.dy+10),
				(2.2*s.dx-20)*sz, s.dy*0.4-20,
				s.lightgreen)
		}
    }
}