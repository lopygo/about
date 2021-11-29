// +build linux

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const modString = "github.com/lopygo/about"

const infoPath = "info"

func main() {

	conf := LoadConfig()

	fmt.Printf("%++v", conf)
	fmt.Println()

	targetOS := string(execShell("go", "env", "GOOS"))
	targetArch := string(execShell("go", "env", "GOARCH"))

	// fmt.
	fmt.Println("target os: ", targetOS)
	fmt.Println("target arch: ", targetArch)

	// create flag info like `-X 'xxx/about.appVersion=v1.2.3'`
	flagInfoStringFunc := func(key, value string) string {
		return fmt.Sprintf(`  -X \"%s/%s.%s=%s\"`, modString, infoPath, key, value)
	}

	//
	now := time.Now()
	buildTimeString := now.Format(time.RFC3339)

	// git
	commitString := string(execShell("git", "rev-parse", "HEAD"))

	// go version
	goVersionString := fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// gcc
	ccCmd := execShell("go", "env", "CC")

	ccCmdList := strings.Split(string(ccCmd), " ")

	ccCmdList = append(ccCmdList, "--version")
	ccCmd2 := execShell(ccCmdList[0], ccCmdList[1:]...)
	ccCmd3 := bytes.SplitN(ccCmd2, []byte{0x0a}, 2)
	gccString := string(ccCmd3[0])

	ldFlats := make([]string, 0)

	ldFlats = append(ldFlats, "-s -w")

	// app info
	ldFlats = append(ldFlats, flagInfoStringFunc("appName", conf.App.Name))
	ldFlats = append(ldFlats, flagInfoStringFunc("appVersion", conf.App.Version))
	ldFlats = append(ldFlats, flagInfoStringFunc("appDescription", conf.App.Description))

	// build env
	ldFlats = append(ldFlats, flagInfoStringFunc("buildGoVersion", goVersionString))
	ldFlats = append(ldFlats, flagInfoStringFunc("gccVersion", gccString))
	ldFlats = append(ldFlats, flagInfoStringFunc("gitCommit", commitString))
	ldFlats = append(ldFlats, flagInfoStringFunc("buildTime", buildTimeString))

	// copyright
	ldFlats = append(ldFlats, flagInfoStringFunc("webSite", conf.Copyright.Website))
	ldFlats = append(ldFlats, flagInfoStringFunc("copyrightStart", fmt.Sprintf("%d", conf.Copyright.StartYear)))
	ldFlats = append(ldFlats, flagInfoStringFunc("copyrightUpdate", now.Format("2006")))

	outputFileName, buildCmdObj := getBuildCmd(conf.App.Name, conf.Build.Source, conf.Build.OutputDir, ldFlats, targetOS, targetArch)

	fmt.Println()
	fmt.Println("build command is : ")
	fmt.Println(buildCmdObj.String())

	//
	dir := filepath.Join(conf.Build.OutputDir)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	filename := filepath.Join(dir, conf.Build.ScriptFile)

	filename, err = filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	content := bytes.NewBufferString(buildCmdObj.String())
	content.WriteString("\necho\n")
	content.WriteString(`echo "show md5 of built bin file"`)
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("md5sum %s", outputFileName))
	content.WriteString("\necho\n")

	err = ioutil.WriteFile(filename, content.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf(`build bash created on location: "%s"`, filename)
	fmt.Println()

	if !conf.Build.Run {
		return
	}
	buildCmd := exec.Command("sh", filename)

	output, err := buildCmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))

}

func getBuildCmd(appName, source, outputDir string, ldFlats []string, osString, archString string) (filename string, cmd *exec.Cmd) {

	extString := ""
	if osString == "windows" {
		extString = ".exe"
	}

	filename = fmt.Sprintf("%s_%s_%s%s", appName, osString, archString, extString)

	filename = filepath.Join(outputDir, filename)
	cmd = exec.Command(
		"go",
		"build",
		"-trimpath",
		"-ldflags",
		fmt.Sprintf(`"%s"`, strings.Join(ldFlats, " \\\n")),
		// fmt.Sprintf(`"%s"`, strings.Join(ldFlats, " ")),
		"\\\n",
		"-o",
		filename,
		"\\\n",
		source,
	)
	return filename, cmd
}
func execShell(name string, args ...string) []byte {
	ccCmd, err := exec.Command(name, args...).Output()
	if err != nil {
		panic(err)
	}

	ccCmd = bytes.TrimRight(ccCmd, "\n")
	ccCmd = bytes.TrimRight(ccCmd, " ")

	return ccCmd
}

type Config struct {
	Build BuildConfig `mapstructure:"build"`

	App AppConfig `mapstructure:"app"`

	Copyright CopyrightConfig `mapstructure:"copyright"`
}

type BuildConfig struct {

	// OutputDir files outout dir
	OutputDir string `mapstructure:"output"`

	// Source source go file used to compile
	Source string `mapstructure:"source"`

	// ScriptFile bash file name
	ScriptFile string `mapstructure:"script"`

	// Run if run script
	Run bool `mapstructure:"run"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`

	Version string `mapstructure:"version"`

	Description string `mapstructure:"description"`
}
type CopyrightConfig struct {
	StartYear uint16 `mapstructure:"start"`

	UpdateYear uint16 `mapstructure:"update"`

	Website string `mapstructure:"website"`
}

func LoadConfig() (conf Config) {

	v := viper.New()

	pflag.String("app.name", "demo", "app name")
	pflag.String("app.version", "0.0.0-test", "app version")
	pflag.String("app.description", "description of app", "build bash file name")

	pflag.String("build.output", ".", "output dir")
	pflag.String("build.source", "main.go", "source go file or dir")
	pflag.String("build.script", "build.sh", "bash file name")
	pflag.Bool("build.run", false, "run script generated")

	pflag.Uint16("copyright.start", 0, "start year of copyright")
	pflag.Uint16("copyright.update", 0, "update year of copyright, no use now.")
	pflag.String("copyright.website", "www.example.com", "your -website")

	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// v.SetDefault("")

	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {

		// 没找到文件，则不用管
		err, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@")
		} else {

			panic(err)
		}

	}

	// fmt.Println(v.Get("copyright.website"))
	err = v.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return
}
