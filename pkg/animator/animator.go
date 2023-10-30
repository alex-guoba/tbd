package animator

import (
	"time"

	"github.com/briandowns/spinner"
)

type Animator struct {
	spinner *spinner.Spinner
}

func New() *Animator {
	return &Animator{
		spinner: spinner.New(spinner.CharSets[4], 100*time.Millisecond),
	}
}

func (s *Animator) Start() {
	s.spinner.Start()
}

func (s *Animator) Stop() {
	s.spinner.Stop()
}
