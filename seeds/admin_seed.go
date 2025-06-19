package seeds

import (
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/services"
)

type TodoSeed struct {
	logger         framework.Logger
	cognitoService services.CognitoAuthService
	todoService    *todo.Service
	env            *framework.Env
}

// NewTodoSeed creates admin seed
func NewTodoSeed(
	logger framework.Logger,
	cognitoService services.CognitoAuthService,
	todoService *todo.Service,
	env *framework.Env,
) TodoSeed {
	return TodoSeed{
		logger:         logger,
		cognitoService: cognitoService,
		todoService:    todoService,
		env:            env,
	}
}

// Run the todo seed
func (s TodoSeed) Setup() {
	s.logger.Info("🌱 seeding todo data...")

	s.logger.Info("Todo user already exists")
}
