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

	oldBox string //保存上次修改的box名称，发生变化时需要根据界面内容设置Gconf的内存数据
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

	//sheet3 成员列表
	MembersScrollBox *vcl.TScrollBox

	Gconf GlobalConfig //保存配置过程的参数情况，可应用、确认，或取消设置
}

var mainForm *TMainForm

func (f *TMainForm) GetConfig() {
	Tgconf = InitKey()
	f.Gconf = Tgconf

	//根据配置文件设置界面元素
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
	//根据界面元素保存配置文件

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
		log.Fatalf("打开日志文件失败：%s:%v\n", Tgconf.DefaultLogFile, err)
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
	//生成对应的conf文件
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

	//退出rclone，下次自动连接默认的box
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
	f.SetCaption("设置智能连接助手")
	f.SetPosition(types.PoScreenCenter)
	f.EnabledMaximize(false)
	f.SetWidth(640)
	f.SetHeight(480)
	// 全局设置提示
	f.SetShowHint(true)

	// 动态创建
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

	if vcl.MessageDlg("是否覆盖式引用"+f.Template.Text()+"号模板?", types.MtInformation, types.MbYes, types.MbNo) == types.MrYes {
		//
		f.Gconf.DefaultTemplate = "rclone.conf." + f.Template.Text() + "box"
		cfg, err := ini.Load(f.Gconf.DefaultTemplate)

		if err != nil {
			Info.Print("%v err->%v\n", "load template：", err)
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
	//发生修改box的切换，需要记录oldBox内容
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
	f.Password.SetText(string(decoded)) //需要解密
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

	//生成对应的conf文件

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
		f.BtnRefresh.SetCaption("刷新(&S)")
		f.LbRefresh.SetCaption("已挂载：(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ")")
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
				vcl.ThreadSync(func() { //非主线程访问ui
					f.LbRefresh.SetCaption("已挂载：(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ");\n" + "总文件数：" + strings.Split(string(line1), ":")[1] + ";\t总字节数：" + strings.Split(string(line2), ":")[1] + ")")
					f.LbRefresh.SetCursor(types.CrDefault)
					f.BtnRefresh.SetCursor(types.CrDefault)
				})
			}
		}()
	} else {
		f.BtnRefresh.SetCaption("重连(&C)")
		f.LbRefresh.SetCaption("未挂载：(" + Tgconf.DefaultBox + ":/" + Tgconf.AllBox[Tgconf.DefaultBox].VirtualDisk + ")")
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
	f.BtnApply.SetCaption("应用(&A)")
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
	f.BtnOk.SetCaption("确定(&O)")
	f.BtnOk.SetFont(ftf)
	f.BtnOk.SetOnClick(f.OnBtnOkClick)

	f.BtnCancel = vcl.NewButton(f)
	f.BtnCancel.SetParent(f)
	f.BtnCancel.SetLeft(472)
	f.BtnCancel.SetTop(408)
	f.BtnCancel.SetHeight(47)
	f.BtnCancel.SetWidth(134)
	f.BtnCancel.SetCaption("取消(&C)")
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
	f.Sheet1.SetCaption("设置存储节点参数")

	f.Sheet2 = vcl.NewTabSheet(f)
	f.Sheet2.SetPageControl(f.Pgc)
	f.Sheet2.SetCaption("设置共享同步参数")

	f.Sheet3 = vcl.NewTabSheet(f)
	f.Sheet3.SetPageControl(f.Pgc)
	f.Sheet3.SetCaption("设置VPN联网参数")

	f.Lb1 = vcl.NewLabel(f)
	f.Lb1.SetParent(f.Sheet1)
	f.Lb1.SetCaption("默认BOX：")
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
	f.Lb2.SetCaption("大脑BOX：")
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
	f.Lb3.SetCaption("引入多节点模板")
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
	f.Lb4.SetCaption("修改节点配置：")
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
	f.BtnNewBox.SetCaption("增加新BOX(&N)")
	f.BtnNewBox.SetOnClick(f.OnBtnNewBoxClick)
	f.BtnDeleteBox = vcl.NewButton(f)
	f.BtnDeleteBox.SetParent(f.Sheet1)
	f.BtnDeleteBox.SetLeft(488)
	f.BtnDeleteBox.SetTop(40)
	f.BtnDeleteBox.SetHeight(30)
	f.BtnDeleteBox.SetWidth(112)
	f.BtnDeleteBox.SetCaption("删除BOX(&D)")
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
	f.Lb5.SetCaption("连接类型：")
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
	f.Lb6.SetCaption("厂商：")
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
	f.Lb7.SetCaption("服务地址：")
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
	f.Lb8.SetCaption("用户：")
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
	f.Lb9.SetCaption("口令：")
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
	f.CbDefaultBox.SetCaption("默认box")
	f.CbDefaultBox.SetLeft(28)
	f.CbDefaultBox.SetTop(80)
	f.CbDefaultBox.SetHeight(24)
	f.CbDefaultBox.SetWidth(82)
	f.CbDefaultBox.SetChecked(true)
	f.CbDefaultBox.SetOnChange(f.OnCbDefaultBoxChange)
	f.CbDefaultBrain = vcl.NewCheckBox(f)
	f.CbDefaultBrain.SetParent(f.GroupBox1)
	f.CbDefaultBrain.SetCaption("大脑box")
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
	f.BtnTestMount.SetCaption("测试挂载(&M)")
	f.BtnTestMount.SetOnClick(f.OnBtnTestMountClick)

	f.Lb10 = vcl.NewLabel(f)
	f.Lb10.SetParent(f.Sheet1)
	f.Lb10.SetLeft(16)
	f.Lb10.SetTop(248)
	f.Lb10.SetHeight(20)
	f.Lb10.SetWidth(105)
	f.Lb10.SetCaption("智能探测间隔：")
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
	f.Lb11.SetCaption("秒")
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
	f.BtnChangePassword.SetCaption("批量改密码(&P)")
	f.BtnChangePassword.SetOnClick(f.OnBtnChangePasswordClick)

	f.Lb12 = vcl.NewLabel(f)
	f.Lb12.SetParent(f.Sheet1)
	f.Lb12.SetLeft(16)
	f.Lb12.SetTop(288)
	f.Lb12.SetHeight(20)
	f.Lb12.SetWidth(75)
	f.Lb12.SetCaption("日志文件：")
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
	f.BtnViewLog.SetCaption("查看日志(&R)")
	f.BtnViewLog.SetOnClick(f.OnBtnViewLogClick)
	f.CbDefaultAutoBoot = vcl.NewCheckBox(f)
	f.CbDefaultAutoBoot.SetParent(f.Sheet1)
	f.CbDefaultAutoBoot.SetCaption("是否开机自启动")
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
	f.BtnRefresh.SetCaption("重连(&C)")
	f.BtnRefresh.SetOnClick(f.OnBtnRefreshClick)
	f.LbRefresh = vcl.NewLabel(f)
	f.LbRefresh.SetParent(f.Sheet1)
	f.LbRefresh.SetLeft(116)
	f.LbRefresh.SetTop(328)
	f.LbRefresh.SetHeight(40)
	f.LbRefresh.SetWidth(405)
	f.LbRefresh.SetCaption("未挂载：(box1:/Z:) 总文件数： 总字节数：")

	//sheet3 初始化控件
	f.VPNConfigGroupBox = vcl.NewGroupBox(f)
	f.VPNConfigGroupBox.SetParent(f.Sheet3)
	f.VPNConfigGroupBox.SetBounds(32, 16, 560, 228)
	f.VPNConfigGroupBox.SetCaption("配置网络")

	f.JoinVPNLabel = vcl.NewLabel(f)
	f.JoinVPNLabel.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNLabel.SetBounds(40, 0, 72, 24)

	f.JoinVPNLabel.SetCaption("管理子网:")

	f.JoinVPNEdit = vcl.NewComboBox(f)
	f.JoinVPNEdit.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNEdit.SetBounds(120, 0, 160, 33)

	f.JoinVPNEdit.SetOnChange(f.onJoinVPNEditChange)

	f.JoinVPNButton = vcl.NewButton(f)
	f.JoinVPNButton.SetParent(f.VPNConfigGroupBox)
	f.JoinVPNButton.SetBounds(312, 0, 96, 33)
	f.JoinVPNButton.SetCaption("添加子网")
	f.JoinVPNButton.SetOnClick(f.onJoinVPNButtonClick)

	f.DeleteVPN = vcl.NewButton(f)
	f.DeleteVPN.SetParent(f.VPNConfigGroupBox)
	f.DeleteVPN.SetBounds(424, 0, 96, 33)
	f.DeleteVPN.SetCaption("退出子网")
	f.DeleteVPN.SetOnClick(f.onDeleteVPNButtonClick)

	f.AutoJoinVPN = vcl.NewCheckBox(f)
	f.AutoJoinVPN.SetParent(f.VPNConfigGroupBox)
	f.AutoJoinVPN.SetBounds(456, 120, 146, 27)
	f.AutoJoinVPN.SetCaption("开机自启动")
	f.AutoJoinVPN.SetOnClick(f.onAutoJoinVPNClick)

	f.VPNStatusLabel = vcl.NewLabel(f)
	f.VPNStatusLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNStatusLabel.SetBounds(40, 120, 80, 24)
	f.VPNStatusLabel.SetCaption("当前状态：")

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
	f.VPNIPLabel.SetCaption("分配IP：")

	f.VPNIPEdit = vcl.NewEdit(f)
	f.VPNIPEdit.SetParent(f.VPNConfigGroupBox)
	f.VPNIPEdit.SetBounds(288, 120, 160, 33)
	f.VPNIPEdit.SetEnabled(false)

	f.VPNUsernameLabel = vcl.NewLabel(f)
	f.VPNUsernameLabel.SetParent(f.VPNConfigGroupBox)
	f.VPNUsernameLabel.SetBounds(40, 40, 72, 24)
	f.VPNUsernameLabel.SetCaption("用户名：")

	f.VPNUsernameEdit = vcl.NewEdit(f)
	f.VPNUsernameEdit.SetParent(f.VPNConfigGroupBox)
	f.VPNUsernameEdit.SetBounds(120, 40, 160, 33)
	var hostname, _ = os.Hostname()
	fmt.Println("hostname:", hostname)
	f.VPNUsernameEdit.SetText(hostname)

	f.DeviceTypeLabel = vcl.NewLabel(f)
	f.DeviceTypeLabel.SetParent(f.VPNConfigGroupBox)
	f.DeviceTypeLabel.SetBounds(40, 79, 72, 24)
	f.DeviceTypeLabel.SetCaption("设备类型:")

	f.DeviceTypeBox = vcl.NewComboBox(f)
	f.DeviceTypeBox.SetParent(f.VPNConfigGroupBox)
	f.DeviceTypeBox.SetBounds(120, 79, 160, 33)
	f.DeviceTypeBox.Items().Add("笔记本电脑")
	f.DeviceTypeBox.Items().Add("台式机")
	f.DeviceTypeBox.Items().Add("虚拟机")
	f.DeviceTypeBox.SetItemIndex(0)

	f.DeviceTypeBox.SetStyle(types.CsDropDownList)

	f.ChangeNameButton = vcl.NewButton(f)
	f.ChangeNameButton.SetParent(f.VPNConfigGroupBox)
	f.ChangeNameButton.SetBounds(312, 56, 96, 33)
	f.ChangeNameButton.SetCaption("保存修改")
	f.ChangeNameButton.SetOnClick(f.onChangeNameButtonClick)

	f.VPNTest = vcl.NewButton(f)
	f.VPNTest.SetParent(f.VPNConfigGroupBox)
	f.VPNTest.SetBounds(424, 56, 96, 33)
	if Tgconf.VPN.Connect {
		f.VPNTest.SetCaption("断开子网")
	} else {
		f.VPNTest.SetCaption("连接子网")
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
	col.SetCaption("节点ID")
	col.SetWidth(100)

	col = lv1.Columns().Add()
	col.SetCaption("用户名")
	col.SetWidth(90)

	col = lv1.Columns().Add()
	col.SetCaption("设备类型")
	col.SetWidth(75)

	col = lv1.Columns().Add()
	col.SetCaption("分配IP")
	col.SetWidth(110)

	col = lv1.Columns().Add()
	col.SetCaption("在线状态")
	col.SetWidth(80)

	col = lv1.Columns().Add()
	col.SetCaption("物理IP")
	col.SetWidth(110)

	//填充子网列表
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

	//显示本机 盒子 vpn的Nodeid
	f.LocalNodeIdLabel = vcl.NewLabel(f)
	f.LocalNodeIdLabel.SetParent(f.VPNConfigGroupBox)
	f.LocalNodeIdLabel.SetBounds(40, 190, 84, 15)
	nodeId := getNodeId()
	if nodeId != "" {
		f.LocalNodeIdLabel.SetCaption(nodeId)
	} else {
		f.LocalNodeIdLabel.SetCaption("未获取到ID")
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
		f.VPNidLabel.SetCaption("未获取到ID")
	}

	vpnstatus, statusresponse := getVPNConnectStatus()
	if vpnstatus == 1 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClLime)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("已连接")
		f.VPNIPEdit.SetText(statusresponse.AssignedAddresses[0])
	} else if vpnstatus == 0 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClYellow)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("连接中")
	} else if vpnstatus == -2 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClRed)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("连接失败")
	} else if vpnstatus == 2 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClGrey)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("未连接")
	} else if vpnstatus == -1 {
		pen.SetColor(colors.ClWhite)
		brush.SetColor(colors.ClBlack)
		f.VPNStatusShape.SetPen(pen)
		f.VPNStatusShape.SetBrush(brush)
		f.VPNStatus.SetCaption("无连接记录")
	}
	if Tgconf.VPN.NetWorkId != "" {
		showAllMembers(Tgconf.VPN.NetWorkId, lv1)
	}

}

