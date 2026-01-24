package o2go

type parseParamsCallbackFunc func(key, value string)

func parseParams(reservedParams map[string]struct{}, paramsData map[string]string, callback parseParamsCallbackFunc) {
	for key, value := range paramsData {
		if _, reserved := reservedParams[key]; reserved {
			continue
		}

		callback(key, value)
	}
}

func baseReservedParams(extraReservedParams []string) map[string]struct{} {
	reservedParams := map[string]struct{}{
		"client_id":     {},
		"client_secret": {},
		"redirect_uri":  {},
		"code":          {},
		"refresh_token": {},
	}

	for _, value := range extraReservedParams {
		reservedParams[value] = struct{}{}
	}

	return reservedParams
}
