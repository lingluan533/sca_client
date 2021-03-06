package main

//https://gitee.com/ying32/govcl/wikis/pages?sort_id=2030600&doc_id=102420
//https://gitee.com/ying32/govcl/wikis/pages?sort_id=2693253&doc_id=102420
//https://gitee.com/ying32/govcl/wikis/pages?sort_id=2782863&doc_id=102420

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
	"github.com/ying32/govcl/vcl/types/keys"
	"golang.org/x/sys/windows/registry"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"syscall"

	"os/exec"
	"os/user"
	"path/filepath"
	"sca_client/datastruct"
	"strings"

	"time"

	"golang.org/x/sys/windows"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
)

type TMainForm struct {
	*vcl.TForm
	ImgList *vcl.TImageList
	ActList *vcl.TActionList

	BtnApply  *vcl.TButton
	BtnOk     *vcl.TButton
	BtnCancel *vcl.TButton

	Pgc    *vcl.TPageControl
	Sheet1 *vcl.TTabSheet
	Sheet2 *vcl.TTabSheet
	Sheet3 *vcl.TTabSheet

	//Sheet1
	Lb1          *vcl.TLabel
	DefaultBox   *vcl.TLabel
	Lb2          *vcl.TLabel
	DefaultBrain *vcl.TLabel
	Lb3          *vcl.TLabel
	Template     *vcl.TComboBox
	Lb4          *vcl.TLabel
	ModifyBox    *vcl.TComboBox
	BtnNewBox    *vcl.TButton
	BtnDeleteBox *vcl.TButton

	GroupBox1      *vcl.TGroupBox
	Lb5            *vcl.TLabel
	Type           *vcl.TLabel
	Lb6            *vcl.TLabel
	Vendor         *vcl.TLabel
	Lb7            *vcl.TLabel
	Url            *vcl.TEdit
	Lb8            *vcl.TLabel
	User           *vcl.TEdit
	Lb9            *vcl.TLabel
	Password       *vcl.TEdit
	CbDefaultBox   *vcl.TCheckBox
	CbDefaultBrain *vcl.TCheckBox
	VirtualDisk    *vcl.TComboBox
	BtnTestMount   *vcl.TButton

	Lb10              *vcl.TLabel
	DefaultInterval   *vcl.TEdit
	AllPassword       *vcl.TEdit
	BtnChangePassword *vcl.TButton

	Lb11              *vcl.TLabel
	Lb12              *vcl.TLabel
	DefaultLogFile    *vcl.TEdit
	BtnViewLog        *vcl.TButton
	CbDefaultAutoBoot *vcl.TCheckBox

	LbRefresh  *vcl.TLabel
	BtnRefresh *vcl.TButton

	oldBox string //?????????????????????box??????????????????????????????????????????????????????Gconf???????????????
	//Sheet2

	//Sheet3

	VPNConfigGroupBox *vcl.TGroupBox
	JoinVPNLabel      *vcl.TLabel
	JoinVPNEdit       *vcl.TComboBox
	JoinVPNButton     *vcl.TButton
	AutoJoinVPN       *vcl.TCheckBox
	DeleteVPN         *vcl.TButton

	VPNStatusLabel *vcl.TLabel
	VPNStatusShape *vcl.TShape
	VPNStatus      *vcl.TLabel

	VPNIPLabel *vcl.TLabel
	VPNIPEdit  *vcl.TEdit

	VPNUsernameLabel *vcl.TLabel
	VPNUsernameEdit  *vcl.TEdit
	DeviceTypeLabel  *vcl.TLabel
	DeviceTypeBox    *vcl.TComboBox
	ChangeNameButton *vcl.TButton
	VPNTest          *vcl.TButton

	LocalDevice  *vcl.TImage
	LocalAndHezi *vcl.TImage
	Hezi         *vcl.TImage
	HeziandVPN   *vcl.TImage
	VPNImage     *vcl.TImage

	LocalNodeIdLabel *vcl.TLabel
	HeziNodeIdLabel  *vcl.TLabel
	VPNidLabel       *vcl.TLabel

	//sheet3 ????????????
	MembersScrollBox *vcl.TScrollBox

	Gconf GlobalConfig //????????????????????????????????????????????????????????????????????????
}

var mainForm *TMainForm

