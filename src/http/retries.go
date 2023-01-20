package http

import "CASBanSyncer/src/http/types"

type retriesOption int

func (ct retriesOption) Apply(o *types.RequestOptions) {
	o.Retries = int(ct)
}

func Retries(count int) RequestOption {
	return retriesOption(count)
}
