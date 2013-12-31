package tvision_test

import (
	. "github.com/bbodi/tvision"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ComboBox", func() {
	var (
		combo *ComboBox
		view  *View
	)

	BeforeEach(func() {
		view, combo = CreateComboBox(10)
	})

	Describe("Handling events", func() {
		Context("when an alphanumeric value was written in an Editable Box", func() {
			It("should clear the event", func() {
				combo.Editable = true
				event := Event{Type: EvKey, Key: 0, Ch: 'a'}
				combo.HandleEvent(&event, view)
				Expect(event.Type).To(Equal(EvNothing))
			})
		})

	})
})
