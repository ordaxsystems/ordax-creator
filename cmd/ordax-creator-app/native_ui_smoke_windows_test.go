//go:build windows

package main

import (
    "os"
    "runtime"
    "testing"
    "unsafe"
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

    // Reproduce the stale-selection bug using synthetic in-memory state only.
    // The test window has no services and never enumerates or opens any disk.
    previousLocale := currentCreatorLocale()
    defer func() {
        refreshInProgress = false
        setCreatorLocale(string(previousLocale))
        stateMu.Lock()
        refreshState = appRefreshState{}
        stateMu.Unlock()
    }()
    fixture := appRefreshState{
        Targets: []physicalTarget{{DiskNumber: 9, PrototypeSafe: true}},
        PhysicalReady: true,
        BackendDirectory: "synthetic-only-not-an-installed-backend",
    }
    send(deviceCombo, cbAddString, 0, uintptr(unsafe.Pointer(utf16Ptr("Synthetic USB — no device"))))
    send(deviceCombo, cbSetCurSel, 0, 0)
    stateMu.Lock()
    refreshState = fixture
    stateMu.Unlock()
    renderSelectionExperience(fixture)
    isEnabled := func() bool {
        value, _, _ := user32.NewProc("IsWindowEnabled").Call(writeButton)
        return value != 0
    }
    if !isEnabled() {
        t.Fatal("synthetic pre-refresh review could not reproduce enabled creation control")
    }
    refreshInProgress = true
    enable(writeButton, false) // beginRefresh's synchronous UI lockout.
    setCreatorLocale(string(creatorLocaleENUS))
    renderCreatorLocaleUI()
    if isEnabled() {
        t.Fatal("changing locale during refresh reenabled a stale USB write action")
    }
    if creatorStatusText != creatorT(msgRefreshSearching) || creatorHintText != creatorT(msgRefreshHint) {
        t.Fatalf("changing locale replaced refresh notice with stale state: %q / %q", creatorStatusText, creatorHintText)
    }
    if _, _, err := selectedPhysicalTarget(); err == nil {
        t.Fatal("a physical target was selectable while the list was being refreshed")
    }
    // A queued terminal completion must not reenable a stale control.
    renderWriteDone()
    if isEnabled() {
        t.Fatal("queued write completion reenabled a target during refresh")
    }
    t.Log("Locale change and queued completion kept write unavailable during synthetic refresh")

}
