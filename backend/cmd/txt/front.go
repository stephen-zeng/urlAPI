package txt

func Request(Format, API, Model, Target, Type, IP, Domain string) (TxtResponse, error) {
	if Type == "sum" {
		return TxtResponse{
			Response: "test",
		}, nil
	} else {
		return genRequest(IP, Domain, Model, API, Target)
	}
}
