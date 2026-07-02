package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Event struct {
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id,omitempty"`
	UserRole  string `json:"user_role,omitempty"`
	Action    string `json:"action"`
	Tool      string `json:"tool,omitempty"`
	Status    string `json:"status"`
	Reason    string `json:"reason,omitempty"`
}

func Log(event Event) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)

	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stdout, "AUDIT {\"action\":\"audit_encode_failed\",\"status\":\"error\",\"reason\":%q}\n", err.Error())
		return
	}

	fmt.Fprintf(os.Stdout, "AUDIT %s\n", payload)
}
