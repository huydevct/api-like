package constants

type (
	// ServiceKind : loại bảo hành
	ServiceKind int
)

// Service kind
const (
	NormalKind     ServiceKind = 1
	WeekKind       ServiceKind = 7
	MonthKind      ServiceKind = 30
	ThreeMonthKind ServiceKind = 90
	FreeKind       ServiceKind = 99
)
