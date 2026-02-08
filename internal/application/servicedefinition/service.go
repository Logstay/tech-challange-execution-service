package servicedefinition

import (
	"context"

	cmds "github.com/Logstay/tech-challange-wherehouse-service/internal/application/servicedefinition/commands"
	qries "github.com/Logstay/tech-challange-wherehouse-service/internal/application/servicedefinition/queries"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type Service struct {
	// commands
	saveServiceDefinitionCommand   *cmds.SaveServiceDefinitionHandler
	updateServiceDefinitionCommand *cmds.UpdateServiceDefinitionHandler
	deleteServiceDefinitionCommand *cmds.DeleteServiceDefinitionHandler

	// queries
	getServiceDefinitionQuery  *qries.GetServiceDefinitionHandler
	listServiceDefinitionQuery *qries.ListServiceDefinitionHandler

	d Dependencies
}

type Dependencies struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func ServiceDefinitionService(d Dependencies) (s *Service) {
	s = &Service{
		d: d,
	}

	s.buildCommands()

	s.buildQueries()

	return
}

// commands
func (s *Service) SaveServiceDefinition(ctx context.Context,
	cmd ports.SaveServiceDefinition) (err error) {
	return s.saveServiceDefinitionCommand.Handle(ctx, cmd)
}

func (s *Service) UpdateServiceDefinition(ctx context.Context,
	cmd ports.UpdateServiceDefinition) (err error) {
	return s.updateServiceDefinitionCommand.Handle(ctx, cmd)
}

func (s *Service) DeleteServiceDefinition(ctx context.Context,
	uid string) (err error) {
	return s.deleteServiceDefinitionCommand.Handle(ctx, uid)
}

// queries
func (s *Service) GetServiceDefinitionByUID(ctx context.Context,
	uid string) (result ports.GetServiceDefinition, err error) {
	return s.getServiceDefinitionQuery.Handle(ctx, uid)
}

func (s *Service) ListServiceDefinition(ctx context.Context,
	filter ports.Filter) (result []ports.GetServiceDefinition, total int64, err error) {
	return s.listServiceDefinitionQuery.Handle(ctx, filter)
}

func (s *Service) buildCommands() {
	s.buildSaveServiceDefinition()

	s.buildUpdateServiceDefinition()

	s.buildDeleteServiceDefinition()
}

func (s *Service) buildSaveServiceDefinition() {
	s.saveServiceDefinitionCommand = &cmds.SaveServiceDefinitionHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildUpdateServiceDefinition() {
	s.updateServiceDefinitionCommand = &cmds.UpdateServiceDefinitionHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildDeleteServiceDefinition() {
	s.deleteServiceDefinitionCommand = &cmds.DeleteServiceDefinitionHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildQueries() {

	s.buildGetServiceDefinition()

	s.buildListServiceDefinition()
}

func (s *Service) buildGetServiceDefinition() {
	s.getServiceDefinitionQuery = &qries.GetServiceDefinitionHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildListServiceDefinition() {
	s.listServiceDefinitionQuery = &qries.ListServiceDefinitionHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}
