/*
 * DCC Causal Bridge for Tetragon (Kernel Space)
 * 
 * This eBPF program implements the Digital Causal Closure (DCC) verification 
 * for critical syscalls intercepted by Tetragon.
 * 
 * It ensures that every sensitive operation (connect, exec, etc.) is backed 
 * by a valid, non-expired Causal Token in the global_dcc_map.
 */

#include <linux/bpf.h>
#include <linux/ptrace.h>
#include <linux/sched.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

char LICENSE[] SEC("license") = "GPL";

/* DCC Token Structure */
struct dcc_token {
    __u64 timestamp;
    __u32 intent_id;
    __u32 age_limit_ns;
    __u8  consumed;
};

/* Global Map for Causal Tokens - shared with userland and Tetragon */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, __u32);   /* Key: PID/TGID */
    __type(value, struct dcc_token);
} global_dcc_map SEC(".maps");

/* 
 * Helper function to verify the causal chain.
 * Returns 1 if the operation is causally closed, 0 if it's an "orphaned" call.
 */
static __always_inline int verify_causality(__u32 pid, __u32 intent_id) {
    struct dcc_token *token;
    __u64 now = bpf_ktime_get_ns();

    token = bpf_map_lookup_elem(&global_dcc_map, &pid);
    if (!token) {
        /* No token found: This is an autonomous/unauthorized API call */
        return 0;
    }

    /* Check for expiration (Causality Window) */
    if (now - token->timestamp > token->age_limit_ns) {
        return 0;
    }

    /* Prevent token reuse if already consumed */
    if (token->consumed) {
        return 0;
    }

    /* Verify if the intent matches the requested operation */
    if (token->intent_id != intent_id) {
        return 0;
    }

    /* Token is valid. Mark as consumed for this atomic operation. */
    token->consumed = 1;
    return 1;
}

/* 
 * Example hook: sys_connect logic closure.
 * This should be called by Tetragon's observer during its connect probe.
 */
SEC("lsm/socket_connect")
int BPF_PROG(dcc_connect_guard, struct socket *sock, struct sockaddr *address, int addrlen) {
    __u32 pid = bpf_get_current_pid_tgid() >> 32;

    /* Intent ID for Network Egress = 0x100 */
    if (!verify_causality(pid, 0x100)) {
        /* Block the orphaned connection */
        return -1; // -EPERM equivalent in LSM
    }

    return 0; // Allow
}
