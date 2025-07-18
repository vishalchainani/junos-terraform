package netconf

import (
	"time"
)

// DriverJunos type is for creating a Junos based driver. Maintains state for session and connection. Implements Driver{}
type DriverJunos struct {
	Timeout   time.Duration // Timeout for SSH timed sessions
	Datastore string        // NETCONF datastore
	Session   *Session      // Session data
}

// New creates a new instance of DriverJunos
func New() *DriverJunos {
	return &DriverJunos{}
}

// SetDatastore sets the target datastore on the data structure
func (d *DriverJunos) SetDatastore(ds string) error {
	d.Datastore = ds
	return nil
}

// Dial function (call this after New())
func (d *DriverJunos) Dial() error {

	var err error

	d.Session, err = lowlevelDial()

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

// DialTimeout NOT IMPLEMENTED. This driver is transactional based and not required.
func (d *DriverJunos) DialTimeout() error {

	return nil
}

// Close function closes the socket
func (d *DriverJunos) Close() error {
	// Close the SSH Session if we have one}

	err := d.Session.Close()

	if err != nil {
		return err
	}

	return nil
}

// Lock the target datastore
func (d *DriverJunos) Lock(ds string) (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodLock(ds))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// Unlock the target datastore
func (d *DriverJunos) Unlock(ds string) (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodUnlock(ds))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// SendRaw sends a raw XML envelope
func (d *DriverJunos) SendRaw(rawxml string) (*RPCReply, error) {
	reply, err := d.Session.Exec(RawMethod(rawxml))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// GetConfig requests the contents of a datastore
func (d *DriverJunos) GetConfig() (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodGetConfig(d.Datastore))

	if err != nil {
		return reply, err
	}

	return reply, nil
}