func fillNetworks(f *TMainForm) {
	//先清空当前控件内容+
	fmt.Println("进入fillNetworks")
	vcl.ThreadSync(func() {
		f.JoinVPNEdit.Clear()
	})
	//1.获取全部网络列表
	// 要遍历的文件夹

	dir := `./networks.d`
	fileinfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("读取文件夹：", err)
	}

	fmt.Println("读取文件夹成功")
	// 遍历这个文件夹
	var allNetworks datastruct.AllNetworks
	for _, fi := range fileinfo {
		// 判断是不是目录
		data, _ := ioutil.ReadFile("./networks.d/" + fi.Name())
		var networkconfig datastruct.Network
		if err := yaml.Unmarshal(data, &networkconfig); err != nil {
			log.Fatalf("unmarshal 1error: %v", err)
		}

		allNetworks = append(allNetworks, networkconfig)
	}
	fmt.Println("遍历文件夹成功")
	//2.填充表单
	var index int
	vcl.ThreadSync(func() {
		for i := 0; i < len(allNetworks); i++ {
			if Tgconf.VPN.NetWorkId != "" && allNetworks[i].NetworkID == Tgconf.VPN.NetWorkId && Tgconf.VPN.Connect {
				index = i
				f.JoinVPNEdit.Items().Add(allNetworks[i].NetworkID + "  √")
			} else {
				f.JoinVPNEdit.Items().Add(allNetworks[i].NetworkID)
			}
		}
		f.JoinVPNEdit.SetItemIndex(int32(index))
	})

	fmt.Println("获取所有子网列表成功！")

}

