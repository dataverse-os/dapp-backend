package model

type UserVsion struct {
	Base
	Did       string `json:"did,omitempty"`
	FsVersion string `json:"fs_version,omitempty"`
}
