package main

type StaticData struct {
	Data map[string]struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"data"`
}

type StaticData2 struct {
	Data map[string]struct {
		Name string `json:"name"`
	} `json:"data"`
}

type Gameflow struct {
	Gamedata struct {
		Queue struct {
			GameMode string `json:"gameMode"`
		} `json:"queue"`
	} `json:"gameData"`
}

type APIData struct {
	Data []struct {
		Role  string `json:"role"`
		Stats struct {
			WR_Shards struct {
				Build []int `json:"build"`
			} `json:"rune_stat_shards"`
			Shards struct {
				Build []int `json:"build"`
			} `json:"most_common_rune_stat_shards"`
			WR_Runes struct {
				Build []int `json:"build"`
			} `json:"runes"`
			Runes struct {
				Build []int `json:"build"`
			} `json:"most_common_runes"`
			WR_Build struct {
				Build []int `json:"build"`
			} `json:"big_item_builds"`
			Build struct {
				Build []int `json:"build"`
			} `json:"most_common_big_item_builds"`
			WR_Starting_items struct {
				Build []int `json:"build"`
			} `json:"starting_items"`
			Starting_items struct {
				Build []int `json:"build"`
			} `json:"most_common_starting_items"`
			Spells struct {
				Build []int `json:"build"`
			} `json:"spells"`
			Games int `json:"games"`
		} `json:"stats"`
	}
}

type RunePayload struct {
	Name      string `json:"name"`
	Primary   int    `json:"primaryStyleId"`
	Runes     []int  `json:"selectedPerkIds"`
	Secondary int    `json:"subStyleId"`
}
