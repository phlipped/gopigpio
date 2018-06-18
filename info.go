package gopigpio

const (
  HW_REVISION = 17
)


func (p Pigpio) HardwareRevision() (uint32, error) {
  cmd := Cmd{
    ID: HW_REVISION,
    P1: 0,
    P2: 0,
  }

  result, _ := p.sendCmd(cmd) // This command cannot fail, so don't worry about the error returned by SnedCmdSimple()

  return uint32(result.res), nil
}
