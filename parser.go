package main

import (
	"bufio"
	"fmt"
	"os"
)

var st_look_ahead_token Token
var st_look_ahead_token_exists int

var paramList map[string]Token //变量列表

func my_get_token(token *Token) {
	if st_look_ahead_token_exists == 1 {
		*token = st_look_ahead_token
		st_look_ahead_token_exists = 0
	} else {
		getToken(token)
	}
}

func unget_token(token *Token) {
	st_look_ahead_token = *token
	st_look_ahead_token_exists = 1
}

func parse_primary_expression() interface{} {
	var token Token
	var value interface{} = 0
	var minus_flages int = 0

	my_get_token(&token)
	if token.kind == SUB_OPERATOR_TOKEN {
		minus_flages = 1
	} else {
		unget_token(&token)
	}

	my_get_token(&token)

	//判断是否声明变量
	if token.kind == STATE_TYPE_TOKEN {
		state_token := token
		my_get_token(&token)
		//获取变量名，放入map
		if token.kind == STATE_TOKEN {
			if state_token.str == "let" {
				token.stateType = LET
			} else if state_token.str == "set" {
				token.stateType = SET
			}
			stk := token
			my_get_token(&token)
			//获取变量类型
			if token.kind == TOKEN_TYPE_TOKEN {
				stk.tokenType = getTokenType(token.str)
				if stk.tokenType == ERRORTYPE {
					fmt.Println("error type : ", token.str)
					os.Exit(1)
				}

				my_get_token(&token)

				//变量后续是否为赋值操作
				if token.kind == ASS_OPERATOR_TOKEN {
					value = parse_expression()
					// tokentype := getTokenType(reflect.TypeOf(value).String())
					// if stk.tokenType != tokentype {
					// 	fmt.Println("The type of variable assignment is not consistent : ", token.str)
					// 	os.Exit(1)
					// }

					value = getValue(stk.tokenType, value, minus_flages)

					if _, ok := paramList[stk.str]; ok {
						fmt.Println("error the variable is existed : ", stk.str)
						os.Exit(1)
					} else {
						stk.value = value
						paramList[stk.str] = stk
					}
				} else {
					unget_token(&token)
					paramList[stk.str] = stk
					return value
				}
			}
		}
	}

	/*获取变量，返回value*/
	if token.kind == STATE_TOKEN {
		//判断是否已经声明
		if tk, ok := paramList[token.str]; ok {
			my_get_token(&token)
			//变量后续是否为赋值操作
			if token.kind == ASS_OPERATOR_TOKEN {
				value = parse_expression()
				value = getValue(tk.tokenType, value, minus_flages)
				tk.value = value
				paramList[tk.str] = tk
				fmt.Println(tk.str, " :: ", tk.value)
			} else {
				unget_token(&token)
				if t, ok := paramList[tk.str]; ok {
					fmt.Println(t.str, " : ", t.value)
					//遗留问题
					if minus_flages == 1 {
						return -t.value.(float32)
					}
					return t.value
				} else {
					fmt.Println("Undeclared variables : ", tk.str)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("an undeclared variable : ", tk.str)
			os.Exit(1)
		}
	}

	//如果是常量
	if token.kind == NUMBER_TOKEN {
		value = token.value.(float64)
	} else if token.kind == LEFT_PAREN_TOKEN {
		value = parse_expression()
		my_get_token(&token)
		if token.kind != RIGHT_PAREN_TOKEN {
			fmt.Println("missing ')' error.")
			os.Exit(1)
		}

	} else {
		unget_token(&token)
	}
	//遗留问题
	if minus_flages == 1 {
		value = -value.(float32)
	}
	return value
}

func parse_term() interface{} {
	var v1 interface{}
	var v2 interface{}
	var token Token

	v1 = parse_primary_expression()
	for {
		my_get_token(&token)
		if token.kind != DIV_OPERATOR_TOKEN && token.kind != MUL_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_primary_expression()
		fmt.Println("kind...", token.kind, "str...", token.str)
		if token.kind == MUL_OPERATOR_TOKEN {
			v1 = v1.(float64) * v2.(float64)
		} else if token.kind == DIV_OPERATOR_TOKEN {
			v1 = v1.(float64) / v2.(float64)
		}
	}
	//fmt.Println("v1...", v1)
	return v1
}

func parse_expression() float64 {
	var v1 interface{}
	var v2 interface{}
	var token Token

	v1 = parse_term()
	for {
		my_get_token(&token)
		if token.kind != ADD_OPERATOR_TOKEN && token.kind != SUB_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_term()
		if token.kind == ADD_OPERATOR_TOKEN {
			v1 = v1.(float64) + v2.(float64)
		} else if token.kind == SUB_OPERATOR_TOKEN {
			v1 = v1.(float64) - v2.(float64)
		} else {
			unget_token(&token)
		}
	}
	return v1.(float64)
}

//获取变量类型
func getTokenType(str string) TokenType {
	if str == "int8" {
		return INT8
	} else if str == "int16" {
		return INT16
	} else if str == "int32" {
		return INT32
	} else if str == "int64" {
		return INT64
	} else if str == "uint8" {
		return UINT8
	} else if str == "uint16" {
		return UINT16
	} else if str == "uint32" {
		return UINT32
	} else if str == "uint64" {
		return UINT64
	} else if str == "bool" {
		return BOOL
	} else if str == "float32" {
		return FLOAT32
	} else if str == "float64" {
		return FLOAT64
	} else if str == "string" {
		return STRING
	} else if str == "rune" {
		return CHAR
	} else {
		return ERRORTYPE
	}
}

//获取value类型
func getValue(t TokenType, value interface{}, minus_flages int) interface{} {
	if t == INT8 {
		if minus_flages == 1 {
			value = -value.(int8)
		} else {
			value = value.(int8)
		}
	} else if t == INT16 {
		if minus_flages == 1 {
			value = -value.(int16)
		} else {
			value = value.(int16)
		}
	} else if t == FLOAT32 {
		if minus_flages == 1 {
			value = -value.(float32)
		} else {
			value = value.(float32)
		}
	} else if t == FLOAT64 {
		if minus_flages == 1 {
			value = -value.(float64)
		} else {
			value = value.(float64)
		}
	}
	return value
}

func parse_line() float64 {
	var value float64

	st_look_ahead_token_exists = 0
	value = parse_expression()

	return value
}

func main() {
	var value float64
	paramList = make(map[string]Token) //变量列表
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Println("please input:")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("There ware errors reading,exiting program.")
			return
		}
		set_line([]rune(input))
		value = parse_line()
		fmt.Println(">>", value)
	}

}
