package app

type Config struct {
	Generators []struct {
		TimeoutS    int `json:"timeout_s"`
		SendPeriodS int `json:"send_period_s"`
		DataSources []struct {
			Id            string `json:"id"`
			InitValue     int    `json:"init_value"`
			MaxChangeStep int    `json:"max_change_step"`
		} `json:"data_sources"`
	} `json:"generators"`
	Queue struct {
		Size int `json:"size"`
	} `json:"queue"`
	Aggregators []struct {
		SubIds           []string `json:"sub_ids"`
		AggregatePeriodS int      `json:"aggregate_period_s"`
	} `json:"aggregators"`
	StorageType int `json:"storage_type"`
}
