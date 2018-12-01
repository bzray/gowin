package app

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"sync"
)

func Gui() {
	Win()
}

type GshMainWindow struct {
	*walk.MainWindow
}

var animal = new(Animal)

func Win() {

	mw := new(GshMainWindow)

	var showAboutBoxAction *walk.Action
	var inTE, outTETop, outTEBon *walk.TextEdit
	var pushBT *walk.PushButton

	if _, err := (MainWindow{
		Title:   "姓氏筛选",
		AssignTo:&mw.MainWindow,
		MenuItems:[]MenuItem{
			Menu{
				Text:"文件",
				Items:[]MenuItem{
					Action{
						Text:"设置",
						OnTriggered:func(){
							if animal.Type==TypeScroll && Start{
								closeChannel<-true
								pushBT.SetText("滚动筛选(开始)")
							}
							if cmd, err := RunAnimalDialog(mw, animal);err!=nil{
								log.Printf("%v", err)
							}else if cmd == walk.DlgCmdOK{
								//outTETop.SetText(fmt.Sprintf("%+v", animal))
								if animal.Type==TypeScroll{
									pushBT.SetText("滚动筛选(开始)")
								}else{
									pushBT.SetText("随机筛选")
									inTE.SetText("")
								}
							}
						},
					},
					Separator{},
					Action{
						Text:"退出",
						OnTriggered:func() { mw.Close() },
					},
				},

			},
			Menu{
				Text:"帮助",
				Items:[]MenuItem{
					Action{
						AssignTo:&showAboutBoxAction,
						Text:"关于",
						OnTriggered:mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},
		ContextMenuItems:[]MenuItem{
			ActionRef{&showAboutBoxAction},
		},
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			//HSplitter{
			Composite{
				Layout:HBox{},
				Children: []Widget{
					TextEdit{AssignTo: &inTE, ReadOnly: true},
					Composite{
						Layout: VBox{
							MarginsZero:true,
							SpacingZero:true,
						},
						Children: []Widget{
							Composite{
								Layout:VBox{
									Margins:Margins{0,0,0,5},
								},
								Children:[]Widget{
									TextEdit{
										AssignTo: &outTETop,
										ReadOnly: true,
										Text:"姓氏",
									},
								},
							},
							Composite{
								Layout:VBox{
									MarginsZero:true,
									SpacingZero:true,
								},
								Children:[]Widget{
									TextEdit{
										AssignTo: &outTEBon,
										ReadOnly: true,
										Text:"解读",
									},
								},
							},
						},
					},
				},
			},
			//},
			PushButton{
				AssignTo:&pushBT,
				Text: "滚动筛选(开始)",
				OnClicked: func() {
					//outTE.SetText(strings.ToUpper(inTE.Text()))

					//字体初始化
					once.Do(func() {
						font, err := walk.NewFont("微软雅黑", 20, walk.FontBold)
						if err == nil {
							outTETop.SetFont(font)
							inTE.SetFont(font)
						}
						//defer font.Dispose()
					})

					if animal.Type==TypeRandom{
						//1.随机生成一个姓氏
						c := GetRandom()
						outTETop.SetText(c.Word+"("+c.Pinyin+")")
						outTEBon.SetText(c.Desc)
					}else{
						//2.滚屏显示
						if Start==true{
							closeChannel<-true
							Start = false
							pushBT.SetText("滚动筛选(开始)")
						}else{
							wordChannel = make(chan Chinese, 5)
							closeChannel = make(chan bool)

							go func() {
								GetScroll(wordChannel, closeChannel)
							}()

							go func() {
								var lword Chinese
								for{
									select{
										case word,ok := <-wordChannel:
											if !ok{
												log.Println("接收结束")
												outTETop.SetText(lword.Word+"("+lword.Pinyin+")")
												outTEBon.SetText(lword.Desc)
												return
											}
											log.Printf("接收:%v\n",word)
											lword = word
											inTE.SetText(word.Word+"("+word.Pinyin+")")

									}
								}
							}()

							Start = true

							pushBT.SetText("滚动筛选(结束)")
						}
					}
				},
			},
		},
	}.Run()); err!=nil{
		log.Fatal(err)
	}
}

var Start bool = false
var wordChannel chan Chinese
var closeChannel chan bool
var once sync.Once

func (mw *GshMainWindow) showAboutBoxAction_Triggered(){
	walk.MsgBox(mw,"关于","随机筛选姓氏\n版本:V1.0\n时间:2017-10-12\n邮箱:liraygo@163.com", walk.MsgBoxIconInformation)
}

type Mode int

const (
	TypeScroll Mode = iota
	TypeRandom
)

type Animal struct {
	Type Mode
}

func RunAnimalDialog(owner walk.Form, animal *Animal) (int, error){
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	return Dialog{
		AssignTo:&dlg,
		Title:"设置",
		DefaultButton:&acceptPB,
		CancelButton:&cancelPB,
		DataBinder:DataBinder{
			AssignTo:       &db,
			DataSource:     animal,
		},
		MinSize:Size{300, 180},
		Layout:VBox{},
		Children:[]Widget{
			Composite{
				Layout:Grid{Columns:2},
				Children:[]Widget{
					RadioButtonGroupBox{
						ColumnSpan: 2,
						Title:      "筛选模式",
						Layout:     HBox{},
						DataMember: "Type",
						Buttons: []RadioButton{
							{Text: "随机筛选", Value: TypeRandom},
							{Text: "滚动选择", Value: TypeScroll},
						},
					},
				},
			},
			Composite{
				Layout:HBox{},
				Children:[]Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确认",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
