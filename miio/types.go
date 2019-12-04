package miio

type deviceCommand struct {
	ID     int64         `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// copied from https://github.com/vkorn/go-miio/blob/master/vacuum.go

// VacError defines possible vacuum error.
type VacError int

const (
	// VacErrorNo describes no errors.
	VacErrorNo VacError = iota
	// VacErrorCharge describes error with charger.
	VacErrorCharge
	// VacErrorFull describes full dust container.
	VacErrorFull
	// VacErrorUnknown describes unknown error
	VacErrorUnknown
)

type vacuumState int

type GetStatusResponse struct {
	Battery       int `mapstructure:"battery"`
	BeginTime     int `mapstructure:"begin_time"`
	CleanTime     int `mapstructure:"clean_time"`
	CleanArea     int `mapstructure:"clean_area"`
	CleanMode     int `mapstructure:"clean_mode"`
	CleanStrategy int `mapstructure:"clean_strategy"`
	CleanTrigger  int `mapstructure:"clean_trigger"`
	BackTrigger   int `mapstructure:"back_trigger"`
	DNDEnabled    int `mapstructure:"dnd_enabled"`
	FanPower      int `mapstructure:"fan_power"`
	MapPresent    int `mapstructure:"map_present"`
	InCleaning    int `mapstructure:"in_cleaning"`
	Completed     int `mapstructure:"completed"`

	ErrorCode int `mapstructure:"error_code"`

	MsgVer int         `mapstructure:"msg_ver"`
	MsgSeq int         `mapstructure:"msg_seq"`
	State  vacuumState `mapstructure:"state"`
}

func (s vacuumState) String() string {
	switch s {
	case 1:
		return "initiating"
	case 2:
		return "sleeping"
	case 3:
		return "waiting"
	case 5:
		return "cleaning"
	case 6:
		return "returning"
	case 8:
		return "charging"
	case 9:
		return "unknown"
	case 10:
		return "paused"
	case 11:
		return "cleaning_spot"
	case 13:
		return "shutting_down"
	case 14:
		return "updating"
	case 15:
		return "docking"
	case 17:
		return "cleaning_zone"
	case 100:
		return "dust_bag_full"
	default:
		return "unknown"
	}
}

type GetConsumableResponse struct {
	MainBrushWorkTime int `mapstructure:"main_brush_work_time"`
	SideBrushWorkTime int `mapstructure:"side_brush_work_time"`
	FilterWorkTime    int `mapstructure:"filter_work_time"`
	SensorDirtyTime   int `mapstructure:"sensor_dirty_time"`
}
