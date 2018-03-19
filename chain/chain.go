package chain

import ()

type Chain []BlockHeader

func (chn Chain) Verify() bool {
	for i := range chn {
		invi := len(chn) - i - 1
		if invi == 0 {
			return true
		} else {
			if !chn[invi].Verify(&chn[invi-1]) {
				return false
			}
		}
	}
	return false
}
