# Quickstart: DCC Causal Enforcement for Tetragon

This guide provides a self-contained environment to verify the **Digital Causal Closure (DCC)** integration for Cilium Tetragon.

## Prerequisites
- **Go** (1.21+)
- **Linux** (for eBPF integration) or **macOS** (for logic demo)

## Step 1: Build the Bridge
```bash
git clone https://github.com/LemonScripter/tetragon-causal-bridge.git
cd tetragon-causal-bridge
make build
```

## Step 2: Run the Automated Proof
Execute the reproduction script to see the "Fail-Closed" logic in action:
```bash
./reproduce.sh
```

## Step 3: Production Deployment
To enable kernel-anchored enforcement, the bridge must be deployed on a BioOS-enabled Linux node. It will automatically attempt to link with the Tetragon eBPF maps at:
`/sys/fs/bpf/tetragon/global_dcc_map`

## Verification Scenarios
- **Verified:** Processes with a valid Causal Token issued by the BioOS kernel are permitted to execute sensitive syscalls.
- **Orphaned:** Processes initiated without a verified hardware intent are blocked at the LSM hook level.

---
*Production-Grade Research Prototype by MetaSpace BioOS Team*
