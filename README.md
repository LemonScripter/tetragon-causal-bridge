# DCC Causal Bridge for Cilium Tetragon

[![Verified](https://img.shields.io/badge/Verified-Tokyo--Node-green)](VERIFICATION.md)
[![Status](https://img.shields.io/badge/Status-Hardened--Prototype-blue)](ROADMAP.md)
[![Project](https://img.shields.io/badge/BioOS-Causal--Security-green)](https://metaspace.bio)
[![DOI](https://img.shields.io/badge/DOI-10.5281%2Fzenodo.20384700-purple)](https://doi.org/10.5281/zenodo.20384700)

## Hardened Architecture: Kernel-Anchored Observability

The **DCC Causal Bridge** transforms Tetragon from a reactive observer into a **Causal Enforcement Engine**. It eliminates the semantic gap by ensuring that every Tetragon event is causally linked to a hardware-verified intent.

### Production-Grade Implementation

- **Direct eBPF Map Interaction:** The bridge utilizes `ebpf.LoadPinnedMap` to query the `global_dcc_map` directly from the BPF filesystem (`/sys/fs/bpf/tetragon/`).
- **Fail-Closed Logic:** If the causal map is unavailable or a PID has no associated token, the bridge triggers a `DCC Violation` alert, facilitating immediate blocking.
- **Hardware-Anchored Truth:** Causal tokens are populated via kernel LSM hooks (`dcc_bridge.bpf.c`) triggered by physical IRQ events.

### Security Guarantees

1. **Temporal Integrity:** Enforces a sub-second causality window (default 500ms).
2. **Atomic Consumption:** Prevents token reuse through atomic kernel-space flags.
3. **Intent Pinning:** Verifies that the specific syscall (e.g., `connect`) matches the authorized intent ID.

### Scientific & Technical Foundation

This implementation is based on the following formal specifications and research:

- **Research Paper:** [The Causal Operating System: Digital Causal Closure for Autonomous Systems](https://doi.org/10.5281/zenodo.20384700) (DOI: 10.5281/zenodo.20384700)
- **Formal Specification:** [BioOS Causal Constitution (PDF)](https://bioos.metaspace.bio/bioos_causal_constitution_en.pdf)

---
*MetaSpace.Bio Logic Project | [metaspace.bio](https://metaspace.bio) | [admin@metaspace.bio](mailto:admin@metaspace.bio)*
