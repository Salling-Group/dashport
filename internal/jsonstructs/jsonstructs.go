// Package jsonstructs holds structs for handling dashboard json
package jsonstructs

// Configuration json type
type Configuration struct {
	APIConfig struct {
		Tenants []struct {
			Env   string `json:"env"`
			Token string `json:"token"`
			URL   string `json:"url"`
		} `json:"Tenants"`
	} `json:"apiConfig"`
}

// ReqParts struct to build request
type ReqParts struct {
	Action string
	Env    string
	Denv   string
	ID     string
	Did    string
	URL    string
	Token  string
	Method string
	Conf   Configuration
}

// DashboardsAll json type
type DashboardsAll struct {
	Dashboards []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Owner string `json:"owner"`
	} `json:"dashboards"`
}

// DashSuccess json type
type DashSuccess struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DashError json type
type DashError struct {
	Error struct {
		Code                 int    `json:"code"`
		Message              string `json:"message"`
		ConstraintViolations []struct {
			Path              string `json:"path"`
			Message           string `json:"message"`
			ParameterLocation string `json:"parameterLocation"`
			Location          string `json:"location"`
		} `json:"constraintViolations"`
	} `json:"error"`
}

// Dashboard json type - id is left out since it must be excluded when creating a new one from a clone
type Dashboard struct {
	Metadata struct {
		ConfigurationVersions []int  `json:"-"`
		ClusterVersion        string `json:"-"`
	} `json:"metadata"`
	ID                string `json:"id,omitempty"`
	DashboardMetadata struct {
		Name           string `json:"name"`
		Shared         bool   `json:"shared"`
		Owner          string `json:"owner"`
		SharingDetails struct {
			LinkShared bool `json:"linkShared"`
			Published  bool `json:"published"`
		} `json:"sharingDetails"`
		DashboardFilter struct {
			Timeframe      string      `json:"timeframe"`
			ManagementZone interface{} `json:"managementZone"`
		} `json:"dashboardFilter"`
	} `json:"dashboardMetadata"`
	Tiles []struct {
		Name       string `json:"name"`
		TileType   string `json:"tileType"`
		Configured bool   `json:"configured"`
		Bounds     struct {
			Top    int `json:"top"`
			Left   int `json:"left"`
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"bounds"`
		TileFilter struct {
			Timeframe      interface{} `json:"timeframe"`
			ManagementZone interface{} `json:"managementZone"`
		} `json:"tileFilter"`
		FilterConfig interface{} `json:"filterConfig,omitempty"`
		AssignedEntities []string    `json:"assignedEntities,omitempty"`
		ChartVisible bool        `json:"chartVisible,omitempty"`
		Markdown     string      `json:"markdown,omitempty"`
	} `json:"tiles"`
}

// DashboardWithID is uesd for updating dashboard where destination id is needed.
type DashboardWithID struct {
	*Dashboard
	ID string `json:"id"`
}
