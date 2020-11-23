package qualys

import (
	"github.com/nortonlifelock/aegis/pkg/log"
	"strings"
)

// LoadAssetGroups loads the asset groups passed in through an int slice from Qualys and returns the Assignment Group
// output from Qualys
func (session *Session) LoadAssetGroups(ids []int) (ags *QSAGListOutput, err error) {

	// API Flags
	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["show_attributes"] = "ALL"

	// If the ids are null then return all assignment groups otherwise pass the ids that were passed in as
	// a comma separated string of asset group ids
	if ids != nil {
		fields["ids"] = strings.Join(intArrayToStringArray(ids), ",")
	}

	ags = &QSAGListOutput{}

	if err = session.post(session.Config.Address()+qsAssetGroup, fields, ags); err != nil {
		session.lstream.Send(log.Errorf(err, "nil response while calling api [%s]", qsAssetGroup))
	}

	return ags, err
}
