package main

import (
	"fmt"
	"github.com/json-iterator/go"
)

var s = `{
    "steps":[
        {
            "type":"if",
            "condition":"7ca6f3ed-c953-4b38-b503-5f29db425ca7",
            "rollback":"a",
            "y":[
                {
                    "type":"pipeline",
                    "steps":"000"
                },
                {
                    "type":"if",
                    "condition":"a == b",
                    "y":[

                    ],
                    "n":[

                    ]
                }
            ],
            "n":[
                {
                    "type":"pipeline",
                    "steps":"111"
                },
                {
                    "type":"if",
                    "condition":"c == d",
                    "y":[
                        {
                            "type":"if",
                            "condition":"e == f",
                            "y":[
                                {
                                    "type":"if",
                                    "condition":"g == h",
                                    "y":[

                                    ],
                                    "n":[

                                    ]
                                }
                            ],
                            "n":[
                                {
                                    "type":"pipeline",
                                    "steps":"222"
                                },
                                {
                                    "type":"if",
                                    "condition":"j == k",
                                    "y":[
                                        {
                                            "type":"for",
                                            "condition":"opNo",
                                            "i_begin":"0",
                                            "i_step":"11",
                                            "i_end":"11",
                                            "loop":[
                                                {
                                                    "type":"pipeline",
                                                    "steps":"333"
                                                }
                                            ]
                                        }
                                    ],
                                    "n":[

                                    ]
                                }
                            ]
                        }
                    ],
                    "n":[

                    ]
                }
            ]
        },
        {
            "type":"switch",
            "condition":"opNo",
            "ops":{
                "a":[
                    {
                        "type":"pipeline",
                        "steps":"444"
                    }
                ],
                "b":[
                    {
                        "type":"pipeline",
                        "steps":"555"
                    }
                ],
                "default":[
                    {
                        "type":"pipeline",
                        "steps":"666"
                    }
                ]
            }
        },
        {
            "type":"for",
            "condition":"opNo",
            "i_begin":"10",
            "i_step":"30",
            "i_end":"40",
            "loop":[
                {
                    "type":"pipeline",
                    "steps":"777"
                }
            ]
        }
    ]
}`

var s2 = `{
"type":"h",
"time":"1"
}`

func main() {
	//Unmarshal(s)
	var m map[string]string
	_ = jsoniter.Unmarshal([]byte(s2), &m)
	fmt.Println(m["type"])
	fmt.Println(m["time"])
}

type Op struct {
	Type     string `json:"type"`
	Rollback string `json:"rollback"`
	// pipeline
	Steps string `json:"steps"`
	// if
	Condition string `json:"condition"` // op_no
	Y         []*Op  `json:"y"`
	N         []*Op  `json:"n"`
	// switch
	// -- Condition string // op_no
	Ops map[string][]*Op `json:"ops"`
	// for
	IBegin string `json:"i_begin"`
	IStep  string `json:"i_step"`
	IEnd   string `json:"i_end"`
	Loop   []*Op  `json:"loop"`
}

type Setps struct {
	Steps []*Op
}

// 解析
func Unmarshal(rule string) {
	st := &Setps{
		Steps: []*Op{},
	}
	if err := jsoniter.Unmarshal([]byte(rule), &st); err != nil {
		fmt.Println("------->", err.Error())
		return
	}
	Parsing(st.Steps)
}
func Parsing(op []*Op) {
	for _, v := range op {
		rollback := v.Rollback
		if rollback != "" {
			fmt.Println("rollback -----> ", rollback)
		}
		switch v.Type {
		case "if":
			condition := v.Condition
			fmt.Println("condition -----> ", condition)
			if len(v.Y) > 0 {
				Parsing(v.Y)
			}
			if len(v.N) > 0 {
				Parsing(v.N)
			}
		case "pipeline": // 顺序
			opNo := v.Steps
			fmt.Println("opNo-----> ", opNo)
		case "switch": // switch
			condition := v.Condition
			fmt.Println("condition -----> ", condition)
			if len(v.Ops["a"]) > 0 {
				Parsing(v.Ops["a"])
			}
			if len(v.Ops["b"]) > 0 {
				Parsing(v.Ops["b"])
			}
			if len(v.Ops["default"]) > 0 {
				Parsing(v.Ops["default"])
			}
		case "for": // for
			iBegin := v.IBegin
			iStep := v.IStep
			iEnd := v.IEnd
			condition := v.Condition
			fmt.Println("condition -----> ", condition)
			fmt.Println("iBegin -----> ", iBegin)
			fmt.Println("iStep -----> ", iStep)
			fmt.Println("iEnd -----> ", iEnd)
			if len(v.Loop) > 0 {
				Parsing(v.Loop)
			}
		}
	}
}

//func Parsing(op []*Op) {
//	for _, v := range op {
//		rollback := v.Rollback
//		fmt.Println("rollback -----> ", rollback)
//
//		switch v.Type {
//		case "if":
//			if len(v.Y) > 0 {
//				for _, v1 := range v.Y {
//					condition := v1.Condition
//					fmt.Println("condition-----> ", condition)
//					switch v1.Type {
//					case "pipeline": // 顺序
//						opNo := v1.Steps
//						fmt.Println("opNo-----> ", opNo)
//					case "if":
//						if len(v1.Y) > 0 {
//							Parsing(v1.Y)
//						}
//						if len(v1.N) > 0 {
//							Parsing(v1.N)
//						}
//					}
//				}
//			}
//			if len(v.N) > 0 {
//				for _, v2 := range v.N {
//					condition := v2.Condition
//					fmt.Println("condition -----> ", condition)
//					switch v2.Type {
//					case "pipeline": // 顺序
//						opNo := v2.Steps
//						fmt.Println("opNo -----> ", opNo)
//					case "if":
//						if len(v2.Y) > 0 {
//							Parsing(v2.Y)
//						}
//						if len(v2.N) > 0 {
//							Parsing(v2.N)
//						}
//					}
//				}
//			}
//		case "switch":
//			condition := v.Condition
//			fmt.Println("condition -----> ", condition)
//			if len(v.Ops["a"]) > 0 {
//				Parsing(v.Ops["a"])
//			}
//			if len(v.Ops["b"]) > 0 {
//				Parsing(v.Ops["b"])
//			}
//			if len(v.Ops["default"]) > 0 {
//				Parsing(v.Ops["default"])
//			}
//		case "for":
//			iBegin := v.IBegin
//			iStep := v.IStep
//			iEnd := v.IEnd
//			condition := v.Condition
//			fmt.Println("condition -----> ", condition)
//			fmt.Println("iBegin -----> ", iBegin)
//			fmt.Println("iStep -----> ", iStep)
//			fmt.Println("iEnd -----> ", iEnd)
//			if len(v.Loop) > 0 {
//				Parsing(v.Loop)
//			}
//		}
//	}
//}
