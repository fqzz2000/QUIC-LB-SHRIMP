# QUIC Version Negotiation Handling

## 1. Why Version Negotiation Must Be Handled

### 1.1 Security and Privacy
**Reference**: Section 6 and Section 21.13
```
Client         LB          Server
  |            |            |
  |--Initial-->|            |  (Unsupported Version)
  |            |--Initial-->|
  |            |<---VN-----|  
  |<---VN-----|            |  (LB must rewrite source IP)
  |            |            |
```

Key reasons:
1. Protect server IP addresses from exposure
2. Maintain server pool privacy
3. Keep consistent client-facing endpoint
4. Prevent direct client-server communication

### 1.2 Protocol Flow
```c
struct version_negotiation {
    uint8_t type;
    uint32_t version;     // Set to 0
    uint8_t dcid_len;
    uint8_t dcid[20];     // From client's Initial
    uint8_t scid_len;
    uint8_t scid[20];     // From client's Initial
    uint32_t versions[];  // Supported versions
};

// LB handling
if (is_version_negotiation(packet)) {
    // Rewrite source IP to LB address
    rewrite_source_address(packet, lb_address);
    // Forward to client
    return forward_to_client(packet);
}
```

## 2. Implementation Requirements

### 2.1 Basic Version Negotiation Handling
```c
static bool handle_version_negotiation(const uint8_t *packet, size_t len) {
    struct version_negotiation_packet *vn = (struct version_negotiation_packet *)packet;
    
    // 1. Validate packet
    if (!is_valid_vn_packet(vn, len))
        return false;
        
    // 2. Keep original Connection IDs
    // (Important: VN must echo client's Connection IDs)
    
    // 3. Rewrite source address to LB address
    rewrite_source_address(packet, get_lb_address());
    
    // 4. Forward to client
    return forward_to_client(packet);
}
```

### 2.2 Connection ID Handling
```c
// Version Negotiation must preserve Connection IDs
static bool validate_vn_connection_ids(const struct version_negotiation_packet *vn,
                                     const struct original_packet *initial) {
    return (
        vn->dcid_len == initial->scid_len &&
        memcmp(vn->dcid, initial->scid, vn->dcid_len) == 0 &&
        vn->scid_len == initial->dcid_len &&
        memcmp(vn->scid, initial->dcid, vn->scid_len) == 0
    );
}
```

### 2.3 Source Address Rewriting
```c
struct lb_address {
    union {
        struct in_addr  v4;
        struct in6_addr v6;
    } addr;
    uint16_t port;
    bool is_ipv6;
};

static void rewrite_version_negotiation_source(
    struct packet *pkt,
    const struct lb_address *lb_addr) {
    
    if (lb_addr->is_ipv6) {
        // Rewrite IPv6 source
        pkt->ip6->ip6_src = lb_addr->addr.v6;
    } else {
        // Rewrite IPv4 source
        pkt->ip->ip_src = lb_addr->addr.v4;
    }
    // Rewrite UDP source port
    pkt->udp->source = lb_addr->port;
}
```

## 3. Processing Flow

1. **Receiving Initial Packet**
```
if (unsupported_version(initial_packet)) {
    // Forward to server to generate VN
    forward_to_server(initial_packet);
}
```

2. **Receiving Version Negotiation**
```
if (is_version_negotiation(packet)) {
    // Must handle to protect server identity
    rewrite_and_forward_vn(packet);
}
```

3. **State Maintenance**
```c
struct vn_state {
    uint8_t client_dcid[20];
    uint8_t client_scid[20];
    uint8_t dcid_len;
    uint8_t scid_len;
};

// Optional: track VN state if needed
static void track_vn_state(const struct initial_packet *initial) {
    struct vn_state state = {
        .dcid_len = initial->dcid_len,
        .scid_len = initial->scid_len,
    };
    memcpy(state.client_dcid, initial->dcid, initial->dcid_len);
    memcpy(state.client_scid, initial->scid, initial->scid_len);
    
    // Store state if needed
    store_vn_state(&state);
}
```

## 4. Security Considerations

1. **Address Privacy**
```c
// Always ensure server address is never exposed
static bool verify_source_address(const struct packet *pkt) {
    return is_lb_address(pkt->source_addr);
}
```

2. **Connection ID Verification**
```c
// Ensure Connection IDs are properly echoed
static bool verify_connection_ids(
    const struct version_negotiation *vn,
    const struct stored_state *state) {
    
    return (
        vn->dcid_len == state->client_scid_len &&
        vn->scid_len == state->client_dcid_len &&
        memcmp(vn->dcid, state->client_scid, vn->dcid_len) == 0 &&
        memcmp(vn->scid, state->client_dcid, vn->scid_len) == 0
    );
}
```

3. **Version List Protection**
```c
// Optionally filter/modify supported versions
static void process_version_list(uint32_t *versions, size_t *count) {
    // Could filter versions based on LB policy
    // Could remove versions not supported by all servers
}
```

You're right to be concerned about server privacy. The load balancer must handle Version Negotiation packets to maintain server anonymity and ensure proper operation of version negotiation while keeping the server pool hidden from clients. This is an important security consideration in QUIC load balancer implementation.