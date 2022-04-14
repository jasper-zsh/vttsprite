package avutil

//#cgo pkg-config: libavutil
//#include <libavutil/error.h>
//#include <stdlib.h>
//static const char *error2string(int code) { return av_err2str(code); }
import "C"
import "errors"

const (
	EAGAIN1 = -11
	EAGAIN2 = -35
)

func ErrorFromCode(code int) error {
	if code >= 0 {
		return nil
	}

	return errors.New(C.GoString(C.error2string(C.int(code))))
}
