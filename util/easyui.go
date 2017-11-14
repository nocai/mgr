package util

type Combobox struct {
	Id   int64  `json:"value"`
	Text string `json:"text"`
}

type TreeNode struct {
	Id         int64       `json:"id"`
	Text       string      `json:"text"`
	State      string      `json:"state"`
	Checked    bool        `json:"checked"`
	Attributes interface{} `json:"attributes"`
	Children   []TreeNode  `json:"children"`
}
