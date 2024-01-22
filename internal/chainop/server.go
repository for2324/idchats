package chainop

import (
	"Open_IM/pkg/metrics"
	"Open_IM/pkg/server"
	tracer "Open_IM/pkg/trace"
	"Open_IM/pkg/xkafka"
	"Open_IM/pkg/xlog"
	"context"
)

type ServerTemp struct {
	ctx    context.Context
	cancel func()
	kc     *xkafka.Consumer
}

func NewServer(groupname string, borkerArray []string, moduleName string, topicName string, listenerFunc xkafka.Listener, metrics metrics.Provider, tracer tracer.Provider) server.Server {
	svr := &ServerTemp{
		ctx:    context.Background(),
		cancel: nil,
	}
	optptr := xkafka.NewDefaultOptions()
	optptr.Name = moduleName
	optptr.Addr = borkerArray
	optptr.Consumer.Group = groupname
	ct, err := xkafka.New(optptr, metrics, tracer)
	if err != nil {
		xlog.CErrorf(err.Error())
		return nil
	}
	ct.Consumer.AddListener(topicName, listenerFunc)
	svr.kc = ct.Consumer
	return svr

}
func (s *ServerTemp) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.kc.Start()
	return nil
}

func (s *ServerTemp) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.kc != nil {
		s.kc.Stop()
	}
	return nil
}
