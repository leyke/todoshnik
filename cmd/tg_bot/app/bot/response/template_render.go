package response

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"todoshnik/cmd/tg_bot/app/bot/response/models"
)

const (
	addSuccessTemplate  = "templates/add_success.tmpl"
	addRequestTemplate  = "templates/add_request.tmpl"
	helpTemplate        = "templates/help.tmpl"
	listEmptyTemplate   = "templates/list_empty.tmpl"
	listSuccessTemplate = "templates/list_success.tmpl"
	startTemplate       = "templates/start.tmpl"
	stateEmptyTemplate  = "templates/state_empty.tmpl"
	stateDoneTemplate   = "templates/state_done.tmpl"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

func Render(tpl string, data any) (string, error) {
	t, err := template.ParseFS(templatesFS, tpl)
	if err != nil {
		return "", fmt.Errorf("парсинг шаблона: %w", err)
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, data); err != nil {
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
		addRequestTemplate,
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