func (f *TMainForm) GetConfig() {
	Tgconf = InitKey()
	f.Gconf = Tgconf

	//????????????????????????????????????
	f.DefaultBox.SetCaption(f.Gconf.DefaultBox)
	f.DefaultBrain.SetCaption(f.Gconf.DefaultBrain)
	if f.Gconf.DefaultTemplate == "rclone.conf.3box" {
		f.Template.SetItemIndex(1)
	}
	if f.Gconf.DefaultTemplate == "rclone.conf.5box" {
		f.Template.SetItemIndex(2)
	}
	if f.Gconf.DefaultTemplate == "rclone.conf.11box" {
		f.Template.SetItemIndex(3)
	}
	f.ModifyBox.Clear()
	for kname, _ := range f.Gconf.AllBox {
		f.ModifyBox.Items().Add(kname)
	}
	f.ModifyBox.SetText(f.Gconf.DefaultBox)
	f.oldBox = f.Gconf.DefaultBox
	f.OnModifyBoxSelect(nil)

	f.DefaultInterval.SetText(strconv.Itoa(f.Gconf.DefaultInterval))
	f.DefaultLogFile.SetText(f.Gconf.DefaultLogFile)
	f.CbDefaultAutoBoot.SetChecked(f.Gconf.DefaultAutoBoot)
	f.OnBtnRefreshClick(nil)

	f.BtnApply.SetEnabled(false)
}
func (f *TMainForm) SetConfig() {
	//????????????????????????????????????

	var allBox AllBox
	allBox.Type = f.Type.Caption()
	allBox.Vendor = f.Vendor.Caption()
	allBox.Url = f.Url.Text()
	allBox.User = f.User.Text()
	allBox.Password = base64.StdEncoding.EncodeToString([]byte(f.Password.Text()))
	allBox.VirtualDisk = f.VirtualDisk.Text()
	f.Gconf.AllBox[f.GroupBox1.Caption()] = allBox

	if f.CbDefaultBox.Checked() {
		f.DefaultBox.SetCaption(f.GroupBox1.Caption())
		f.Gconf.DefaultBox = f.DefaultBox.Caption()
	}
	if f.CbDefaultBrain.Checked() {
		f.DefaultBrain.SetCaption(f.GroupBox1.Caption())
		f.Gconf.DefaultBrain = f.DefaultBrain.Caption()
	}
	u, _ := url.Parse(f.Gconf.AllBox[f.DefaultBox.Caption()].Url)
	f.Gconf.DefaultIpAddr = strings.Split(u.Host, ":")[0]
	f.Gconf.DefaultInterval, _ = strconv.Atoi(f.DefaultInterval.Text())
	f.Gconf.DefaultLogFile = f.DefaultLogFile.Text()
	logFile, err := os.OpenFile(f.Gconf.DefaultLogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("???????????????????????????%s:%v\n", Tgconf.DefaultLogFile, err)
	}
	Info = log.New(io.MultiWriter(logFile, os.Stderr), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	f.Gconf.DefaultAutoBoot = f.CbDefaultAutoBoot.Checked()

	key, _, _ := registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if f.CbDefaultAutoBoot.Checked() {
		fabs, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		key.SetStringValue("sca", fabs+"\\sca_client.exe")
	} else {
		key.SetStringValue("sca", "")
	}
	Tgconf = f.Gconf
	UpdateKey(Tgconf)
}

func (f *TMainForm) SaveIni() {
	//???????????????conf??????
	cfg := ini.Empty()
	var rconf string
	user, err := user.Current()
	if nil == err {
		rconf = user.HomeDir + "/.config/rclone/rclone.conf"
	} else {
		rconf = os.Getenv("USERPROFILE") + "/.config/rclone/rclone.conf"
	}

	for kname, allBox := range Tgconf.AllBox {
		_, err = cfg.NewSection(kname)
		_, err = cfg.Section(kname).NewKey("type", allBox.Type)
		_, err = cfg.Section(kname).NewKey("vendor", allBox.Vendor)
		_, err = cfg.Section(kname).NewKey("url", allBox.Url)
		_, err = cfg.Section(kname).NewKey("user", allBox.User)
		_, err = cfg.Section(kname).NewKey("pass", "")
	}
	err = cfg.SaveToIndent(rconf, "\t")
	for kname, allBox := range Tgconf.AllBox {
		decode, _ := base64.StdEncoding.DecodeString(allBox.Password)
		cmd := exec.Command("rclone", "config", "password", kname, "pass", string(decode))
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
	}

	//??????rclone??????????????????????????????box
	cmd := exec.Command("taskkill", "/im", "rclone.exe", "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	//MobaXterm.ini

	cfg, err = ini.Load("MobaXterm.ini")
	if err != nil {
		cfg = ini.Empty()
	}
	if cfg.Section("Bookmarks") == nil {
		_, err = cfg.NewSection("Bookmarks")
		_, err = cfg.Section("Bookmarks").NewKey("SubRep", "")
		_, err = cfg.Section("Bookmarks").NewKey("ImgNum", "42")
		_, err = cfg.Section("Bookmarks").NewKey("box1", "#109#0%192.168.216.197%22%%%-1%-1%%%22%%0%-1%0%%%-1%0%0%0%%1080%%0%0%1#MobaFont%10%0%0%0%15%236,236,236%30,30,30%180,180,192%0%-1%0%%xterm%-1%-1%_Std_Colors_0_%80%24%0%1%-1%<none>%%0#0# #-1")

		for kname, allBox := range Tgconf.AllBox {
			u, _ := url.Parse(allBox.Url)
			_, err = cfg.Section("Bookmarks").NewKey(kname, "#109#0%"+strings.Split(u.Host, ":")[0]+"%22%%%-1%-1%%%22%%0%-1%0%%%-1%0%0%0%%1080%%0%0%1#MobaFont%10%0%0%0%15%236,236,236%30,30,30%180,180,192%0%-1%0%%xterm%-1%-1%_Std_Colors_0_%80%24%0%1%-1%<none>%%0#0# #-1")
		}
	} else {
		cfg.Section("Bookmarks").Key("ImgNum").SetValue("42")
		for kname, allBox := range Tgconf.AllBox {
			u, _ := url.Parse(allBox.Url)
			if cfg.Section("Bookmarks").HasKey(kname) {
				cfg.Section("Bookmarks").Key(kname).SetValue("#109#0%" + strings.Split(u.Host, ":")[0] + "%22%%%-1%-1%%%22%%0%-1%0%%%-1%0%0%0%%1080%%0%0%1#MobaFont%10%0%0%0%15%236,236,236%30,30,30%180,180,192%0%-1%0%%xterm%-1%-1%_Std_Colors_0_%80%24%0%1%-1%<none>%%0#0# #-1")

			} else {
				_, err = cfg.Section("Bookmarks").NewKey(kname, "#109#0%"+strings.Split(u.Host, ":")[0]+"%22%%%-1%-1%%%22%%0%-1%0%%%-1%0%0%0%%1080%%0%0%1#MobaFont%10%0%0%0%15%236,236,236%30,30,30%180,180,192%0%-1%0%%xterm%-1%-1%_Std_Colors_0_%80%24%0%1%-1%<none>%%0#0# #-1")
			}
		}
	}

	err = cfg.SaveToIndent("MobaXterm.ini", "\t")

}

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("????????????????????????")
	f.SetPosition(types.PoScreenCenter)
	f.EnabledMaximize(false)
	f.SetWidth(640)
	f.SetHeight(480)
	// ??????????????????
	f.SetShowHint(true)

	// ????????????
	f.initComponents()

}

func (f *TMainForm) OnBtnApplyClick(sender vcl.IObject) {
	//vcl.Application.CreateForm(&Form1)
	//Form1.Show()
	f.SetCursor(types.CrHourGlass)
	f.BtnApply.SetCursor(types.CrHourGlass)
	f.SetConfig()
	f.SaveIni()
	f.BtnApply.SetEnabled(false)
	f.SetCursor(types.CrDefault)
	f.BtnApply.SetCursor(types.CrDefault)
}

func (f *TMainForm) OnBtnOkClick(sender vcl.IObject) {
	if !f.BtnApply.Enabled() {
		f.Close()
		return
	}
	f.SetCursor(types.CrHourGlass)
	f.BtnOk.SetCursor(types.CrHourGlass)
	vcl.Application.ProcessMessages()
	f.SetConfig()
	f.SaveIni()
	f.BtnApply.SetEnabled(false)
	f.SetCursor(types.CrDefault)
	f.BtnOk.SetCursor(types.CrDefault)
	f.Close()
}

func (f *TMainForm) OnBtnCancelClick(sender vcl.IObject) {
	f.BtnApply.SetEnabled(false)
	f.Close()
}

func (f *TMainForm) OnTemplateChange(sender vcl.IObject) {

	if vcl.MessageDlg("?????????????????????"+f.Template.Text()+"??????????", types.MtInformation, types.MbYes, types.MbNo) == types.MrYes {
		//
		f.Gconf.DefaultTemplate = "rclone.conf." + f.Template.Text() + "box"
		cfg, err := ini.Load(f.Gconf.DefaultTemplate)

		if err != nil {
			Info.Print("%v err->%v\n", "load template???", err)
			return
		}
		for k := range f.Gconf.AllBox {
			delete(f.Gconf.AllBox, k)
		}

		for _, v := range cfg.Sections() {
			if len(v.KeyStrings()) == 0 {
				continue
			}

			var allBox AllBox
			allBox.Type = v.Key("type").Value()
			allBox.Vendor = v.Key("vendor").Value()
			allBox.Url = v.Key("url").Value()
			allBox.User = v.Key("user").Value()
			allBox.Password = base64.StdEncoding.EncodeToString([]byte("password"))
			allBox.VirtualDisk = "Z:"
			f.Gconf.AllBox[v.Name()] = allBox
			/*fmt.Println(v.Name(),len(v.KeyStrings()),v.KeyStrings(),v.Key("type").Value())

			_, ok := f.Gconf.AllBox[v.Name()]
			if !ok {
				var allBox AllBox
				allBox.Type = v.Key("type").Value()
				allBox.Vendor = v.Key("vendor").Value()
				allBox.Url = v.Key("url").Value()
				allBox.User = v.Key("user").Value()
				allBox.Password = base64.StdEncoding.EncodeToString( []byte("password") )
				allBox.VirtualDisk = "Z:"
				f.Gconf.AllBox[v.Name()] = allBox
			}

			*/

		}
		f.ModifyBox.Clear()
		for kname, _ := range f.Gconf.AllBox {
			f.ModifyBox.Items().Add(kname)
		}
		_, ok := f.Gconf.AllBox[f.Gconf.DefaultBox]
		if !ok {
			f.Gconf.DefaultBox = f.ModifyBox.Items().Strings(0)
		}
		_, ok = f.Gconf.AllBox[f.Gconf.DefaultBrain]
		if !ok {
			f.Gconf.DefaultBrain = f.ModifyBox.Items().Strings(0)
		}
		f.ModifyBox.SetText(f.Gconf.DefaultBox)
		f.BtnApply.SetEnabled(true)
		f.OnModifyBoxSelect(sender)
	}
}

func (f *TMainForm) OnModifyBoxSelect(sender vcl.IObject) {
	//????????????box????????????????????????oldBox??????
	if f.oldBox != f.ModifyBox.Text() && f.ModifyBox.Items().IndexOf(f.oldBox) >= 0 {
		var allBox AllBox
		allBox.Type = f.Type.Caption()
		allBox.Vendor = f.Vendor.Caption()
		allBox.Url = f.Url.Text()
		allBox.User = f.User.Text()
		allBox.Password = base64.StdEncoding.EncodeToString([]byte(f.Password.Text()))
		allBox.VirtualDisk = f.VirtualDisk.Text()
		f.Gconf.AllBox[f.oldBox] = allBox
	}
	f.GroupBox1.SetCaption(f.ModifyBox.Text())
	f.Url.SetText(f.Gconf.AllBox[f.ModifyBox.Text()].Url)
	f.User.SetText(f.Gconf.AllBox[f.ModifyBox.Text()].User)
	decoded, _ := base64.StdEncoding.DecodeString(f.Gconf.AllBox[f.ModifyBox.Text()].Password)
	f.Password.SetText(string(decoded)) //????????????
	if f.ModifyBox.Text() == f.Gconf.DefaultBox {
		f.CbDefaultBox.SetChecked(true)
	} else {
		f.CbDefaultBox.SetChecked(false)
	}
	if f.ModifyBox.Text() == f.Gconf.DefaultBrain {
		f.CbDefaultBrain.SetChecked(true)
	} else {
		f.CbDefaultBrain.SetChecked(false)
	}
	f.VirtualDisk.SetItemIndex(f.VirtualDisk.Items().IndexOf(f.Gconf.AllBox[f.ModifyBox.Text()].VirtualDisk))

	f.oldBox = f.ModifyBox.Text()
}
func (f *TMainForm) OnModifyBoxKeyPress(sender vcl.IObject, key *types.Char) {
	if int(*key) == keys.VkReturn {
		if f.ModifyBox.Items().IndexOf(f.ModifyBox.Text()) < 0 {
			f.ModifyBox.Items().Add(f.ModifyBox.Text())
			f.ModifyBox.SetItemIndex(f.ModifyBox.Items().Count() - 1)
			var allBox AllBox
			allBox.Type = "root"
			allBox.Vendor = "other"
			allBox.Url = "http://192.168.216." + strconv.Itoa(197+int(f.ModifyBox.Items().Count())-1) + ":800"
			allBox.User = "root"
			allBox.Password = base64.StdEncoding.EncodeToString([]byte("password"))
			allBox.VirtualDisk = "Z:"
			f.Gconf.AllBox[f.ModifyBox.Text()] = allBox
			f.BtnApply.SetEnabled(true)

		}
		f.OnModifyBoxSelect(sender)
	}

}
func (f *TMainForm) OnBtnNewBoxClick(sender vcl.IObject) {
	if f.ModifyBox.Items().IndexOf("box"+strconv.Itoa(int(f.ModifyBox.Items().Count())+1)) < 0 {
		f.ModifyBox.Items().Add("box" + strconv.Itoa(int(f.ModifyBox.Items().Count())+1))
		f.ModifyBox.SetItemIndex(f.ModifyBox.Items().Count() - 1)
		var allBox AllBox
		allBox.Type = "root"
		allBox.Vendor = "other"
		allBox.Url = "http://192.168.216." + strconv.Itoa(197+int(f.ModifyBox.Items().Count())-1) + ":800"
		allBox.User = "root"
		allBox.Password = base64.StdEncoding.EncodeToString([]byte("password"))
		allBox.VirtualDisk = "Z:"
		f.Gconf.AllBox[f.ModifyBox.Text()] = allBox
		f.BtnApply.SetEnabled(true)

		f.OnModifyBoxSelect(sender)
	}

}
func (f *TMainForm) OnBtnDeleteBoxClick(sender vcl.IObject) {
	if f.ModifyBox.Items().Count() > 1 {
		delete(f.Gconf.AllBox, f.ModifyBox.Text())
		cidx := f.ModifyBox.ItemIndex()
		f.ModifyBox.Items().Delete(cidx)
		if cidx > 0 {
			cidx = cidx - 1
		}
		f.ModifyBox.SetItemIndex(cidx)

		if f.ModifyBox.Items().IndexOf(f.DefaultBox.Caption()) < 0 {
			f.DefaultBox.SetCaption(f.ModifyBox.Text())
			f.Gconf.DefaultBox = f.DefaultBox.Caption()
		}
		if f.ModifyBox.Items().IndexOf(f.DefaultBrain.Caption()) < 0 {
			f.DefaultBrain.SetCaption(f.ModifyBox.Text())
			f.Gconf.DefaultBrain = f.DefaultBrain.Caption()
		}
		f.BtnApply.SetEnabled(true)
		f.OnModifyBoxSelect(sender)

	}

}

func (f *TMainForm) OnBtnTestMountClick(sender vcl.IObject) {
	cmd := exec.Command("taskkill", "/im", "rclone.exe", "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
	//open.Run("cmd /b rclone mount "+f.GroupBox1.Caption()+":/ " +f.VirtualDisk.Text()+ " --cache-dir .\\"+f.GroupBox1.Caption()+" --vfs-cache-mode writes")

	//???????????????conf??????

	var rconf string
	user, err := user.Current()
	if nil == err {
		rconf = user.HomeDir + "/.config/rclone/rclone.conf"
	} else {
		rconf = os.Getenv("USERPROFILE") + "/.config/rclone/rclone.conf"
	}
	cfg, err := ini.Load(rconf)
	if err != nil {
		cfg = ini.Empty()
	}
	if cfg.Section(f.GroupBox1.Caption()) == nil {
		_, err = cfg.NewSection(f.GroupBox1.Caption())
		_, err = cfg.Section(f.GroupBox1.Caption()).NewKey("type", f.Type.Caption())
		_, err = cfg.Section(f.GroupBox1.Caption()).NewKey("url", f.Url.Text())
		_, err = cfg.Section(f.GroupBox1.Caption()).NewKey("vendor", f.Vendor.Caption())
		_, err = cfg.Section(f.GroupBox1.Caption()).NewKey("user", f.User.Text())
		_, err = cfg.Section(f.GroupBox1.Caption()).NewKey("pass", "")
	} else {
		cfg.Section(f.GroupBox1.Caption()).Key("type").SetValue(f.Type.Caption())
		cfg.Section(f.GroupBox1.Caption()).Key("url").SetValue(f.Url.Text())
		cfg.Section(f.GroupBox1.Caption()).Key("vendor").SetValue(f.Vendor.Caption())
		cfg.Section(f.GroupBox1.Caption()).Key("user").SetValue(f.User.Text())
		cfg.Section(f.GroupBox1.Caption()).Key("pass").SetValue("")
	}

	err = cfg.SaveToIndent(rconf, "\t")
	cmd = exec.Command("rclone", "config", "password", f.GroupBox1.Caption(), "pass", f.Password.Text())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
	go func() {
		// --attr-timeout and --dir-cache-time
		cmd = exec.Command("cmd.exe", "/c", "start", "rclone", "mount", f.GroupBox1.Caption()+":/", f.VirtualDisk.Text(), "--cache-dir", ".\\"+f.GroupBox1.Caption(), "--vfs-cache-mode", "writes", "--attr-timeout", strconv.Itoa(Tgconf.DefaultInterval)+"s", "--dir-cache-time", strconv.Itoa(Tgconf.DefaultInterval)+"s")
		Info.Println(cmd.Args)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		cmd.Run()

	}()
}
func (f *TMainForm) OnBtnViewLogClick(sender vcl.IObject) {

	open.Run(f.DefaultLogFile.Text())
}

//TODO
func (f *TMainForm) OnBtnChangePasswordClick(sender vcl.IObject) {

}
func (f *TMainForm) OnBtnRefreshClick(sender vcl.IObject) {
	//
	f.SetCursor(types.CrHourGlass)
	if UrlStatus("800") {
		isr := IsExeRuning("rclone.exe", "rclone.exe")
		if !isr {
			MountBox()
		}
		f.BtnRefresh.SetCaption("??????(&S)")
		f.LbRefresh.SetCaption("????????????(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ")")
		f.LbRefresh.SetCursor(types.CrHourGlass)
		f.BtnRefresh.SetCursor(types.CrHourGlass)
		go func() {
			cmd := exec.Command("rclone", "size", Tgconf.DefaultBox+":/")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				//fmt.Println(cmd.ProcessState.Pid())
				sr := strings.NewReader(out.String())
				reader := bufio.NewReader(sr)
				line1, _, _ := reader.ReadLine()
				line2, _, _ := reader.ReadLine()
				vcl.ThreadSync(func() { //??????????????????ui
					f.LbRefresh.SetCaption("????????????(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ");\n" + "???????????????" + strings.Split(string(line1), ":")[1] + ";\t???????????????" + strings.Split(string(line2), ":")[1] + ")")
					f.LbRefresh.SetCursor(types.CrDefault)
					f.BtnRefresh.SetCursor(types.CrDefault)
				})
			}
		}()
	} else {
		f.BtnRefresh.SetCaption("??????(&C)")
		f.LbRefresh.SetCaption("????????????(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ")")
		UnmountBox()
	}

	f.SetCursor(types.CrDefault)
}

