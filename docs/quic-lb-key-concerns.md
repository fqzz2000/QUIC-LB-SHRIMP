# Essential QUIC Protocol Elements for Load Balancer Implementation

## 1. Critical Components to Handle

### 1.1 Connection ID Handling
**Draft Reference**: Section 5.1
- Most important part for load balancing
- Used for routing decisions
- Must handle variable lengths (0-20 bytes)
- Format determined by server

Key points:
```
Connection ID Format:
+--------+------------------+------------+
| Config |    Server ID     |   Nonce    |
| Bits   |                  |            |
+--------+------------------+------------+
```

### 1.2 Packet Types to Parse
**Draft Reference**: Section 17

Only need to handle these packet types:
1. **Initial Packets** (0x00)
   - First packet in connection
   - Contains Destination Connection ID
   - Need for initial routing

2. **Short Header Packets** (0x40)
   - Most common during connection
   - Contains only Destination Connection ID
   - Simpler to parse

Can ignore:
- Handshake packets
- Retry packets
- Version negotiation packets
- Payload contents
- Other header fields

### 1.3 Header Parsing
**Draft Reference**: Section 17.2-17.3

Only need to extract:
1. Packet type (first byte)
2. Connection ID length (for Initial packets)
3. Connection ID

Example packet format checking:
```c
// First byte pattern check
#define LONG_HEADER_MASK  0x80
#define SHORT_HEADER_MASK 0x40

static inline bool is_initial_packet(uint8_t first_byte) {
    return (first_byte & LONG_HEADER_MASK) && 
           ((first_byte & 0x30) == 0x00);
}

static inline bool is_short_header(uint8_t first_byte) {
    return !(first_byte & LONG_HEADER_MASK) && 
           (first_byte & SHORT_HEADER_MASK);
}
```

## 2. What to Ignore

1. **Protocol Negotiation**
   - Version negotiation
   - Transport parameters
   - Crypto handshake

2. **Stream Management**
   - Stream IDs
   - Flow control
   - Stream state

3. **Packet Contents**
   - Frame types
   - Frame contents
   - Payload encryption/decryption

4. **Connection Management**
   - Connection state
   - Idle timeout
   - Error handling

5. **Recovery and Congestion**
   - Loss detection
   - Congestion control
   - ACK handling

## 3. Essential Operations

### 3.1 Connection ID Extraction
```c
struct cid_info {
    uint8_t data[20];  // Max CID length
    uint8_t len;
};

static bool extract_cid(const uint8_t *packet, size_t len, 
                       struct cid_info *cid) {
    if (len < 1) return false;
    
    uint8_t first_byte = packet[0];
    size_t offset;
    
    if (is_initial_packet(first_byte)) {
        // Skip version (4 bytes)
        offset = 5;
        // DCID length is at offset 5
        if (len < 6) return false;
        cid->len = packet[5];
        offset = 6;
    } else if (is_short_header(first_byte)) {
        // Fixed CID length based on config
        offset = 1;
        cid->len = get_config_cid_length();
    } else {
        return false;
    }
    
    if (len < offset + cid->len) return false;
    memcpy(cid->data, packet + offset, cid->len);
    return true;
}
```

### 3.2 Server ID Decoding
```c
static bool decode_server_id(const struct cid_info *cid,
                           uint32_t *server_id) {
    if (cid->len < MIN_CID_LENGTH) return false;
    
    // Extract server ID based on config
    size_t server_id_len = get_config_server_id_length();
    size_t offset = get_config_rotation_bits();
    
    *server_id = extract_bits(cid->data, offset, server_id_len);
    return true;
}
```

### 3.3 Basic Packet Forwarding
```c
static int forward_packet(const void *packet, size_t len,
                        uint32_t server_id) {
    struct server_info *server = get_server(server_id);
    if (!server) return -1;
    
    // Forward packet to server
    return redirect_packet(packet, len, server->address);
}
```

## 4. Key Configuration Parameters

```yaml
load_balancer:
  # Only essential parameters
  connection_id:
    min_length: 8
    max_length: 20
    server_id_length: 8
    config_rotation_bits: 2

  routing:
    # Simple routing config
    server_map:
      - id: 1
        address: "10.0.0.1:443"
      - id: 2
        address: "10.0.0.2:443"
```

## 5. Minimal Processing Flow

```
1. Receive UDP datagram
2. Quick check if it's QUIC (first byte pattern)
3. Extract Connection ID
4. Decode server ID from Connection ID
5. Forward packet to server
```

## 6. Error Cases to Handle

Only need to handle:
1. Invalid packet format
2. Missing/invalid Connection ID
3. Unknown server ID
4. Forwarding failures

Can ignore:
1. Protocol errors
2. Crypto errors
3. Stream errors
4. Connection errors

## 7. Optional Optimizations

1. **Connection Tracking**
```c
struct conn_track_entry {
    uint32_t server_id;
    uint64_t last_seen;
};

// Track connections for faster routing
static struct conn_track_entry *track_connection(
    const struct cid_info *cid) {
    return map_lookup_or_create(conn_track_map, cid);
}
```

2. **Server Health Checking**
```c
struct server_health {
    bool healthy;
    uint64_t last_check;
};

// Basic health checking
static bool is_server_healthy(uint32_t server_id) {
    struct server_health *health = get_server_health(server_id);
    return health && health->healthy;
}
```

This focused approach lets you implement a functional QUIC load balancer without getting bogged down in protocol complexity. Would you like me to explain any of these aspects in more detail?

Remember:
1. Only need to parse packet headers minimally
2. Focus on Connection ID handling
3. Simple packet forwarding
4. Ignore complex protocol features

The goal is fast, efficient routing based primarily on Connection ID.