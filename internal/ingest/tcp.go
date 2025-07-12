package ingest

import (
	"bufio"
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type TCPServer struct {
	addr     string
	sink     Sink
	listener net.Listener
	wg       sync.WaitGroup
	active   atomic.Int32
	lines    atomic.Uint64
}

func NewTCPServer(addr string, sink Sink) *TCPServer {
	return &TCPServer{addr: addr, sink: sink}
}

func (s *TCPServer) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		_ = ln.Close()
	}()
	s.wg.Add(1)
	go s.acceptLoop(ctx)
	return nil
}

func (s *TCPServer) acceptLoop(ctx context.Context) {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if ctx.Err() != nil {
				return
			}
			continue
		}
		s.active.Add(1)
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			defer s.active.Add(-1)
			defer c.Close()
			s.handleConn(c)
		}(conn)
	}
}

func (s *TCPServer) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		entry := logentry.Entry{
			Service: "tcp",
			Message: line,
			Raw:     line,
		}
		if err := s.sink(entry); err != nil {
			return
		}
		s.lines.Add(1)
	}
}

func (s *TCPServer) Wait() {
	s.wg.Wait()
}

func (s *TCPServer) ActiveClients() int32 {
	return s.active.Load()
}

func (s *TCPServer) LinesIngested() uint64 {
	return s.lines.Load()
}
