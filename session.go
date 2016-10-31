package yiyun

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Manager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}
type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
	IsOverTime(maxLifeTime int64) bool
	UpdateTime()
}

var provides = make(map[string]Provider)

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provider
}

func init() {
	Info("register memory session provider")
	Register("memory", &memoryProvider{data: make(map[string]Session)})
}

var globalSessions *Manager

//GetGlobalSessions  Then, initialize the session manager
func GetGlobalSessions() *Manager {
	if globalSessions == nil {
		var err error
		globalSessions, err = NewManager("memory", "gosessionid", int64(time.Minute*20))
		if err != nil {
			panic("session manager start error." + err.Error())
		}
		go globalSessions.GC()
	}
	return globalSessions
}

//NewManager create new manager
func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func (manager *Manager) sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//GetSession START
func (manager *Manager) GetSession(rs *fasthttp.Response, re *fasthttp.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie := re.Header.Cookie(manager.cookieName)
	//Debug("session cookie id:", string(cookie))
	if cookie == nil {
		return nil
	}
	sid, _ := url.QueryUnescape(string(cookie))
	session, _ = manager.provider.SessionRead(sid)
	return
}

//CreateSession 创建session
func (manager *Manager) CreateSession(rs *fasthttp.Response, re *fasthttp.Request) (session Session) {
	sid := manager.sessionID()
	session, _ = manager.provider.SessionInit(sid)
	cookie := fasthttp.AcquireCookie()
	cookie.SetHTTPOnly(true)
	//必须设置成“／”才能使session全站生效
	cookie.SetPath("/")
	cookie.SetExpire(time.Now().Add(time.Duration(manager.maxlifetime)))
	cookie.SetKey(manager.cookieName)
	cookie.SetValue(url.QueryEscape(sid))
	rs.Header.SetCookie(cookie)
	return
}

//SessionDestroy sessionid
func (manager *Manager) SessionDestroy(rs *fasthttp.Response, re *fasthttp.Request) {
	v := re.Header.Cookie(manager.cookieName)
	if v != nil {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionDestroy(string(v))
	cookie := fasthttp.AcquireCookie()
	cookie.SetHTTPOnly(true)
	cookie.SetExpire(time.Now())
	cookie.SetKey(manager.cookieName)
	cookie.SetValue(string(v))
	rs.Header.SetCookie(cookie)
}

//GC GC
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}

type sessionList struct {
	next   *sessionList
	pre    *sessionList
	me     Session
	isHead bool
	isEnd  bool
}
type memoryProvider struct {
	data     map[string]Session
	dataList *sessionList
}

func (m *memoryProvider) SessionInit(sid string) (Session, error) {
	if m.data == nil {
		return nil, fmt.Errorf("memoryProvider is not init")
	}
	s := m.createMemorySession(sid)
	m.data[sid] = s
	//链表头部
	if m.dataList == nil {
		m.dataList = &sessionList{
			me:     s,
			isHead: true,
			isEnd:  true,
		}
	} else {
		if m.dataList.isHead {
			m.dataList.pre = &sessionList{
				next:   m.dataList,
				me:     s,
				isHead: true,
			}
			m.dataList.isHead = false
			m.dataList = m.dataList.pre
		}
	}
	return s, nil
}
func (m *memoryProvider) SessionRead(sid string) (Session, error) {
	if m.data == nil {
		return nil, fmt.Errorf("memoryProvider is not init")
	}
	if m.data[sid] == nil {
		return nil, fmt.Errorf("session is not exit")
	}
	//更新session时间
	m.data[sid].UpdateTime()
	return m.data[sid], nil
}
func (m *memoryProvider) SessionDestroy(sid string) error {
	if m.data == nil {
		return fmt.Errorf("memoryProvider is not init")
	}
	if m.dataList == nil {
		return fmt.Errorf("memoryProvider no info")
	}
	if m.data[sid] == nil {
		return fmt.Errorf("session is not exit")
	}
	delete(m.data, sid)
	sl := m.dataList
	if sl.me.SessionID() == sid {
		dl := m.dataList.next
		dl.pre = nil
		dl.isHead = true
		m.dataList = dl
	} else {
		for {
			if sl.next == nil {
				break
			}
			sl := sl.next
			if sl.me.SessionID() == sid {
				sl.pre.next = sl.next
				sl.next.pre = sl.pre
				break
			}
		}
	}
	return nil
}
func (m *memoryProvider) SessionGC(maxLifeTime int64) {
	if m.dataList == nil {
		return
	}
	var sl sessionList
	sl = *m.dataList
	if sl.me.IsOverTime(maxLifeTime) && sl.isHead {
		Info("session gc close headSession :", sl.me.SessionID())
		if sl.next != nil {
			sl.next.pre = nil
			sl.next.isHead = true
		}
		m.dataList = sl.next
		delete(m.data, sl.me.SessionID())
	}
	for {
		if sl.next == nil {
			break
		}
		sl = *sl.next
		if sl.me.IsOverTime(maxLifeTime) && !sl.isEnd {
			if sl.pre != nil {
				sl.pre.next = sl.next
			}
			if sl.next != nil {
				sl.next.pre = sl.pre
			}
			delete(m.data, sl.me.SessionID())
			Info("session gc close :", sl.me.SessionID())
		} else if sl.me.IsOverTime(maxLifeTime) && sl.isEnd {
			if sl.pre != nil {
				sl.pre.next = nil
			}
			delete(m.data, sl.me.SessionID())
			Info("session gc close endSession :", sl.me.SessionID())
		}
	}
}

func (m *memoryProvider) createMemorySession(sid string) *memorySession {
	return &memorySession{
		data:       make(map[interface{}]interface{}, 0),
		createTime: time.Now(),
		sessionID:  sid,
	}
}

type memorySession struct {
	data       map[interface{}]interface{}
	createTime time.Time
	updateTime time.Time
	sessionID  string
}

var SessionTimeOutError error

func (m *memorySession) IsOverTime(maxLifeTime int64) bool {
	return time.Now().After(m.updateTime.Add(time.Duration(maxLifeTime)))
}
func (m *memorySession) UpdateTime() {
	m.updateTime = time.Now()
}
func (m *memorySession) Set(key, value interface{}) error {
	if m.data == nil {
		return fmt.Errorf("session data is nil")
	}
	m.data[key] = value
	return nil
} //set session value
//IsSessionTimeOutError 判断是否为session过时错误
func IsSessionTimeOutError(err error) bool {
	if err.Error() == "session time out" {
		return true
	}
	return false
}
func (m *memorySession) Get(key interface{}) interface{} {
	return m.data[key]
} //get session value
func (m *memorySession) Delete(key interface{}) error {
	if m.data[key] == nil {
		return fmt.Errorf("key is not exit")
	}
	delete(m.data, key)
	return nil
} //delete session value
func (m *memorySession) SessionID() string {
	return m.sessionID
} //back current sessionID
