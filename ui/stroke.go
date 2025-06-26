package ui

type StrokeDraggable interface {
	MoveBy(dx, dy float64)
	MoveTo(x, y float64)
}

type Stroke struct {
	// initX and initY represents the position when dragging starts.
	initX float64
	initY float64

	prevX float64
	prevY float64

	// currentX and currentY represents the current position
	currentX float64
	currentY float64

	Released bool

	DropEventTaken bool

	DraggingObject StrokeDraggable
}

func NewStroke(cx float64, cy float64, s StrokeDraggable) *Stroke {
	return &Stroke{
		//source:   source,
		initX:          cx,
		initY:          cy,
		prevX:          cx,
		prevY:          cy,
		currentX:       cx,
		currentY:       cy,
		DraggingObject: s,
		DropEventTaken: false,
	}
}

func (s *Stroke) Update(cx, cy float64) {
	if s.Released {
		return
	}
	s.prevX = s.currentX
	s.prevY = s.currentY

	s.currentX = cx
	s.currentY = cy

	s.DraggingObject.MoveBy(s.PositionDiff())
}

func (s *Stroke) Release() {
	s.Released = true
}

func (s *Stroke) PositionDiff() (float64, float64) {
	dx := s.currentX - s.prevX
	dy := s.currentY - s.prevY
	return dx, dy
}
