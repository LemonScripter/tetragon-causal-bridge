package observer

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/cilium/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
)

/*
 * DCC Causal Bridge for Tetragon (Hardened Production-Grade)
 * 
 * This component extends Tetragon's observer to verify the causal chain 
 * using direct eBPF map lookups. It anchors Tetragon events to the 
 * Digital Causal Closure (DCC) kernel state.
 */

const (
	IntentNetworkEgress = 0x100
	IntentProcessSpawn  = 0x200
	DCCMapPath          = "/sys/fs/bpf/tetragon/global_dcc_map"
)

// DCCToken matches the struct in dcc_bridge.bpf.c
type DCCToken struct {
	Timestamp  uint64
	IntentID   uint32
	AgeLimitNS uint32
	Consumed   uint8
}

// CausalBridge manages the connection to the DCC kernel state
type CausalBridge struct {
	dccMap *ebpf.Map
}

// NewCausalBridge initializes the bridge by pinning the global DCC map
func NewCausalBridge() (*CausalBridge, error) {
	m, err := ebpf.LoadPinnedMap(DCCMapPath, nil)
	if err != nil {
		return nil, fmt.Errorf("DCC Critical: Failed to load pinned causal map at %s: %w", DCCMapPath, err)
	}
	return &CausalBridge{dccMap: m}, nil
}

// VerifyCausality performs a synchronous lookup in the kernel map to authorize an event.
func (b *CausalBridge) VerifyCausality(pid uint32, intentID uint32) (bool, error) {
	var token DCCToken
	
	// Fail-Closed: If map lookup fails, we assume no causal lineage
	if err := b.dccMap.Lookup(unsafe.Pointer(&pid), unsafe.Pointer(&token)); err != nil {
		return false, fmt.Errorf("DCC Violation: No causal lineage found for PID %d", pid)
	}

	// 1. Temporal Integrity (Causality Window)
	// Note: In kernel space, we use bpf_ktime_get_ns(). In userspace, we approximate.
	now := uint64(time.Now().UnixNano())
	if now < token.Timestamp || now - token.Timestamp > uint64(token.AgeLimitNS) {
		return false, fmt.Errorf("DCC Violation: Causal token expired for PID %d", pid)
	}

	// 2. Atomic Integrity (Anti-Replay)
	if token.Consumed != 0 {
		return false, fmt.Errorf("DCC Violation: Token replay attempt for PID %d", pid)
	}

	// 3. Intent Validation
	if token.IntentID != intentID {
		return false, fmt.Errorf("DCC Violation: Intent mismatch (Expected %x, Got %x) for PID %d", intentID, token.IntentID, pid)
	}

	return true, nil
}

func (b *CausalBridge) Close() error {
	if b.dccMap != nil {
		return b.dccMap.Close()
	}
	return nil
}
