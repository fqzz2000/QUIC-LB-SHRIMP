package lb

import (
	"net"
	"sync"

	"github.com/fqzz2000/QUIC-LB-SHRIMP/pkg/packet"
)

// LoadBalancer represents the main QUIC load balancer structure
type LoadBalancer struct {
	// Configuration
	listenAddr string
	backends   []string

	// Runtime state
	listener net.PacketConn
	mu       sync.RWMutex
	running  bool

	// Packet processing
	packetProcessor *packet.HeaderParser
}

// InitLoadBalancer creates and initializes a new LoadBalancer instance
func InitLoadBalancer(listenAddr string, backends []string) (*LoadBalancer, error) {
	lb := &LoadBalancer{
		listenAddr: listenAddr,
		backends:   backends,
		running:    false,
	}
	return lb, nil
}

// Start begins the load balancer operations
func (lb *LoadBalancer) Start() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.running {
		return nil
	}

	listener, err := net.ListenPacket("udp", lb.listenAddr)
	if err != nil {
		return err
	}

	lb.listener = listener
	lb.running = true

	return nil
}

// ReadPacket reads a single QUIC packet from the UDP listener
func (lb *LoadBalancer) ReadPacket() ([]byte, net.Addr, error) {
	// Buffer size for QUIC packets (typical MTU size)
	buffer := make([]byte, 1500)

	n, addr, err := lb.listener.ReadFrom(buffer)
	if err != nil {
		return nil, nil, err
	}

	// Return only the bytes that were actually read
	return buffer[:n], addr, nil
}

// ExtractCID extracts the Connection ID from a QUIC packet
// Returns the CID as a byte slice and an error if extraction fails
func (lb *LoadBalancer) ExtractCID(packet []byte) ([]byte, error) {
	return lb.packetProcessor.HeaderParser.ExtractCID(packet)
}

// Shutdown gracefully stops the load balancer
func (lb *LoadBalancer) Shutdown() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if !lb.running {
		return nil
	}

	if err := lb.listener.Close(); err != nil {
		return err
	}

	lb.running = false
	return nil
}
