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
	decoder := json.NewDecoder(src)
	encoder := json.NewEncoder(dst)

	// Decode
	var agentMsg AgentMessage
	if err := decoder.Decode(&agentMsg); err != nil {
		return fmt.Errorf("failed to decode request json: %w", err)
	}

	// Handle
	srvMsg, err := srv.HandleAgent(ctx, agentMsg)
	if err != nil {
		if encodeErr := encoder.Encode(&ServerMessage{}); encodeErr != nil {
			return fmt.Errorf("failed to handle agent: %w, failed to encode reply: %v", err, encodeErr)
		}
		return fmt.Errorf("failed to handle agent request: %w", err)
	}

	// Encode
	if err := encoder.Encode(srvMsg); err != nil {
		return fmt.Errorf("failed to encode response json: %w", err)
	}

	return nil
}
