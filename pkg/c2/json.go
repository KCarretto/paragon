package c2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// HandleJSON provides wraps HandleAgent by decoding a JSON AgentMessage from src and encoding a
// JSON ServerMessage to dst.
func (srv Server) HandleJSON(ctx context.Context, src io.Reader, dst io.Writer) error {
	// Decode
	decoder := json.NewDecoder(src)
	var agentMsg AgentMessage
	if err := decoder.Decode(&agentMsg); err != nil {
		return fmt.Errorf("failed to decode request json: %w", err)
	}

	// Handle
	srvMsg, err := srv.HandleAgent(ctx, agentMsg)
	if err != nil {
		return fmt.Errorf("failed to handle agent request: %w", err)
	}

	// Encode
	encoder := json.NewEncoder(dst)
	if err := encoder.Encode(srvMsg); err != nil {
		return fmt.Errorf("failed to encode response json: %w", err)
	}

	return nil
}
