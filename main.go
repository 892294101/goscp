package main

import (
	"github.com/892294101/goscp/dec"
	"github.com/892294101/goscp/gossh"
	"github.com/alecthomas/kingpin/v2"
	"github.com/povsister/scp"
	"log"
	"os"
	"strings"
)

func main() {

	app := kingpin.New("qnscp", "a versatile command-line tool").Version("1.0.0")
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	// 定义一个全局标志
	uplaodAcd := app.Flag("acd", "create a non-existent directory on the target host").Default("enable").Enum("enable", "disable")

	// 定义一个子命令：serve
	uploadCmd := app.Command("upload", "upload files locally to the target server")
	uploadHost := uploadCmd.Flag("host", "enter the target host address").Required().String()
	uploadPort := uploadCmd.Flag("port", "enter the target host port").Default("22").Uint()
	uploadHostUser := uploadCmd.Flag("user", "Remote target username").Required().String()
	uploadHostPass := uploadCmd.Flag("password", "Remote target username").Required().String()
	uploadTargetLocation := uploadCmd.Flag("dest", "upload destination directory").Required().String()
	uploadLocalLocation := uploadCmd.Flag("local", "local directory or file").Required().String()

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Printf("%s\n\n", err)
		os.Exit(1)
	}

	// 根据不同的子命令执行不同的操作
	switch command {
	case uploadCmd.FullCommand():
		Upload(&dec.UploadStruct{LocalLocation: *uploadLocalLocation, Host: *uploadHost, Port: *uploadPort, HostUser: *uploadHostUser, HostPass: *uploadHostPass, TargetLocation: *uploadTargetLocation, CreateDir: *uplaodAcd})
	default:
		app.Usage([]string{command}) // 显示指定命令的帮助信息
	}

	return
}

func Upload(u *dec.UploadStruct) {
	log.Printf("Start uploading: local %v ------>> remote %v/%v@%v:%v", u.LocalLocation, u.HostUser, u.HostPass, u.Host, u.TargetLocation)
	scpClient, err := scp.NewClient(u.Host, scp.NewSSHConfigFromPassword(u.HostUser, u.HostPass), &scp.ClientOption{})
	checkErr(scpClient, err)
	defer scpClient.Close()
	uploadTo(scpClient, u)
}

func uploadTo(scpClient *scp.Client, us *dec.UploadStruct) {
	sc, err := gossh.NewSSHClient(us)
	checkErr(scpClient, err)

	err = sc.CreateDir(us.CreateDir)
	checkErr(scpClient, err)
	defer sc.Close()
	LocalLocation := strings.Split(us.LocalLocation, " ")
	for _, filedir := range LocalLocation {
		if !isExist(filedir) {
			log.Printf("skip file or direcotry not exist: %v\n", filedir)
			continue
		}
		switch {
		case isFile(filedir):
			log.Printf("Start uplaod file: %v\n", filedir)
			err = scpClient.CopyFileToRemote(filedir, us.TargetLocation, &scp.FileTransferOption{})
			checkErr(scpClient, err)
			log.Printf("end uplaod file: %v\n", filedir)
		case isDir(filedir):
			log.Printf("Start uplaod directory: %v\n", filedir)
			err = scpClient.CopyDirToRemote(filedir, us.TargetLocation, &scp.DirTransferOption{})
			checkErr(scpClient, err)
			log.Printf("end uplaod directory: %v\n", filedir)
		}
	}
	log.Printf("All uploaded directories or files completed\n")
}

func checkErr(client *scp.Client, err error) {
	if err != nil {
		log.Printf("error: %s\n\n", err)
		client.Close()
		os.Exit(1)
	}
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

/*	// 定义另一个子命令：config
	configCmd := app.Command("config", "Manage configuration")
	setConfigCmd := configCmd.Command("set", "Set a configuration value")
	setConfigKey := setConfigCmd.Arg("key", "Configuration key").Required().String()
	setConfigValue := setConfigCmd.Arg("value", "Configuration value").Required().String()

	setConfigCmd2 := configCmd.Command("set2", "Set a configuration value")
	_ = setConfigCmd2.Arg("key", "Configuration key").Required().String()

	case setConfigCmd.FullCommand():
	fmt.Printf("Setting configuration: %s=%s\n", *setConfigKey, *setConfigValue)
*/