func (f *TMainForm) OnCbDefaultBoxChange(sender vcl.IObject) {
	if f.CbDefaultBox.Checked() {
		f.DefaultBox.SetCaption(f.GroupBox1.Caption())
		f.Gconf.DefaultBox = f.GroupBox1.Caption()
		f.BtnApply.SetEnabled(true)
	}
}
func (f *TMainForm) OnCbDefaultBrainChange(sender vcl.IObject) {
	if f.CbDefaultBrain.Checked() {
		f.DefaultBrain.SetCaption(f.GroupBox1.Caption())
		f.Gconf.DefaultBrain = f.GroupBox1.Caption()
		f.BtnApply.SetEnabled(true)
	} else {
		f.DefaultBrain.SetCaption("")
		f.Gconf.DefaultBrain = ""
		f.BtnApply.SetEnabled(true)
	}

}

func (f *TMainForm) OnUrlChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}
func (f *TMainForm) OnUserChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}
func (f *TMainForm) OnPasswordChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}
func (f *TMainForm) OnVirtualDiskChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}
func (f *TMainForm) OnDefaultIntervalChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}
func (f *TMainForm) OnDefaultLogFileChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}

func (f *TMainForm) OnCbDefaultAutoBootChange(sender vcl.IObject) {

	f.BtnApply.SetEnabled(true)
}

var lv1 *vcl.TListView

