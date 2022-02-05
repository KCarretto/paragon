package agent

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AgentVersion represents a unique agent version string to identify the agent.
const AgentVersion string = "paragon-v0.1.1"

// collectMetadata about the system and update the agent with the gathered information.
func (agent *Agent) collectMetadata() {
	if agent == nil {
		return
	}
	if agent.Log == nil {
		agent.Log = zap.NewNop()
	}

	// Generate session ID once
	if agent.Metadata.SessionID == "" {
		sessionUUID, err := uuid.NewRandom()
		if err != nil {
			agent.Log.Error("Failed to generate session ID", zap.Error(err))
			agent.Metadata.SessionID = fmt.Sprintf("session_id_err_%d", time.Now().Unix())
		} else {
			agent.Metadata.SessionID = sessionUUID.String()
		}
	}

	// Generate machine UUID once
	if agent.Metadata.MachineUUID == "" {
		machineUUID, err := machineid.ID()
		if err != nil {
			agent.Log.Warn("Failed to collect machine UUID", zap.Error(err))
		} else {
			agent.Metadata.MachineUUID = machineUUID
		}
		if agent.machineidPrefix != "" {
			agent.Metadata.MachineUUID = fmt.Sprintf("%s%s", agent.machineidPrefix, agent.Metadata.MachineUUID)
		}
	}

	// Always update hostname
	hostname, err := os.Hostname()
	if err != nil {
		agent.Log.Error("Failed to collect machine hostname", zap.Error(err))
		hostname = ""
	}
	if hostname != "" {
		agent.Metadata.Hostname = hostname
	}

	// Get primary interface
	iface := routedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if iface == nil {
		return
	}

	// Always update primary MAC address
	if mac := iface.HardwareAddr.String(); mac != "" {
		agent.Metadata.PrimaryMAC = mac
	}

	// Always update primary IP
	ipAddrs, err := iface.Addrs()
	if err != nil || len(ipAddrs) < 1 {
		agent.Log.Error("Failed to collect machine primary IP", zap.Error(err))
		return
	}

	for _, ipAddr := range ipAddrs {
		ip, _, err := net.ParseCIDR(ipAddr.String())
		if err != nil {
			agent.Log.Error("Failed to collect machine primary IP", zap.Error(err))
			continue
		}
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.To4() == nil {
			continue
		}
		agent.Metadata.PrimaryIP = ip.String()
		break
	}
	// agent.transport.TaskExecutor = agent.Metadata
}

// isMulticastCapable reports whether ifi is an IP multicast-capable
// network interface. Network must be "ip", "ip4" or "ip6".
func isMulticastCapable(network string, ifi *net.Interface) (net.IP, bool) {
	switch network {
	case "ip", "ip4", "ip6":
	default:
		return nil, false
	}
	if ifi == nil || ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagMulticast == 0 {
		return nil, false
	}
	return hasRoutableIP(network, ifi)
}

// routedInterface returns a network interface that can route IP
// traffic and satisfies flags. It returns nil when an appropriate
// network interface is not found. Network must be "ip", "ip4" or
// "ip6".
func routedInterface(network string, flags net.Flags) *net.Interface {
	switch network {
	case "ip", "ip4", "ip6":
	default:
		return nil
	}
	ift, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ift {
		if ifi.Flags&flags != flags {
			continue
		}
		if _, ok := hasRoutableIP(network, &ifi); !ok {
			continue
		}
		return &ifi
	}
	return nil
}

func hasRoutableIP(network string, ifi *net.Interface) (net.IP, bool) {
	ifat, err := ifi.Addrs()
	if err != nil {
		return nil, false
	}
	for _, ifa := range ifat {
		switch ifa := ifa.(type) {
		case *net.IPAddr:
			if ip := routableIP(network, ifa.IP); ip != nil {
				return ip, true
			}
		case *net.IPNet:
			if ip := routableIP(network, ifa.IP); ip != nil {
				return ip, true
			}
		}
	}
	return nil, false
}

func routableIP(network string, ip net.IP) net.IP {
	if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsGlobalUnicast() {
		return nil
	}
	switch network {
	case "ip4":
		if ip := ip.To4(); ip != nil {
			return ip
		}
	case "ip6":
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() { // addressing scope of the loopback address depends on each implementation
			return nil
		}
		if ip := ip.To16(); ip != nil && ip.To4() == nil {
			return ip
		}
	default:
		if ip := ip.To4(); ip != nil {
			return ip
		}
		if ip := ip.To16(); ip != nil {
			return ip
		}
	}
	return nil
}
