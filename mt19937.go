// Package mt19937 implements Mersenne Twister PRNG
package mt19937

const (
	n      = 312
	m      = 156
	hiMask = 0xFFFFFFFF80000000 // Most significant 33 bits
	loMask = 0x000000007FFFFFFF // Least significant 31 bits
)

// A Source represents a source of uniformly-distributed pseudo-random int64 values in the range [0, 1<<63).
type Source struct {
	state  [n]uint64 // The array for the state vector
	index  int
	seeded bool
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
func (s *Source) Seed(seed int64) {
	s.state[0] = uint64(seed)
	for i := 1; i < n; i++ {
		s.state[i] = 6364136223846793005*(s.state[i-1]^(s.state[i-1]>>62)) + uint64(i)
	}
	s.seeded = true
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

// SeedByArray ...
func (s *Source) SeedByArray(keys []uint64) {
	s.Seed(19650218)
	i, j, k := 1, 0, max(n, len(keys))
	for ; k > 0; k-- {
		s.state[i] = (s.state[i] ^ ((s.state[i-1] ^ (s.state[i-1] >> 62)) * 3935559000370003845)) + keys[j] + uint64(j)
		if i++; i >= n {
			s.state[0], i = s.state[n-1], 1
		}
		if j++; j >= len(keys) {
			j = 0
		}
	}
	for k = n - 1; k > 0; k-- {
		s.state[i] = (s.state[i] ^ ((s.state[i-1] ^ (s.state[i-1] >> 62)) * 2862933555777941757)) - uint64(i)
		if i++; i >= n {
			s.state[0], i = s.state[n-1], 1
		}
	}
	s.state[0] = 1 << 63
}

func (s *Source) generate() {
	if !s.seeded {
		s.Seed(5489) // default initial seed
	}
	mag := [2]uint64{0, 0xB5026F5AA96619E9}
	var i int
	for ; i < n-m; i++ {
		x := (s.state[i] & hiMask) | (s.state[i+1] & loMask)
		s.state[i] = s.state[i+m] ^ (x >> 1) ^ mag[int(x&1)]
	}
	for ; i < n-1; i++ {
		x := (s.state[i] & hiMask) | (s.state[i+1] & loMask)
		s.state[i] = s.state[i+(m-n)] ^ (x >> 1) ^ mag[int(x&1)]
	}
	x := (s.state[n-1] & hiMask) | (s.state[0] & loMask)
	s.state[n-1] = s.state[m-1] ^ (x >> 1) ^ mag[int(x&1)]
}

// Uint64 returns a non-negative pseudo-random 64-bit integer as an uint64.
func (s *Source) Uint64() uint64 {
	i := s.index
	if i == 0 || i >= n { // generate NN words at one time
		s.generate()
		i = 0
	}
	x := s.state[i]
	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)
	s.index = i + 1
	return x
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Float64A returns a non-negative pseudo-random 64-bit float
func (s *Source) Float64A() float64 {
	return float64(s.Uint64()>>11) * (1.0 / 9007199254740991.0)
}

// Float64B returns a non-negative pseudo-random 64-bit float
func (s *Source) Float64B() float64 {
	return float64(s.Uint64()>>11) * (1.0 / 9007199254740992.0)
}

// Float64C returns a non-negative pseudo-random 64-bit float
func (s *Source) Float64C() float64 {
	return (float64(s.Uint64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}
