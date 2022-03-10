package parser

import "strconv"

const _MessageType_name = "ISAUTHENTICATEDPUBLISHREMOVETOKENSETTOKENEVENTACKRECEIVE"

var _MessageType_index = [...]uint8{0, 15, 22, 33, 41, 46, 56}

func (i MessageType) String() string {
	if i < 0 || i >= MessageType(len(_MessageType_index)-1) {
		return "MessageType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MessageType_name[_MessageType_index[i]:_MessageType_index[i+1]]
}
