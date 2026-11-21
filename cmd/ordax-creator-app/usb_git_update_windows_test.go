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

func TestOwnerUSBTrustedHotplugRequiresUniqueAuthorizedSerial(t *testing.T) {
    home := t.TempDir()
    creatorRepo := filepath.Join(home, "GitHub", "ordax-creator")
    osRepo := filepath.Join(home, "GitHub", "ordax-os")
    tools := filepath.Join(creatorRepo, "tools")
    if err := os.MkdirAll(tools, 0700); err != nil { t.Fatal(err) }
    if err := os.MkdirAll(filepath.Join(osRepo, ".git"), 0700); err != nil { t.Fatal(err) }
    python := filepath.Join(home, "python.exe")
    script := filepath.Join(tools, "owner_local_sync.py")
    deliver := filepath.Join(tools, "owner_usb_git_delivery.py")
    for _, path := range []string{python, script, deliver} {
        if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil { t.Fatal(err) }
    }
    install := filepath.Join(home, "install", "OrdaX-Creator-Owner-Prototype")
    executable := filepath.Join(install, "current", "OrdaX-Creator", "OrdaX-Creator.exe")
    task := ownerUSBTaskInfo{
        Executable: python,
        Arguments: "\"" + script + "\" --install-root \"" + install +
            "\" --stage-owner-dev-usb-serial 6699E36B",
        State: "Ready", Concurrency: "IgnoreNew",
    }
    state := appRefreshState{PhysicalReady: true, Targets: []physicalTarget{testOwnerUSBTarget()}}
    target, ok := ownerUSBTrustedTargetFromTask(state, task, executable)
    if !ok || target.DeviceSerial != "6699E36B" {
        t.Fatalf("previously authorized Owner USB not found: %+v %v", target, ok)
    }
    state.Targets = append(state.Targets, testOwnerUSBTarget())
    if _, ok := ownerUSBTrustedTargetFromTask(state, task, executable); ok {
        t.Fatal("ambiguous physical USB was auto-staged")
    }
    state.Targets = state.Targets[:1]
    task.Arguments = strings.Replace(task.Arguments, "6699E36B", "UNTRUSTED", 1)
    if _, ok := ownerUSBTrustedTargetFromTask(state, task, executable); ok {
        t.Fatal("USB without previous local serial authorization was auto-staged")
    }
    state.PhysicalReady = false
    if _, ok := ownerUSBTrustedTargetFromTask(state, task, executable); ok {
        t.Fatal("unready physical device backend authorized hotplug staging")
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

// A staged source bundle must never be presented as an activated system update.
func TestOwnerUSBComponentSummaryDoesNotClaimExecutableInstallation(t *testing.T) {
    for _, result := range []string{"CURRENT", "STAGED_NOT_ACTIVATED"} {
        view := ownerUSBComponentSummary(ownerUSBUpdateStatus{
            Result: result, SourceCommit: strings.Repeat("a", 40),
        })
        for _, expected := range []string{"Git", "Surface", "Base", "aaaaaaaaaaaa"} {
            if !strings.Contains(view, expected) {
                t.Fatalf("missing %q in component evidence %q", expected, view)
            }
        }
    }
    for _, input := range []ownerUSBUpdateStatus{
        {Result: "AVAILABLE", SourceCommit: strings.Repeat("a", 40)},
        {Result: "CURRENT", SourceCommit: "invalid"},
        {Result: ""},
    } {
        if got := ownerUSBComponentSummary(input); got != "" {
            t.Fatalf("unverified transport generated evidence: %q", got)
        }
    }
}

func TestOwnerUSBStagingPresentationNeverClaimsOSActivated(t *testing.T) {
    for _, locale := range []creatorLocale{creatorLocalePTBR, creatorLocaleENUS} {
        current := strings.ToLower(creatorMessageFor(locale, msgUSBUpdateCurrent, nil))
        staged := strings.ToLower(creatorMessageFor(locale, msgUSBUpdateComplete, nil))
        if locale == creatorLocalePTBR {
            if !strings.Contains(current, "não verificado") || !strings.Contains(staged, "não atualizado") {
                t.Fatalf("Portuguese Git staging must not be reported as OS activation: %q / %q", current, staged)
            }
        } else if !strings.Contains(current, "unverified") || !strings.Contains(staged, "not updated") {
            t.Fatalf("English Git staging must not be reported as OS activation: %q / %q", current, staged)
        }
    }
}

func TestOwnerUSBGitStageMatchesPinnedOSUpdateCoverage(t *testing.T) {
    root := os.Getenv("ORDAX_OS_CONTRACTS_DIR")
    if root == "" {
        t.Skip("requires the pinned read-only OS contract checkout")
    }
    data, err := os.ReadFile(filepath.Join(root, "device-update-coverage.json"))
    if err != nil { t.Fatal(err) }
    var contract struct {
        Schema string `json:"$schema"`
        Layers []struct {
            ID string `json:"id"`
            Delivery string `json:"development_delivery"`
            ActivationImplemented *bool `json:"activation_implemented"`
            RuntimeOneShot *bool `json:"runtime_one_shot_activation_wiring_implemented"`
            SelfUpdateGap *bool `json:"self_update_gap"`
        } `json:"layers"`
    }
    if err := json.Unmarshal(data, &contract); err != nil { t.Fatal(err) }
    if contract.Schema != "prototype-ordax.device-update-coverage/1" {
        t.Fatalf("unsupported OS coverage authority %q", contract.Schema)
    }
    layers := map[string]struct {
        ID string
        Delivery string
        ActivationImplemented *bool
        RuntimeOneShot *bool
        SelfUpdateGap *bool
    }{}
    for _, layer := range contract.Layers {
        if _, found := layers[layer.ID]; found { t.Fatalf("duplicate OS layer %q", layer.ID) }
        layers[layer.ID] = layer
    }
    rootfs, ok := layers["development-rootfs"]
    if !ok || rootfs.Delivery != "exact-commit-development-base-rootfs-candidate" ||
        rootfs.ActivationImplemented == nil || *rootfs.ActivationImplemented ||
        rootfs.RuntimeOneShot == nil || *rootfs.RuntimeOneShot {
        t.Fatal("Creator Git-only staging cannot claim Base activation while OS rootfs activation is unavailable")
    }
    kernel, ok := layers["kernel-and-initramfs"]
    if !ok || kernel.RuntimeOneShot == nil || *kernel.RuntimeOneShot {
        t.Fatal("Creator Git staging cannot claim kernel/initramfs activation")
    }
    loader, ok := layers["uefi-fallback-loader"]
    if !ok || loader.SelfUpdateGap == nil || !*loader.SelfUpdateGap {
        t.Fatal("Creator Git staging cannot claim UEFI loader self-update")
    }
    for _, locale := range []creatorLocale{creatorLocalePTBR, creatorLocaleENUS} {
        action := strings.ToLower(creatorMessageFor(locale, msgUSBUpdateAction, nil))
        if !strings.Contains(action, "git") {
            t.Fatalf("Git-only Creator action must identify its real scope: %q", action)
        }
    }
}