//JoinVPNButton处理函数  加入子网
func (f *TMainForm) onJoinVPNButtonClick(sender vcl.IObject) {
	vpnID := f.JoinVPNEdit.Text()
	fmt.Println(vpnID)
	if len(vpnID) < 16 {
		vcl.ShowMessage("VPN地址长度不足16位，请修改！")
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
		fmt.Println("序列化失败")
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
			vcl.ShowMessage("连接失败，错误信息请查看日志！")
			return
		} else {
			result, _ := ioutil.ReadAll(resp.Body)
			fmt.Print(string(result))
			err = json.Unmarshal(result, &joinNetworkResponse)
			//连接成功后：
			//1.修改配置文件VPN部分
			Tgconf.VPN.NetWorkId = vpnID
			Tgconf.VPN.Connect = true
			Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
			UpdateKey(Tgconf)
			Tgconf = InitKey()
			//2.同步更新下方的成员列表信息
			showAllMembers(vpnID, lv1)
			//3.保存子网配置文件
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
				vcl.ShowMessage("添加子网成功！")
				//刷新子网列表
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

//断开/连接网络的方法  VPNTest
func (f *TMainForm) onVPNTestClick(sender vcl.IObject) {
	if f.JoinVPNEdit.Text()== "" {
		vcl.ShowMessage("请选择或添加子网！")
		return
	}
	//TODO  由于主UI线程如果被协程阻塞的话会失去响应 ，所以如果想抽取函数并根据函数的返回值来执行程序的话，会使主ui线程失去响应 目前是在协程中执行全部的逻辑 执行完主线程不对协程的处理结果进行响应
	//if f.VPNTest.Caption() == "连接子网" {
		//执行连接选定的子网逻辑
		//1.是否需要断开现有的 现在是否有连接的
		//先断开
		//if a, _ := getVPNConnectStatus(); a == 1 {
		//	//下线
		//	fmt.Println("先下线**********************************")
		//	localDisConnectVPN(Tgconf.VPN.NetWorkId)
		//}
		//2.断开之后，连接目前选中的
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
		//		//连接成功 ： 更新配置文件
		//		Tgconf.VPN.NetWorkId = connectresponse.ID
		//		Tgconf.VPN.Connect = true
		//		Tgconf.VPN.NetWorkName = connectresponse.Name
		//		UpdateKey(Tgconf)
		//		Tgconf = InitKey()
		//		//连接成功 刷新子网列表
		//		//	fillNetworks(mainForm)
		//		//连接成功 刷新成员列表
		//		fmt.Println("和额呵呵集合接口就能看就不能控件并")
		//		showAllMembers(connectresponse.ID, lv1)
		//		//连接成功 刷新分配的ip
		//
		//		vcl.ThreadSync(func() {
		//			//mainForm.VPNIPEdit.SetText(connectresponse.AssignedAddresses[0])
		//			f.VPNTest.SetCaption("断开子网")
		//		})
		//	} else {
		//		vcl.ShowMessage("连接失败，请重试！")
		//	}
		//
		//} else if f.VPNTest.Caption() == "断开子网" {
		//	//执行断开当前已连接的网络的逻辑
		//	//1.断开子网
		//
		//	fmt.Println("断开子网返回状态：", localDisConnectVPN(Tgconf.VPN.NetWorkId))
		//	//2.成功后修改配置文件
		//	Tgconf.VPN.Connect = false
		//	UpdateKey(Tgconf)
		//	Tgconf = InitKey()
		//	fmt.Println("到这里 ")
		//	//3.刷新显示
		//	//fillNetworks(mainForm)
		//
		//
		//}
		////连接
		if f.VPNTest.Caption() == "连接子网"  {
			if Tgconf.VPN.Connect {
				localDisConnectVPN(Tgconf.VPN.NetWorkId)
			}


			//TODO 抽取上线 下线3 函数出来
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
				fmt.Println("序列化失败")
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
					fmt.Println("ChangeVPNStatus反序列化joinNetworkResponse失败：=", err)
				}
				//连接成功后：
				//1.修改配置文件VPN部分
				Tgconf.VPN.NetWorkId = networkid
				Tgconf.VPN.Connect = true
				Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()
				//2.成功后实时刷新子网成员列表
				showAllMembers(joinNetworkResponse.ID, lv1)
				//3.成功后 修改子网列表 修改按钮提示

				vcl.ThreadSync(func() {
					mainForm.VPNTest.SetCaption("断开子网")
				})
				fillNetworks(mainForm)

			}()
		}
		//断开
	if f.VPNTest.Caption() == "断开子网"  {
			//下线
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
				//连接成功后：
				//1.修改配置文件VPN部分
				//Tgconf.VPN.NetWorkId = joinNetworkResponse.ID
				Tgconf.VPN.Connect = false
				//Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
				fmt.Println(Tgconf)
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()

				//3.成功后 修改子网列表 修改按钮提示

				vcl.ThreadSync(func() {
					//修改分配ip为空
					f.VPNIPEdit.SetText("")
					mainForm.VPNTest.SetCaption("连接子网")
				})
				fillNetworks(mainForm)
			}()

		}
	}
