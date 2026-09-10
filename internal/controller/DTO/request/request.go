package request

import "mime/multipart"

type Registration struct {
	Token    string `json:"token" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

type UploadFile struct {
	Meta       UploadMeta            `form:"meta" json:"meta"`
	JSON       interface{}           `form:"json" json:"json,omitempty"`
	File       multipart.File        `form:"file" json:"-"`
	FileHeader *multipart.FileHeader `form:"file" json:"-"`
}

type UploadMeta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

type DocsListQuery struct {
	Token string `form:"token"` // опционально
	// Key   string `form:"key"`   // имя колонки для фильтрации
	// Value string `form:"value"` // значение фильтра
	Limit int `form:"limit"` // кол-во документов
}
