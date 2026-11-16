//go:build windows

package main

import (
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func testOwnerUSBTarget() physicalTarget {
    return physicalTarget{
        DiskNumber: 1, DeviceSerial: "6699E36B", PhysicalDiskBytes: 8 << 30,
        ConfirmationToken: strings.Repeat("a", 64), PrototypeSafe: true,
    }
}

func TestOwnerUSBUpdateRequiresDeliberatePhysicalSelection(t *testing.T) {
    device := testOwnerUSBTarget()
    state := appRefreshState{PhysicalReady: true, Targets: []physicalTarget{device}}
    if _, err := ownerUSBSelectedTarget(state, -1); err == nil {
        t.Fatal("missing explicit selection was accepted")
    }
    if _, err := ownerUSBSelectedTarget(state, 1); err == nil {
        t.Fatal("out-of-range selection was accepted")
    }
    if got, err := ownerUSBSelectedTarget(state, 0); err != nil || got.DiskNumber != 1 {
        t.Fatalf("validated selected target unavailable: %+v %v", got, err)
    }
    device.SystemDisk = true
    state.Targets[0] = device
    if _, err := ownerUSBSelectedTarget(state, 0); err == nil {
        t.Fatal("system disk accepted")
    }
    device.SystemDisk = false
    device.ConfirmationToken = "bad"
    state.Targets[0] = device
    if _, err := ownerUSBSelectedTarget(state, 0); err == nil {
        t.Fatal("unverified physical selection accepted")
    }
    state.PhysicalReady = false
    if _, err := ownerUSBSelectedTarget(state, 0); err == nil {
        t.Fatal("backend unavailable but USB update allowed")
    }
}

func TestOwnerUSBTaskReusesCanonicalInstalledTransport(t *testing.T) {
    home := t.TempDir()
    github := filepath.Join(home, "GitHub")
    creator := filepath.Join(github, "ordax-creator")
    osRepo := filepath.Join(github, "ordax-os")
    tools := filepath.Join(creator, "tools")
    if err := os.MkdirAll(tools, 0700); err != nil { t.Fatal(err) }
    if err := os.MkdirAll(filepath.Join(osRepo, ".git"), 0700); err != nil { t.Fatal(err) }
    for _, p := range []string{filepath.Join(tools, "owner_local_sync.py"),
        filepath.Join(tools, "owner_usb_git_delivery.py"), filepath.Join(home, "python.exe")} {
        if err := os.WriteFile(p, []byte("test"), 0600); err != nil { t.Fatal(err) }
    }
    app := filepath.Join(home, "install", "OrdaX-Creator-Owner-Prototype", "current",
        "OrdaX-Creator", "OrdaX-Creator.exe")
    script := filepath.Join(tools, "owner_local_sync.py")
    install := filepath.Join(home, "install", "OrdaX-Creator-Owner-Prototype")
    task := ownerUSBTaskInfo{
        Executable: filepath.Join(home, "python.exe"),
        Arguments: `"` + script + `" --install-root "` + install + `" --stage-owner-dev-usb-serial 6699E36B`,
        State: "Ready", Concurrency: "IgnoreNew",
    }
    value, err := ownerUSBTransportFromTask(task, testOwnerUSBTarget(), app)
    if err != nil || value.Repo != osRepo || value.Script != filepath.Join(tools, "owner_usb_git_delivery.py") {
        t.Fatalf("canonical Owner Dev transport rejected: %+v %v", value, err)
    }
    task.Concurrency = "Parallel"
    if _, err := ownerUSBTransportFromTask(task, testOwnerUSBTarget(), app); err == nil {
        t.Fatal("parallel scheduled task accepted")
    }
    task.Concurrency = "IgnoreNew"
    task.Arguments = strings.Replace(task.Arguments, "6699E36B", "UNTRUSTED", 1)
    if _, err := ownerUSBTransportFromTask(task, testOwnerUSBTarget(), app); err == nil {
        t.Fatal("unselected physical serial was accepted")
    }
    task.Arguments = strings.Replace(task.Arguments, "UNTRUSTED", "6699E36B", 1)
    task.State = "Running"
    if _, err := ownerUSBTransportFromTask(task, testOwnerUSBTarget(), app); err == nil {
        t.Fatal("running scheduled task accepted")
    }
}

func TestOwnerUSBUpdateAcceptsOnlyVerifiedStagingReceipts(t *testing.T) {
    commit := strings.Repeat("a", 40)
    for _, status := range []string{"CURRENT", "STAGED_NOT_ACTIVATED"} {
        receipt := ownerUSBReceipt{Status: status, SourceCommit: commit}
        if status == "STAGED_NOT_ACTIVATED" {
            receipt.SHA256 = strings.Repeat("b", 64)
            receipt.Bytes = 28163491
        }
        data, err := json.Marshal(receipt)
        if err != nil { t.Fatal(err) }
        result, err := parseOwnerUSBReceipt(data)
        if err != nil || result.Status != status {
            t.Fatalf("expected verified %s receipt: %+v %v", status, result, err)
        }
    }
    cases := []string{
        "{}", `{"status":"SKIP_TRUSTED_USB_NOT_PRESENT","source_commit":"` + commit + `"}`,
        `{"status":"STAGED_NOT_ACTIVATED","source_commit":"` + commit + `","sha256":"bad","bytes":10}`,
        `{"status":"CURRENT","source_commit":"` + commit + `"} {"extra":true}`,
    }
    for _, text := range cases {
        if _, err := parseOwnerUSBReceipt([]byte(text)); err == nil {
            t.Fatalf("malformed or unauthorized receipt accepted: %q", text)
        }
    }
}
