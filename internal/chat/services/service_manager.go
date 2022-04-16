package services

import (
	"sync"
)

type ServiceManagerService struct {
	conversationService ConversationService
}

var serviceManagerOnce sync.Once
var serviceManagerInstance ServiceManagerService

func NewServiceManager(conversationService ConversationService) ServiceManagerService {
	serviceManagerOnce.Do(func() { // <-- atomic, does not allow repeating
		serviceManagerInstance = ServiceManagerService{
			conversationService: conversationService,
		}
	})

	return serviceManagerInstance
}
