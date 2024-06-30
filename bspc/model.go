package bspc

const (
	LayoutTypeTiled   LayoutType = "tiled"
	LayoutTypeMonocle LayoutType = "monocle"

	SplitTypeHorizontal SplitType = "horizontal"
	SplitTypeVertical   SplitType = "vertical"

	DirectionTypeUp    DirectionType = "north"
	DirectionTypeDown  DirectionType = "south"
	DirectionTypeLeft  DirectionType = "west"
	DirectionTypeRight DirectionType = "east"

	StateTypeTiled       StateType = "tiled"
	StateTypePseudoTiled StateType = "pseudo_tiled"
	StateTypeFloating    StateType = "floating"
	StateTypeFullscreen  StateType = "fullscreen"

	FlagTypeHidden  FlagType = "hidden"
	FlagTypeSticky  FlagType = "sticky"
	FlagTypePrivate FlagType = "private"
	FlagTypeLocked  FlagType = "locked"
	FlagTypeMarked  FlagType = "marked"
	FlagTypeUrgent  FlagType = "urgent"

	RelativePositionTypeAbove RelativePositionType = "above"
	RelativePositionTypeBelow RelativePositionType = "below"

	LayerTypeAbove  LayerType = "above"
	LayerTypeNormal LayerType = "normal"
	LayerTypeBelow  LayerType = "below"

	PointerActionTypeMove         PointerActionType = "move"
	PointerActionTypeResizeCorner PointerActionType = "resize_corner"
	PointerActionTypeResizeSide   PointerActionType = "resize_side"

	PointerActionStateTypeBegin PointerActionStateType = "begin"
	PointerActionStateTypeEnd   PointerActionStateType = "end"
)

type (
	LayoutType             string
	SplitType              string
	DirectionType          string
	StateType              string
	FlagType               string
	RelativePositionType   string
	LayerType              string
	PointerActionType      string
	PointerActionStateType string
)

type Monitor struct {
	// RandRID is the monitor's ID in the RandR tool.
	RandRID          int       `json:"randrId"`
	ID               ID        `json:"id"`
	Name             string    `json:"name"`
	Wired            bool      `json:"wired"`
	StickyCount      int       `json:"stickyCount"`
	WindowGap        int       `json:"windowGap"`
	BorderWidth      int       `json:"borderWidth"`
	FocusedDesktopID ID        `json:"focusedDesktopId"`
	Padding          padding   `json:"padding"`
	Rectangle        rectangle `json:"rectangle"`
	Desktops         []Desktop `json:"desktops"`
}

// Desktop contains all the info regarding a given desktop.
type Desktop struct {
	Name          string     `json:"name"`
	ID            ID         `json:"id"`
	Layout        LayoutType `json:"layout"`
	UserLayout    LayoutType `json:"userLayout"`
	WindowGap     int        `json:"windowGap"`
	BorderWidth   int        `json:"borderWidth"`
	FocusedNodeID ID         `json:"focusedNodeId"`
	Padding       padding    `json:"padding"`
	Root          Node       `json:"root"`
}

// Node contains all the info regarding a given node.
type Node struct {
	ID          ID             `json:"id"`
	SplitType   SplitType      `json:"SplitType"`
	SplitRatio  float64        `json:"splitRatio"`
	Vacant      bool           `json:"vacant"`
	Hidden      bool           `json:"hidden"`
	Sticky      bool           `json:"sticky"`
	Private     bool           `json:"private"`
	Locked      bool           `json:"locked"`
	Marked      bool           `json:"marked"`
	Preselect   *NodePreselect `json:"presel"`
	Rectangle   rectangle      `json:"rectangle"`
	Constraints constraints    `json:"constraints"`
	FirstChild  *Node          `json:"firstChild"`
	SecondChild *Node          `json:"secondChild"`
	Client      *NodeClient    `json:"client"`
}

// NodePreselect contains all the infor regarding a node's preselection state.
type NodePreselect struct {
	SplitDirection DirectionType `json:"splitDir"`
	SplitRatio     float64       `json:"splitRatio"`
}

// NodeClient contains all the info regarding a node's client. The program it contains.
type NodeClient struct {
	ClassName         string    `json:"className"`
	InstanceName      string    `json:"instanceName"`
	BorderWidth       int       `json:"borderWidth"`
	State             StateType `json:"state"`     // TODO: Add validation for this in the GetState method
	LastState         StateType `json:"lastState"` // TODO: Add validation for this in the GetState method
	Layer             LayerType `json:"layer"`     // TODO: Add validation for this in the GetState method
	LastLayer         LayerType `json:"lastLayer"` // TODO: Add validation for this in the GetState method
	Urgent            bool      `json:"urgent"`
	Shown             bool      `json:"shown"`
	TiledRectangle    rectangle `json:"tiledRectangle"`
	FloatingRectangle rectangle `json:"floatingRectangle"`
}

type rectangle struct {
	X      int `json:"x"`
	Y      int `json:"Y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type constraints struct {
	MinWidth  int `json:"min_width"`
	MinHeight int `json:"min_height"`
}

type padding struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}
