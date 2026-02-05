package email

import "context"

type EmailSender interface {
	SendEmail(ctx context.Context, e Email) (err error)
}

type Email struct {
	From             string
	CostumerName     string
	Subject          string
	To               string
	Logo             string
	SelectedTemplate string
	TapLink          string
	WorkOrderStatus  string
	WorkOrderCode    string
	WorkOrderType    string
}

type InternalParams struct {
	Email
	Tmap      map[string]string
	AssetPath string
}
