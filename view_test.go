package tvision_test

import (
	//	. "github.com/bbodi/tvision"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("ComboBox", func() {
	var ()

	BeforeEach(func() {
	})

	Describe("Drawing View", func() {
		Context("when a view without own DrawBuffer are Executed in anotherView", func() {
			It("should use the parent's DrawBuffer", func() {
				/*view := new(View)
				view.Widget = new(ComboBox)
				view.SetDrawBuffer(new(DrawBuffer))
				child := new(View)
				child.Widget = new(ComboBox)
				child.SetDrawBuffer(nil)
				view.ExecuteView(0, 0, child)*/
				// shouldnt throw NPE
			})
		})

	})
})
