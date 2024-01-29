package types

type Order struct {
	RequestId  string     `json:"requestId"`
	Timestamp  int64      `json:"timestamp"`
	FleetId    string     `json:"fleetId"`
	DeviceUUID string     `json:"deviceUUID"`
	Intent     string     `json:"intent"`
	Event      CloudEvent `json:"event"`
}
type CloudEvent struct {
	Name  string          `json:"name"`
	Value CloudEventValue `json:"value"`
}
type CloudEventValue struct {
	ResponseStatus int    `json:"responseStatus"`
	ResponseMsg    string `json:"responseMsg"`
	ExecStatus     int    `json:"execStatus"`
	ExecMsg        string `json:"execMsg"`
}
