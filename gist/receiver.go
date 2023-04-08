package gist

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

// 对 *Runtime 类型实现 UnMarshalJSON 方法
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unQuotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(unQuotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	*r = Runtime(i)
	return nil
}
/*
	$ curl -d '{"title": "Moana", "runtime": "107 mins"}' 10.0.1.17:4000/v1/movies
	  {Title:Moana Year:0 Runtime:107 Genres:[]}
	Runtime 自动实现上述 *Runtime实现的 UnMarshalJSON 方法，将请求体的"107 mins"序列化为107(int32)
	如果要手动实现一个值类型的方法，对于传入的值(r Runtime)的修改就不那么方便了
*/
func decoderExample(w http.ResponseWriter, r *http.Request) {
	//匿名结构体
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime Runtime 	 `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	_ := json.NewDecoder(r.Body).Decode(&input)
	//	...
}
