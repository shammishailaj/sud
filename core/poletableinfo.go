package core

import (
	"errors"
	"strings"

	"github.com/crazyprograms/sud/corebase"
)

const PrefixTable = "t_"

type PoleTableInfo struct {
	Configuration string
	TableName     string
	PoleName      string
	PoleInfo      corebase.IPoleInfo
}

func (pti *PoleTableInfo) FullPoleName(TableName string, PoleName string) string {
	return TableName[len(PrefixTable):] + "." + PoleName
}
func (pti *PoleTableInfo) FromPoleInfo(pi corebase.IPoleInfo) error {
	pti.Configuration = pi.GetConfigurationName()
	s := strings.Split(pi.GetPoleName(), ".")
	s2 := strings.Split(pi.GetRecordType(), ".")
	if len(s) < len(s2) {
		return errors.New("PTI:pole name error")
	}
	for i := 0; i < len(s2); i++ {
		if s[i] != s2[i] {
			return errors.New("PTI:pole name error")
		}
	}
	pti.PoleName = s[len(s)-1]
	pti.TableName = PrefixTable + strings.Join(s[0:len(s)-1], ".")
	pti.PoleInfo = pi
	return nil
}
