package consumable

import (
	"context"

	cmds "github.com/Logstay/tech-challange-wherehouse-service/internal/application/consumable/commands"
	qries "github.com/Logstay/tech-challange-wherehouse-service/internal/application/consumable/queries"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type Service struct {
	// commands
	saveConsumableCommand   *cmds.SaveConsumableHandler
	updateConsumableCommand *cmds.UpdateConsumableHandler
	deleteConsumableCommand *cmds.DeleteConsumableHandler

	// queries
	getConsumableQuery  *qries.GetConsumableHandler
	listConsumableQuery *qries.ListConsumableHandler

	d Dependencies
}

type Dependencies struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func ConsumableService(d Dependencies) (s *Service) {
	s = &Service{
		d: d,
	}

	s.buildCommands()

	s.buildQueries()

	return
}

// commands
func (s *Service) SaveConsumable(ctx context.Context, cmd ports.SaveConsumable) (err error) {
	return s.saveConsumableCommand.Handle(ctx, cmd)
}

func (s *Service) UpdateConsumable(ctx context.Context, cmd ports.UpdateConsumable) (err error) {
	return s.updateConsumableCommand.Handle(ctx, cmd)
}

func (s *Service) DeleteConsumable(ctx context.Context, uid string) (err error) {
	return s.deleteConsumableCommand.Handle(ctx, uid)
}

// queries
func (s *Service) GetConsumableByUID(ctx context.Context, uid string) (result ports.GetConsumable, err error) {
	return s.getConsumableQuery.Handle(ctx, uid)
}

func (s *Service) ListConsumable(ctx context.Context, filter ports.Filter) (result []ports.GetConsumable, total int64, err error) {
	return s.listConsumableQuery.Handle(ctx, filter)
}

func (s *Service) buildCommands() {
	s.buildSaveConsumable()

	s.buildUpdateConsumable()

	s.buildDeleteConsumable()
}

func (s *Service) buildSaveConsumable() {
	s.saveConsumableCommand = &cmds.SaveConsumableHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildUpdateConsumable() {
	s.updateConsumableCommand = &cmds.UpdateConsumableHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildDeleteConsumable() {
	s.deleteConsumableCommand = &cmds.DeleteConsumableHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildQueries() {

	s.buildGetConsumable()

	s.buildListConsumable()
}

func (s *Service) buildGetConsumable() {
	s.getConsumableQuery = &qries.GetConsumableHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}

func (s *Service) buildListConsumable() {
	s.listConsumableQuery = &qries.ListConsumableHandler{
		Repository: s.d.Repository,
		Tracer:     s.d.Tracer,
	}
}
