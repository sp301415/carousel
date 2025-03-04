package carousel

import (
	"math"
	"runtime"
	"sync"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

// GenEvaluationKey samples a new evaluation key for bootstrapping.
//
// This can take a long time.
// Use [*Encryptor.GenEvaluationKeyParallel] for better key generation performance.
func (e *Encryptor) GenEvaluationKey() EvaluationKey {
	return EvaluationKey{
		BlindRotateKey:  e.GenBlindRotateKey(),
		KeySwitchKey:    e.GenKeySwitchKeyForBootstrap(),
		AutomorphismKey: e.GenAutomorphismKeyForBootstrap(),
	}
}

// GenEvaluationKeyParallel samples a new evaluation key for bootstrapping in parallel.
func (e *Encryptor) GenEvaluationKeyParallel() EvaluationKey {
	return EvaluationKey{
		BlindRotateKey:  e.GenBlindRotateKeyParallel(),
		KeySwitchKey:    e.GenKeySwitchKeyForBootstrapParallel(),
		AutomorphismKey: e.GenAutomorphismKeyForBootstrapParallel(),
	}
}

// GenBlindRotateKey samples a new bootstrapping key.
//
// This can take a long time.
// Use [*Encryptor.GenBlindRotateKeyParallel] for better key generation performance.
func (e *Encryptor) GenBlindRotateKey() BlindRotateKey {
	brk := NewBlindRotateKey(e.Parameters)

	for i := 0; i < e.Parameters.lweDimension; i++ {
		e.buffer.ptRGSW.Clear()
		e.buffer.ptRGSW.Coeffs[0] = e.SecretKey.LWEKey.Value[i]
		for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
			e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, e.Parameters.blindRotateParameters.BaseQ(k), e.buffer.ctRLWE.Value[0])
			e.EncryptRLWEBody(e.buffer.ctRLWE)
			e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, brk.Value[i].Value[0].Value[k])
		}

		e.PolyEvaluator.ScalarMulAddPolyAssign(e.SecretKey.RLWEKey.Value, e.SecretKey.LWEKey.Value[i], e.buffer.ptRGSW)
		for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
			e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, e.Parameters.blindRotateParameters.BaseQ(k), e.buffer.ctRLWE.Value[0])
			e.EncryptRLWEBody(e.buffer.ctRLWE)
			e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, brk.Value[i].Value[0].Value[k])
		}
	}

	return brk
}

// GenBlindRotateKeyParallel samples a new bootstrapping key in parallel.
func (e *Encryptor) GenBlindRotateKeyParallel() BlindRotateKey {
	brk := NewBlindRotateKey(e.Parameters)

	chunkCount := num.Min(runtime.NumCPU(), int(math.Sqrt(float64(e.Parameters.lweDimension))))

	encryptorPool := make([]*Encryptor, chunkCount)
	for i := range encryptorPool {
		encryptorPool[i] = e.ShallowCopy()
	}

	jobs := make(chan int)
	go func() {
		defer close(jobs)
		for i := 0; i < e.Parameters.lweDimension; i++ {
			jobs <- i
		}
	}()

	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for c := 0; c < chunkCount; c++ {
		go func(idx int) {
			eIdx := encryptorPool[idx]
			for i := range jobs {
				eIdx.buffer.ptRGSW.Clear()
				eIdx.buffer.ptRGSW.Coeffs[0] = e.SecretKey.LWEKey.Value[i]
				for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
					eIdx.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, e.Parameters.blindRotateParameters.BaseQ(k), e.buffer.ctRLWE.Value[0])
					eIdx.EncryptRLWEBody(e.buffer.ctRLWE)
					eIdx.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, brk.Value[i].Value[0].Value[k])
				}

				eIdx.PolyEvaluator.ScalarMulAddPolyAssign(e.SecretKey.RLWEKey.Value, e.SecretKey.LWEKey.Value[i], e.buffer.ptRGSW)
				for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
					eIdx.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, e.Parameters.blindRotateParameters.BaseQ(k), e.buffer.ctRLWE.Value[0])
					eIdx.EncryptRLWEBody(e.buffer.ctRLWE)
					eIdx.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, brk.Value[i].Value[0].Value[k])
				}
			}
			wg.Done()
		}(c)
	}
	wg.Wait()

	return brk
}

