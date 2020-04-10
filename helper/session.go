package helper

type session map[string]map[string]string

// Session 全局内存储存
// 推荐结构为： map[namespace]map[uid]value
var Session session = session{}

// AuthToken namespace
const AuthToken = "AuthToken"

func (v session) Set(namespace string, name string, value string) {
	scopeData := v[namespace]
	if scopeData == nil {
		v[namespace] = map[string]string{}
	}

	scopeData[name] = value
}

func (v session) Get(namespace string, name string) string {
	scopeData := v[namespace]
	if scopeData == nil {
		return ""
	}

	value := scopeData[name]
	if len(value) == 0 {
		return ""
	}

	return value
}

func (v session) Del(namespace string, name string) {
	v.Set(namespace, name, "")
}

func (v session) Has(namespace string, name string) bool {
	return len(v.Get(namespace, name)) == 0
}
