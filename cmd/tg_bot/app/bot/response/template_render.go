package response

import (
	"bytes"
	"fmt"
	"text/template"

	"todoshnik/cmd/tg_bot/app/bot/response/models"
)

const (
	addSuccessTemplate  = "cmd/tg_bot/app/bot/response/templates/add_success.tmpl"
	addRequestTemplate  = "cmd/tg_bot/app/bot/response/templates/add_request.tmpl"
	helpTemplate        = "cmd/tg_bot/app/bot/response/templates/help.tmpl"
	listEmptyTemplate   = "cmd/tg_bot/app/bot/response/templates/list_empty.tmpl"
	listSuccessTemplate = "cmd/tg_bot/app/bot/response/templates/list_success.tmpl"
	startTemplate       = "cmd/tg_bot/app/bot/response/templates/start.tmpl"
	stateEmptyTemplate  = "cmd/tg_bot/app/bot/response/templates/state_empty.tmpl"
	stateDoneTemplate   = "cmd/tg_bot/app/bot/response/templates/state_done.tmpl"
)

func Render(tpl string, data any) (string, error) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		return "", fmt.Errorf("парсинг шаблона: %w", err)
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("рендер шаблона: %w", err)
	}

	return buf.String(), nil
}

func RenderAddSuccess(title string) (string, error) {
	return Render(
		addSuccessTemplate,
		models.AddSuccessData{
			Title: title,
		},
	)
}

func RenderAddRequest() (string, error) {
	return Render(
		addSuccessTemplate,
		models.AddRequestData{},
	)
}

func RenderStart() (string, error) {
	return Render(
		startTemplate,
		models.StartData{},
	)
}

func RenderHelp() (string, error) {
	return Render(
		helpTemplate,
		models.HelpData{},
	)
}

func RenderListEmpty() (string, error) {
	return Render(
		listEmptyTemplate,
		models.TaskListEmptyData{},
	)
}

func RenderListSuccess() (string, error) {
	return Render(
		listSuccessTemplate,
		models.TaskListSuccessData{},
	)
}

func RenderStateEmpty() (string, error) {
	return Render(
		stateEmptyTemplate,
		models.StateEmptyData{},
	)
}

func RenderStateDone() (string, error) {
	return Render(
		stateDoneTemplate,
		models.StateDoneData{},
	)
}