//}
func (f *TMainForm) onDeleteVPNButtonClick(sender vcl.IObject) {
	//1.如果当前正在连接或者已连接该子网 则需要先退出
	if mainForm.JoinVPNEdit.Text() == "" {
		vcl.ShowMessage("请选择要退出的子网！")
		return
	}
	var removenetworkid = mainForm.JoinVPNEdit.Text()
	if strings.HasSuffix(removenetworkid, "  √") {
		removenetworkid = strings.TrimSuffix(removenetworkid, "  √")
	}
	status, statusresponse := getVPNConnectStatus()
	fmt.Println(statusresponse.ID ,"  ",removenetworkid)
	if statusresponse.ID == removenetworkid {
		fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&退出子网",status)
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

				//1.修改配置文件VPN部分
				Tgconf.VPN.NetWorkId = ""
				Tgconf.VPN.Connect = false
				Tgconf.VPN.NetWorkName = ""
				fmt.Println(Tgconf)
				UpdateKey(Tgconf)
				time.Sleep(1 * time.Second)
				Tgconf = InitKey()
				vcl.ThreadSync(func() {
					//修改分配ip为空
					f.VPNIPEdit.SetText("")
				})

				//删除相应的配置文件
				err = os.Remove("./networks.d/" + networkid + ".yaml")
				if err != nil {
					vcl.ShowMessage("退出子网" + networkid + "失败！")
				} else {
					vcl.ShowMessage("退出子网" + networkid + "成功！")
					//重新生成子网列表
					fillNetworks(f)
				}

			}()
		} else {
			//修改分配ip为空
			f.VPNIPEdit.SetText("")
fmt.Println("////////////////////////////////////////////////////")
			//删除相应的配置文件
			err := os.Remove("./networks.d/" + removenetworkid + ".yaml")

			if err != nil {
				vcl.ShowMessage("退出子网" + removenetworkid + "失败！")
			} else {
				vcl.ShowMessage("退出子网" + removenetworkid + "成功！")
				//重新生成子网列表
				fillNetworks(f)
			}
		}
	} else {
		//删除相应的配置文件
		err := os.Remove("./networks.d/" + removenetworkid + ".yaml")

		if err != nil {
			fmt.Println("////////////////////////77777777////////////////////////////")
			vcl.ShowMessage("退出子网" + removenetworkid + "失败！")
		} else {
			vcl.ShowMessage("退出子网" + removenetworkid + "成功！")
			//重新生成子网列表
			fillNetworks(f)
		}
	}
}
//处理修改用户名和设备类型的函数
func (f *TMainForm) onChangeNameButtonClick(sender vcl.IObject) {
	username := f.VPNUsernameEdit.Text()
	devicetype := f.DeviceTypeBox.Text()
	//central api 修改成员的信息
	//POST   https://my.zerotier.com/api/v1/network/{networkID}/member/{memberID}

	status, statusresponse := getVPNConnectStatus()
	if status != 1 {
		vcl.ShowMessage("当前未连接到任何网络，请连接后设置！")
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
			vcl.ShowMessage("与远程服务器建立连接失败！")
			fmt.Print("err=", err)
			return
		}
		req.Header.Add("Authorization", "Bearer "+"kocvbWu3ByI3SnVZ3LnZGNtMgzTKXr1p")
		resp, err := client.Do(req)
		if err != nil {
			vcl.ShowMessage("与远程服务器建立连接失败！")
			fmt.Print("err=", err)
			return
		}
		result, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(result))
		err = json.Unmarshal(result, &memberinfo)
		if err != nil {
			fmt.Print("Unmarshalerr=", err)
			vcl.ShowMessage("修改失败！")
			return
		}
		//验证是否修改成功
		if memberinfo.Name == username && memberinfo.Description == devicetype {

			vcl.ShowMessage("修改用户名和设备类型成功！")
			//重新填充用户信息列表
			showAllMembers(statusresponse.ID, lv1)
		} else {
			vcl.ShowMessage("修改用户名和设备类型失败！")
		}
	}()

}

