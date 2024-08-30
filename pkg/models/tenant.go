package models

type Tenant struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	ApiKey         string                 `json:"api_key"`
	ActiveSources  []string               `json:"active_sources"`
	Configurations map[string]interface{} `json:"configurations"`
}