func (f *TMainForm) initComponents() {
	f.ImgList = vcl.NewImageList(f)

	if vcl.Application.Icon().Handle() != 0 {
		f.ImgList.AddIcon(vcl.Application.Icon())
	}

	f.ActList = vcl.NewActionList(f)
	f.ActList.SetImages(f.ImgList)

	f.BtnApply = vcl.NewButton(f)
	f.BtnApply.SetParent(f)
	f.BtnApply.SetLeft(152)
	f.BtnApply.SetTop(408)
	f.BtnApply.SetHeight(47)
	f.BtnApply.SetWidth(134)
	f.BtnApply.SetCaption("??????(&A)")
	var ftf *vcl.TFont = vcl.AsFont(f.BtnApply.Font())
	ftf.SetSize(11)
	ftf.SetHeight(-18)
	f.BtnApply.SetFont(ftf)
	f.BtnApply.SetOnClick(f.OnBtnApplyClick)

	f.BtnOk = vcl.NewButton(f)
	f.BtnOk.SetParent(f)
	f.BtnOk.SetLeft(312)
	f.BtnOk.SetTop(408)
	f.BtnOk.SetHeight(47)
	f.BtnOk.SetWidth(134)
	f.BtnOk.SetCaption("??????(&O)")
	f.BtnOk.SetFont(ftf)
	f.BtnOk.SetOnClick(f.OnBtnOkClick)

	f.BtnCancel = vcl.NewButton(f)
	f.BtnCancel.SetParent(f)
	f.BtnCancel.SetLeft(472)
	f.BtnCancel.SetTop(408)
	f.BtnCancel.SetHeight(47)
	f.BtnCancel.SetWidth(134)
	f.BtnCancel.SetCaption("??????(&C)")
	f.BtnCancel.SetFont(ftf)
	f.BtnCancel.SetOnClick(f.OnBtnCancelClick)

	f.Pgc = vcl.NewPageControl(f)
	f.Pgc.SetParent(f)
	f.Pgc.SetLeft(1)
	f.Pgc.SetTop(0)
	f.Pgc.SetHeight(393)
	f.Pgc.SetWidth(638)
	//f.Pgc.SetAlign(types.AlClient)

	f.Sheet1 = vcl.NewTabSheet(f)
	f.Sheet1.SetPageControl(f.Pgc)
	f.Sheet1.SetCaption("????????????????????????")

	f.Sheet2 = vcl.NewTabSheet(f)
	f.Sheet2.SetPageControl(f.Pgc)
	f.Sheet2.SetCaption("????????????????????????")

	f.Sheet3 = vcl.NewTabSheet(f)
	f.Sheet3.SetPageControl(f.Pgc)
	f.Sheet3.SetCaption("??????VPN????????????")

	f.Lb1 = vcl.NewLabel(f)
	f.Lb1.SetParent(f.Sheet1)
	f.Lb1.SetCaption("??????BOX???")
	f.Lb1.SetLeft(16)
	f.Lb1.SetTop(8)
	f.Lb1.SetHeight(20)
	f.Lb1.SetWidth(76)
	f.DefaultBox = vcl.NewLabel(f)
	f.DefaultBox.SetParent(f.Sheet1)
	f.DefaultBox.SetCaption("box1")
	f.DefaultBox.SetLeft(96)
	f.DefaultBox.SetTop(8)
	f.DefaultBox.SetHeight(20)
	f.DefaultBox.SetWidth(76)

	f.Lb2 = vcl.NewLabel(f)
	f.Lb2.SetParent(f.Sheet1)
	f.Lb2.SetCaption("??????BOX???")
	f.Lb2.SetLeft(192)
	f.Lb2.SetTop(8)
	f.Lb2.SetHeight(20)
	f.Lb2.SetWidth(76)
	f.DefaultBrain = vcl.NewLabel(f)
	f.DefaultBrain.SetParent(f.Sheet1)
	f.DefaultBrain.SetCaption("box1")
	f.DefaultBrain.SetLeft(272)
	f.DefaultBrain.SetTop(8)
	f.DefaultBrain.SetHeight(20)
	f.DefaultBrain.SetWidth(76)

	f.Lb3 = vcl.NewLabel(f)
	f.Lb3.SetParent(f.Sheet1)
	f.Lb3.SetCaption("?????????????????????")
	f.Lb3.SetLeft(372)
	f.Lb3.SetTop(8)
	f.Lb3.SetHeight(20)
	f.Lb3.SetWidth(105)
	f.Template = vcl.NewComboBox(f)
	f.Template.SetParent(f.Sheet1)
	f.Template.Items().Add("1")
	f.Template.Items().Add("3")
	f.Template.Items().Add("5")
	f.Template.Items().Add("11")
	f.Template.SetItemIndex(0)
	f.Template.SetStyle(types.CsDropDownList)
	f.Template.SetLeft(488)
	f.Template.SetTop(8)
	f.Template.SetItemHeight(20)
	f.Template.SetWidth(112)
	f.Template.SetOnChange(f.OnTemplateChange)

	f.Lb4 = vcl.NewLabel(f)
	f.Lb4.SetParent(f.Sheet1)
	f.Lb4.SetCaption("?????????????????????")
	f.Lb4.SetLeft(16)
	f.Lb4.SetTop(40)
	f.Lb4.SetHeight(20)
	f.Lb4.SetWidth(105)
	f.ModifyBox = vcl.NewComboBox(f)
	f.ModifyBox.SetParent(f.Sheet1)
	f.ModifyBox.Items().Add("box1")
	f.ModifyBox.SetItemIndex(0)
	f.ModifyBox.SetStyle(types.CsDropDown)
	f.ModifyBox.SetLeft(116)
	f.ModifyBox.SetTop(40)
	f.ModifyBox.SetItemHeight(20)
	f.ModifyBox.SetWidth(215)
	f.ModifyBox.SetOnKeyPress(f.OnModifyBoxKeyPress)
	f.ModifyBox.SetOnSelect(f.OnModifyBoxSelect)
	f.BtnNewBox = vcl.NewButton(f)
	f.BtnNewBox.SetParent(f.Sheet1)
	f.BtnNewBox.SetLeft(372)
	f.BtnNewBox.SetTop(40)
	f.BtnNewBox.SetHeight(30)
	f.BtnNewBox.SetWidth(112)
	f.BtnNewBox.SetCaption("?????????BOX(&N)")
	f.BtnNewBox.SetOnClick(f.OnBtnNewBoxClick)
	f.BtnDeleteBox = vcl.NewButton(f)
	f.BtnDeleteBox.SetParent(f.Sheet1)
	f.BtnDeleteBox.SetLeft(488)
	f.BtnDeleteBox.SetTop(40)
	f.BtnDeleteBox.SetHeight(30)
	f.BtnDeleteBox.SetWidth(112)
	f.BtnDeleteBox.SetCaption("??????BOX(&D)")
	f.BtnDeleteBox.SetOnClick(f.OnBtnDeleteBoxClick)

	f.GroupBox1 = vcl.NewGroupBox(f)
	f.GroupBox1.SetParent(f.Sheet1)
	f.GroupBox1.SetLeft(8)
	f.GroupBox1.SetTop(80)
	f.GroupBox1.SetHeight(144)
	f.GroupBox1.SetWidth(615)
	f.GroupBox1.SetCaption(f.ModifyBox.Text())
	f.Lb5 = vcl.NewLabel(f)
	f.Lb5.SetParent(f.GroupBox1)
	f.Lb5.SetLeft(28)
	f.Lb5.SetTop(10)
	f.Lb5.SetHeight(20)
	f.Lb5.SetWidth(75)
	f.Lb5.SetCaption("???????????????")
	f.Type = vcl.NewLabel(f)
	f.Type.SetParent(f.GroupBox1)
	f.Type.SetLeft(104)
	f.Type.SetTop(10)
	f.Type.SetHeight(20)
	f.Type.SetWidth(57)
	f.Type.SetCaption("webdav")
	f.Lb6 = vcl.NewLabel(f)
	f.Lb6.SetParent(f.GroupBox1)
	f.Lb6.SetLeft(176)
	f.Lb6.SetTop(10)
	f.Lb6.SetHeight(20)
	f.Lb6.SetWidth(45)
	f.Lb6.SetCaption("?????????")
	f.Vendor = vcl.NewLabel(f)
	f.Vendor.SetParent(f.GroupBox1)
	f.Vendor.SetLeft(224)
	f.Vendor.SetTop(10)
	f.Vendor.SetHeight(20)
	f.Vendor.SetWidth(40)
	f.Vendor.SetCaption("other")
	f.Lb7 = vcl.NewLabel(f)
	f.Lb7.SetParent(f.GroupBox1)
	f.Lb7.SetLeft(280)
	f.Lb7.SetTop(10)
	f.Lb7.SetHeight(20)
	f.Lb7.SetWidth(75)
	f.Lb7.SetCaption("???????????????")
	f.Url = vcl.NewEdit(f)
	f.Url.SetParent(f.GroupBox1)
	f.Url.SetLeft(360)
	f.Url.SetTop(10)
	f.Url.SetHeight(28)
	f.Url.SetWidth(244)
	f.Url.SetText("http://192.168.216.197:800")
	f.Url.SetOnChange(f.OnUrlChange)
	f.Lb8 = vcl.NewLabel(f)
	f.Lb8.SetParent(f.GroupBox1)
	f.Lb8.SetLeft(28)
	f.Lb8.SetTop(40)
	f.Lb8.SetHeight(20)
	f.Lb8.SetWidth(45)
	f.Lb8.SetCaption("?????????")
	f.User = vcl.NewEdit(f)
	f.User.SetParent(f.GroupBox1)
	f.User.SetLeft(104)
	f.User.SetTop(40)
	f.User.SetHeight(28)
	f.User.SetWidth(100)
	f.User.SetText("root")
	f.User.SetOnChange(f.OnUserChange)
	f.Lb9 = vcl.NewLabel(f)
	f.Lb9.SetParent(f.GroupBox1)
	f.Lb9.SetLeft(280)
	f.Lb9.SetTop(40)
	f.Lb9.SetHeight(20)
	f.Lb9.SetWidth(45)
	f.Lb9.SetCaption("?????????")
	f.Password = vcl.NewEdit(f)
	f.Password.SetParent(f.GroupBox1)
	f.Password.SetLeft(360)
	f.Password.SetTop(40)
	f.Password.SetHeight(28)
	f.Password.SetWidth(112)
	f.Password.SetText("password")
	f.Password.SetPasswordChar('#')
	f.Password.SetOnChange(f.OnPasswordChange)
	f.CbDefaultBox = vcl.NewCheckBox(f)
	f.CbDefaultBox.SetParent(f.GroupBox1)
	f.CbDefaultBox.SetCaption("??????box")
	f.CbDefaultBox.SetLeft(28)
	f.CbDefaultBox.SetTop(80)
	f.CbDefaultBox.SetHeight(24)
	f.CbDefaultBox.SetWidth(82)
	f.CbDefaultBox.SetChecked(true)
	f.CbDefaultBox.SetOnChange(f.OnCbDefaultBoxChange)
	f.CbDefaultBrain = vcl.NewCheckBox(f)
	f.CbDefaultBrain.SetParent(f.GroupBox1)
	f.CbDefaultBrain.SetCaption("??????box")
	f.CbDefaultBrain.SetLeft(182)
	f.CbDefaultBrain.SetTop(80)
	f.CbDefaultBrain.SetHeight(24)
	f.CbDefaultBrain.SetWidth(82)
	f.CbDefaultBrain.SetChecked(true)
	f.CbDefaultBrain.SetOnChange(f.OnCbDefaultBrainChange)
	f.VirtualDisk = vcl.NewComboBox(f)
	f.VirtualDisk.SetParent(f.GroupBox1)
	f.VirtualDisk.Items().Add("Z:")
	f.VirtualDisk.Items().Add("Y:")
	f.VirtualDisk.Items().Add("X:")
	f.VirtualDisk.Items().Add("W:")
	f.VirtualDisk.Items().Add("V:")
	f.VirtualDisk.Items().Add("U:")
	f.VirtualDisk.Items().Add("T:")
	f.VirtualDisk.Items().Add("S:")
	f.VirtualDisk.Items().Add("R:")
	f.VirtualDisk.Items().Add("Q:")
	f.VirtualDisk.Items().Add("P:")
	f.VirtualDisk.Items().Add("O:")
	f.VirtualDisk.Items().Add("N:")
	f.VirtualDisk.Items().Add("M:")
	f.VirtualDisk.Items().Add("L:")
	f.VirtualDisk.Items().Add("K:")
	f.VirtualDisk.Items().Add("J:")
	f.VirtualDisk.Items().Add("I:")
	f.VirtualDisk.SetItemIndex(0)
	f.VirtualDisk.SetStyle(types.CsDropDownList)
	f.VirtualDisk.SetLeft(280)
	f.VirtualDisk.SetTop(80)
	f.VirtualDisk.SetItemHeight(20)
	f.VirtualDisk.SetWidth(45)
	f.VirtualDisk.SetOnChange(f.OnVirtualDiskChange)
	f.BtnTestMount = vcl.NewButton(f)
	f.BtnTestMount.SetParent(f.GroupBox1)
	f.BtnTestMount.SetLeft(360)
	f.BtnTestMount.SetTop(80)
	f.BtnTestMount.SetHeight(30)
	f.BtnTestMount.SetWidth(112)
	f.BtnTestMount.SetCaption("????????????(&M)")
	f.BtnTestMount.SetOnClick(f.OnBtnTestMountClick)

	f.Lb10 = vcl.NewLabel(f)
	f.Lb10.SetParent(f.Sheet1)
	f.Lb10.SetLeft(16)
	f.Lb10.SetTop(248)
	f.Lb10.SetHeight(20)
	f.Lb10.SetWidth(105)
	f.Lb10.SetCaption("?????????????????????")
	f.DefaultInterval = vcl.NewEdit(f)
	f.DefaultInterval.SetParent(f.Sheet1)
	f.DefaultInterval.SetLeft(116)
	f.DefaultInterval.SetTop(248)
	f.DefaultInterval.SetHeight(28)
	f.DefaultInterval.SetWidth(105)
	f.DefaultInterval.SetText("30")
	f.DefaultInterval.SetNumbersOnly(true)
	f.DefaultInterval.SetOnChange(f.OnDefaultIntervalChange)
	f.Lb11 = vcl.NewLabel(f)
	f.Lb11.SetParent(f.Sheet1)
	f.Lb11.SetLeft(232)
	f.Lb11.SetTop(248)
	f.Lb11.SetHeight(20)
	f.Lb11.SetWidth(15)
	f.Lb11.SetCaption("???")
	f.AllPassword = vcl.NewEdit(f)
	f.AllPassword.SetParent(f.Sheet1)
	f.AllPassword.SetLeft(372)
	f.AllPassword.SetTop(248)
	f.AllPassword.SetHeight(28)
	f.AllPassword.SetWidth(112)
	f.AllPassword.SetText("password")
	f.AllPassword.SetPasswordChar('*')
	f.BtnChangePassword = vcl.NewButton(f)
	f.BtnChangePassword.SetParent(f.Sheet1)
	f.BtnChangePassword.SetLeft(488)
	f.BtnChangePassword.SetTop(248)
	f.BtnChangePassword.SetHeight(30)
	f.BtnChangePassword.SetWidth(112)
	f.BtnChangePassword.SetCaption("???????????????(&P)")
	f.BtnChangePassword.SetOnClick(f.OnBtnChangePasswordClick)

	f.Lb12 = vcl.NewLabel(f)
	f.Lb12.SetParent(f.Sheet1)
	f.Lb12.SetLeft(16)
	f.Lb12.SetTop(288)
	f.Lb12.SetHeight(20)
	f.Lb12.SetWidth(75)
	f.Lb12.SetCaption("???????????????")
	f.DefaultLogFile = vcl.NewEdit(f)
	f.DefaultLogFile.SetParent(f.Sheet1)
	f.DefaultLogFile.SetLeft(116)
	f.DefaultLogFile.SetTop(288)
	f.DefaultLogFile.SetHeight(28)
	f.DefaultLogFile.SetWidth(105)
	f.DefaultLogFile.SetText("sca_client.log")
	f.DefaultLogFile.SetOnChange(f.OnDefaultLogFileChange)
	f.BtnViewLog = vcl.NewButton(f)
	f.BtnViewLog.SetParent(f.Sheet1)
	f.BtnViewLog.SetLeft(232)
	f.BtnViewLog.SetTop(288)
	f.BtnViewLog.SetHeight(30)
	f.BtnViewLog.SetWidth(112)
	f.BtnViewLog.SetCaption("????????????(&R)")
	f.BtnViewLog.SetOnClick(f.OnBtnViewLogClick)
	f.CbDefaultAutoBoot = vcl.NewCheckBox(f)
	f.CbDefaultAutoBoot.SetParent(f.Sheet1)
	f.CbDefaultAutoBoot.SetCaption("?????????????????????")
	f.CbDefaultAutoBoot.SetLeft(496)
	f.CbDefaultAutoBoot.SetTop(288)
	f.CbDefaultAutoBoot.SetHeight(24)
	f.CbDefaultAutoBoot.SetWidth(129)
	f.CbDefaultAutoBoot.SetChecked(true)
	f.CbDefaultAutoBoot.SetOnChange(f.OnCbDefaultAutoBootChange)

	f.BtnRefresh = vcl.NewButton(f)
	f.BtnRefresh.SetParent(f.Sheet1)
	f.BtnRefresh.SetLeft(16)
	f.BtnRefresh.SetTop(328)
	f.BtnRefresh.SetHeight(30)
	f.BtnRefresh.SetWidth(90)
	f.BtnRefresh.SetCaption("??????(&C)")
	f.BtnRefresh.SetOnClick(f.OnBtnRefreshClick)
	f.LbRefresh = vcl.NewLabel(f)
	f.LbRefresh.SetParent(f.Sheet1)
	f.LbRefresh.SetLeft(116)
	f.LbRefresh.SetTop(328)
	f.LbRefresh.SetHeight(40)
	f.LbRefresh.SetWidth(405)
	f.LbRefresh.SetCaption("????????????(box1:/Z:) ??????????????? ???????????????")

	//sheet3 ???????????????
	f.VPNConfigGroupBox = vcl.NewGroupBox(f)
	f.VPNConfigGroupBox.SetParent(f.Sheet3)
	f.VPNConfigGroupBox.SetBounds(32, 16, 560, 228)
	f.VPNConfigGroupBox.SetCaption("????????????")

	f.JoinVPNLabel = vcl.NewLabel(f)
	f.JoinVPNLabel.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNLabel.SetBounds(40, 0, 72, 24)

	f.JoinVPNLabel.SetCaption("????????????:")

	f.JoinVPNEdit = vcl.NewComboBox(f)
	f.JoinVPNEdit.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNEdit.SetBounds(120, 0, 160, 33)

	f.JoinVPNEdit.SetOnChange(f.onJoinVPNEditChange)

	f.JoinVPNButton = vcl.NewButton(f)
	f.JoinVPNButton.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNButton.SetBounds(312, 0, 96, 33)
	f.JoinVPNButton.SetCaption("????????????")
	f.JoinVPNButton.SetOnClick(f.onJoinVPNButtonClick)

	f.DeleteVPN = vcl.NewButton(f)
	f.DeleteVPN.SetParent(f.VPNConfigGroupBox)
	f.DeleteVPN.SetBounds(424, 0, 96, 33)
	f.DeleteVPN.SetCaption("????????????")
	f.DeleteVPN.SetOnClick(f.onDeleteVPNButtonClick)

	f.AutoJoinVPN = vcl.NewCheckBox(f)
	f.AutoJoinVPN.SetParent(f.VPNConfigGroupBox)
	f.AutoJoinVPN.SetBounds(456, 120, 146, 27)
	f.AutoJoinVPN.SetCaption("???????????????")
	f.AutoJoinVPN.SetOnClick(f.onAutoJoinVPNClick)

	f.VPNStatusLabel = vcl.NewLabel(f)
	f.VPNStatusLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNStatusLabel.SetBounds(40, 120, 80, 24)
	f.VPNStatusLabel.SetCaption("???????????????")

	f.VPNStatusShape = vcl.NewShape(f)
	f.VPNStatusShape.SetParent(f.VPNConfigGroupBox)
	f.VPNStatusShape.SetShape(types.StCircle)
	f.VPNStatusShape.SetBounds(120, 120, 24, 20)

	var brush = vcl.NewBrush()
	var pen = vcl.NewPen()

	f.VPNStatus = vcl.NewLabel(f)
	f.VPNStatus.SetParent(f.VPNConfigGroupBox)
	f.VPNStatus.SetBounds(152, 120, 45, 24)

	f.VPNIPLabel = vcl.NewLabel(f)
	f.VPNIPLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNIPLabel.SetBounds(232, 120, 72, 24)
	f.VPNIPLabel.SetCaption("??????IP???")

	f.VPNIPEdit = vcl.NewEdit(f)
	f.VPNIPEdit.SetParent(f.VPNConfigGroupBox)
	f.VPNIPEdit.SetBounds(288, 120, 160, 33)
	f.VPNIPEdit.SetEnabled(false)

	f.VPNUsernameLabel = vcl.NewLabel(f)
	f.VPNUsernameLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNUsernameLabel.SetBounds(40, 40, 72, 24)
	f.VPNUsernameLabel.SetCaption("????????????")

	f.VPNUsernameEdit = vcl.NewEdit(f)
	f.VPNUsernameEdit.SetParent(f.VPNConfigGroupBox)
	f.VPNUsernameEdit.SetBounds(120, 40, 160, 33)
	var hostname, _ = os.Hostname()
	fmt.Println("hostname:", hostname)
	f.VPNUsernameEdit.SetText(hostname)

	f.DeviceTypeLabel = vcl.NewLabel(f)
	f.DeviceTypeLabel.SetParent(f.VPNConfigGroupBox)
	f.DeviceTypeLabel.SetBounds(40, 79, 72, 24)
	f.DeviceTypeLabel.SetCaption("????????????:")

	f.DeviceTypeBox = vcl.NewComboBox(f)
	f.DeviceTypeBox.SetParent(f.VPNConfigGroupBox)
	f.DeviceTypeBox.SetBounds(120, 79, 160, 33)
	f.DeviceTypeBox.Items().Add("???????????????")
	f.DeviceTypeBox.Items().Add("?????????")
	f.DeviceTypeBox.Items().Add("?????????")
	f.DeviceTypeBox.SetItemIndex(0)

	f.DeviceTypeBox.SetStyle(types.CsDropDownList)

	f.ChangeNameButton = vcl.NewButton(f)
	f.ChangeNameButton.SetParent(f.VPNConfigGroupBox)
	f.ChangeNameButton.SetBounds(312, 56, 96, 33)
	f.ChangeNameButton.SetCaption("????????????")
	f.ChangeNameButton.SetOnClick(f.onChangeNameButtonClick)

	f.VPNTest = vcl.NewButton(f)
	f.VPNTest.SetParent(f.VPNConfigGroupBox)
	f.VPNTest.SetBounds(424, 56, 96, 33)
	if Tgconf.VPN.Connect {
		f.VPNTest.SetCaption("????????????")
	} else {
		f.VPNTest.SetCaption("????????????")
	}

	f.VPNTest.SetOnClick(f.onVPNTestClick)

	lv1 = vcl.NewListView(f)
	lv1.SetParent(f.Sheet3)
	lv1.SetBounds(16, 248, 600, 110)
	lv1.SetRowSelect(true)
	lv1.SetReadOnly(true)
	lv1.SetViewStyle(types.VsReport)
	lv1.SetGridLines(true)
	//lv1.SetColumnClick(false)
	lv1.SetHideSelection(false)

	col := lv1.Columns().Add()
	col.SetCaption("??????ID")
	col.SetWidth(100)

	col = lv1.Columns().Add()
	col.SetCaption("?????????")
	col.SetWidth(90)

	col = lv1.Columns().Add()
	col.SetCaption("????????????")
	col.SetWidth(75)

	col = lv1.Columns().Add()
	col.SetCaption("??????IP")
	col.SetWidth(110)

	col = lv1.Columns().Add()
	col.SetCaption("????????????")
	col.SetWidth(80)

	col = lv1.Columns().Add()
	col.SetCaption("??????IP")
	col.SetWidth(110)

	//??????????????????
	fillNetworks(f)

	f.LocalDevice = vcl.NewImage(f)
	f.LocalDevice.SetParent(f.VPNConfigGroupBox)
	f.LocalDevice.SetBounds(55, 151, 50, 50)
	pictue := vcl.NewPicture()
	pictue.LoadFromFile("./images/local.jpg")
	f.LocalDevice.SetPicture(pictue)

	f.LocalAndHezi = vcl.NewImage(f)
	f.LocalAndHezi.SetParent(f.VPNConfigGroupBox)
	f.LocalAndHezi.SetBounds(160, 151, 50, 50)
	pictuer := vcl.NewPicture()
	pictuer.LoadFromFile("./images/unlinked.jpg")
	f.LocalAndHezi.SetPicture(pictuer)

	f.Hezi = vcl.NewImage(f)
	f.Hezi.SetParent(f.VPNConfigGroupBox)
	f.Hezi.SetBounds(264, 154, 50, 50)
	pictuer1 := vcl.NewPicture()
	pictuer1.LoadFromFile("./images/hezi.jpg")
	f.Hezi.SetPicture(pictuer1)

	f.HeziandVPN = vcl.NewImage(f)
	f.HeziandVPN.SetParent(f.VPNConfigGroupBox)
	f.HeziandVPN.SetBounds(365, 151, 50, 50)
	pictuer2 := vcl.NewPicture()
	pictuer2.LoadFromFile("./images/unlinked.jpg")
	f.HeziandVPN.SetPicture(pictuer2)
	f.HeziandVPN.SetOnClick(f.onHeziandVPNClick)

	f.VPNImage = vcl.NewImage(f)
	f.VPNImage.SetParent(f.VPNConfigGroupBox)
	f.VPNImage.SetBounds(456, 151, 50, 50)
	pictuer3 := vcl.NewPicture()
	pictuer3.LoadFromFile("./images/net.jpg")
	f.VPNImage.SetPicture(pictuer3)

	//???????????? ?????? vpn???Nodeid
	f.LocalNodeIdLabel = vcl.NewLabel(f)
	f.LocalNodeIdLabel.SetParent(f.VPNConfigGroupBox)
	f.LocalNodeIdLabel.SetBounds(40, 190, 84, 15)
	nodeId := getNodeId()
	if nodeId != "" {
		f.LocalNodeIdLabel.SetCaption(nodeId)
	} else {
		f.LocalNodeIdLabel.SetCaption("????????????ID")
	}

	f.HeziNodeIdLabel = vcl.NewLabel(f)
	f.HeziNodeIdLabel.SetParent(f.VPNConfigGroupBox)
	f.HeziNodeIdLabel.SetBounds(256, 190, 84, 15)

	f.VPNidLabel = vcl.NewLabel(f)
	f.VPNidLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNidLabel.SetBounds(420, 190, 132, 15)

	if Tgconf.VPN.NetWorkId != "" {
		f.VPNidLabel.SetCaption(Tgconf.VPN.NetWorkId)
	} else {
		f.VPNidLabel.SetCaption("????????????ID")
	}

	vpnstatus, statusresponse := getVPNConnectStatus()
	if vpnstatus == 1 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClLime)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("?????????")
		f.VPNIPEdit.SetText(statusresponse.AssignedAddresses[0])
	} else if vpnstatus == 0 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClYellow)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("?????????")
	} else if vpnstatus == -2 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClRed)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("????????????")
	} else if vpnstatus == 2 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClGrey)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("?????????")
	} else if vpnstatus == -1 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClBlack)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("???????????????")
	}
	if Tgconf.VPN.NetWorkId != "" {
		showAllMembers(Tgconf.VPN.NetWorkId, lv1)
	}

}

