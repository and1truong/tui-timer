package sound

import (
	"context"
	"fmt"
	"os/exec"
)

// Player is the interface for playing sounds.
type Player interface {
	PlayBeep(ctx context.Context) error
	PlayVoice(ctx context.Context, voice, message string) error
	PlayFile(ctx context.Context, path string) error
}

// MacPlayer implements Player using macOS commands.
type MacPlayer struct{}

func NewMacPlayer() *MacPlayer {
	return &MacPlayer{}
}

func (p *MacPlayer) PlayBeep(_ context.Context) error {
	// Print BEL character for system beep
	fmt.Print("\a")
	return nil
}

func (p *MacPlayer) PlayVoice(ctx context.Context, voice, message string) error {
	cmd := exec.CommandContext(ctx, "say", "-v", voice, message)
	return cmd.Run()
}

func (p *MacPlayer) PlayFile(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "afplay", path)
	return cmd.Run()
}

// NoopPlayer does nothing (for testing or when sound is disabled).
type NoopPlayer struct{}

func (p *NoopPlayer) PlayBeep(_ context.Context) error                     { return nil }
func (p *NoopPlayer) PlayVoice(_ context.Context, _, _ string) error       { return nil }
func (p *NoopPlayer) PlayFile(_ context.Context, _ string) error           { return nil }
