package server

import (
	"context"
	"fmt"
	"net"
	"slices"
	"sync"
	"vitalsign-publisher/common"
	"vitalsign-publisher/config"
	"vitalsign-publisher/protos"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var conf = config.GetConfig()

func ServerStart(vsp *VSP, port int, ch chan bool) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		color.Red("ServerStart: FAIL - listen port failed: %v", err)
	}

	s := grpc.NewServer()
	protos.RegisterVitalSignPublishServer(s, vsp)
	color.Green("ServerStart: SUCCESS - server is serving at %v", lis.Addr())
	ch <- true

	if err := s.Serve(lis); err != nil {
		color.Red("ServerStart: FAIL - serve failed: %v", err)
	}
}

type VSP struct {
	protos.UnimplementedVitalSignPublishServer
	MuRpn     sync.Mutex
	MuPatient sync.Mutex
	RPNs      []*protos.RPN
	Patients  []*protos.Patient
}

func (vsp *VSP) CheckRPNs(ctx context.Context, in *protos.VoidRequest) (*protos.RPNs, error) {
	color.Cyan("%v CheckRPNs: %+v", common.TimeNow(), vsp.RPNs)
	return nil, status.Errorf(codes.Unimplemented, "method CheckRPNs not implemented")
}

func (vsp *VSP) CheckPatients(ctx context.Context, in *protos.VoidRequest) (*protos.Patients, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPatients not implemented")
}

func (vsp *VSP) RegisterRPN(ctx context.Context, in *protos.RPN) (*protos.Msg, error) {
	// Generate MQTT topic
	topic := fmt.Sprintf("%s%s", conf.Setting.RpnBase, in.Id)

	// Check if RPN is serving or not
	vsp.MuRpn.Lock()
	if slices.ContainsFunc(vsp.RPNs, func(r *protos.RPN) bool {
		if r.Id == in.Id && r.Hospital == in.Hospital {
			return true
		}
		return false
	}) {
		vsp.MuRpn.Unlock()
		color.Yellow("%v RegisterRPN: RPN {%+v} has been serving", common.TimeNow(), in)
		return &protos.Msg{Status: false, Topic: topic, Msg: "Fail: already serving"}, nil
	} else {
		vsp.RPNs = append(vsp.RPNs, in)
		vsp.MuRpn.Unlock()
		color.Cyan("%v RegisterRPN: {%+v}", common.TimeNow(), in)
	}
	return &protos.Msg{Status: true, Topic: topic, Msg: "Success: RPN start serving"}, nil
}

func (vsp *VSP) UnregisterRPN(ctx context.Context, in *protos.RPN) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterRPN not implemented")
}
func (vsp *VSP) RegisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPatient not implemented")
}
func (vsp *VSP) UnregisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterPatient not implemented")
}
