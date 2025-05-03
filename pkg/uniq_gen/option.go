package uniq_gen

import "github.com/yitter/idgenerator-go/idgen"

type UniqGenOption func(op *idgen.IdGeneratorOptions)

// WithGenMethod 雪花计算方法,（1-漂移算法|2-传统算法），默认1
func WithGenMethod(method int) UniqGenOption {
	return func(op *idgen.IdGeneratorOptions) {
		op.Method = uint16(method)
	}
}

// WithBaseTime 基础时间（ms单位），不能超过当前系统时间
func WithBaseTime(baseTime int64) UniqGenOption {
	return func(op *idgen.IdGeneratorOptions) {
		op.BaseTime = baseTime
	}
}

// WithWorkerIdBitLength 机器码位长，默认值6，取值范围 [1, 15]（要求：序列数位长+机器码位长不超过22）
func WithWorkerIdBitLength(workerIdBitLength byte) UniqGenOption {
	return func(op *idgen.IdGeneratorOptions) {
		op.WorkerIdBitLength = workerIdBitLength
	}
}

// WithSeqBitLength 序列数位长，默认值6，取值范围 [3, 21]（要求：序列数位长+机器码位长不超过22）
func WithSeqBitLength(seqBitLength byte) UniqGenOption {
	return func(op *idgen.IdGeneratorOptions) {
		op.SeqBitLength = seqBitLength
	}
}

// WithTopOverCostCount 最大漂移次数（含），默认2000，推荐范围500-10000（与计算能力有关）
func WithTopOverCostCount(topOverCostCount int) UniqGenOption {
	return func(op *idgen.IdGeneratorOptions) {
		op.TopOverCostCount = uint32(topOverCostCount)
	}
}