func fillNetworks(f *TMainForm) {
	//???????????????????????????+
	fmt.Println("??????fillNetworks")
	vcl.ThreadSync(func() {
		f.JoinVPNEdit.Clear()
	})
	//1.????????????????????????
	// ?????????????????????

	dir := `./networks.d`
	fileinfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("??????????????????", err)
	}

	fmt.Println("?????????????????????")
	// ?????????????????????
	var allNetworks datastruct.AllNetworks
	for _, fi := range fileinfo {
		// ?????????????????????
		data, _ := ioutil.ReadFile("./networks.d/" + fi.Name())
		var networkconfig datastruct.Network
		if err := yaml.Unmarshal(data, &networkconfig); err != nil {
			log.Fatalf("unmarshal 1error: %v", err)
		}

		allNetworks = append(allNetworks, networkconfig)
	}
	fmt.Println("?????????????????????")
	//2.????????????
	var index int
	vcl.ThreadSync(func() {
		for i := 0; i < len(allNetworks); i++ {
			if Tgconf.VPN.NetWorkId != "" && allNetworks[i].NetworkID == Tgconf.VPN.NetWorkId && Tgconf.VPN.Connect {
				index = i
				f.JoinVPNEdit.Items().Add(allNetworks[i].NetworkID + "  ???")
			} else {
				f.JoinVPNEdit.Items().Add(allNetworks[i].NetworkID)
			}
		}
		f.JoinVPNEdit.SetItemIndex(int32(index))
	})

	fmt.Println("?????????????????????????????????")

}

