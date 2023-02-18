package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/mafik/pulseaudio"
	"github.com/visualfc/atk/tk"
)

type PulseAudioOutput struct {
	Output   pulseaudio.Output
	Name     string
	IsActive bool
}

type OutputsByName []PulseAudioOutput

func (a OutputsByName) Len() int           { return len(a) }
func (a OutputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a OutputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func getPulseAudioOutputs(client *pulseaudio.Client) []PulseAudioOutput {
	outs, activeIndex, err := client.Outputs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := make([]PulseAudioOutput, len(outs))
	for i, out := range outs {
		name := fmt.Sprintf("%s (%s)", out.CardName, out.PortName)
		result[i] = PulseAudioOutput{out, name, i == activeIndex}
	}
	sort.Sort(OutputsByName(result))
	return result
}

type Window struct {
	*tk.Window

	client *pulseaudio.Client

	soundOutputs       []PulseAudioOutput
	soundOutputButtons []*tk.RadioButton
}

func NewWindow(client *pulseaudio.Client) *Window {
	outputs := getPulseAudioOutputs(client)
	mw := &Window{client: client, soundOutputs: outputs}

	mw.Window = tk.RootWindow()
	vbox := tk.NewVPackLayout(mw)

	rgroup := tk.NewRadioGroup()
	for i, v := range mw.soundOutputs {
		btn := rgroup.AddNewRadio(mw, v.Name, i)
		if v.IsActive {
			btn.SetChecked(true)
		}
		mw.soundOutputButtons = append(mw.soundOutputButtons, btn)
		vbox.AddWidget(btn, tk.PackAttrFillBoth(), tk.PackAttrExpand(true))
	}
	mw.BindKeyEvent(mw.OnKeyEvent)
	return mw
}

func (w *Window) OnKeyEvent(e *tk.KeyEvent) {
	if e.KeyCode == 36 { // enter
		for index, rbutton := range w.soundOutputButtons {
			if rbutton.IsFocus() {
				w.soundOutputs[index].Output.Activate()
				os.Exit(0)
			}
		}
	} else if e.KeyCode == 9 { // esc
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
		mw := NewWindow(client)
		mw.SetTitle("Sound Output Switcher")
		mw.SetSizeN(500, 400)
		mw.Center(nil)
		mw.ShowNormal()
	})
}
