package info

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Info struct {
	// AppName name of your app
	AppName string

	// WebSite url of your website
	WebSite string

	// AppDescription a AppDescription of your app
	AppDescription string

	// AppVersion your add AppVersion
	AppVersion string

	// BuildTime time of building your app
	BuildTime string

	// GoVersion cmd "go version" of building env
	BuildGoVersion string

	// GCCVersion cmd "`go env CC` --version"
	GCCVersion string

	// git rev-parse HEAD
	GITCommit string

	// GoVersion cmd "go version" of building env
	GoVersion string

	// CopyrightStart copyright year start
	CopyrightStart uint16

	// CopyrightUpdate copyright year update, with CopyrightStart, can generate like xxxx-xxxx
	CopyrightUpdate uint16
}

func (p *Info) copyRightString() string {
	y := p.copyRightYears()

	if y != "" {
		y = y + " "
	}

	return fmt.Sprintf("Copyright %s%s ALL Rights Reserved", y, webSite)
}

func (p *Info) copyRightYears() string {

	yearList := []uint16{}

	if p.CopyrightStart < p.CopyrightUpdate {
		//
		yearList = append(yearList, p.CopyrightStart)
		yearList = append(yearList, p.CopyrightUpdate)

	} else {
		p.CopyrightUpdate = p.CopyrightStart

		yearList = append(yearList, p.CopyrightUpdate)
	}

	if p.CopyrightUpdate == 0 {
		return ""
	}

	if p.CopyrightStart == 0 {
		return fmt.Sprintf("%d", p.CopyrightUpdate)
	}

	stringList := make([]string, 0)
	for _, v := range yearList {
		stringList = append(stringList, fmt.Sprintf("%d", v))
	}

	return strings.Join(stringList, "-")
}

func (p *Info) gccVersionString() string {
	g := ""
	if p.GCCVersion != "" {
		g = fmt.Sprintf(`gcc version  : %s`, p.GCCVersion)
	}

	return g
}

func (p *Info) ShowAll() {

	// copyright
	fmt.Printf(
		`%s for %s
%s

%s

git commit   : %s
build time   : %s
go compile   : %s
go runtime   : %s
%s

%s
Powered by %s
`,
		p.AppName, p.WebSite,
		p.AppVersion,
		p.AppDescription,
		p.GITCommit,
		p.BuildTime,
		p.BuildGoVersion,
		p.GoVersion,
		p.gccVersionString(),
		p.copyRightString(),
		p.WebSite,
	)
}

func FromInput() (p Info, err error) {
	p.AppName = appName
	p.WebSite = webSite
	p.AppDescription = appDescription
	p.AppVersion = appVersion
	p.BuildTime = buildTime
	p.BuildGoVersion = buildGoVersion
	p.GCCVersion = gccVersion
	p.GITCommit = gitCommit
	p.GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	p.CopyrightStart = 0
	p.CopyrightUpdate = 0

	copyrightStart, err := strconv.ParseUint(copyrightStart, 10, 32)
	if err != nil {
		copyrightStart = 0
	}

	copyrightUpdate, err := strconv.ParseUint(copyrightUpdate, 10, 32)
	if err != nil {
		copyrightUpdate = 0
	}

	p.CopyrightStart = uint16(copyrightStart)
	p.CopyrightUpdate = uint16(copyrightUpdate)
	return
}
