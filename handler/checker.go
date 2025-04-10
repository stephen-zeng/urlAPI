package handler

import "urlAPI/request"

func (web) checker(r *request.Request) {
	r.Security.General.GeneralChecker()
	r.Security.WebImg.FunctionChecker(&r.Security.General)
}

func (img) checker(r *request.Request) {
	r.Security.General.GeneralChecker()
	r.Security.ImgGen.APIChecker(&r.Security.General)
	r.Security.ImgGen.FunctionChecker(&r.Security.General)
}

func (rand) checker(r *request.Request) {
	r.Security.General.GeneralChecker()
	r.Security.Rand.FunctionChecker(&r.Security.General)
	r.Security.Rand.APIChecker(&r.Security.General)
}

func (txt) checker(r *request.Request) {
	r.Security.General.GeneralChecker()
	r.Security.TxtGen.APIChecker(&r.Security.General)
	r.Security.TxtGen.FunctionChecker(&r.Security.General)
}
