package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var SrcPath *string
var DestPath *string
var IncludeDates *bool

const ResourcesFolder string = "resources"

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

type FileInfo struct {
	name          string
	metaIndex     int
	metaId        string
	metaParentId  string
	metaType      int //1:Article 2:Folder 4:Resource 5:Tag 6:Tag-Article association
	metaFileExt   string
	metaCreatedAt string
	metaUpdatedAt string
	metaNoteId    string
	metaTagId     string
}

func (fi FileInfo) getValidName() string {
	r := strings.NewReplacer(
		"*", ".",
		"\"", "''",
		"\\", "-",
		"/", "_",
		"<", ",",
		">", ".",
		":", ";",
		"|", "-",
		"?", "!")
	return r.Replace(fi.name)
}

type Folder struct {
	*FileInfo
	parent *Folder
}

func (f Folder) getPath() string {
	return path.Join(*DestPath, f.getRelativePath())
}

func (f Folder) getRelativePath() string {
	if f.parent == nil {
		return f.getValidName()
	} else {
		return path.Join(f.parent.getRelativePath(), f.getValidName())
	}
}

type Article struct {
	*FileInfo
	folder  *Folder
	content string
	prefix  string
}

func (a Article) getPath() string {
	return fmt.Sprintf("%s.md", path.Join(a.folder.getPath(), a.getValidName()))
}
func (a Article) save() {
	filePath := a.getPath()
	dirName := path.Dir(filePath)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		CheckError(err)
	}
	// optional meta info to sort by time like Joplin
	meta := a.FileInfo
	if meta != nil {
		// if we have accumulated tags, add the tag category label
		if len(a.prefix) > 0 {
			a.prefix = "Tags:\n" + a.prefix
		}
		// add the dates and Joplin ID only if not suppressed
		if *IncludeDates {
			if meta.metaCreatedAt != "" && meta.metaUpdatedAt != "" && meta.metaId != "" {
				a.prefix += fmt.Sprintf("created: %v\nupdated: %v\njoplin_id: %v\n",
					meta.metaCreatedAt, meta.metaUpdatedAt, meta.metaId,
				)
			}
		}
		if len(a.prefix) > 0 {
			a.prefix = "---\n" + a.prefix + "---\n"
		}
	}

	err := os.WriteFile(filePath, []byte(a.prefix+a.content), 0644)
	CheckError(err)

	// optionally change mtime and atime
	// 2021-07-10T02:10:03.850Z
	if meta != nil && meta.metaCreatedAt != "" && meta.metaUpdatedAt != "" {
		updatedAt, err := time.Parse(time.RFC3339, meta.metaUpdatedAt)
		CheckError(err)
		err = os.Chtimes(filePath, updatedAt, updatedAt)
		CheckError(err)
		// the following is macOS-specific, changes the creation date
		createdAt, err := time.Parse(time.RFC3339, meta.metaCreatedAt)
		CheckError(err)
		cmd := exec.Command("SetFile", "-d", createdAt.Format("01/02/2006 15:04:05"), filePath)
		err = cmd.Run()
		CheckError(err)
	}
}

type Resource struct {
	*FileInfo
}

func (r Resource) getFileName() string {
	var fileName string
	if /*len(r.metaFileExt) > 0*/ false {
		fileName = fmt.Sprintf("%s.%s", r.metaId, r.metaFileExt)
	} else {
		resPath := path.Join(*SrcPath, "resources")
		c, err := os.ReadDir(resPath)
		CheckError(err)
		for _, entry := range c {
			if entry.IsDir() {
				continue
			}
			if strings.Index(entry.Name(), r.metaId) >= 0 {
				fileName = entry.Name()
				break
			}
		}
	}
	return fileName
}
