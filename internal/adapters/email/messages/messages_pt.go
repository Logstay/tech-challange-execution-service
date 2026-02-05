package messages

var MessagesPT = map[string]string{
	"EmailHi": `Olá, {{ .CostumerName }} <br />temos atualizações sobre sua Ordem de Serviço!`,
	"EmailInviteMessage": `Sua Ordem de Serviço <span style="color: #3b82f6;">{{ .From }} Código: {{ .WorkOrderCode }}</span> teve o status atualizado para <br/> <span style="color: #1e40af;">"{{ .WorkOrderStatus }}"</span> em nosso sistema.
								<p>Para acompanhar o progresso completo e ver todos os detalhes, acesse o portal do cliente clicando no botão abaixo.<p/>`,
	"EmailAccessPortalButton": `<a href="{{.TapLink}}" class="f-fallback button" style="color: #ffffff; background: #3b82f6;" target="_blank">{{.WorkOrderType}}</a> <br />`,
	"EmailFooter":             `<span style="font-size: 16px; color: #242424;">Atenciosamente,<br />Equipe Tech Challenge.</span>`,
	"EmailSubjectFrom":        `Ordem de Serviço #%s - Atualização de Status`,
}
