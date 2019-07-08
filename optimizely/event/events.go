package event

type EventContext struct {
	Revision        string  `json:"revision"`
	AccountID       string  `json:"account_id"`
	ClientVersion   string  `json:"client_version"`
	ProjectID       string  `json:"project_id"`
	ClientName      string  `json:"client_name"`
	AnonymizeIP     bool    `json:"anonymize_ip"`
}

type UserEvent struct {
	Timestamp int64 `json:"timestamp"`
	Uuid      string `json:"uuid"`
	EventContext EventContext
	VisitorID string
	Impression *ImpressionEvent
	Conversion *ConversionEvent
}

type ImpressionEvent struct {
	EntityID  string `json:"entity_id"`
	Key       string `json:"key"`
	Attributes 	 []VisitorAttribute
	VariationID  string `json:"variation_id"`
	CampaignID   string `json:"campaign_id"`
	ExperimentID string `json:"experiment_id"`
}

type ConversionEvent struct {
	EntityID  string `json:"entity_id"`
	Key       string `json:"key"`
	Attributes 	 []VisitorAttribute
	Tags      map[string]interface{} `json:"tags,omitempty"`
	Revenue   int `json:"revenue,omitempty"`
	Value     float32 `json:"value,omitempty"`
}


// Context about the event
type EventBatch struct {
	Revision        string  `json:"revision"`
	AccountID       string  `json:"account_id"`
	ClientVersion   string  `json:"client_version"`
	Visitors        []Visitor `json:"visitors"`
	ProjectID       string  `json:"project_id"`
	ClientName      string  `json:"client_name"`
	AnonymizeIP     bool    `json:"anonymize_ip"`
	EnrichDecisions bool    `json:"enrich_decisions"`
}

type Visitor struct {
	Attributes []VisitorAttribute `json:"attributes"`
	Snapshots  []Snapshot         `json:"snapshots"`
	VisitorID  string             `json:"visitor_id"`
}

type VisitorAttribute struct {
	Value         interface{} `json:"value"`
	Key           string	`json:"key"`
	AttributeType string	`json:"type"`
	EntityID      string	`json:"entity_id"`
}

type Snapshot struct {
	Decisions []Decision `json:"decisions"`
 	Events []DispatchEvent `json:"events"`
}

type Decision struct {
	VariationID  string `json:"variation_id"`
	CampaignID   string `json:"campaign_id"`
	ExperimentID string `json:"experiment_id"`
}

type DispatchEvent struct {
	EntityID  string `json:"entity_id"`
	Key       string `json:"key"`
	Timestamp int64 `json:"timestamp"`
	Uuid      string `json:"uuid"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
	Revenue   int `json:"revenue,omitempty"`
	Value     float32 `json:"value,omitempty"`
}