package tvision

var (
	TextColor          = ColorWhite
	FocusedTextColor   = ColorWhite
	TitleColor         = ColorWhite
	FocusedTitleColor  = ColorYellow
	FrameColor         = ColorBlue
	FocusedFrameColor  = ColorBlue
	BorderColor        = ColorWhite
	FocusedBorderColor = ColorYellow

	cItemFgColor       = ColorWhite
	FocusedItemFgColor = ColorBlack
	cItemBgColor       = ColorBlack
	FocusedItemBgColor = ColorWhite

	ComboFrameColor        = ColorGreen
	FocusedComboFrameColor = ColorRed

	EnabledBorderActionColor  = ColorGreen
	DisabledBorderActionColor = ColorGray
)

func (self *View) ComboFrameColor() Color {
	if self.focused {
		return FocusedComboFrameColor
	} else {
		return ComboFrameColor
	}
}

func (self *View) TextColor() Color {
	if self.focused {
		return FocusedTitleColor
	} else {
		return TitleColor
	}
}

func (self *View) TitleColor() Color {
	if self.focused {
		return FocusedTitleColor
	} else {
		return TitleColor
	}
}

func ItemFgColor(selected bool) Color {
	if selected {
		return FocusedItemFgColor
	} else {
		return cItemFgColor
	}
}

func ItemBgColor(selected bool) Color {
	if selected {
		return FocusedItemBgColor
	} else {
		return cItemBgColor
	}
}

func (self *View) FrameColor() Color {
	if self.focused {
		return FocusedFrameColor
	} else {
		return FrameColor
	}
}

func BorderActionColor(enabled bool) Color {
	if enabled {
		return EnabledBorderActionColor
	} else {
		return DisabledBorderActionColor
	}
}

func (self *View) BorderColor() Color {
	if self.focused {
		return FocusedBorderColor
	} else {
		return BorderColor
	}
}
