package http

import "CASBanSyncer/src/http/types"

type noInstantRead bool

func (ct noInstantRead) Apply(o *types.RequestOptions) {
	o.NoInstantRead = bool(ct)
}

func NoInstantRead() RequestOption {
	return noInstantRead(true)
}