//JoinVPNButton????????????  ????????????
func (f *TMainForm) onJoinVPNButtonClick(sender vcl.IObject) {
	vpnID := f.JoinVPNEdit.Text()
	fmt.Println(vpnID)
	if len(vpnID) < 16 {
		vcl.ShowMessage("VPN??????????????????16??????????????????")
		return
	}
	var joinNetworkResponse datastruct.JoinNetworkResponse
	var joinNetworkRequest datastruct.JoinNetworkRequest
	joinNetworkRequest.AllowManaged = true
	joinNetworkRequest.AllowDNS = true
	joinNetworkRequest.AllowGlobal = true
	joinNetworkRequest.AllowDefault = true

	req, err := json.Marshal(joinNetworkRequest)
	if err != nil {
		fmt.Println("???????????????")
	}
	go func() {
		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://127.0.0.1:9993/network/"+vpnID, bytes.NewBuffer(req))
		if err != nil {
			fmt.Print("err=", err)
		}
		req.Header.Add("X-ZT1-Auth", getVpnToken())
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("err=", err)
			vcl.ShowMessage("?????????????????????????????????????????????")
			return
		} else {
			result, _ := ioutil.ReadAll(resp.Body)
			fmt.Print(string(result))
			err = json.Unmarshal(result, &joinNetworkResponse)
			//??????????????????
			//1.??????????????????VPN??????
			Tgconf.VPN.NetWorkId = vpnID
			Tgconf.VPN.Connect = true
			Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
			UpdateKey(Tgconf)
			Tgconf = InitKey()
			//2.???????????????????????????????????????
			showAllMembers(vpnID, lv1)
			//3.????????????????????????
			var network datastruct.Network
			network.NetworkID = vpnID
			network.UserName = f.VPNUsernameEdit.Text()
			network.DeviceType = f.DeviceTypeBox.Text()
			network.NetworkName = joinNetworkResponse.Name
			out, err := yaml.Marshal(network)
			if err = ioutil.WriteFile("./networks.d/"+vpnID+".yaml", out, 0666); err != nil {
				fmt.Println("Writefile Error =", err)
				return
			} else {
				vcl.ShowMessage("?????????????????????")
				//??????????????????
				fillNetworks(f)
			}

		}
	}()
}

