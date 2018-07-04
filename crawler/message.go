package crawler

const (
	LOG_MESSAGE_CLOSE = iota
	LOG_MESSAGE_FATAL
	LOG_MESSAGE_ERROR
	LOG_MESSAGE_WARN
	LOG_MESSAGE_INFO
	LOG_MESSAGE_DEBUG
)

type LogService interface {
	Log(message LogMessage)
	Close()
}

type LogMessage interface {
	GetLevel() int
	GetMessage() string
	GetExtraData() ExtraDataMessage
}

type ExtraDataMessage struct {
	Url          string
	ResponseCode int
	Domain       string
}

type CrawlerMessage struct {
	Message string
	ExtraData ExtraDataMessage
}

type InfoMessage struct {
	CrawlerMessage
}

type WarnMessage struct {
	CrawlerMessage
}

type ErrorMessage struct {
	CrawlerMessage
}

type DebugMessage struct {
	CrawlerMessage
}

type CloseMessage struct {
	CrawlerMessage
}

func (m CrawlerMessage) GetLevel() int {
	return LOG_MESSAGE_CLOSE
}

func (m CrawlerMessage) GetMessage() string {
	return m.Message
}

func (m CrawlerMessage) GetExtraData() ExtraDataMessage {
	return m.ExtraData
}


func (m InfoMessage) GetLevel() int {
	return LOG_MESSAGE_INFO
}

func (m WarnMessage) GetLevel() int {
	return LOG_MESSAGE_WARN
}

func (m ErrorMessage) GetLevel() int {
	return LOG_MESSAGE_ERROR
}

func (m DebugMessage) GetLevel() int {
	return LOG_MESSAGE_DEBUG
}

func Info(msg string, data ...ExtraDataMessage) LogMessage{
	m := InfoMessage{}
	m.Message = msg

	if len(data) > 0 {
		m.ExtraData = data[0]
	}


	return m
}

func Warn(msg string, data ...ExtraDataMessage) LogMessage {
	m := InfoMessage{}
	m.Message = msg

	if len(data) > 0 {
		m.ExtraData = data[0]
	}

	return m
}

func Error(msg string, data ...ExtraDataMessage) LogMessage {
	m := ErrorMessage{}
	m.Message = msg

	if len(data) > 0 {
		m.ExtraData = data[0]
	}


	return m
}

func Debug(msg string, data ...ExtraDataMessage) LogMessage {
	m := DebugMessage{}
	m.Message = msg

	if len(data) > 0 {
		m.ExtraData = data[0]
	}

	return m
}