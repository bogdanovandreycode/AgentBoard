package httpapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// Method indices are fixed by the IFileOpenDialog/IFileDialog/IShellItem COM
// interfaces in the Windows SDK (shobjidl_core.h).
const (
	comRelease              = 2
	dialogShow              = 3
	dialogSetOptions        = 9
	dialogGetOptions        = 10
	dialogSetFolder         = 12
	dialogSetTitle          = 17
	dialogGetResult         = 20
	shellItemGetDisplayName = 5
)

type shellCOMObject struct{ methods *[29]uintptr }

func (object *shellCOMObject) call(method int, args ...uintptr) uintptr {
	arguments := append([]uintptr{uintptr(unsafe.Pointer(object))}, args...)
	result, _, _ := syscall.SyscallN(object.methods[method], arguments...)
	runtime.KeepAlive(object)
	return result
}

func shellResult(operation string, result uintptr) error {
	if int32(uint32(result)) < 0 {
		return fmt.Errorf("%s failed (HRESULT 0x%08X)", operation, uint32(result))
	}
	return nil
}

var folderDialogMutex sync.Mutex

func selectWindowsFolder(ctx context.Context, initial string) (string, error) {
	folderDialogMutex.Lock()
	defer folderDialogMutex.Unlock()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	// Shell dialogs require an STA; all COM operations stay on this OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ole := syscall.NewLazyDLL("ole32.dll")
	initResult, _, _ := ole.NewProc("CoInitializeEx").Call(0, 2)
	if err := shellResult("CoInitializeEx", initResult); err != nil {
		return "", err
	}
	defer ole.NewProc("CoUninitialize").Call()
	classID := syscall.GUID{Data1: 0xDC1C5A9C, Data2: 0xE88A, Data3: 0x4DDE, Data4: [8]byte{0xA5, 0xA1, 0x60, 0xF8, 0x2A, 0x20, 0xAE, 0xF7}}
	interfaceID := syscall.GUID{Data1: 0xD57C7288, Data2: 0xD4AD, Data3: 0x4768, Data4: [8]byte{0xBE, 0x02, 0x9D, 0x96, 0x95, 0x32, 0xD9, 0x60}}
	var dialog *shellCOMObject
	result, _, _ := ole.NewProc("CoCreateInstance").Call(uintptr(unsafe.Pointer(&classID)), 0, 1, uintptr(unsafe.Pointer(&interfaceID)), uintptr(unsafe.Pointer(&dialog)))
	if err := shellResult("Create folder dialog", result); err != nil {
		return "", err
	}
	defer dialog.call(comRelease)
	var options uint32
	if err := shellResult("Read dialog options", dialog.call(dialogGetOptions, uintptr(unsafe.Pointer(&options)))); err != nil {
		return "", err
	}
	// FOS_PICKFOLDERS | FOS_FORCEFILESYSTEM | FOS_PATHMUSTEXIST | FOS_DONTADDTORECENT.
	if err := shellResult("Set folder mode", dialog.call(dialogSetOptions, uintptr(options|0x20|0x40|0x800|0x02000000))); err != nil {
		return "", err
	}
	title, _ := syscall.UTF16PtrFromString("Select project folder")
	if err := shellResult("Set dialog title", dialog.call(dialogSetTitle, uintptr(unsafe.Pointer(title)))); err != nil {
		return "", err
	}
	runtime.KeepAlive(title)
	if initial != "" {
		initialPath, err := syscall.UTF16PtrFromString(initial)
		if err != nil {
			return "", err
		}
		itemID := syscall.GUID{Data1: 0x43826D1E, Data2: 0xE718, Data3: 0x42EE, Data4: [8]byte{0xBC, 0x55, 0xA1, 0xE2, 0x61, 0xC3, 0x7B, 0xFE}}
		var folder *shellCOMObject
		result, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("SHCreateItemFromParsingName").Call(uintptr(unsafe.Pointer(initialPath)), 0, uintptr(unsafe.Pointer(&itemID)), uintptr(unsafe.Pointer(&folder)))
		runtime.KeepAlive(initialPath)
		if shellResult("Read initial folder", result) == nil {
			result = dialog.call(dialogSetFolder, uintptr(unsafe.Pointer(folder)))
			folder.call(comRelease)
			if err := shellResult("Set initial folder", result); err != nil {
				return "", err
			}
		}
	}
	// The HTTP server has no browser HWND. Do not borrow the foreground
	// window: it may belong to an unrelated application or another client.
	result = dialog.call(dialogShow, 0)
	if uint32(result) == 0x800704C7 {
		return "", nil
	} // ERROR_CANCELLED.
	if err := shellResult("Show folder dialog", result); err != nil {
		return "", err
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var item *shellCOMObject
	if err := shellResult("Read selected folder", dialog.call(dialogGetResult, uintptr(unsafe.Pointer(&item)))); err != nil {
		return "", err
	}
	defer item.call(comRelease)
	var name *uint16
	// SIGDN_FILESYSPATH includes Unicode and long filesystem paths.
	if err := shellResult("Read folder path", item.call(shellItemGetDisplayName, 0x80058000, uintptr(unsafe.Pointer(&name)))); err != nil {
		return "", err
	}
	defer ole.NewProc("CoTaskMemFree").Call(uintptr(unsafe.Pointer(name)))
	var units []uint16
	for index := 0; index < 32768; index++ {
		unit := *(*uint16)(unsafe.Add(unsafe.Pointer(name), uintptr(index)*2))
		if unit == 0 {
			break
		}
		units = append(units, unit)
	}
	path := syscall.UTF16ToString(units)
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		return "", fmt.Errorf("selected folder is unavailable")
	}
	return filepath.Clean(path), nil
}
