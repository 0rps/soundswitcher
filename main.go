package main

import (
	"fmt"
	"os"

	"github.com/mafik/pulseaudio"
	"github.com/visualfc/atk/tk"
)

const (
	btnEscapeCode = 9
	btnEnterCode  = 36

	mwWidth  = 500
	mwHeight = 400
)

type Window struct {
	*tk.Window

	radioGroup *tk.RadioGroup
}

func NewWindow(audioOutputs []PulseAudioOutput) *Window {
	mw := &Window{
		Window:     tk.NewWindow(),
		radioGroup: tk.NewRadioGroup(),
	}

	vbox := tk.NewVPackLayout(mw)
	for _, v := range audioOutputs {
		btn := mw.radioGroup.AddNewRadio(mw, v.Name, v)
		if v.IsActive {
			btn.SetChecked(true)
			btn.SetFocus()
		}
		vbox.AddWidget(btn, tk.PackAttrFillBoth(), tk.PackAttrExpand(true))
	}

	mw.BindKeyEvent(mw.OnKeyEvent)
	mw.OnClose(func() bool {
		os.Exit(0)
		return false
	})
	return mw
}

func (w *Window) OnKeyEvent(e *tk.KeyEvent) {
	if e.KeyCode == btnEnterCode {
		for _, radioBtn := range w.radioGroup.RadioList() {
			if radioBtn.IsFocus() {
				soundOutput := w.radioGroup.RadioData(radioBtn).(PulseAudioOutput)
				soundOutput.Output.Activate()
				os.Exit(0)
			}
		}
	} else if e.KeyCode == btnEscapeCode {
		os.Exit(0)
	}
}

func main() {
	client, err := pulseaudio.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer client.Close()
	tk.MainLoop(func() {
		mw := NewWindow(getPulseAudioOutputs(client))
		mw.SetTitle("Sound Output Switcher")
		mw.SetSizeN(mwWidth, mwHeight)
		mw.Center(nil)
		mw.ShowNormal()
	})
}
