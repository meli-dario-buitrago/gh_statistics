package opsgenie

import (
	"fmt"
	"time"
)

type Collaborator struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	EMail string `json:"email"`
}

type OnCallUsersResponse struct {
	Main   OpsGenieUser
	BackUp OpsGenieUser
}

type OpsGenieUser struct {
	EMail string
}

type OnCallData struct {
	Data []OnCallDetail `json:"data"`
}

type OnCallDetail struct {
	Participants []OnCallParticipants `json:"onCallParticipants"`
}

type OnCallParticipants struct {
	Name string `json:"name"`
}

func getYesterdayDate() string {
	t2 := time.Now().AddDate(0, 0, -1)
	fmt.Println("===== date =====\n", t2.Format(time.RFC3339))
	return t2.Format(time.RFC3339)
}
