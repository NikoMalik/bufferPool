package bufferPool

import (
	"sync"

	lowlevelfunctions "github.com/NikoMalik/low-level-functions"
)

const (
	numPools = 0x04
	size     = 0x04
)

func createAllocFunc(size int) func() interface{} {
	return func() interface{} {
		return lowlevelfunctions.MakeNoZero(size)
	}
}

var (
	pool     [numPools]*sync.Pool
	poolSize [numPools]int //An array of pool sizes. These sizes will correspond to the number of bytes each pool will manage
)

func Init() {

	for i := 0; i < numPools; i++ {

		pool[i] = &sync.Pool{
			New: createAllocFunc(1 << (i + 11)),
		}
		poolSize[i] = 1 << (i + 11)
		// fmt.Printf("poolSize[%d] = %d\n", i, poolSize[i])
		//poolSize[0] = 2048
		//poolSize[1] = 4096
		//poolSize[2] = 8192
		//poolSize[3] = 16384

	}
}

// GetPool returns a sync.Pool that can be used to allocate a byte slice
// of a given size. If the size is larger than the maximum size of the
// sync.Pool, it will return nil.
func GetPool(size int) *sync.Pool {
	if pool[0] == nil {
		Init()
	}
	// Iterate over the poolSize array and find the first sync.Pool that
	// can be used to allocate a byte slice of the given size.
	for idx := range poolSize {
		if size <= poolSize[idx] {
			return pool[idx]
		}
	}

	// If no sync.Pool can be used to allocate a byte slice of the given
	// size, return nil.
	return nil

	/*
		Returns a pool that can be used to allocate a byte slice of the given size.
		If the pool is not initialized, it is called to initialize it. If there is no suitable pool, returns nil.

	*/
}

// Allocate returns a byte slice of the given size. If the size is larger
// than the maximum size of the sync.Pool, it will allocate a new byte slice
// using the make function.
func Allocate(size int) []byte {
	// Get the sync.Pool that can be used to allocate a byte slice of the
	// given size. If no sync.Pool can be used to allocate a byte slice of the
	// given size, GetPool will return nil.
	pool := GetPool(size)
	if pool != nil {
		// Get a byte slice from the sync.Pool.
		buf := pool.Get().([]byte)
		// If the byte slice is larger than the requested size, return a slice
		// of the requested size.
		if len(buf) >= size {
			return buf[:size]
		}
		// If the byte slice is not larger than the requested size, return the
		// byte slice itself.
		return buf
	}
	// If no sync.Pool can be used to allocate a byte slice of the given size,
	// allocate a new byte slice using the make function.
	return lowlevelfunctions.MakeNoZero(size)

	/*
		Allocates a byte slice of the desired size. First tries to get a slice from the pool. If the slice from the pool is greater than or equal to the desired size,
		returns a slice of the desired size.
		If the slice is smaller than the desired size, returns a full slice. If there is no suitable pool, allocates memory using the function


	*/
}

func Free(buf []byte) {
	size := cap(buf)
	buf = buf[0:size]
	for idx := range poolSize {
		if len(buf) == poolSize[idx] {

			pool[idx].Put(buf)
			return
		}
	}

	/*
		Returns the byte slice back to the pool.
		Trims the slice to its capacity and places it in a suitable pool. If the slice size does not match the pools size, the slice is not placed back into the pool.

	*/
}
