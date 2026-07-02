package agent

import (
	"errors"
	"sync"
	"time"

	"agente-bancario/mcp"
	"agente-bancario/policy"
	"agente-bancario/tools"
)

type PendingAction struct {
	ID        string
	UserID    string
	ToolCall  mcp.ToolCall
	CreatedAt time.Time
}

type Orchestrator struct {
	mu             sync.Mutex
	pendingActions map[string]PendingAction
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		pendingActions: map[string]PendingAction{},
	}
}

func HandleToolCall(user policy.AuthenticatedUser, toolcall mcp.ToolCall) (any, error) {
	return NewOrchestrator().HandleToolCall(user, toolcall)
}

func (o *Orchestrator) HandleToolCall(user policy.AuthenticatedUser, toolcall mcp.ToolCall) (any, error) {
	if toolcall.Name == "create_pix" {
		pendingToolCall := copyToolCallWithoutConfirmation(toolcall)
		pendingAction := PendingAction{
			ID:        "pending-" + user.ID,
			UserID:    user.ID,
			ToolCall:  pendingToolCall,
			CreatedAt: time.Now(),
		}

		o.mu.Lock()
		o.pendingActions[user.ID] = pendingAction
		o.mu.Unlock()

		return "pix_requires_confirmation", nil
	}

	responseTool, err := tools.ExecuteTool(user, toolcall)
	if err != nil {
		return nil, err
	}

	return responseTool, nil
}

func (o *Orchestrator) ConfirmPendingAction(user policy.AuthenticatedUser) (any, error) {
	o.mu.Lock()
	pendingAction, ok := o.pendingActions[user.ID]
	if !ok {
		o.mu.Unlock()
		return nil, errors.New("no_pending_action")
	}
	delete(o.pendingActions, user.ID)
	o.mu.Unlock()

	confirmedToolCall := addConfirmation(pendingAction.ToolCall)

	responseTool, err := tools.ExecuteTool(user, confirmedToolCall)
	if err != nil {
		return nil, err
	}

	return responseTool, nil
}

func copyToolCallWithoutConfirmation(toolcall mcp.ToolCall) mcp.ToolCall {
	arguments := map[string]string{}
	for key, value := range toolcall.Arguments {
		if key == "confirmed" {
			continue
		}
		arguments[key] = value
	}

	return mcp.ToolCall{
		Name:      toolcall.Name,
		Arguments: arguments,
	}
}

func addConfirmation(toolcall mcp.ToolCall) mcp.ToolCall {
	arguments := map[string]string{}
	for key, value := range toolcall.Arguments {
		arguments[key] = value
	}
	arguments["confirmed"] = "true"

	return mcp.ToolCall{
		Name:      toolcall.Name,
		Arguments: arguments,
	}
}
