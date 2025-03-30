package security

import (
	"github.com/pkg/errors"
	"urlAPI/util"
)

var (
	txt  = []string{"openai", "alibaba", "deepseek", "otherapi"}
	img  = []string{"openai", "alibaba"}
	rand = []string{"github", "gitee"}
)

//func (*General) APIChecker() {}

func (info *TxtGen) APIChecker(general *General) error {
	if !(util.ListChecker(&txt, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
		return errors.WithStack(errors.New(general.Info))
	}
	return nil
}

func (info *TxtSum) APIChecker(general *General) error {
	if !(util.ListChecker(&txt, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
		return errors.WithStack(errors.New(general.Info))
	}
	return nil
}

func (info *ImgGen) APIChecker(general *General) error {
	if !(util.ListChecker(&img, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
		return errors.WithStack(errors.New(general.Info))
	}
	return nil
}

func (info *Rand) APIChecker(general *General) error {
	if !(util.ListChecker(&rand, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
		return errors.WithStack(errors.New(general.Info))
	}
	return nil
}
