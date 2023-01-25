package graph

import "github.com/atjhoendz/notpushcation-service/internal/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver :nodoc:
type Resolver struct {
	MessageProcessorUsecase model.MessageProcessorUsecase
	ThreadUsecase           model.ThreadUsecase
	LiveBlogPostUsecase     model.LiveBlogPostUsecase
}
