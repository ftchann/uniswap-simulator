package result

type Snapshot struct {
	Timestamp int    `json:"timestamp"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
	AmountUSD string `json:"amountUSD"`
	Price     string `json:"price"`
}
type Save struct {
	UpdateInterval int         `json:"update_interval"`
	StartAmount    string      `json:"start_amount"`
	StartTime      int         `json:"start_time"`
	EndTime        int         `json:"end_time"`
	Results        []RunResult `json:"results"`
}

type RunResult struct {
	ParameterA int `json:"parameterA"`
	//ParameterB    int    `json:"parameterB"`
	//MultiplierX10  int    `json:"multiplierX10"`
	//HistoryWindow  int    `json:"history_window"`
	EndAmount               string  `json:"end_amount"`
	StandardDeviationHourly float64 `json:"standard_deviation_hourly"`
	StandardDeviationDaily  float64 `json:"standard_deviation_daily"`
	StandardDeviationWeekly float64 `json:"standard_deviation_weekly"`
}
