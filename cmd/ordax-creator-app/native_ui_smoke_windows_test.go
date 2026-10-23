//go:build windows

package main

import (
    "os"
    "runtime"
    "testing"
)

// This smoke test creates the same native Win32 UI controls as production,
// but NEVER starts device discovery, refresh, updates, UAC or a writer.
// It exercises Windows DLL imports and actual control creation on a Windows runner.
func TestWindowsNativeUIInitializesWithoutAnyPhysicalIO(t *testing.T) {
    if os.Getenv("ORDAX_CREATOR_WINDOWS_UI_SMOKE") != "1" {
        t.Skip("native Win32 smoke runs only in explicitly guarded Windows CI")
    }
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()
    setCreatorLocale(string(creatorSourceLocale))
    createMainWindow(false)
    if mainWindow == 0 {
        t.Fatal("Win32 main window was not created")
    }
    defer func() {
        hwnd := mainWindow
        if hwnd != 0 {
            destroyed, _, err := procDestroyWindow.Call(hwnd)
            if destroyed == 0 {
                t.Errorf("DestroyWindow failed: %v", err)
            }
        }
        mainWindow = 0
    }()

    for name, handle := range map[string]uintptr{
        "locale": localeCombo,
        "device": deviceCombo,
        "update": updateButton,
        "refresh": refreshButton,
        "write": writeButton,
        "boot-help": bootHelpButton,
    } {
        if handle == 0 {
            t.Errorf("missing native %s control", name)
        }
    }
    if writeButton != 0 {
        enabled, _, _ := user32.NewProc("IsWindowEnabled").Call(writeButton)
        if enabled != 0 {
            t.Fatal("physical write button enabled without any signed physical candidate")
        }
    }
    stateMu.Lock()
    state := refreshState
    stateMu.Unlock()
    if state.PhysicalReady || state.BackendDirectory != "" || len(state.Targets) > 0 || state.PhysicalCandidate != nil {
        t.Fatal("Win32 startup smoke unexpectedly acquired a physical candidate or target")
    }
    t.Log("Win32 controls created with physical write disabled; no enumeration or update started")
}
