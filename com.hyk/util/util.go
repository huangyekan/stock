package util


var Layout = "20060102" //时间格式
var Layout_2 = "2006-01-02"
var Layout_3 = "2006/01/02"
var Layout_4 = "01/02/2006"

type Util struct {

}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}
