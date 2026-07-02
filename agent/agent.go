package agent

import (
	"errors"
	"sync"
	"time"

	"agente-bancario/audit"
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
	audit.Log(audit.Event{
		UserID:   user.ID,
		UserRole: string(user.Role),
		Action:   "tool_call_received",
		Tool:     toolcall.Name,
		Status:   "received",
	})

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

		audit.Log(audit.Event{
			UserID:   user.ID,
			UserRole: string(user.Role),
			Action:   "pending_action_created",
			Tool:     toolcall.Name,
			Status:   "pending_confirmation",
			Reason:   "pix_requires_confirmation",
		})

		return "pix_requires_confirmation", nil
	}

	responseTool, err := tools.ExecuteTool(user, toolcall)
	if err != nil {
		audit.Log(audit.Event{
			UserID:   user.ID,
			UserRole: string(user.Role),
			Action:   "tool_call_finished",
			Tool:     toolcall.Name,
			Status:   "failed",
			Reason:   err.Error(),
		})
		return nil, err
	}

	audit.Log(audit.Event{
		UserID:   user.ID,
		UserRole: string(user.Role),
		Action:   "tool_call_finished",
		Tool:     toolcall.Name,
		Status:   "success",
	})

	return responseTool, nil
}

func (o *Orchestrator) ConfirmPendingAction(user policy.AuthenticatedUser) (any, error) {
	o.mu.Lock()
	pendingAction, ok := o.pendingActions[user.ID]
	if !ok {
		o.mu.Unlock()
		audit.Log(audit.Event{
			UserID:   user.ID,
			UserRole: string(user.Role),
			Action:   "pending_action_confirmation",
			Status:   "failed",
			Reason:   "no_pending_action",
		})
		return nil, errors.New("no_pending_action")
	}
	delete(o.pendingActions, user.ID)
	o.mu.Unlock()

	audit.Log(audit.Event{
		UserID:   user.ID,
		UserRole: string(user.Role),
		Action:   "pending_action_confirmation",
		Tool:     pendingAction.ToolCall.Name,
		Status:   "confirmed",
	})

	confirmedToolCall := addConfirmation(pendingAction.ToolCall)

	responseTool, err := tools.ExecuteTool(user, confirmedToolCall)
	if err != nil {
		audit.Log(audit.Event{
			UserID:   user.ID,
			UserRole: string(user.Role),
			Action:   "pending_action_execution",
			Tool:     pendingAction.ToolCall.Name,
			Status:   "failed",
			Reason:   err.Error(),
		})
		return nil, err
	}

	audit.Log(audit.Event{
		UserID:   user.ID,
		UserRole: string(user.Role),
		Action:   "pending_action_execution",
		Tool:     pendingAction.ToolCall.Name,
		Status:   "success",
	})

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
