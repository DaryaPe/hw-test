package mocks

type Logger struct{}

func (m Logger) Debug(_ ...interface{}) {}

func (m Logger) Info(_ ...interface{}) {}

func (m Logger) Warn(_ ...interface{}) {}

func (m Logger) Error(_ ...interface{}) {}

func (m Logger) Debugf(_ string, _ ...interface{}) {}

func (m Logger) Infof(_ string, _ ...interface{}) {}

func (m Logger) Warnf(_ string, _ ...interface{}) {}

func (m Logger) Errorf(_ string, _ ...interface{}) {}

func (m Logger) Debugw(_ string, _ ...interface{}) {}

func (m Logger) Infow(_ string, _ ...interface{}) {}

func (m Logger) Warnw(_ string, _ ...interface{}) {}

func (m Logger) Errorw(_ string, _ ...interface{}) {}

func (m Logger) Print(_ ...interface{}) {}

func (m Logger) Printf(_ string, _ ...interface{}) {}
