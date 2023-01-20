package http

import "CASBanSyncer/src/http/types"

type RequestOption interface {
	Apply(o *types.RequestOptions)
}

type MultiPartOption interface {
	Apply(o *types.MultiPartInfo)
}
