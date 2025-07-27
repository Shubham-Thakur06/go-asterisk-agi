package agi

import (
	"fmt"
	"strings"
)

// MRCPRecog represents an MRCP recognition request
type MRCPRecog struct {
	Grammar     string
	Timeout     int
	Options     map[string]string
	ResultVar   string
	CompletionC string
}

// MRCPSynth represents an MRCP synthesis request
type MRCPSynth struct {
	Text      string
	Options   map[string]string
	ResultVar string
}

// SpeechCreate creates a speech object
func (s *AgiSession) SpeechCreate() error {
	_, err := s.execute("SPEECH CREATE")
	return err
}

// SpeechDestroy destroys a speech object
func (s *AgiSession) SpeechDestroy() error {
	_, err := s.execute("SPEECH DESTROY")
	return err
}

// SpeechLoadGrammar loads a grammar
func (s *AgiSession) SpeechLoadGrammar(grammar, path string) error {
	_, err := s.execute(fmt.Sprintf("SPEECH LOAD GRAMMAR %s %s", grammar, path))
	return err
}

// SpeechUnloadGrammar unloads a grammar
func (s *AgiSession) SpeechUnloadGrammar(grammar string) error {
	_, err := s.execute(fmt.Sprintf("SPEECH UNLOAD GRAMMAR %s", grammar))
	return err
}

// SpeechActivateGrammar activates a grammar
func (s *AgiSession) SpeechActivateGrammar(grammar string) error {
	_, err := s.execute(fmt.Sprintf("SPEECH ACTIVATE GRAMMAR %s", grammar))
	return err
}

// SpeechDeactivateGrammar deactivates a grammar
func (s *AgiSession) SpeechDeactivateGrammar(grammar string) error {
	_, err := s.execute(fmt.Sprintf("SPEECH DEACTIVATE GRAMMAR %s", grammar))
	return err
}

// SpeechRecognize performs speech recognition
func (s *AgiSession) SpeechRecognize(r *MRCPRecog) error {
	if r == nil {
		return fmt.Errorf("recognition request cannot be nil")
	}

	cmd := strings.Builder{}
	cmd.WriteString("SPEECH RECOGNIZE ")

	if r.Grammar != "" {
		cmd.WriteString(r.Grammar)
		cmd.WriteString(" ")
	}

	if r.Timeout > 0 {
		cmd.WriteString(fmt.Sprintf("%d ", r.Timeout))
	}

	if r.Options != nil {
		for k, v := range r.Options {
			cmd.WriteString(fmt.Sprintf("%s=%s ", k, v))
		}
	}

	if r.ResultVar != "" {
		cmd.WriteString(fmt.Sprintf("'%s' ", r.ResultVar))
	}

	if r.CompletionC != "" {
		cmd.WriteString(fmt.Sprintf("'%s'", r.CompletionC))
	}

	_, err := s.execute(strings.TrimSpace(cmd.String()))
	return err
}

// SpeechSet sets a speech engine setting
func (s *AgiSession) SpeechSet(name, value string) error {
	_, err := s.execute(fmt.Sprintf("SPEECH SET %s %s", name, value))
	return err
}

// SpeechSpeak performs speech synthesis
func (s *AgiSession) SpeechSpeak(synth *MRCPSynth) error {
	if synth == nil {
		return fmt.Errorf("synthesis request cannot be nil")
	}

	cmd := strings.Builder{}
	cmd.WriteString("SPEECH SYNTHESIZE ")
	cmd.WriteString(fmt.Sprintf("'%s' ", synth.Text))

	if synth.Options != nil {
		for k, v := range synth.Options {
			cmd.WriteString(fmt.Sprintf("%s=%s ", k, v))
		}
	}

	if synth.ResultVar != "" {
		cmd.WriteString(fmt.Sprintf("'%s'", synth.ResultVar))
	}

	_, err := s.execute(strings.TrimSpace(cmd.String()))
	return err
}
