package config

import "fmt"

type Config struct {
	MaxWindows                int               `toml:"max_windows"`
	NameMaxLength             int               `toml:"name_max_length"`
	NamePadding               int               `toml:"name_padding"`
	SeparatorString           string            `toml:"separator_string"`
	EmptyDesktopString        string            `toml:"empty_desktop_string"`
	ActiveWindowLeftClick     string            `toml:"active_window_left_click"`
	ActiveWindowMiddleClick   string            `toml:"active_window_middle_click"`
	ActiveWindowRightClick    string            `toml:"active_window_right_click"`
	ActiveWindowScrollUp      string            `toml:"active_window_scroll_up"`
	ActiveWindowScrollDown    string            `toml:"active_window_scroll_down"`
	InactiveWindowLeftClick   string            `toml:"inactive_window_left_click"`
	InactiveWindowMiddleClick string            `toml:"inactive_window_middle_click"`
	InactiveWindowRightClick  string            `toml:"inactive_window_right_click"`
	InactiveWindowScrollUp    string            `toml:"inactive_window_scroll_up"`
	InactiveWindowScrollDown  string            `toml:"inactive_window_scroll_down"`
	HiddenWindowLeftClick     string            `toml:"hidden_window_left_click"`
	HiddenWindowMiddleClick   string            `toml:"hidden_window_middle_click"`
	HiddenWindowRightClick    string            `toml:"hidden_window_right_click"`
	HiddenWindowScrollUp      string            `toml:"hidden_window_scroll_up"`
	HiddenWindowScrollDown    string            `toml:"hidden_window_scroll_down"`
	ActiveWindowFgColor       string            `toml:"active_window_fg_color"`
	ActiveWindowBgColor       string            `toml:"active_window_bg_color"`
	ActiveWindowUlColor       string            `toml:"active_window_ul_color"`
	InactiveWindowFgColor     string            `toml:"inactive_window_fg_color"`
	InactiveWindowBgColor     string            `toml:"inactive_window_bg_color"`
	InactiveWindowUlColor     string            `toml:"inactive_window_ul_color"`
	HiddenWindowFgColor       string            `toml:"hidden_window_fg_color"`
	HiddenWindowBgColor       string            `toml:"hidden_window_bg_color"`
	HiddenWindowUlColor       string            `toml:"hidden_window_ul_color"`
	SeparatorFgColor          string            `toml:"separator_fg_color"`
	SeparatorBgColor          string            `toml:"separator_bg_color"`
	SeparatorUlColor          string            `toml:"separator_ul_color"`
	EmptyDesktopFgColor       string            `toml:"empty_desktop_fg_color"`
	EmptyDesktopBgColor       string            `toml:"empty_desktop_bg_color"`
	EmptyDesktopUlColor       string            `toml:"empty_desktop_ul_color"`
	Flags                     Flags             `toml:"flags"`
	IgnoredClasses            []string          `toml:"ignored_classes"`
	WindowNicknames           map[string]string `toml:"window_nicknames"`
}

type Flags struct {
	StickyFlag  string `toml:"sticky"`
	LockedFlag  string `toml:"locked"`
	PrivateFlag string `toml:"private"`
	MarkedFlag  string `toml:"marked"`
}

func (c Config) GetBgColor(isActive, isHidden bool) string {
	color := c.InactiveWindowBgColor
	if isActive && c.ActiveWindowBgColor != "" {
		color = c.ActiveWindowBgColor
	}
	if isHidden && c.HiddenWindowBgColor != "" {
		color = c.HiddenWindowBgColor
	}

	return color
}

func (c Config) GetFgColor(isActive, isHidden bool) string {
	color := c.InactiveWindowFgColor
	if isActive && c.ActiveWindowFgColor != "" {
		color = c.ActiveWindowFgColor
	}
	if isHidden && c.HiddenWindowFgColor != "" {
		color = c.HiddenWindowFgColor
	}

	return color
}

func (c Config) GetUlColor(isActive, isHidden bool) string {
	color := c.InactiveWindowUlColor
	if isActive && c.ActiveWindowUlColor != "" {
		color = c.ActiveWindowUlColor
	}
	if isHidden && c.HiddenWindowUlColor != "" {
		color = c.HiddenWindowUlColor
	}

	return color
}

func (c Config) GetActionLeftClick(isActive, isHidden bool) string {
	if isActive {
		return c.ActiveWindowLeftClick
	}
	if isHidden {
		return c.HiddenWindowLeftClick
	}
	return c.InactiveWindowLeftClick
}

func (c Config) GetActionRightClick(isActive, isHidden bool) string {
	if isActive {
		return c.ActiveWindowRightClick
	}
	if isHidden {
		return c.HiddenWindowRightClick
	}
	return c.InactiveWindowRightClick
}

func (c Config) GetActionMiddleClick(isActive, isHidden bool) string {
	if isActive {
		return c.ActiveWindowMiddleClick
	}
	if isHidden {
		return c.HiddenWindowMiddleClick
	}
	return c.InactiveWindowMiddleClick
}

func (c Config) GetActionScrollUp(isActive, isHidden bool) string {
	if isActive {
		return c.ActiveWindowScrollUp
	}
	if isHidden {
		return c.HiddenWindowScrollUp
	}
	return c.InactiveWindowScrollUp
}

func (c Config) GetActionScrollDown(isActive, isHidden bool) string {
	if isActive {
		return c.ActiveWindowScrollDown
	}
	if isHidden {
		return c.HiddenWindowScrollDown
	}
	return c.InactiveWindowScrollDown
}

func FormatStringColors(string, bgColor, fgColor, ulColor string) string {
	if bgColor != "" {
		string = fmt.Sprintf("%%{B%v}%v%%{B-}", bgColor, string)
	}
	if fgColor != "" {
		string = fmt.Sprintf("%%{F%v}%v%%{F-}", fgColor, string)
	}
	if ulColor != "" {
		string = fmt.Sprintf("%%{u%v}%%{+u}%v%%{-u}", ulColor, string)
	}

	return string
}
