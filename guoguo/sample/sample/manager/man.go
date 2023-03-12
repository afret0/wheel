package manager

import (
	"context"
	"github.com/sirupsen/logrus"
	"sample/source"
)

var m *Manager

type Manager struct {
	logger *logrus.Logger
}

func (m *Manager) Sample(ctx context.Context) {
	m.logger.Infoln("succeed")
}

func GetManager() *Manager {
	if m != nil {
		return m
	}

	m = new(Manager)
	m.logger = source.GetLogger()

	return m
}
