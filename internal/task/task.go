package task

type Task struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Retries int                    `json:"retries"`
}
