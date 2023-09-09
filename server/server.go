package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"vitalsign-publisher/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ServerStart(vsp *VSP, port int, ch chan bool) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protos.RegisterVitalSignPublishServer(s, vsp)

	log.Printf("server listening at %v", lis.Addr())
	ch <- true
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type VSP struct {
	protos.UnimplementedVitalSignPublishServer
	Mutex    sync.Mutex
	RPNs     []*protos.RPN
	Patients []*protos.Patient
}

func (s *VSP) CheckRPNs(ctx context.Context, in *protos.VoidRequest) (*protos.RPNs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRPNs not implemented")
}
func (s *VSP) CheckPatients(ctx context.Context, in *protos.VoidRequest) (*protos.Patients, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPatients not implemented")
}
func (s *VSP) RegisterRPN(ctx context.Context, in *protos.RPN) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterRPN not implemented")
	// s.Mutex.Lock()
	// s.RPNs = append(s.RPNs, in)
	// s.Mutex.Unlock()
	// fmt.Println(in, s.RPNs)
	// return &protos.Msg{Status: true, Topic: "rpn_topoc", Msg: "rpn_msg"}, nil
}
func (s *VSP) UnregisterRPN(ctx context.Context, in *protos.RPN) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterRPN not implemented")
}
func (s *VSP) RegisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPatient not implemented")
}
func (s *VSP) UnregisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterPatient not implemented")
}
