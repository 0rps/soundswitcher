package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/mafik/pulseaudio"
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
