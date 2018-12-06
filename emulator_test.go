package x86_emulator

import (
	"testing"
	"bytes"
	"io"
	"os"
	"io/ioutil"
	)

type machineCode []byte

func rawHeader() machineCode {
	// 32 bytes + 3 bytes (since initIP is set as 0x0003)
	return []byte{
		0x4d, 0x5a, 0x2b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0x01, 0xff, 0xff, 0x01, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x03, 0x00, 0x02, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
	}
}

func TestParseHeaderSignature(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := [2]byte{'M', 'Z'}
	if actual.exSignature != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exSignature)
	}
}

func TestParseHeaderSize(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(2)
	if actual.exHeaderSize != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exHeaderSize)
	}
}

func TestParseHeaderInitSS(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(0x0001)
	if actual.exInitSS != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitSS)
	}
}

func TestParseHeaderInitSP(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(0x1000)
	if actual.exInitSP != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitSP)
	}
}

func TestParseHeaderInitIP(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(0x0003)
	if actual.exInitIP != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitIP)
	}
}

func TestParseHeaderInitCS(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(0x0002)
	if actual.exInitCS != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitCS)
	}
}

// relocation

func rawHeaderWithRelocation() machineCode {
	// 48 bytes
	return []byte{
		0x4d, 0x5a, 0x4f, 0x00, 0x01, 0x00, 0x01, 0x00, 0x03, 0x00, 0x01, 0x01, 0xff, 0xff, 0x02, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func TestParseHeaderRelocationItems(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeaderWithRelocation())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(1)
	if actual.relocationItems != expected {
		t.Errorf("expected %v but actual %v", expected, actual.relocationItems)
	}
}

func TestParseHeaderRelocationOffset(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeaderWithRelocation())
	actual, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := word(0x0020)
	if actual.relocationTableOffset != expected {
		t.Errorf("expected %v but actual %v", expected, actual.relocationTableOffset)
	}
}

// intialize

