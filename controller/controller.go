package controller

import (
	"idtp/values"
)

func ProcessDataAs(
	data []byte,
	isFirstTime bool,
	config *values.Configuration) values.ProcessAs {

	if isFirstTime {
		if len(data) == 1 {
			return values.ProcessAs{
				Type:      values.PT_ERROR,
				PreRcCode: values.RC_MISSING_PAYLOAD,
			}
		}

		if data[0] == values.CC_ONE_TIME_CONN {
			if config.OperationMode != values.OP_MODE_STRICT {
				return values.ProcessAs{
					Type:      values.PT_REQUEST,
					PreRcCode: 0,
				}
			}

			return values.ProcessAs{
				Type:      values.PT_ERROR,
				PreRcCode: values.RC_ONE_TIME_CONNECTION_NOT_ALLOWED,
			}
		}

		if data[0] == values.CC_PERSISTENT_CONN {
			return values.ProcessAs{
				Type:      values.PT_CREATE_CONNECTION,
				PreRcCode: 0,
			}
		}

		return values.ProcessAs{
			Type:      values.PT_ERROR,
			PreRcCode: values.RC_UNKNOWN_CONNECTION_CODE,
		}
	}

	if len(data) == 1 {
		if data[0] == values.CC_PING {
			return values.ProcessAs{
				Type:      values.PT_PING,
				PreRcCode: 0,
			}
		}

		if data[0] == values.CC_DISCONNECTION {
			return values.ProcessAs{
				Type:      values.PT_DISCONNECTION,
				PreRcCode: 0,
			}
		}

		return values.ProcessAs{
			Type:      values.PT_ERROR,
			PreRcCode: values.RC_UNKNOWN_CONNECTION_CODE,
		}
	}

	return values.ProcessAs{
		Type:      values.PT_REQUEST,
		PreRcCode: 0,
	}
}
