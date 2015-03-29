package goTunTap

/*
#include"tapUtils.h"
*/
import "C"
import (
	"os"
	"errors"
	"syscall"
	"sync"
	"io"
)

type TapManager struct {
	sock C.int
	name string
	up bool
}

func CheckBridge(bridge string) error {
	if(C.CheckBridge(C.CString(bridge)) < 0) {
		return errors.New("Bridge does not exist")
	}
	return nil
}

func CreateBridge(name string) error {
	if CheckBridge(name) == nil {
		return nil
	}
	if(C.CreateBridge(C.CString(name)) < 0) {
		return errors.New("Failed to create bridge")
	}
	return nil
}

func DeleteBridge(name string) error {
	if(C.DeleteBridge(C.CString(name)) < 0) {
		return errors.New("Failed to delete bridge")
	}
	return nil
}

func AddTapToBridge(bridge, tap string) error {
	if(C.AddTapToBridge(C.CString(bridge), C.CString(tap)) < 0) {
		return errors.New("Failed to add tap to bridge")
	}
	return nil
}

func RemoveTapFromBridge(bridge, tap string) error {
	if(C.RemoveTapFromBridge(C.CString(bridge), C.CString(tap)) < 0) {
		return errors.New("Failed to remove tap from bridge")
	}
	return nil
}

func CreateTap(name string) (*TapManager, error) {
	tap := &TapManager{0, name, false};
	if os.Geteuid() != 0 {
		return nil, errors.New("Must execute as ROOT")
	}
	return tap, tap.Start()
}

func (t *TapManager) Start() error {
	sock := C.StartTap(C.CString(t.name))
	if sock <= 0 {
		return errors.New("Failed to create tap")
	}
	t.sock = sock
	return nil
}

func (t *TapManager) Stop() error {
	if C.StopTap(t.sock, C.CString(t.name)) != 0 {
		return errors.New("Failed to destroy tap")
	}
	return nil
}

func (t *TapManager) Read(b []byte) (int, error) {
	var n int
	var err error
	if t == nil {
		return 0, os.ErrInvalid
	}
	if len(b) > 0 {
		n, err = syscall.Read(int(t.sock), b)
		if n == 0 && err == nil {
			return n, io.EOF
		}
		if err != nil {
			return 0, err
		}
	}
	return n, err
}


func (t *TapManager) Write(b []byte) (int, error) {
	var n int
	var err error
	if t == nil {
		return 0, os.ErrInvalid
	}
	if len(b) > 0 {
		n, err = syscall.Write(int(t.sock), b)
		if n == 0 && err == nil {
			return n, io.EOF
		}
		if err != nil {
			return 0, err
		}
	}
	return n, err
}

//Relay relays data from a reader/Writer.  Like a net.Conn
//to and from the tap
func (t* TapManager) Relay(c io.ReadWriter) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	//relay from outside conn to the tap
	go func(rdr io.Reader, wtr io.Writer, wg *sync.WaitGroup) {
		defer wg.Done()
		io.Copy(wtr, rdr)
		
	}(c, t, &wg)

	//relay from tap to outside
	go func(rdr io.Reader, wtr io.Writer, wg *sync.WaitGroup) {
		defer wg.Done()
		io.Copy(wtr, rdr)
	}(t, c, &wg)
	wg.Wait()
	return nil 
}

func (t* TapManager) AddToBridge(bridge string) error {
	err := CheckBridge(bridge)
	if err != nil {
		if C.CreateBridge(C.CString(bridge)) != 0 {
			return errors.New("Failed to create bridge")
		}
	}
	if C.AddTapToBridge(C.CString(bridge), C.CString(t.name)) != 0 {
		return errors.New("Failed to add tap to bridge")
	}
	return nil
}

func (t* TapManager) RemoveFromBridge(bridge string) error {
	if C.RemoveTapFromBridge(C.CString(bridge), C.CString(t.name)) != 0 {
		return errors.New("Failed to add tap to bridge")
	}
	return nil
}
