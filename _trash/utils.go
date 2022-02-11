package utils

import "fmt"

func Split(dataShards int, data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("short data")
	}
	dataLen := len(data)
	// Calculate number of bytes per data shard.
	perShard := (len(data) + dataShards - 1) / dataShards

	if cap(data) > len(data) {
		data = data[:cap(data)]
	}

	// Only allocate memory if necessary
	var padding []byte
	if len(data) < (dataShards * perShard) {
		// calculate maximum number of full shards in `data` slice
		fullShards := len(data) / perShard
		padding = make([]byte, dataShards*perShard-perShard*fullShards)
		copy(padding, data[perShard*fullShards:])
		data = data[0 : perShard*fullShards]
	} else {
		for i := dataLen; i < dataLen+dataShards; i++ {
			data[i] = 0
		}
	}

	// Split into equal-length shards.
	dst := make([][]byte, dataShards)
	i := 0
	for ; i < len(dst) && len(data) >= perShard; i++ {
		dst[i] = data[:perShard:perShard]
		data = data[perShard:]
	}

	for j := 0; i+j < len(dst); j++ {
		dst[i+j] = padding[:perShard:perShard]
		padding = padding[perShard:]
	}

	return dst, nil
}
