package model

import (
	"encoding/xml"
)

type Properties struct {
	XMLName		xml.Name `xml:"DAV: propfind"`
	XmlnsD		string	 `xml:"xmlns:d,attr"`
	XmlnsOC		string	 `xml:"xmlns:oc,attr"`
	XmlnsNC		string	 `xml:"xmlns:nc,attr"`
	Prop		Property `xml:"DAV: prop"`
}

type Property struct {
	XMLName				xml.Name	`xml:"DAV: prop"`
	LastModified		string  	`xml:"DAV: getlastmodified"`
	Etag				string  	`xml:"DAV: getetag"`
	ContentType			string  	`xml:"DAV: getcontenttype"`
	ResourceType		string  	`xml:"DAV: resourcetype"`
	ContentLength		string  	`xml:"DAV: getcontentlength"`
	HasPreview			string		`xml:"http://nextcloud.org/ns has-preview"`
	FileId				string		`xml:"http://owncloud.org/ns fileid"`
	Permissions			string		`xml:"http://owncloud.org/ns permissions"`
	Size				string		`xml:"http://owncloud.org/ns size"`
	Favorite			string		`xml:"http://owncloud.org/ns favorite"`
	CommentsUnread		string		`xml:"http://owncloud.org/ns comments-unread"`
	OwnerDisplayName	string		`xml:"http://owncloud.org/ns owner-display-name"`
	ShareTypes			string		`xml:"http://owncloud.org/ns share-types"`
}

type PropResponse struct {
	XMLName		xml.Name	`xml:"response"`
	Href		string		`xml:"href"`
	Properties	[]Property	`xml:"propstat>prop"`
}

type MultiStatusResponse struct {
	XMLName		xml.Name		`xml:"multistatus"`
	Responses	[]PropResponse	`xml:"response"`
}

func CreateProperties() *Properties{
	return &Properties{
		XmlnsD: "DAV:",
		XmlnsNC: "http://nextcloud.org/ns",
		XmlnsOC: "http://owncloud.org/ns",
		Prop: Property{},
	}
}

