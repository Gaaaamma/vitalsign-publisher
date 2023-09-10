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
	vsp.MuRpn.Lock()
	defer vsp.MuRpn.Unlock()
	return &protos.RPNs{Rpns: vsp.RPNs}, nil
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
		msg := fmt.Sprintf("Fail: %s is serving already", in.Id)
		return &protos.Msg{Status: false, Topic: topic, Msg: msg}, nil
	} else {
		vsp.RPNs = append(vsp.RPNs, in)
		vsp.MuRpn.Unlock()
		color.Cyan("%v RegisterRPN: {%+v}", common.TimeNow(), in)
		msg := fmt.Sprintf("Success: start serving RPN %s", in.Id)
		return &protos.Msg{Status: true, Topic: topic, Msg: msg}, nil
	}
}

func (vsp *VSP) UnregisterRPN(ctx context.Context, in *protos.RPN) (*protos.Msg, error) {
	// Check if RPN is serving or not
	vsp.MuRpn.Lock()
	if index := slices.IndexFunc(vsp.RPNs, func(r *protos.RPN) bool {
		if r.Id == in.Id && r.Hospital == in.Hospital {
			return true
		}
		return false
	}); index != -1 {
		// Find index of rpn => Remove it from list
		vsp.RPNs = slices.Delete(vsp.RPNs, index, index+1)
		vsp.MuRpn.Unlock()
		color.Cyan("%v UnregisterRPN: {%+v}", common.TimeNow(), in)
		msg := fmt.Sprintf("Success: unregister %s from serving list", in.Id)
		return &protos.Msg{Status: true, Topic: "-", Msg: msg}, nil
	} else {
		// Not found
		vsp.MuRpn.Unlock()
		color.Yellow("%v UnregisterRPN: RPN {%+v} doesn't exist in list", common.TimeNow(), in)
		msg := fmt.Sprintf("Fail: %s doesn't exist in list", in.Id)
		return &protos.Msg{Status: false, Topic: "-", Msg: msg}, nil
	}
}

func (vsp *VSP) RegisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	// Generate MQTT topic
	topic := fmt.Sprintf("%s%s", conf.Setting.PatientBase, in.Id)

	// Check if patient is serving or not
	vsp.MuPatient.Lock()
	if slices.ContainsFunc(vsp.Patients, func(p *protos.Patient) bool {
		if p.Id == in.Id && p.Hospital == in.Hospital {
			return true
		}
		return false
	}) {
		vsp.MuPatient.Unlock()
		color.Yellow("%v RegisterPatient: Patient {%+v} has been serving", common.TimeNow(), in)
		msg := fmt.Sprintf("Fail: %s is serving already", in.Id)
		return &protos.Msg{Status: false, Topic: topic, Msg: msg}, nil
	} else {
		vsp.Patients = append(vsp.Patients, in)
		vsp.MuPatient.Unlock()
		color.Cyan("%v RegisterPatient: {%+v}", common.TimeNow(), in)
		msg := fmt.Sprintf("Success: start serving Patient %s", in.Id)
		return &protos.Msg{Status: true, Topic: topic, Msg: msg}, nil
	}
}

func (vsp *VSP) UnregisterPatient(ctx context.Context, in *protos.Patient) (*protos.Msg, error) {
	// Check if Patient is serving or not
	vsp.MuPatient.Lock()
	if index := slices.IndexFunc(vsp.Patients, func(p *protos.Patient) bool {
		if p.Id == in.Id && p.Hospital == in.Hospital {
			return true
		}
		return false
	}); index != -1 {
		// Find index of patient => Remove it from list
		vsp.Patients = slices.Delete(vsp.Patients, index, index+1)
		vsp.MuPatient.Unlock()
		color.Cyan("%v UnregisterPatient: {%+v}", common.TimeNow(), in)
		msg := fmt.Sprintf("Success: unregister %s from serving list", in.Id)
		return &protos.Msg{Status: true, Topic: "-", Msg: msg}, nil
	} else {
		// Not found
		vsp.MuPatient.Unlock()
		color.Yellow("%v UnregisterPatient: Patient {%+v} doesn't exist in list", common.TimeNow(), in)
		msg := fmt.Sprintf("Fail: %s doesn't exist in list", in.Id)
		return &protos.Msg{Status: false, Topic: "-", Msg: msg}, nil
	}
}