// GenKeySwitchKeyForBootstrap samples a new keyswitch key LWELargeKey -> LWEKey,
// used for bootstrapping.
//
// This can take a long time.
// Use [*Encryptor.GenKeySwitchKeyForBootstrapParallel] for better key generation performance.
func (e *Encryptor) GenKeySwitchKeyForBootstrap() LWEKeySwitchKey {
	skIn := LWESecretKey{Value: e.SecretKey.LWELargeKey.Value[e.Parameters.lweDimension:]}
	ksk := NewKeySwitchKeyForBootstrap(e.Parameters)

	for i := 0; i < ksk.InputLWEDimension(); i++ {
		for j := 0; j < e.Parameters.keySwitchParameters.level; j++ {
			ksk.Value[i].Value[j].Value[0] = skIn.Value[i] << e.Parameters.keySwitchParameters.LogBaseQ(j)

			e.UniformSampler.SampleVecAssign(ksk.Value[i].Value[j].Value[1:])
			ksk.Value[i].Value[j].Value[0] += -vec.Dot(ksk.Value[i].Value[j].Value[1:], e.SecretKey.LWEKey.Value)
			ksk.Value[i].Value[j].Value[0] += e.GaussianSampler.Sample(e.Parameters.LWEStdDev())
		}
	}

	return ksk
}

// GenKeySwitchKeyForBootstrapParallel samples a new keyswitch key LWELargeKey -> LWEKey in parallel,
// used for bootstrapping.
func (e *Encryptor) GenKeySwitchKeyForBootstrapParallel() LWEKeySwitchKey {
	skIn := LWESecretKey{Value: e.SecretKey.LWELargeKey.Value[e.Parameters.lweDimension:]}
	ksk := NewKeySwitchKeyForBootstrap(e.Parameters)

	workSize := ksk.InputLWEDimension() * e.Parameters.keySwitchParameters.level
	chunkCount := num.Min(runtime.NumCPU(), int(math.Sqrt(float64(workSize))))

	encryptorPool := make([]*Encryptor, chunkCount)
	for i := range encryptorPool {
		encryptorPool[i] = e.ShallowCopy()
	}

	jobs := make(chan [2]int)
	go func() {
		defer close(jobs)
		for i := 0; i < ksk.InputLWEDimension(); i++ {
			for j := 0; j < e.Parameters.keySwitchParameters.level; j++ {
				jobs <- [2]int{i, j}
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for c := 0; c < chunkCount; c++ {
		go func(idx int) {
			eIdx := encryptorPool[idx]
			for jobs := range jobs {
				i, j := jobs[0], jobs[1]
				ksk.Value[i].Value[j].Value[0] = skIn.Value[i] << eIdx.Parameters.keySwitchParameters.LogBaseQ(j)
				eIdx.UniformSampler.SampleVecAssign(ksk.Value[i].Value[j].Value[1:])
				ksk.Value[i].Value[j].Value[0] += -vec.Dot(ksk.Value[i].Value[j].Value[1:], eIdx.SecretKey.LWEKey.Value)
				ksk.Value[i].Value[j].Value[0] += eIdx.GaussianSampler.Sample(eIdx.Parameters.LWEStdDev())
			}
			wg.Done()
		}(c)
	}
	wg.Wait()

	return ksk
}

// GenAutomorphismKeyForBootstrap samples a new automorphism key for bootstrapping.
//
// This can take a long time.
// Use [*Encryptor.GenAutomorphismKeyForBootstrapParallel] for better key generation performance.
func (e *Encryptor) GenAutomorphismKeyForBootstrap() []RLWEKeySwitchKey {
	skPermute := NewRLWESecretKey(e.Parameters)
	atk := make([]RLWEKeySwitchKey, e.Parameters.polyDegree)
	for i := 0; i < e.Parameters.polyDegree; i++ {
		e.PolyEvaluator.PermutePolyAssign(e.SecretKey.RLWEKey.Value, i, skPermute.Value)
		atk[i] = e.GenGLWEKeySwitchKey(skPermute, e.Parameters.keySwitchParameters)
	}
	return atk
}

// GenAutomorphismKeyForBootstrapParallel samples a new automorphism key for bootstrapping in parallel.
func (e *Encryptor) GenAutomorphismKeyForBootstrapParallel() []RLWEKeySwitchKey {
	atk := make([]RLWEKeySwitchKey, e.Parameters.polyDegree)

	chunkCount := num.Min(runtime.NumCPU(), int(math.Sqrt(float64(e.Parameters.polyDegree))))

	encryptorPool := make([]*Encryptor, chunkCount)
	skPermutePool := make([]RLWESecretKey, chunkCount)
	for i := range encryptorPool {
		encryptorPool[i] = e.ShallowCopy()
		skPermutePool[i] = NewRLWESecretKey(e.Parameters)
	}

	jobs := make(chan int)
	go func() {
		defer close(jobs)
		for i := 0; i < e.Parameters.polyDegree; i++ {
			jobs <- i
		}
	}()

	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for c := 0; c < chunkCount; c++ {
		go func(idx int) {
			eIdx := encryptorPool[idx]
			skPermute := skPermutePool[idx]
			for i := range jobs {
				eIdx.PolyEvaluator.PermutePolyAssign(eIdx.SecretKey.RLWEKey.Value, i, skPermute.Value)
				atk[i] = eIdx.GenGLWEKeySwitchKey(skPermute, eIdx.Parameters.keySwitchParameters)
			}
			wg.Done()
		}(c)
	}
	wg.Wait()

	return atk
}
