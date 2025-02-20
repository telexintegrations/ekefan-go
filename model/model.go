package model

type IntegrationConfig struct {
	Data IntegrationData `json:"data"`
}

type IntegrationDescriptions struct {
	AppDescription  string `json:"app_description"`
	AppLogo         string `json:"app_logo"`
	AppName         string `json:"app_name"`
	AppURL          string `json:"app_url"`
	BackgroundColor string `json:"background_color"`
}

type IntegrationDates struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IntegrationOutputs struct {
	Label string `json:"label"`
	Value bool   `json:"value"`
}

type IntegrationKeyFeatures struct {
	KeyFeatures []string `json:"key_features"`
}

type IntegrationPermissions struct {
	MonitoringUser MonitoringUser `json:"monitoring_user"`
}
type MonitoringUser struct {
	AlwaysOnline bool   `json:"always_online"`
	DisplayName  string `json:"display_name"`
}

type IntegrationSettings struct {
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Default  string   `json:"default"`
	Options  []string `json:"options,omitempty"`
}

type IntegrationData struct {
	Date                IntegrationDates        `json:"date"`
	Descriptions        IntegrationDescriptions `json:"descriptions"`
	IntegrationCategory string                  `json:"integration_category"`
	IntegrationType     string                  `json:"integration_type"`
	IsActive            bool                    `json:"is_active"`
	Output              []IntegrationOutputs    `json:"output"`
	Author              string                  `json:"author"`
	KeyFeatures         []string                `json:"key_features"`
	Permissions         IntegrationPermissions  `json:"permissions"`
	Settings            []IntegrationSettings   `json:"settings"`
	TickURL             string                  `json:"tick_url"`
	TargetURL           string                  `json:"target_url"`
}

type TelexRequestPayload struct {
	ChannelID string                `json:"channel_id"`
	ReturnURL string                `json:"return_url"`
	Settings  []IntegrationSettings `json:"settings"`
}

type TelexResponsePayload struct {
	Message   string `json:"message"`
	Username  string `json:"username"` // the name of your integration
	EventName string `json:"event_name"`
	Status    string `json:"status"`
}

type TelexErrMsg struct {
	ErrMsg string `json:"err_msg"`
}