func (f *TMainForm) onAutoJoinVPNClick(sender vcl.IObject) {

	if f.AutoJoinVPN.Checked() == true {
		changeZerotierStartMode("auto")

	} else {
		changeZerotierStartMode("demand")
	}
}

//??????/?????????????????????  VPNTest
func (f *TMainForm) onVPNTestClick(sender vcl.IObject) {
	if f.JoinVPNEdit.Text()== "" {
		vcl.ShowMessage("???????????????????????????")
		return
	}
	//TODO  ?????????UI???????????????????????????????????????????????? ??????????????????????????????????????????????????????????????????????????????????????????ui?????????????????? ?????????????????????????????????????????? ?????????????????????????????????????????????????????????
	//if f.VPNTest.Caption() == "????????????" {
		//?????????????????????????????????
		//1.??????????????????????????? ????????????????????????
		//?????????
		//if a, _ := getVPNConnectStatus(); a == 1 {
		//	//??????
		//	fmt.Println("?????????**********************************")
		//	localDisConnectVPN(Tgconf.VPN.NetWorkId)
		//}
		//2.????????????????????????????????????
		//connectchan := make(chan int)
		//connectresponsechan := make(chan datastruct.JoinNetworkResponse)

		//	status, connectresponse := localConnectVPN(mainForm.JoinVPNEdit.Text())
		//	//go func() {
		//	//	time.Sleep(4 * time.Second)
		//	//	connectchan <- 100
		//	//
		//	//}()
		//
		//	if status == 1 {
		//		//???????????? ??? ??????????????????
		//		Tgconf.VPN.NetWorkId = connectresponse.ID
		//		Tgconf.VPN.Connect = true
		//		Tgconf.VPN.NetWorkName = connectresponse.Name
		//		UpdateKey(Tgconf)
		//		Tgconf = InitKey()
		//		//???????????? ??????????????????
		//		//	fillNetworks(mainForm)
		//		//???????????? ??????????????????
		//		fmt.Println("???????????????????????????????????????????????????")
		//		showAllMembers(connectresponse.ID, lv1)
		//		//???????????? ???????????????ip
		//
		//		vcl.ThreadSync(func() {
		//			//mainForm.VPNIPEdit.SetText(connectresponse.AssignedAddresses[0])
		//			f.VPNTest.SetCaption("????????????")
		//		})
		//	} else {
		//		vcl.ShowMessage("???????????????????????????")
		//	}
		//
		//} else if f.VPNTest.Caption() == "????????????" {
		//	//?????????????????????????????????????????????
		//	//1.????????????
		//
		//	fmt.Println("???????????????????????????", localDisConnectVPN(Tgconf.VPN.NetWorkId))
		//	//2.???????????????????????????
		//	Tgconf.VPN.Connect = false
		//	UpdateKey(Tgconf)
		//	Tgconf = InitKey()
		//	fmt.Println("????????? ")
		//	//3.????????????
		//	//fillNetworks(mainForm)
		//
		//
		//}
		////??????
		if f.VPNTest.Caption() == "????????????"  {
			if Tgconf.VPN.Connect {
				localDisConnectVPN(Tgconf.VPN.NetWorkId)
			}


			//TODO ???????????? ??????3 ????????????
			var joinNetworkResponse datastruct.JoinNetworkResponse
			var joinNetworkRequest datastruct.JoinNetworkRequest
			var networkid string
			networkid = mainForm.JoinVPNEdit.Text()
			fmt.Println("***************" + networkid)
			joinNetworkRequest.AllowManaged = true
			joinNetworkRequest.AllowDNS = true
			joinNetworkRequest.AllowGlobal = true
			joinNetworkRequest.AllowDefault = true
			req, err := json.Marshal(joinNetworkRequest)
			if err != nil {
				fmt.Println("???????????????")
			}
			go func() {
				client := &http.Client{}
				req, err := http.NewRequest("POST", "http://127.0.0.1:9993/network/"+networkid, bytes.NewBuffer(req))
				if err != nil {
					fmt.Print("err=", err)
				}
				req.Header.Add("X-ZT1-Auth", getVpnToken())
				resp, err := client.Do(req)
				if err != nil {
					fmt.Print("err=", err)
				}
				result, _ := ioutil.ReadAll(resp.Body)
				fmt.Print(string(result))
				err = json.Unmarshal(result, &joinNetworkResponse)
				if err != nil {
					fmt.Println("ChangeVPNStatus????????????joinNetworkResponse?????????=", err)
				}
				//??????????????????
				//1.??????????????????VPN??????
				Tgconf.VPN.NetWorkId = networkid
				Tgconf.VPN.Connect = true
				Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()
				//2.???????????????????????????????????????
				showAllMembers(joinNetworkResponse.ID, lv1)
				//3.????????? ?????????????????? ??????????????????

				vcl.ThreadSync(func() {
					mainForm.VPNTest.SetCaption("????????????")
				})
				fillNetworks(mainForm)

			}()
		}
		//??????
	if f.VPNTest.Caption() == "????????????"  {
			//??????
			go func() {
				client := &http.Client{}
				req, err := http.NewRequest("DELETE", "http://127.0.0.1:9993/network/"+Tgconf.VPN.NetWorkId, nil)
				if err != nil {
					fmt.Print("err=", err)
				}
				req.Header.Add("X-ZT1-Auth", getVpnToken())
				resp, err := client.Do(req)
				if err != nil {
					fmt.Print("err=", err)
				}
				result, _ := ioutil.ReadAll(resp.Body)
				fmt.Print(string(result))
				//??????????????????
				//1.??????????????????VPN??????
				//Tgconf.VPN.NetWorkId = joinNetworkResponse.ID
				Tgconf.VPN.Connect = false
				//Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
				fmt.Println(Tgconf)
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()

				//3.????????? ?????????????????? ??????????????????

				vcl.ThreadSync(func() {
					//????????????ip??????
					f.VPNIPEdit.SetText("")
					mainForm.VPNTest.SetCaption("????????????")
				})
				fillNetworks(mainForm)
			}()

		}
	}
