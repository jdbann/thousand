package server

import (
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/middleware"
)

func (s *Server) setupRoutes() {
	s.setup.Do(func() {
		middleware.RequestLogger(s.mux, s.logger.Named("server"))
		middleware.MethodOverride(s.mux)
		middleware.RedirectSlashes(s.mux)

		handlers.Assets(s.mux, s.assets)

		handlers.Root(s.mux)

		handlers.NewCharacter(s.mux, s.logger, s.renderer, s.repository)
		handlers.CreateCharacter(s.mux, s.logger, s.repository)

		handlers.NewExperience(s.mux, s.logger, s.renderer, s.repository)
		handlers.CreateExperience(s.mux, s.logger, s.repository)

		handlers.NewMark(s.mux, s.logger, s.renderer, s.repository)
		handlers.CreateMark(s.mux, s.logger, s.repository)

		handlers.NewResource(s.mux, s.logger, s.renderer, s.repository)
		handlers.CreateResource(s.mux, s.logger, s.repository)

		handlers.NewSkill(s.mux, s.logger, s.renderer, s.repository)
		handlers.CreateSkill(s.mux, s.logger, s.repository)

		handlers.ListVampires(s.mux, s.logger, s.renderer, s.repository)
		handlers.NewVampire(s.mux, s.logger, s.renderer)
		handlers.CreateVampire(s.mux, s.logger, s.repository)
		handlers.ShowVampire(s.mux, s.logger, s.renderer, s.repository)
	})
}