func rawHeaderForTestInitilization() []byte {
	return []byte{
		0x4d, 0x5a, 0x71, 0x00, 0x01, 0x00, 0x01, 0x00, 0x03, 0x00, 0x01, 0x01, 0xff, 0xff, 0x05, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x0c, 0x00, 0x03, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func TestInitialization(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeaderForTestInitilization())
	header, _, err := parseHeader(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	intHandlers := make(intHandlers)
	state := newState(header, intHandlers)

	// check CS
	expectedCS := word(0x0003)
	if state.cs != expectedCS {
		t.Errorf("expected %v but actual %v", expectedCS, state.cs)
	}

	// check IP
	expectedIP := word(0x000c)
	if state.ip != expectedIP {
		t.Errorf("expected %v but actual %v", expectedIP, state.ip)
	}

	// check SS
	expectedSS := word(0x0005)
	if state.ss != expectedSS {
		t.Errorf("expected %v but actual %v", expectedSS, state.ss)
	}

	// check SP
	expectedSP := word(0x1000)
	if state.sp != expectedSP {
		t.Errorf("expected %v but actual %v", expectedSP, state.sp)
	}
}

// decode

func TestDecodeInstInt(t *testing.T) {
	// int 21
	var reader io.Reader = bytes.NewReader([]byte{0xcd, 0x21})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instInt{operand: 0x21}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeMovAX(t *testing.T) {
	// mov ax,1
	var reader io.Reader = bytes.NewReader([]byte{0xb8, 0x01, 0x00})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instMov{dest: AX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeMovCX(t *testing.T) {
	// mov cx,1
	var reader io.Reader = bytes.NewReader([]byte{0xb9, 0x01, 0x00})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instMov{dest: CX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeShlAX(t *testing.T) {
	// shl ax,1
	var reader io.Reader = bytes.NewReader([]byte{0xc1, 0xe0, 0x01})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instShl{register: AX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeShlCX(t *testing.T) {
	// shl cx,1
	var reader io.Reader = bytes.NewReader([]byte{0xc1, 0xe1, 0x01})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instShl{register: CX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeAddAX(t *testing.T) {
	// add ax,1
	var reader io.Reader = bytes.NewReader([]byte{0x83, 0xc0, 0x01})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instAdd{dest: AX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeAddCX(t *testing.T) {
	// add ax,1
	var reader io.Reader = bytes.NewReader([]byte{0x83, 0xc1, 0x01})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instAdd{dest: CX, imm: 0x0001}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeSubAXImm(t *testing.T) {
	// sub ax,2
	var reader io.Reader = bytes.NewReader([]byte{0x83, 0xec, 0x02})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instSub{dest: SP, imm: 0x0002}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeMovDs(t *testing.T) {
	// mov ds,ax
	var reader io.Reader = bytes.NewReader([]byte{0x8e, 0xd8})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instMovSRegReg{dest: DS, src: AX}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeMovAh(t *testing.T) {
	// mov ah,09h
	var reader io.Reader = bytes.NewReader([]byte{0xb4, 0x09})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instMovB{dest: AH, imm: 0x09}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeLeaDx(t *testing.T) {
	// lea dx,msg
	var reader io.Reader = bytes.NewReader([]byte{0x8d, 0x16, 0x02, 0x00}) // 0b00010110
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instLea{dest: DX, address: 0x0002}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodePushGeneralRegisters(t *testing.T) {
	// push ax, cx, dx, bx, sp, bp, si, di
	var readers = []io.Reader{
		bytes.NewReader([]byte{0x50}),
		bytes.NewReader([]byte{0x51}),
		bytes.NewReader([]byte{0x52}),
		bytes.NewReader([]byte{0x53}),
		bytes.NewReader([]byte{0x54}),
		bytes.NewReader([]byte{0x55}),
		bytes.NewReader([]byte{0x56}),
		bytes.NewReader([]byte{0x57}),
	}
	var expected = []instPush{
		instPush{src: AX},
		instPush{src: CX},
		instPush{src: DX},
		instPush{src: BX},
		instPush{src: SP},
		instPush{src: BP},
		instPush{src: SI},
		instPush{src: DI},
	}

	for i := 0; i < len(readers); i++ {
		actual, _, err := decodeInst(readers[i])
		if err != nil {
			t.Errorf("%+v", err)
		}
		if actual != expected[i] {
			t.Errorf("expected %v but actual %v", expected[i], actual)
		}
	}
}

func TestDecodePopGeneralRegisters(t *testing.T) {
	// pop ax, cx, dx, bx, sp, bp, si, di
	var readers = []io.Reader{
		bytes.NewReader([]byte{0x58}),
		bytes.NewReader([]byte{0x59}),
		bytes.NewReader([]byte{0x5a}),
		bytes.NewReader([]byte{0x5b}),
		bytes.NewReader([]byte{0x5c}),
		bytes.NewReader([]byte{0x5d}),
		bytes.NewReader([]byte{0x5e}),
		bytes.NewReader([]byte{0x5f}),
	}
	var expected = []instPop{
		instPop{dest: AX},
		instPop{dest: CX},
		instPop{dest: DX},
		instPop{dest: BX},
		instPop{dest: SP},
		instPop{dest: BP},
		instPop{dest: SI},
		instPop{dest: DI},
	}

	for i := 0; i < len(readers); i++ {
		actual, _, err := decodeInst(readers[i])
		if err != nil {
			t.Errorf("%+v", err)
		}
		if actual != expected[i] {
			t.Errorf("expected %v but actual %v", expected[i], actual)
		}
	}
}

func TestDecodeCall(t *testing.T) {
	// call rel16
	var reader io.Reader = bytes.NewReader([]byte{0xe8, 0xdc, 0xff})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instCall{rel: -36} // 0xffdc
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeRet(t *testing.T) {
	// ret (near return)
	var reader io.Reader = bytes.NewReader([]byte{0xc3})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instRet{}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

func TestDecodeMovWithDisp(t *testing.T) {
	// mov ax,[bp+4]
	var reader io.Reader = bytes.NewReader([]byte{0x8b, 0x46, 0x04})
	actual, _, err := decodeInst(reader)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expected := instMovRegMemBP{dest: AX, disp: 4}
	if actual != expected {
		t.Errorf("expected %v but actual %v", expected, actual)
	}
}

// run

func (code machineCode) withMov() machineCode {
	// mov ax,1
	mov := []byte{0xb8, 0x01, 0x00}
	return append(code, mov...)
}

func (code machineCode) withInt21_4c() machineCode {
	// mov ah, 4c
	// int 21h
	int21 := []byte{0xb4, 0x4c, 0xcd, 0x21}
	return append(code, int21...)
}

func rawHeaderForRunExe() machineCode {
	return []byte{
		0x4d, 0x5a, 0x2b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0x01, 0xff, 0xff, 0x01, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}


func TestRunExe(t *testing.T) {
	_, _, err := RunExe(bytes.NewReader(rawHeaderForRunExe().withMov().withInt21_4c()))
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestInt21_4c_ax(t *testing.T) {
	// exit status 1
	b := rawHeaderForRunExe()
	b = append(b, []byte{0xb8, 0x4c, 0x00}...) // mov ax,4ch
	b = append(b, []byte{0xc1, 0xe0, 0x08}...) // shl ax,8
	b = append(b, []byte{0x83, 0xc0, 0x01}...) // add ax,01h
	b = append(b, []byte{0xcd, 0x21}...)       // int 21h

	intHandlers := make(intHandlers)

	actual, err := runExeWithCustomIntHandlers(bytes.NewReader(b), intHandlers)
	if err != nil {
		t.Errorf("%+v", err)
	}

	if actual.ax != 0x4c01 {
		t.Errorf("ax is illegal: %04x", actual.ax)
	}
	if actual.al() != 0x01 {
		t.Errorf("al is illegal: %02x", actual.al())
	}
	if actual.ah() != 0x4c {
		t.Errorf("ah is illegal: %02x", actual.ah())
	}
	if actual.exitCode != 0x01 {
		t.Errorf("exitCode is expected to be %02x but %02x", 0x01, actual.exitCode)
	}
}

func TestInt21_4c_cx(t *testing.T) {
	// exit status 1
	b := rawHeaderForRunExe()
	b = append(b, []byte{0xb9, 0x4c, 0x00}...) // mov cx,4ch
	b = append(b, []byte{0xc1, 0xe1, 0x08}...) // shl cx,8
	b = append(b, []byte{0x83, 0xc1, 0x01}...) // add cx,01h
	b = append(b, []byte{0x8b, 0xc1}...)       // mov ax,cx
	b = append(b, []byte{0xcd, 0x21}...)       // int 21h

	intHandlers := make(intHandlers)

	actual, err := runExeWithCustomIntHandlers(bytes.NewReader(b), intHandlers)
	if err != nil {
		t.Errorf("%+v", err)
	}

	if actual.ax != 0x4c01 {
		t.Errorf("ax is illegal: %04x", actual.ax)
	}
	if actual.al() != 0x01 {
		t.Errorf("al is illegal: %02x", actual.al())
	}
	if actual.ah() != 0x4c {
		t.Errorf("ah is illegal: %02x", actual.ah())
	}
	if actual.exitCode != 0x01 {
		t.Errorf("exitCode is expected to be %02x but %02x", 0x01, actual.exitCode)
	}
}

// for print systemcall

func rawHeaderForTestInt21_09() machineCode {
	// assume that there is one relocation
	// (1) relocation items
	// (2) item of relocation table
	return []byte{
	//                                      <--(1)--->
		0x4d, 0x5a, 0x4f, 0x00, 0x01, 0x00, 0x01, 0x00, 0x03, 0x00, 0x01, 0x01, 0xff, 0xff, 0x02, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	//  <--(2)--->
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func TestInt21_09(t *testing.T) {
	b := rawHeaderForTestInt21_09()
	b = append(b, []byte{0xb8, 0x01, 0x00}...) // mov ax,seg msg
	b = append(b, []byte{0x8e, 0xd8}...) // mov ds,ax
	b = append(b, []byte{0xb4, 0x09}...) // mov ah,09h
	b = append(b, []byte{0x8d, 0x16, 0x02, 0x00}...) // lea dx,msg
	b = append(b, []byte{0xcd, 0x21}...) // int 21h
	b = append(b, []byte{0xb8, 0x00, 0x4c}...) // mov ax,4c00h
	b = append(b, []byte{0xcd, 0x21}...) // int 21h
	b = append(b, []byte("Hello world!$")...)

	tempFile, err := ioutil.TempFile("", "TestInt21_09")
	if err != nil {
		t.Errorf("%+v", err)
	}
	defer os.Remove(tempFile.Name())

	intHandlers := make(intHandlers)
	intHandlers[0x09] = func(s *state, m *memory) error {
		originalStdout := os.Stdout
		os.Stdout = tempFile
		err = intHandler09(s, m)
		os.Stdout = originalStdout
		return err
	}

	_, err = runExeWithCustomIntHandlers(bytes.NewReader(b), intHandlers)
	if err != nil {
		t.Errorf("%+v", err)
	}

	tempFile.Sync()
	tempFile.Seek(0, 0)
	output, err := ioutil.ReadAll(tempFile)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(output) != "Hello world!" {
		t.Errorf("expect output \"%s\" but \"%s\"", "Hello world!", string(output))
	}

	if err = tempFile.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}

func rawHeaderForTestPush() machineCode {
	return []byte{
		0x4d, 0x5a, 0x2b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0x01, 0xff, 0xff, 0x01, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func TestPushAndPop(t *testing.T) {
	b := rawHeaderForTestPush()
	b = append(b, []byte{0xb8, 0x35, 0x10}...) // mov ax, 0x1035
	b = append(b, []byte{0xb9, 0x36, 0x20}...) // mov cx, 0x2036
	b = append(b, []byte{0x50}...) // push ax
	b = append(b, []byte{0x51}...) // push cx
	b = append(b, []byte{0x5b}...) // pop bx
	b = append(b, []byte{0x5a}...) // pop dx
	b = append(b, []byte{0xb8, 0x00, 0x4c}...) // mov ax,4c00h
	b = append(b, []byte{0xcd, 0x21}...) // int 21h

	intHandlers := make(intHandlers)

	actual, err := runExeWithCustomIntHandlers(bytes.NewReader(b), intHandlers)
	if err != nil {
		t.Errorf("%+v", err)
	}

	if actual.dx != 0x1035 {
		t.Errorf("expect 0x%04x but 0x%04x", 0x1035, actual.dx)
	}
	if actual.bx != 0x2036 {
		t.Errorf("expect 0x%04x but 0x%04x", 0x2036, actual.bx)
	}
}

// RunExe with sample file

func TestRunExeWithSampleFcall(t *testing.T) {
	file, err := os.Open("sample/fcall.exe")
	if err != nil {
		t.Errorf("%+v", err)
	}
	exitCode, _, err := RunExe(file)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if exitCode != 7 {
		t.Errorf("expect exitCode to be %d but actual %d", 7, exitCode)
	}
}

func TestRunExeWithSampleFcallp(t *testing.T) {
	file, err := os.Open("sample/fcallp.exe")
	if err != nil {
		t.Errorf("%+v", err)
	}
	exitCode, _, err := RunExe(file)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if exitCode != 7 {
		t.Errorf("expect exitCode to be %d but actual %d", 7, exitCode)
	}
}

func TestRunExeWithSampleCmain(t *testing.T) {
	file, err := os.Open("sample/cmain.exe")
	if err != nil {
		t.Errorf("%+v", err)
	}
	exitCode, _, err := RunExe(file)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if exitCode != 0 {
		t.Errorf("expect exitCode to be %d but actual %d", 0, exitCode)
	}
}