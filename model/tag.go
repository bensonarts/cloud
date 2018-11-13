package model

import (
	"encoding/xml"
)

// Tag type that represents a system tag
type Tag struct {
	CanAssign 		bool	`json:"canAssign"`
	UserAssignable  bool	`json:"UserAssignable"`
	UserVisible	    bool	`json:"userVisible"`
	Name 			string	`json:"name"`
}

type SystemTagProperties struct {
	XMLName		xml.Name 			`xml:"DAV: propfind"`
	XmlnsD		string	 			`xml:"xmlns:d,attr"`
	XmlnsOC		string	 			`xml:"xmlns:oc,attr"`
	XmlnsNC		string	 			`xml:"xmlns:nc,attr"`
	Prop		SystemTagProperty 	`xml:"DAV: prop"`
}

type SystemTagProperty struct {
	Id				string		`xml:"http://owncloud.org/ns id"`
	DisplayName		string		`xml:"http://owncloud.org/ns display-name"`
	UserVisible		string		`xml:"http://owncloud.org/ns user-visible"`
	UserAssignable	string		`xml:"http://owncloud.org/ns user-assignable"`
	CanAssign		string		`xml:"http://owncloud.org/ns can-assign"`
}

type MultiStatusTagResponse struct {
	XMLName		xml.Name			`xml:"multistatus"`
	Responses	[]TagPropResponse	`xml:"response"`
}

type TagPropResponse struct {
	XMLName		xml.Name			`xml:"response"`
	Href		string				`xml:"href"`
	Properties	[]SystemTagProperty	`xml:"propstat>prop"`
}

func CreateTagProperties() *SystemTagProperties{
	return &SystemTagProperties{
		XmlnsD: "DAV:",
		XmlnsNC: "http://nextcloud.org/ns",
		XmlnsOC: "http://owncloud.org/ns",
		Prop: SystemTagProperty{},
	}
}

