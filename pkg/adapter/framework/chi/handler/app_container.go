package handler

type AppContainer struct {
	User    UserHandler
	Session SessionHandler
}

func NewAppContainer(userHandler UserHandler, sessionHandler SessionHandler) *AppContainer {
	return &AppContainer{
		User:    userHandler,
		Session: sessionHandler,
	}
}
