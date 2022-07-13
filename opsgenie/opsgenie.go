package opsgenie

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	headerAcceptParam      = "Accept"
	headerContentTypeParam = "Content-Type"
	headerAuthorization    = "Authorization"
	headerJsonValue        = "application/json"
	authCodeParam          = "code"
)

var httpClient = http.Client{}

func GetYesterdayOnCallUsers(accessToken string, channel chan OnCallUsersResponse) {
	var response OnCallUsersResponse

	onCallUrl := "https://api.opsgenie.com/v2/schedules/on-calls?date=" + getYesterdayDate()
	var onCallData OnCallData
	doRequest(onCallUrl, accessToken, &onCallData)
	response = OnCallUsersResponse{
		Main:   OpsGenieUser{EMail: onCallData.Data[0].Participants[0].Name},
		BackUp: OpsGenieUser{EMail: onCallData.Data[0].Participants[1].Name},
	}

	channel <- response
}

func doRequest(url, accessToken string, target interface{}) error {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(headerContentTypeParam, headerJsonValue)
	request.Header.Set(headerAuthorization, "GenieKey "+accessToken)
	response, err := httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("::: Error in HTTP request: %v", err)
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
