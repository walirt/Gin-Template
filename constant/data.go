package constant

type BaseObject struct {
	Guid     string `json:"guid"`
	Tag      string `json:"tag"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	NodeType int    `json:"node_type"`
}

type PointObject struct {
	BaseObject
	PointType  int     `json:"node_type"`
	Unit       string  `json:"unit"`
	StatusMap  string  `json:"status_map"`
	AlarmLevel int     `json:"alarm_level"`
	AlarmType  int     `json:"alarm_type"`
	Period     int     `json:"period"`
	Percentage float64 `json:"percentage"`
	AbsValue   float64 `json:"abs_value"`
	AoBound    string  `json:"ao_bound"`
}

type DeviceObject struct {
	BaseObject
	DeviceType int           `json:"device_type"`
	Nodes      []PointObject `json:"nodes"`
}

type SpaceObject struct {
	BaseObject
	SpaceType int            `json:"space_type"`
	Nodes     []DeviceObject `json:"nodes"`
}
