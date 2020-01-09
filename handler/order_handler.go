package handler

import (
	"encoding/json"
	"fmt"
	"github.com/onmpw/JYGO/model"
	"server/client"
	"server/parser"
)

type OrderHandler struct {
}

type User struct {
	Id        int
	Name      string
	Mobile    string
	Address   string
	GoodsName string
	OrderInfo string
	Oid       string
}

func (u *User) TableName() string {
	return "user_info"
}

func (handle *OrderHandler) Init() error {
	model.Init()
	model.RegisterModel(new(User))
	return nil
}

// Handle 实现 ContractHandler 接口 Handle
// 处理接收到的订单信息
func (handle *OrderHandler) Handle(c *client.Client) bool {
	str, err := parser.GetString(c.Data)
	if err != nil {
		c.Err = err
		c.ErrorResponse()
		return false
	}

	var data map[string]interface{}

	if err = json.Unmarshal([]byte(str), &data); err != nil {
		c.Err = err
		c.ErrorResponse()
		return false
	}
	var user User
	user.Name = convertToString(data["Name"])
	user.Mobile = convertToString(data["Mobile"])
	user.Address = convertToString(data["Address"])
	user.GoodsName = convertToString(data["GoodsName"])
	user.OrderInfo = convertToString(data["OrderInfo"])
	user.Oid = convertToString(data["Oid"])

	count := model.Read(new(User)).Filter("oid", user.Oid).Count()

	if count == 0 {
		lastInsertId, err := model.Add(user)

		if err != nil {
			c.Err = err
			c.ErrorResponse()
			return false
		}
		str = fmt.Sprintf("处理成功，添加Id:%d", lastInsertId)
	} else {
		str = fmt.Sprintf("订单已经存在，无须再处理！")
	}

	c.Result = str
	c.Response()

	return true
}

func convertToString(data interface{}) string {
	return fmt.Sprintf("%s", data)
}
