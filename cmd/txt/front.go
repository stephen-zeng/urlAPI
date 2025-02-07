package txt

func Request(API, Model, Target, Type, IP, Domain, Regen string) (TxtResponse, error) {
	if Type == "sum" {
		return TxtResponse{
			Response: "test",
		}, nil
	} else {
		return genRequest(IP, Domain, Model, API, Target, Regen)
	}
}
