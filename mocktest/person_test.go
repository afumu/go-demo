package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"test-demo/mocktest/equipment"
	"testing"
)

func TestPerson_DayLife(t *testing.T) {
	type fields struct {
		name  string
		phone equipment.Phone
	}

	mockCtl := gomock.NewController(t)
	mockPhone := equipment.NewMockPhone(mockCtl)

	// 设置mockPhone对象的接口方法返回值
	mockPhone.EXPECT().ZhiHu().Return(true)
	mockPhone.EXPECT().WeiXin().Return(true)
	mockPhone.EXPECT().WangZhe().Return(true)

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"case1", fields{"iphone6s", equipment.NewIphone6s()}, true},
		{"case2", fields{"mocked phone", mockPhone}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Person{
				name:  tt.fields.name,
				phone: tt.fields.phone,
			}
			assert.Equal(t, tt.want, x.DayLife())
		})
	}
}
