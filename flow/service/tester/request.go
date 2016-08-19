package tester

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("tester")

// RequestProcessor processes request objects and invokes the corresponding
// flow Manager methods
type RequestProcessor struct {
	runner action.Runner
}

// NewRequestProcessor creates a new RequestProcessor
func NewRequestProcessor() *RequestProcessor {

	var rp RequestProcessor
	rp.runner = runner.NewDirectRunner()

	return &rp
}

// StartFlow handles a StartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a StartRequest.
func (rp *RequestProcessor) StartFlow(startRequest *StartRequest) (code int, retData interface{}, err error) {

	execOptions := &flowinst.ExecOptions{Interceptor: startRequest.Interceptor, Patch: startRequest.Patch}

	attrs := startRequest.Attrs

	dataLen := len(startRequest.Data)

	// attrs, not supplied so attempt to create attrs from Data
	if attrs == nil && dataLen > 0 {
		attrs = make([]*data.Attribute, 0, dataLen)

		for k, v := range startRequest.Data {

			//todo handle error
			t, _ := data.GetType(v)
			attrs = append(attrs, data.NewAttribute(k, t, v))
		}
	}

	ao := flowinst.RunOptions{Op: flowinst.AoStart, ExecOptions: execOptions}

	action := action.Get(flowinst.ActionType)
	ctx := trigger.NewContext(context.Background(), attrs)

	return rp.runner.Run(ctx, action, startRequest.FlowURI, ao)
}

// RestartFlow handles a RestartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a RestartRequest.
func (rp *RequestProcessor) RestartFlow(restartRequest *RestartRequest) (code int, retData interface{}, err error) {

	execOptions := &flowinst.ExecOptions{Interceptor: restartRequest.Interceptor, Patch: restartRequest.Patch}
	ao := flowinst.RunOptions{Op: flowinst.AoRestart, InitialState: restartRequest.InitialState, ExecOptions: execOptions}

	ctx := context.Background()

	if restartRequest.Data != nil {

		log.Debugf("Updating flow attrs: %v", restartRequest.Data)
		attrs := make([]*data.Attribute, len(restartRequest.Data))

		for k, v := range restartRequest.Data {
			attrs = append(attrs, data.NewAttribute(k, data.ANY, v))
		}

		ctx = trigger.NewContext(context.Background(), attrs)
	}

	action := action.Get(flowinst.ActionType)
	return rp.runner.Run(ctx, action, restartRequest.InitialState.FlowURI, ao)
}

// ResumeFlow handles a ResumeRequest for a FlowInstance.  This will
// queue a RestartRequest.
func (rp *RequestProcessor) ResumeFlow(resumeRequest *ResumeRequest) (code int, retData interface{}, err error) {

	execOptions := &flowinst.ExecOptions{Interceptor: resumeRequest.Interceptor, Patch: resumeRequest.Patch}
	ao := flowinst.RunOptions{Op: flowinst.AoRestart, InitialState: resumeRequest.State, ExecOptions: execOptions}

	ctx := context.Background()

	if resumeRequest.Data != nil {

		log.Debugf("Updating flow attrs: %v", resumeRequest.Data)
		attrs := make([]*data.Attribute, len(resumeRequest.Data))

		for k, v := range resumeRequest.Data {
			attrs = append(attrs, data.NewAttribute(k, data.ANY, v))
		}

		ctx = trigger.NewContext(context.Background(), attrs)
	}

	action := action.Get(flowinst.ActionType)
	return rp.runner.Run(ctx, action, resumeRequest.State.FlowURI, ao)
}

// StartRequest describes a request for starting a FlowInstance
type StartRequest struct {
	FlowURI     string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Attrs       []*data.Attribute      `json:"attrs"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}

// RestartRequest describes a request for restarting a FlowInstance
// todo: can be merged into StartRequest
type RestartRequest struct {
	InitialState *flowinst.Instance     `json:"initialState"`
	Data         map[string]interface{} `json:"data"`
	Interceptor  *support.Interceptor   `json:"interceptor"`
	Patch        *support.Patch         `json:"patch"`
}

// ResumeRequest describes a request for resuming a FlowInstance
//todo: Data for resume request should be directed to wating task
type ResumeRequest struct {
	State       *flowinst.Instance     `json:"state"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
}
