package player

import "go.uber.org/zap"

type PlayerService struct {
	log *zap.Logger
}

func NewPlayerService(log *zap.Logger) *PlayerService {
	return &PlayerService{
		log: log,
	}
}

func (ps *PlayerService) Build() {

}
