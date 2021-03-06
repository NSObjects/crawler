package util

//
//const (
//	ERROR = iota
//	WARNING
//	INFO
//	DEBUG
//	STATISTICS
//)
//
//var levelStr = []string{
//	"ERROR",
//	"WARNING",
//	"INFO",
//	"DEBUG",
//	"STATISTICS",
//}
//
//var (
//	eIP     string
//	AppName string
//)
//
//func init() {
//	if ip, err := externalIP(); err == nil {
//		eIP = ip
//	}
//}
//
//func Debugln(level int, args ...interface{}) {
//	_, file, line, _ := runtime.Caller(1)
//	Println(file, line, DEBUG, level, args)
//}
//
//func Infoln(level int, args ...interface{}) {
//	_, file, line, _ := runtime.Caller(1)
//	Println(file, line, INFO, level, args)
//}
//
//func Warningln(level int, args ...interface{}) {
//	_, file, line, _ := runtime.Caller(1)
//	Println(file, line, WARNING, level, args)
//}
//
//func Errorln(level int, args ...interface{}) {
//	_, file, line, _ := runtime.Caller(1)
//	Println(file, line, ERROR, level, args)
//}
//
//func Statistics(level int, args ...interface{}) {
//	_, file, line, _ := runtime.Caller(1)
//	Println(file, line, STATISTICS, level, args)
//}
//
//func Println(file string, line int, errorType int, level int, args ...interface{}) {
//
//	etype := levelStr[errorType]
//	msg := lastPath(file, "/")
//	msg += ":"
//	msg += fmt.Sprintf("%d", line)
//	msg += " "
//	msg += fmt.Sprintln(args...)
//
//	if level > 3 {
//		SendLog(msg, errorType)
//	} else {
//		logMsg := log{Message: msg, Ip: eIP}
//		t := time.Now()
//		logMsg.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d:%02d\n",
//			t.Year(), t.Month(), t.Day(),
//			t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
//		logMsg.Type = etype
//		logMsg.AppName = AppName
//
//		if errorType == DEBUG {
//			fmt.Println(logMsg.Time, msg)
//		} else {
//
//			//logs := make([]log, 0)
//			//logs = append(logs, logMsg)
//			//go sendLog(logs)
//		}
//	}
//}
//
//func SendDingdingMsg(msg dingDingMsg, errorType int) {
//	if errorType == ERROR {
//		msg.At.AtMobiles = append(msg.At.AtMobiles, "13802426870")
//		msg.At.IsAtAll = false
//
//	}
//	jsonMsg, err := json.Marshal(&msg)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	body := bytes.NewBuffer(jsonMsg)
//	// Create client
//	client := &http.Client{}
//	token := "3d3a5891d6e002fa7d9ed6d0e6d3c12f1b851c84cd852b32943b9418db9f0753"
//	if errorType == ERROR {
//		token = "7437afe7d4b7db5eb1255d0e9ce75113f0357064c214f062493f8763d9b77862"
//	}
//
//	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
//	// Create request
//	req, err := http.NewRequest("POST", url, body)
//	// Headers
//	req.Header.Add("Content-Type", "application/json; charset=utf-8")
//	parseFormErr := req.ParseForm()
//	if parseFormErr != nil {
//		fmt.Println(parseFormErr)
//	}
//	// Fetch Request
//	_, err = client.Do(req)
//	if err != nil {
//		fmt.Println("Failure : ", err)
//		return
//	}
//
//}
//
//func SendLog(msgs string, errorType int) {
//	t := time.Now()
//	m := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d:%02d\n",
//		t.Year(), t.Month(), t.Day(),
//		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()) + msgs
//	_, file, line, _ := runtime.Caller(1)
//	content := fmt.Sprintf("%s:%d:%s", file, line, m)
//
//	msg := dingDingMsg{Msgtype: "text", Text: text{Content: content}}
//	SendDingdingMsg(msg, errorType)
//}
//
//type log struct {
//	Type    string `json:"type"`
//	Message string `json:"message"`
//	Time    string `json:"time"`
//	Ip      string `json:"ip"`
//	AppName string `json:"app_name"`
//	Level   int    `json:"level"`
//}
//
//type logJSON struct {
//	Data []log `json:"data"`
//}
//
//type dingDingMsg struct {
//	Msgtype string `json:"msgtype"`
//	Text    text   `json:"text"`
//	At      struct {
//		AtMobiles []string `json:"atMobiles"`
//		IsAtAll   bool     `json:"isAtAll"`
//	} `json:"at"`
//}
//
//type text struct {
//	Content string `json:"content"`
//}
//
//func lastPath(s string, sep string) string {
//	ts := strings.Split(s, sep)
//	return ts[len(ts)-1]
//}
//
//func sendLog(log []log) {
//	if len(log) <= 0 {
//		return
//	}
//	logjson := logJSON{Data: log}
//	body, err := json.Marshal(&logjson)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	client := &http.Client{}
//
//	req, err := http.NewRequest("POST", "http://108.61.162.82:9528/api/log", bytes.NewBuffer(body))
//
//	req.Header.Add("Content-Type", "application/json; charset=utf-8")
//
//	resp, err := client.Do(req)
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	if err != nil {
//		fmt.Println("Failure : ", err)
//	}
//}
