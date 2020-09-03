package config

type LiveConfig struct {
	APIBaseURL string `json:"api_base_url"`
	APIVersion string `json:"api_version"`
	DynoUser   string `json:"dyno_user"`
	DynoPw     string
	Data       []struct {
		Redirect struct {
			RedirectFrom string `json:"redirect_from"`
			PathMap      []struct {
				From string `json:"from"`
				To   string `json:"to"`
			} `json:"path_map"`
			RedirectTo  string `json:"redirect_to"`
			Promotable  bool   `json:"promotable"`
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"redirect"`
		Cases []struct {
			Call     string `json:"call"`
			Expected string `json:"expected"`
		} `json:"cases"`
	} `json:"data"`
}
