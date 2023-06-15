package util

import "math"

const (
	Width4K  = 3840            // 4K resolution
	Height4K = 2160            // 4K resolution
	ToRad    = math.Pi / 180.0 // Degrees to Radians
)

// Renderer priorities
const (
	TestCardBasePriority   = 1100 // TestCards background
	TestCardLowerPriority  = 1125 // TestCards mid-ground
	TestCardUpperPriority  = 1175 // TestCards mid-ground
	TestCardTopPriority    = 1199 // TestCards foreground
	DialBackgroundPriority = 2000 // Clock dial background
	ClockPriority          = 2040 // Countdown Clock
	DialForegroundPriority = 2050 // Clock dial foreground
	TitlePriority          = 2100 // Title
	BoxPriority            = 3000 // Box and derivatives if not implicitly defined
)
