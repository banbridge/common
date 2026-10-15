package logs

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
)

const (
	// IPUnknown represents unknown ip
	// 32 * 0
	IPUnknown = "00000000000000000000000000000000"
)

const (
	CTXKeyLogID  = "CTX_LOGID" // logid key
	HEADERLogKey = "X-TT-Logid"
)

const (
	version    = "02"
	length     = 53
	maxRandNum = 1 << 12
)

var (
	localIP      string
	defaultLogID LogID
)

func init() {
	ip, err := getLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	localIP = ipToBytes(ip)

	defaultLogID = NewLogID()
}

// LogID represents a logID generator
type LogID struct{}

// NewLogID create a new LogID instance
func NewLogID() LogID {
	return LogID{}
}

// GenLogID return a new logID string
func (l LogID) GenLogID() string {
	r := fastrand.Uint32n(maxRandNum)
	now := time.Now()
	sb := strings.Builder{}
	sb.Grow(length)
	sb.WriteString(now.Format("20060102150405"))
	sb.WriteString(version)
	sb.WriteString(strconv.FormatInt(now.UnixMilli(), 10))
	sb.WriteString(string(localIP))
	sb.WriteString(fmt.Sprintf("%03X", r))
	return sb.String()
}

// GenLogID return a new logID
func GenLogID() string {
	return defaultLogID.GenLogID()
}

func getLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("no non-loopback IP address found")
}

func ipToBytes(ip net.IP) string {
	ipBytes := ip.To4()
	if ipBytes == nil {
		return IPUnknown
	}
	// 将每个字节转换为十六进制字符串并拼接
	var hexParts []string
	for _, b := range ipBytes {
		hexParts = append(hexParts, fmt.Sprintf("%02X", b))
	}
	return strings.Join(hexParts, "")
}

// SetLogID set logid to context.Context / 在 context.Context 中设置 logid
func SetLogID(ctx context.Context, logid string) context.Context {
	return context.WithValue(ctx, CTXKeyLogID, logid)
}

// CtxLogID get logid from context.Context / 从 context.Context 中获取 logid
//
// return (val, ok), if ok == false, then val must be empty string
func CtxLogID(ctx context.Context) (string, bool) {
	return getStringFromContext(ctx, CTXKeyLogID)
}

// LogIDDefault get logid from context.Context / 从 context.Context 中获取 logid and not return bool
//
// if logid not found in ctx, then return default logid: '-'
func LogIDDefault(ctx context.Context) string {
	val, _ := CtxLogID(ctx)
	if val == "" {
		return "-"
	}
	return val
}

func getStringFromContext(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(key)
	if v == nil {
		return "", false
	}

	switch v := v.(type) {
	case string:
		return v, true
	case *string:
		if v == nil {
			return "", false
		}
		return *v, true
	}
	return "", false
}
