package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/htenjo/gh_statistics/config"
	"github.com/htenjo/gh_statistics/github"
)

var (
	DangerStyle   = "danger"
	PrimaryStyle  = "primary"
	helloMessage  = "<!here> :sunny: Buenos días equipo, *¡¡es hora del daily!!*, nuestros facilitadores hoy son:"
	collaborators = []github.Collaborator{
		{Name: "Mary", ID: "U01CENY63K2"},
		{Name: "Angie", ID: "U0370NNMWDA"},
		{Name: "Dario", ID: "U01KGU0TSHZ"},
		{Name: "Pierre", ID: "U01NYBUUV7D"},
		{Name: "Steven", ID: "U01ML2N7QM6"},
	}
)

func NewPlainTextBlock(text string) PlanTextBlock {
	return PlanTextBlock{
		Type: "plain_text",
		Text: text,
	}
}

func NewHeader(text string) HeaderBlock {
	return HeaderBlock{
		Type: "header",
		Text: NewPlainTextBlock(text),
	}
}

func NewActions() ButtonSection {
	return ButtonSection{
		Type:     "actions",
		Elements: []ButtonBlock{},
	}
}

func SendSlackMessage(messageTitle string, prInfo *[]github.RepoPR) {
	message := WebhookMessage{}
	header := NewHeader(messageTitle)
	actions := NewActions()

	redPrs := getButtonsByFlag(prInfo, github.Red)
	yellowPrs := getButtonsByFlag(prInfo, github.Yellow)
	GreenPrs := getButtonsByFlag(prInfo, github.Green)

	actions.Elements = append(actions.Elements, getButtonElements(&redPrs)...)
	actions.Elements = append(actions.Elements, getButtonElements(&yellowPrs)...)
	actions.Elements = append(actions.Elements, getButtonElements(&GreenPrs)...)
	maxNotifications := math.Min(float64(20), float64(len(actions.Elements)))
	actions.Elements = actions.Elements[0:int(maxNotifications)]

	message.Blocks = append(message.Blocks, header, actions)
	byteResponse, _ := json.MarshalIndent(message, "", "  ")
	log.Printf("%v", string(byteResponse))
	sendNotification(byteResponse)
}

func getButtonsByFlag(prInfo *[]github.RepoPR, flag github.PrReviewFlag) []github.PullRequestDetail {
	var redPrs []github.PullRequestDetail

	for _, pr := range *prInfo {
		for _, info := range pr.Prs {
			if info.ReviewFlag == flag {
				redPrs = append(redPrs, info)
			}
		}
	}

	return redPrs
}

func getButtonElements(prDetails *[]github.PullRequestDetail) []ButtonBlock {
	var buttonBlocks []ButtonBlock

	for _, pr := range *prDetails {
		buttonBlock := ButtonBlock{
			Type: "button",
			Text: PlanTextBlock{
				Type: "plain_text",
				Text: pr.Title,
			},
			Url: pr.HtmlUrl,
		}

		if pr.ReviewFlag == github.Red {
			buttonBlock.Style = &DangerStyle
		} else if pr.ReviewFlag == github.Green {
			buttonBlock.Style = &PrimaryStyle
		}

		buttonBlocks = append(buttonBlocks, buttonBlock)
	}

	return buttonBlocks
}

func sendNotification(message []byte) {
	resp, err := http.Post(config.SlackWebhookUrl(), "application/json", bytes.NewReader(message))

	if err != nil {
		log.Fatal(err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(bodyText))
}

func SendDailyReminderMessage() {
	message := WebhookMessage{}
	helloSection := NewMarkDownSection(helloMessage)
	gopherbotsBlock := NewFieldsSection(getTodayGopherbots(config.GetCollaboratorIndex(len(collaborators))))
	linksSection := NewMarkDownSection("Nuestros enlaces...")
	linksBlock := NewFieldsSection(getLinks())

	message.Blocks = append(message.Blocks, helloSection, gopherbotsBlock, linksSection, linksBlock)
	byteResponse, _ := json.MarshalIndent(message, "", "  ")
	sendNotification(byteResponse)
}

func getTodayGopherbots(collaboratorIndex int) []PlanTextBlock {
	alternateIndex := (collaboratorIndex + 2) % len(collaborators)
	fmt.Println("El colaborador es:", collaborators[collaboratorIndex].Name)
	var todayGopherbotsBlock []PlanTextBlock
	todayGopherbotsBlock = append(todayGopherbotsBlock,
		NewMarkDownBlock(fmt.Sprintf("*Presentador:*\n<@%s>", collaborators[collaboratorIndex].ID)),
		NewMarkDownBlock(fmt.Sprintf("*Suplente:*\n<@%s>", collaborators[alternateIndex].ID)),
	)

	return todayGopherbotsBlock
}

func getLinks() []PlanTextBlock {
	var linksBlock []PlanTextBlock

	linksBlock = append(linksBlock,
		NewMarkDownBlock(":jira: *<https://mercadolibre.atlassian.net/jira/software/projects/SDK/boards/3802|Tablero Jira>* :jira:"),
		NewMarkDownBlock(":meet: *<https://meet.google.com/mbs-hgob-fdf|Meet>* :meet:"),
	)

	return linksBlock
}

func NewMarkDownBlock(text string) PlanTextBlock {
	return PlanTextBlock{
		Type: "mrkdwn",
		Text: text,
	}
}

func NewMarkDownSection(text string) HeaderBlock {
	return HeaderBlock{
		Type: "section",
		Text: NewMarkDownBlock(text),
	}
}

func NewFieldsSection(fields []PlanTextBlock) FieldBlock {
	return FieldBlock{
		Type:   "section",
		Fields: fields,
	}
}
