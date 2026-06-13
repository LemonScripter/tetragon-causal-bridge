# DCC Causal Bridge for Cilium Tetragon

## Overview
The **DCC Causal Bridge** is a professional extension for Cilium Tetragon that introduces **Digital Causal Closure (DCC)** to the eBPF-based runtime security ecosystem. While traditional security tools focus on *Identity* (Who is calling?), the DCC Bridge enforces *Causality* (What caused this call?).

By closing the **Semantic Gap** between user intent and kernel-level execution, this bridge prevents "orphaned" or "autonomous" API calls—critical for securing AI agents and autonomous workloads.

## The Problem: The Semantic Gap
Current runtime security solutions (including Tetragon's standard TracingPolicy) suffer from a Semantic Gap. If an authorized binary (e.g., a Python-based AI agent) is compromised, it can initiate outbound connections to allowed endpoints. The kernel allows these calls because the identity matches the policy, even if the call was initiated without a valid causal trigger (e.g., without a human command or an authenticated event).

## The Solution: Digital Causal Closure (DCC)
The DCC Causal Bridge implements a **Causal Logic Gate** at the kernel level:
1. **Causal Tokens:** Every sensitive operation (connect, exec, bpf) must be preceded by a valid, non-expired Causal Token.
2. **Hard Enforcement:** If a syscall is intercepted without a verified causal chain in the `global_dcc_map`, the operation is blocked with `-EPERM`.
3. **Atomic Consumption:** Tokens are single-use, preventing replay attacks and ensuring that every action is explicitly authorized.

## Scientific & Theoretical Background
The implementation is based on the following research and formal specifications:

*   **Scientific Paper:** [The Causal Operating System: Digital Causal Closure for Autonomous Systems](https://doi.org/10.5281/zenodo.20384700) (DOI: 10.5281/zenodo.20384700)
*   **Formal Specification:** [BioOS Causal Constitution (PDF)](https://bioos.metaspace.bio/bioos_causal_constitution_en.pdf)

These documents provide the mathematical proof for Digital Causal Closure and the transition from traditional access control to Causal Enforcement.

## Architecture
The bridge consists of two main components:
*   **Kernel Space (`dcc_bridge.bpf.c`)**: An eBPF/LSM module that verifies tokens in the `global_dcc_map` during syscall execution.
*   **User Space (`causal_bridge.go`)**: A Go extension for the Tetragon observer to manage token lifecycle and flag causal violations.

## Verification
A simulation-based test suite is included to verify the bridge's logic:
```bash
python tests/verify_bridge.py
```
**Validated Scenarios:**
*   Blocking orphaned calls (No Token).
*   Allowing verified calls (Fresh Token).
*   Blocking stale/expired intent.
*   Preventing token reuse (Atomic Closure).

## Upstreaming to Tetragon
To integrate this bridge into your Tetragon deployment, refer to the proposed `TracingPolicy` extension in the `causal_shadow_analysis.yaml` and the `matchCausal` selector implementation details in the source directory.

---
*Created by MetaSpace BioOS | Strategic Integration with Cilium Tetragon*
