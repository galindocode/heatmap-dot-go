package heatmap

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()

	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}

	if builder.config == nil {
		t.Fatal("Builder config is nil")
	}

	// Check defaults
	if builder.config.Width != 800 {
		t.Errorf("Expected default width 800, got %d", builder.config.Width)
	}

	if builder.config.Height != 600 {
		t.Errorf("Expected default height 600, got %d", builder.config.Height)
	}
}

func TestBuilderSize(t *testing.T) {
	hm, err := NewBuilder().
		Size(1920, 1080).
		AddPoint(100, 100, 1.0).
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if hm.config.Width != 1920 || hm.config.Height != 1080 {
		t.Errorf("Expected size 1920x1080, got %dx%d", hm.config.Width, hm.config.Height)
	}
}

func TestBuilderMaxValue(t *testing.T) {
	hm, err := NewBuilder().
		Size(800, 600).
		MaxValue(100).
		AddPoint(100, 100, 50).
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if hm.config.MaxValue == nil || *hm.config.MaxValue != 100 {
		t.Error("MaxValue not set correctly")
	}
}

func TestBuilderMustBuild_Success(t *testing.T) {
	hm := NewBuilder().
		Size(800, 600).
		AddPoint(100, 100, 1.0).
		MustBuild()

	if hm == nil {
		t.Error("MustBuild() returned nil")
	}
}