//}
func (f *TMainForm) onDeleteVPNButtonClick(sender vcl.IObject) {
	//1.???????????????????????????????????????????????? ??????????????????
	if mainForm.JoinVPNEdit.Text() == "" {
		vcl.ShowMessage("??????????????????????????????")
		return
	}
	var removenetworkid = mainForm.JoinVPNEdit.Text()
	if strings.HasSuffix(removenetworkid, "  ???") {
		removenetworkid = strings.TrimSuffix(removenetworkid, "  ???")
	}
	status, statusresponse := getVPNConnectStatus()
	fmt.Println(statusresponse.ID ,"  ",removenetworkid)
	if statusresponse.ID == removenetworkid {
		fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&????????????",status)
		if status == 1 {
			networkid := statusresponse.ID
			go func() {
				client := &http.Client{}
				req, err := http.NewRequest("DELETE", "http://127.0.0.1:9993/network/"+statusresponse.ID, nil)
				if err != nil {
					fmt.Print("err=", err)
				}
				req.Header.Add("X-ZT1-Auth", getVpnToken())
				resp, err := client.Do(req)
				if err != nil {
					fmt.Print("err=", err)
				}
				result, _ := ioutil.ReadAll(resp.Body)
				fmt.Print(string(result))

				//1.??????????????????VPN??????
				Tgconf.VPN.NetWorkId = ""
				Tgconf.VPN.Connect = false
				Tgconf.VPN.NetWorkName = ""
				fmt.Println(Tgconf)
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()
				vcl.ThreadSync(func() {
					//????????????ip??????
					f.VPNIPEdit.SetText("")
				})

				//???????????????????????????
				err = os.Remove("./networks.d/" + networkid + ".yaml")
				if err != nil {
					vcl.ShowMessage("????????????" + networkid + "?????????")
				} else {
					vcl.ShowMessage("????????????" + networkid + "?????????")
					//????????????????????????
					fillNetworks(f)
				}

			}()
		} else {
			//????????????ip??????
			f.VPNIPEdit.SetText("")
fmt.Println("////////////////////////////////////////////////////")
			//???????????????????????????
			err := os.Remove("./networks.d/" + removenetworkid + ".yaml")

			if err != nil {
				vcl.ShowMessage("????????????" + removenetworkid + "?????????")
			} else {
				vcl.ShowMessage("????????????" + removenetworkid + "?????????")
				//????????????????????????
				fillNetworks(f)
			}
		}
	} else {
		//???????????????????????????
		err := os.Remove("./networks.d/" + removenetworkid + ".yaml")

		if err != nil {
			fmt.Println("////////////////////////77777777////////////////////////////")
			vcl.ShowMessage("????????????" + removenetworkid + "?????????")
		} else {
			vcl.ShowMessage("????????????" + removenetworkid + "?????????")
			//????????????????????????
			fillNetworks(f)
		}
	}
}
//?????????????????????????????????????????????
func (f *TMainForm) onChangeNameButtonClick(sender vcl.IObject) {
	username := f.VPNUsernameEdit.Text()
	devicetype := f.DeviceTypeBox.Text()
	//central api ?????????????????????
	//POST   https://my.zerotier.com/api/v1/network/{networkID}/member/{memberID}

	status, statusresponse := getVPNConnectStatus()
	if status != 1 {
		vcl.ShowMessage("??????????????????????????????????????????????????????")
		return
	}

	go func() {
		var modifymember datastruct.ModifyMemberRequest
		var memberinfo datastruct.ModifyMemberResponse
		var capabilities = make([]int, 0)
		var ipAssignments = make([]string, 0)
		var tags = make([][]int, 0)
		modifymember.Name = username
		modifymember.Description = devicetype

		modifymember.Config.Capabilities = capabilities
		modifymember.Config.IPAssignments = ipAssignments
		modifymember.Config.Tags = tags

		req1, err := json.Marshal(modifymember)
		client := &http.Client{}
		address := getNodeId()
		fmt.Println(string(req1))
		url := "https://my.zerotier.com/api/v1/network/" + statusresponse.ID + "/member/" + address
		fmt.Println("url:", url)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(req1))
		if err != nil {
			vcl.ShowMessage("???????????????????????????????????????")
			fmt.Print("err=", err)
			return
		}
		req.Header.Add("Authorization", "Bearer "+"kocvbWu3ByI3SnVZ3LnZGNtMgzTKXr1p")
		resp, err := client.Do(req)
		if err != nil {
			vcl.ShowMessage("???????????????????????????????????????")
			fmt.Print("err=", err)
			return
		}
		result, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(result))
		err = json.Unmarshal(result, &memberinfo)
		if err != nil {
			fmt.Print("Unmarshalerr=", err)
			vcl.ShowMessage("???????????????")
			return
		}
		//????????????????????????
		if memberinfo.Name == username && memberinfo.Description == devicetype {

			vcl.ShowMessage("???????????????????????????????????????")
			//??????????????????????????????
			showAllMembers(statusresponse.ID, lv1)
		} else {
			vcl.ShowMessage("???????????????????????????????????????")
		}
	}()

}

//????????????????????????????????????
func (f *TMainForm) onHeziandVPNClick(sender vcl.IObject) {
	if Ips() == "127.0.0.1" {
		vcl.ShowMessage("?????????????????????????????????")
		return
	}
	status, _ := getHeziVPNConnectStatus()
	if status == 1 { //????????????????????????vpn ????????????????????????vpn??????
		fmt.Println("?????????????????? 1 ")
		heziDisconnectVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -1 { //????????????????????????vpn ????????????????????????vpn??????
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -2 { //?????????????????????????????????vpn ,??????????????????
		vcl.ShowMessage("??????????????????VPN??????!")
	} else if status == 0 { //????????????????????????vpn ????????????????????????vpn??????
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -3 { //????????????????????????vpn ????????????
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	}
}

//???????????????????????????????????????
func (f *TMainForm) onJoinVPNEditChange(sender vcl.IObject) {
	//1.??????????????????
	var selected string
	selected = f.JoinVPNEdit.Text()
	if strings.HasSuffix(selected, "  ???") {
		selected = strings.TrimSuffix(selected, "  ???")
		f.VPNTest.SetCaption("????????????")
	} else {
		f.VPNTest.SetCaption("????????????")
	}
	showAllMembers(selected, lv1)

	//2.??????/??????????????????

}

//????????????????????????cmd???????????????????????????????????????
func changeZerotierStartMode(mode string) {
	verb := "runas"
	//exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	//args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString("cmd")
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString("/c sc config ZeroTierOneService start= " + mode)

	var showCmd int32 = 0 //0?????????????????????  1?????????????????????????????????

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func showAllMembers(vpnid string, lv1 *vcl.TListView) {
	//????????????????????????

	//????????????
	lv1.Clear()

	//1.????????????api ????????????????????????
	var memberlist datastruct.MemberListResponse
	var err error

	memberlist, err = getAllVPNUsersById(vpnid)
	if err == nil {
		lv1.Items().BeginUpdate()
		for _, v := range memberlist {
			item := lv1.Items().Add()
			// ????????????Caption???????????????
			item.SetCaption(v.NodeID)
			item.SubItems().Add(v.Name)
			item.SubItems().Add(v.Description)
			if len(v.Config.IPAssignments) > 0 {
				item.SubItems().Add(v.Config.IPAssignments[0])
			} else {
				item.SubItems().Add("")
			}

			if v.Online == true {
				item.SubItems().Add("??????")
			} else {
				//timestamp := strconv.FormatInt(v.LastOnline/1e3, 10)
				//fmt.Println(timestamp)
				ts := time.Unix(v.LastOnline/1e3, 0).Format("2006-01-02")
				//var t,_ =time.Parse("2006-01-02 15:04:05",timestamp)
				item.SubItems().Add(ts)
			}

			item.SubItems().Add(v.PhysicalAddress)
		}

		lv1.Items().EndUpdate()
	}
}
