package helper

type session map[string]map[string]string

// Session 全局内存储存
// 推荐结构为： map[namespace]map[uid]value
var Session = session{}

// AuthTokenNamespace namespace
const AuthTokenNamespace = "AuthToken"

func (v session) Set(namespace string, key string, value string) {
	scopeData := v[namespace]
	if scopeData == nil {
		scopeData = map[string]string{}
		v[namespace] = scopeData
	}

	scopeData[key] = value
}

func (v session) Get(namespace string, key string) string {
	scopeData := v[namespace]
	if scopeData == nil {
		return ""
	}

	return scopeData[key]
}

func (v session) Del(namespace string, key string) {
	scopeData := v[namespace]
	if scopeData == nil {
		return
	}
	v.Set(namespace, key, "")
}

func (v session) ClearByValue(namespace string, value string) {
	scopeData := v[namespace]
	if scopeData == nil {
		return
	}
	for k, v := range scopeData {
		if v == value {
			delete(scopeData, k)
		}
	}
	// v[namespace] = scopeData
}

func (v session) SearchByValue(namespace string, value string) string {
	scopeData := v[namespace]
	if scopeData == nil {
		return ""
	}
	for k, v := range scopeData {
		if v == value {
			return k
		}
	}
	return ""
}

func (v session) Has(namespace string, key string) bool {
	return v.Get(namespace, key) != ""
}
