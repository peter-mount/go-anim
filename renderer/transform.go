package renderer

func (c *context) Rotate(angle float64) Context {
	c.Gc().Rotate(angle)
	return c
}

func (c *context) Translate(tx, ty float64) Context {
	c.Gc().Translate(tx, ty)
	return c
}

func (c *context) Scale(sx, sy float64) Context {
	c.Gc().Scale(sx, sy)
	return c
}
