package controllers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Minnek-Digital-Studio/cominnek/config"
	"github.com/Minnek-Digital-Studio/cominnek/controllers/bridge"
	github_controller "github.com/Minnek-Digital-Studio/cominnek/controllers/github"
	"github.com/Minnek-Digital-Studio/cominnek/controllers/loading"
	"github.com/Minnek-Digital-Studio/cominnek/pkg/shell"
	"github.com/fatih/color"
)

var currentVersion = config.Public.Version
var allOk = true
var osName = runtime.GOOS
var maxToCheck = 10

func CheckUpdates(printMessage bool) bool {
	latestVersion := github_controller.GetLatestVersion()

	if currentVersion != latestVersion {
		if printMessage {
			fmt.Print("\n\n")
			color.HiYellow("🎉🎉🎉 A new version of cominnek is available! 🎉🎉🎉")
			fmt.Println(color.MagentaString(currentVersion), "→ ", color.GreenString(latestVersion))
			fmt.Print("\n")
			fmt.Println("Run", color.HiGreenString("'cominnek update'"), "to update or download the latest version from:")
			color.HiBlue("https://github.com/Minnek-Digital-Studio/cominnek/releases/latest/")
		}

		return true
	}

	return false
}

func checkDistToUnmount(mountOut string, firstNumber int, lastNumber int) string {
	str := "/dev/disk" + fmt.Sprint(firstNumber) + "s" + fmt.Sprint(lastNumber)
	disk := strings.Contains(mountOut, str)

	if disk {
		return str
	}

	if lastNumber == 5 && firstNumber < maxToCheck {
		return checkDistToUnmount(mountOut, firstNumber+1, 1)
	}

	if lastNumber == 5 && firstNumber == maxToCheck {
		return ""
	}

	return checkDistToUnmount(mountOut, firstNumber, lastNumber+1)
}

func getMountedDisc(mountOut string, name string, num int) string {
	str := "/Volumes/" + name
	var preTest = false

	if num == 0 {
		preTest = strings.Contains(mountOut, str+" ")
	}

	if num == maxToCheck {
		return str
	}

	if num > 0 {
		str = str + " " + fmt.Sprint(num)
	}

	mounted := strings.Contains(mountOut, str)

	if mounted && !preTest {
		return str
	}

	return getMountedDisc(mountOut, name, num+1)
}

func mountDisk(route string, name string) (string, string) {
	err, out, _ := shell.Out("hdiutil attach " + route)

	if err != nil {
		fmt.Println(err)
		allOk = false
	}

	disk := checkDistToUnmount(out, 1, 1)
	mounted := getMountedDisc(out, name, 0)

	return disk, mounted
}

func installUpdates(route string, fileName string) {
	if osName == "windows" {
		shell.ExecuteCommand(`Start-Process -FilePath "`+route+`" -Argument "/silent" -PassThru`, false)

		if allOk {
			color.HiBlue("\n🎉🎉🎉 cominnek " + github_controller.GetLatestVersion() + " has been downloaded successfully! 🎉🎉🎉")
		}
	}

	if osName == "darwin" {
		name := strings.Split(fileName, ".dmg")[0]
		loading.Start("📦 Installing " + color.HiGreenString(name))

		disk, mounted := mountDisk(route, name)
		shell.ExecuteCommand("cd "+mounted+"; bash installer.sh", false)
		shell.ExecuteCommand("hdiutil detach "+disk, false)

		loading.Stop()

		if allOk {
			color.HiBlue("\n🎉🎉🎉 cominnek " + github_controller.GetLatestVersion() + " has been updated successfully! 🎉🎉🎉")
		}
	}
}

func Update() {
	if !CheckUpdates(false) {
		fmt.Println("🥳🎈 You are using the latest version of cominnek")
		return
	}

	fileName := github_controller.GetLatestFileName()
	url := "https://github.com/Minnek-Digital-Studio/cominnek/releases/latest/download/" + fileName
	route := bridge.DownloadFromURL(url, fileName)

	installUpdates(route, fileName)
}