//处理盒子联网按钮点击事件
func (f *TMainForm) onHeziandVPNClick(sender vcl.IObject) {
	if Ips() == "127.0.0.1" {
		vcl.ShowMessage("请先使用网线连接盒子！")
		return
	}
	status, _ := getHeziVPNConnectStatus()
	if status == 1 { //表示盒子已经连接vpn 此时执行盒子断开vpn逻辑
		fmt.Println("得道盒子状态 1 ")
		heziDisconnectVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -1 { //表示盒子已经断开vpn 此时执行盒子连接vpn逻辑
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -2 { //此时表示未指定连接哪个vpn ,弹框提醒一下
		vcl.ShowMessage("请指定或添加VPN网络!")
	} else if status == 0 { //表示盒子正在连接vpn 此时执行盒子重连vpn逻辑
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	} else if status == -3 { //表示盒子未连接过vpn 需要连接
		heziJoinVPN(Tgconf.VPN.NetWorkId, Tgconf.DefaultIpAddr)
	}
}

//处理管理子网的下垃框的事件
func (f *TMainForm) onJoinVPNEditChange(sender vcl.IObject) {
	//1.成员列表刷新
	var selected string
	selected = f.JoinVPNEdit.Text()
	if strings.HasSuffix(selected, "  √") {
		selected = strings.TrimSuffix(selected, "  √")
		f.VPNTest.SetCaption("断开子网")
	} else {
		f.VPNTest.SetCaption("连接子网")
	}
	showAllMembers(selected, lv1)

	//2.连接/断开按钮修改

}

//调用管理员权限的cmd执行修改服务启动方式的操作
func changeZerotierStartMode(mode string) {
	verb := "runas"
	//exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	//args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString("cmd")
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString("/c sc config ZeroTierOneService start= " + mode)

	var showCmd int32 = 0 //0表示窗口不显示  1表示执行时窗口正常显示

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func showAllMembers(vpnid string, lv1 *vcl.TListView) {
	//填充设备列表信息

	//首先清空
	lv1.Clear()

	//1.调用中心api 获取所有用户信息
	var memberlist datastruct.MemberListResponse
	var err error

	memberlist, err = getAllVPNUsersById(vpnid)
	if err == nil {
		lv1.Items().BeginUpdate()
		for _, v := range memberlist {
			item := lv1.Items().Add()
			// 第一列为Caption属性所管理
			item.SetCaption(v.NodeID)
			item.SubItems().Add(v.Name)
			item.SubItems().Add(v.Description)
			if len(v.Config.IPAssignments) > 0 {
				item.SubItems().Add(v.Config.IPAssignments[0])
			} else {
				item.SubItems().Add("")
			}

			if v.Online == true {
				item.SubItems().Add("在线")
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
