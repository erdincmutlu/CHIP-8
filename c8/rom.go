package c8

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"unicode"
)

const (
	memorySize          = 0x1000
	programCounterStart = 0x200
	screenMemoryStart   = 0x100
	clearScreen         = "\033[H\033[2J"
)

var memory = make([]byte, memorySize)

type registerStruct struct {
	v           []byte
	index       uint16 // Memory address, 16 bit register
	progCounter uint16 // Instruction pointer
	delayTimer  byte
	soundTimer  byte
}

// Regs are registers
var regs = registerStruct{
	v:           make([]byte, 16),
	progCounter: programCounterStart,
}

type stackPtr struct {
	v           []byte
	index       uint16
	progCounter uint16
}

var stack []stackPtr
var keyPressed byte

// ReadROM will read the ROM
func ReadROM(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fmt.Printf("File is %d bytes\n", fileInfo.Size())

	reader := bufio.NewReader(file)
	_, err = reader.Read(memory[programCounterStart:])
	initSprites()
	return err
}

// RunROM will run the ROM
func RunROM() error {
	for regs.progCounter < memorySize {
		b1 := uint16(memory[regs.progCounter])
		b2 := uint16(memory[regs.progCounter+1])
		val := b1<<8 + b2
		fmt.Printf("At instruction 0x%X: 0x%x ", regs.progCounter, val)
		switch {
		case val == 0x0000:
			// 0000 NOP No Operation
			fmt.Printf("No operation\n")

		case val == 0x00E0:
			// 00E0	Display	disp_clear()	Clears the screen.
			pixels = [pixelsVertically][pixelsHorizontally]byte{}
			fmt.Printf("Clear the screen\n")

		case val == 0x00EE:
			// 00EE	Flow	return;	Returns from a subroutine.
			if len(stack) == 0 {
				fmt.Printf("Return from a subroutine => Cannot return as stack is empty\n")
			} else {
				regs.v = stack[len(stack)-1].v
				regs.index = stack[len(stack)-1].index
				regs.progCounter = stack[len(stack)-1].progCounter
				stack = stack[:len(stack)-1]
				fmt.Printf("Return from a subroutine\n")
			}

		case val <= 0x0FFF:
			// 0NNN	Call		Calls RCA 1802 program at address NNN. Not necessary for most ROMs.
			regs.progCounter = val & 0xFFF
			for i := range regs.v {
				regs.v[i] = 0
			}
			fmt.Printf("Call RCA 1802 program at address 0x%X\n", val&0xFFF)

		case val >= 0x1000 && val <= 0x1FFF:
			// 1NNN	Flow	goto NNN;	Jumps to address NNN.
			regs.progCounter = val&0xFFF - 2
			fmt.Printf("Jump to address at 0x%X\n", val&0xFFF)

		case val >= 0x2000 && val <= 0x2FFF:
			// 2NNN	Flow	*(0xNNN)()	Calls subroutine at NNN.
			regs.progCounter = val&0xFFF - 2
			stack = append(stack, stackPtr{v: regs.v, index: regs.index, progCounter: regs.progCounter})
			fmt.Printf("Call subroutine at 0x%X\n", val&0xFFF)

		case val >= 0x3000 && val <= 0x3FFF:
			// 3XNN	Cond	if(Vx==NN)	Skips the next instruction if VX equals NN.
			// (Usually the next instruction is a jump to skip a code block)
			fmt.Printf("Skip next instr if V%X (val:%d) equals 0x%X", b1&0x0F, regs.v[b1&0x0F], b2)
			if regs.v[b1&0x0F] == byte(b2) {
				regs.progCounter += 2
				fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
			} else {
				fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
			}

		case val >= 0x4000 && val <= 0x4FFF:
			// 4XNN	Cond	if(Vx!=NN)	Skips the next instruction if VX doesn't equal NN.
			// (Usually the next instruction is a jump to skip a code block)
			fmt.Printf("Skip next instr if V%X (val:%d) doesn't equal 0x%X", b1&0x0F, regs.v[b1&0x0F], b2)
			if regs.v[b1&0x0F] != byte(b2) {
				regs.progCounter += 2
				fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
			} else {
				fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
			}

		case val >= 0x5000 && val <= 0x5FFF:
			switch b2 & 0xF {
			case 0x0:
				// 5XY0	Cond	if(Vx==Vy)	Skips the next instruction if VX equals VY.
				// (Usually the next instruction is a jump to skip a code block)
				fmt.Printf("Skip next instr if V%X (val:%d) equals V%X (val:%d)",
					b1&0x0F, regs.v[b1&0x0F], b2&0xF0>>4, regs.v[b2&0xF0>>4])
				if regs.v[b1&0x0F] == regs.v[b2&0xF0>>4] {
					regs.progCounter += 2
					fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
				} else {
					fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
				}
			default:
				fmt.Printf("-------------> Unknown statement !!! Data:0x%X\n", val)
				return nil
			}

		case val >= 0x6000 && val <= 0x6FFF:
			// 6XNN	Const	Vx = NN	Sets VX to NN.
			regs.v[b1&0x0F] = byte(b2)
			fmt.Printf("Set V%X to 0x%X (%d)\n", b1&0x0F, b2, b2)

		case val >= 0x7000 && val <= 0x7FFF:
			// 7XNN	Const	Vx += NN	Adds NN to VX. (Carry flag is not changed)
			regs.v[b1&0x0F] = regs.v[b1&0x0F] + byte(b2)
			fmt.Printf("Add %d to V%X (final val:%d)\n", b2, b1&0x0F, regs.v[b1&0x0F])

		case val >= 0x8000 && val <= 0x8FFF:
			switch b2 & 0xF {
			case 0x0:
				// 8XY0	Assign	Vx=Vy	Sets VX to the value of VY.
				regs.v[b1&0x0F] = regs.v[b2&0xF0>>4]
				fmt.Printf("Set V%X to value of V%X (val:%d)\n",
					b1&0x0F, b2&0xF0>>4, regs.v[b1&0x0F])
			case 0x1:
				// 8XY1	BitOp	Vx=Vx|Vy	Sets VX to VX or VY. (Bitwise OR operation)
				regs.v[b1&0x0F] = regs.v[b1&0x0F] ^ regs.v[b2&0xF0>>4]
				fmt.Printf("Set V%X to bitwise V%X or V%X (Final val:%d)\n",
					b1&0x0F, b1&0x0F, b2&0xF0>>4, regs.v[b1&0x0F])
			case 0x2:
				// 8XY2	BitOp	Vx=Vx&Vy	Sets VX to VX and VY. (Bitwise AND operation)
				regs.v[b1&0x0F] = regs.v[b1&0x0F] & regs.v[b2&0xF0>>4]
				fmt.Printf("Set V%X to bitwise V%X and V%X (Final val:%d)\n",
					b1&0x0F, b1&0x0F, b2&0xF0>>4, regs.v[b1&0x0F])
			case 0x3:
				// 8XY3	BitOp	Vx=Vx^Vy	Sets VX to VX xor VY.
				regs.v[b1&0x0F] = regs.v[b1&0x0F] ^ regs.v[b2&0xF0>>4]
				fmt.Printf("Set V%X to bitwise V%X xor V%X (Final val:%d)\n",
					b1&0x0F, b1&0x0F, b2&0xF0>>4, regs.v[b1&0x0F])
			case 0x4:
				// 8XY4	Math	Vx += Vy	Adds VY to VX. VF is set to 1 when
				// there's a carry, and to 0 when there isn't.
				total := int(regs.v[b1&0x0F]) + int(regs.v[b2&0xF0>>4])
				if total >= 256 {
					total -= 256
					regs.v[0xF] = 1
				} else {
					regs.v[0xF] = 0
				}
				fmt.Printf("Set V%X to V%X (val:%d) + V%X (val:%d) (Final val:%d)",
					b1&0x0F, b1&0x0F, regs.v[b1&0x0F], b2&0xF0>>4, regs.v[b2&0xF0>>4], total)
				if regs.v[0xF] == 1 {
					fmt.Printf(" => CARRY OVER\n")
				} else {
					fmt.Printf(" => NOT CARRY OVER\n")
				}
				regs.v[b1&0x0F] = byte(total)
			case 0x5:
				// 8XY5	Math	Vx -= Vy	VY is subtracted from VX.
				// VF is set to 0 when there's a borrow, and 1 when there isn't.
				sub := int(regs.v[b1&0x0F]) - int(regs.v[b2&0xF0>>4])
				if sub < 0 {
					sub += 256
					regs.v[0xF] = 1
				} else {
					regs.v[0xF] = 0
				}
				fmt.Printf("Set V%X to V%X (val:%d) - V%X (val:%d) (Final val:%d)",
					b1&0x0F, b1&0x0F, regs.v[b1&0x0F], b2&0xF0>>4, regs.v[b2&0xF0>>4], sub)
				if regs.v[0xF] == 1 {
					fmt.Printf(" => BORROWED\n")
				} else {
					fmt.Printf(" => NOT BORROWED\n")
				}
				regs.v[b1&0x0F] = byte(sub)
			case 0x6:
				// 8XY6	BitOp	Vx>>=1	Stores the least significant bit of VX in VY and then shifts VX to the right by 1
				oldVal := regs.v[b1&0x0F]
				regs.v[b2&0xF0>>4] = regs.v[b1&0x0F] & 1
				regs.v[b1&0x0F] = regs.v[b1&0x0F] >> 1
				fmt.Printf("Store the least significant bit of V%X (val:%d) in V%X then shift V%X to the right by 1 (final val:%d)\n",
					b1&0x0F, oldVal, b2&0xF0>>4, regs.v[b2&0xF0>>4], regs.v[b1&0x0F])
			case 0x7:
				// 8XY7	Math	Vx=Vy-Vx	Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
				sub := int(regs.v[b2&0xF0>>4]) - int(regs.v[b1&0x0F])
				if sub < 0 {
					sub += 256
					regs.v[0xF] = 1
				} else {
					regs.v[0xF] = 0
				}
				fmt.Printf("Set V%X to V%X (val:%d) - V%X (val:%d) (Final val:%d)",
					b1&0x0F, b2&0xF0>>4, regs.v[b2&0xF0>>4], b1&0x0F, regs.v[b1&0x0F], sub)
				if regs.v[0xF] == 1 {
					fmt.Printf(" => BORROWED\n")
				} else {
					fmt.Printf(" => NOT BORROWED\n")
				}
				regs.v[b1&0x0F] = byte(sub)
			case 0xE:
				// 8XYE	BitOp	Vx<<=1	Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
				oldVal := regs.v[b1&0x0F]
				regs.v[b2&0xF0>>4] = regs.v[b1&0x0F] & 0x80
				regs.v[b1&0x0F] = regs.v[b1&0x0F] << 1
				fmt.Printf("Store the most significant bit of V%X (val:%d) in V%X then shift V%X to the left by 1 (final val:%d)\n",
					b1&0x0F, oldVal, b2&0xF0>>4, regs.v[b2&0xF0>>4], regs.v[b1&0x0F])
			default:
				fmt.Printf("-------------> Unknown statement !!! Data:0x%X\n", val)
				return nil
			}

		case val >= 0x9000 && val <= 0x9FFF:
			// 9XY0	Cond	if(Vx!=Vy)	Skips the next instruction if VX doesn't equal VY.
			// (Usually the next instruction is a jump to skip a code block)
			fmt.Printf("Skip next instr if V%X (val:%d) doesn't equal to V%X (val:%d)",
				b1&0xF, regs.v[b1&0xF], b2&0xF0>>4, regs.v[b2&0xF0>>4])
			if regs.v[b1&0x0F] != regs.v[b2&0xF0>>4] {
				regs.progCounter += 2
				fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
			} else {
				fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
			}

		case val >= 0xA000 && val <= 0xAFFF:
			// ANNN	MEM	I = NNN	Sets I to the address NNN.
			regs.index = val & 0xFFF
			fmt.Printf("Set I (memory pointer) to 0x%X (%d)\n", val&0xFFF, val&0xFFF)

		case val >= 0xB000 && val <= 0xBFFF:
			// BNNN	Flow	PC=V0+NNN	Jumps to the address NNN plus V0.
			regs.progCounter = val&0xFFF + uint16(regs.v[0])
			fmt.Printf("Jump to address 0x%X plus V0 (val:0x%X) (Final add:0x%X)\n",
				val&0xFFF, regs.v[0], regs.progCounter)

		case val >= 0xC000 && val <= 0xCFFF:
			// CXNN	Rand	Vx=rand()&NN	Sets VX to the result of a bitwise and operation
			// on a random number (Typically: 0 to 255) and NN.
			regs.v[b1&0xF] = byte(b2 & uint16(rand.Int31n(256)))
			fmt.Printf("Set V%X random value 0x%X (%d)\n", b1&0xF, regs.v[b1&0xF], regs.v[b1&0xF])

		case val >= 0xD000 && val <= 0xDFFF:
			// DXYN	Disp	draw(Vx,Vy,N)	Draws a sprite at coordinate (VX, VY) that
			// has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is
			// read as bit-coded starting from memory location I; I value doesn’t change
			// after the execution of this instruction. As described above, VF is set to
			// 1 if any screen pixels are flipped from set to unset when the sprite is drawn,
			// and to 0 if that doesn’t happen
			//  Display resolution is 64×32 pixels
			var valSlice []byte
			X := b1 & 0xF
			Y := (b2 & 0xF0) >> 4
			height := byte(b2 & 0xF)
			for i := byte(0); i < height; i++ {
				value := byte(memory[regs.index+uint16(i)])
				digits := getDigits(value)
				visiblePixels := int(math.Min(float64(pixelsHorizontally-regs.v[X]), 8))
				for j := 0; j < visiblePixels; j++ {
					if digits[j] {
						pixels[regs.v[Y]+i][regs.v[X]+byte(j)] = 1
					} else {
						pixels[regs.v[Y]+i][regs.v[X]+byte(j)] = 0
					}
				}
				valSlice = append(valSlice, value)
			}

			fmt.Printf("Draw a sprite at coor (V%X:%d, V%X:%d) width 8 pixels height %d pixels (valSlice:%d)\n",
				X, regs.v[X], Y, regs.v[Y], height, valSlice)
			drawScreen()

		case val >= 0xE000 && val <= 0xEFFF:
			switch b2 {
			case 0x9E:
				// EX9E	KeyOp	if(key()==Vx)	Skips the next instruction if the key stored in VX is pressed.
				// (Usually the next instruction is a jump to skip a code block)
				fmt.Printf("Skip instruction if key %d is pressed", b1&0xF)
				if keyPressed == regs.v[b1&0x0F] {
					regs.progCounter += 2
					fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
				} else {
					fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
				}
			case 0xA1:
				//EXA1	KeyOp	if(key()!=Vx)	Skips the next instruction if the key stored in VX
				// isn't pressed. (Usually the next instruction is a jump to skip a code block)
				fmt.Printf("Skip instruction if key %d is not pressed", b1&0xF)
				if keyPressed != regs.v[b1&0x0F] {
					regs.progCounter += 2
					fmt.Printf(" ==> SKIP NEXT INTRUCTION\n")
				} else {
					fmt.Printf(" ==> DO NOT SKIP NEXT INTRUCTION\n")
				}
			default:
				fmt.Printf("-------------> Unknown statement !!! Data:0x%X\n", val)
				return nil
			}

		case val >= 0xF000:
			switch b2 {
			case 0x0:
				// Stop
				fmt.Printf("Stop\n")
				return nil
			case 0x7:
				//FX07	Timer	Vx = get_delay()	Sets VX to the value of the delay timer.
				regs.v[b1&0xF] = getDelay()
				fmt.Printf("Set V%X to the value of the delay timer\n", b1&0xF)
			case 0xA:
				//FX0A	KeyOp	Vx = get_key()	A key press is awaited, and then stored in VX.
				// (Blocking Operation. All instruction halted until next key event)
				fmt.Printf("A key press is awaited, and then stored in V%X Blocking operation\n", b1&0xF)
				keyPressed, _ = readKeyboard()
				regs.v[b1&0xF] = keyPressed
				fmt.Printf("Set V%X as %d\n", b1&0xF, regs.v[b1&0xF])
			case 0x15:
				// FX15	Timer	delay_timer(Vx)	Sets the delay timer to VX.
				setDelayTimer(byte(b1) & 0xF)
				fmt.Printf("Set the delay timer to V%X (val:%d)\n", b1&0xF, regs.v[b1&0xF])
			case 0x18:
				// FX18	Sound	sound_timer(Vx)	Sets the sound timer to VX.
				setSoundTimer(byte(b1) & 0xF)
				fmt.Printf("Set the sound timer to V%X (val:%d)\n", b1&0xF, regs.v[b1&0xF])
			case 0x1E:
				// FX1E	MEM	I +=Vx	Adds VX to I
				regs.index = regs.index + uint16(regs.v[b1&0xF])
				fmt.Printf("Add V%X (val:%d) to I (final val:%d)\n", b1&0xF, regs.v[b1&0xF], regs.index)
			case 0x29:
				// FX29	MEM	I=sprite_addr[Vx]	Sets I to the location of the sprite for the character in VX.
				// Characters 0-F (in hexadecimal) are represented by a 4x5 font.
				// All sprites are 5 bytes long, so the location of the specified sprite
				// is its index multiplied by 5.
				regs.index = screenMemoryStart + uint16(regs.v[b1&0xF])*5
				fmt.Printf("Set I to the location (0x%X) of the sprite for the character in V%X (val:%d)\n",
					regs.index, b1&0xF, regs.v[b1&0xF])
			case 0x33:
				// FX33	BCD	set_BCD(Vx);
				// *(I+0)=BCD(3);  *(I+1)=BCD(2);  *(I+2)=BCD(1);
				// Stores the binary-coded decimal representation of VX, with the most significant of
				// three digits at the address in I, the middle digit at I plus 1, and the least significant
				// digit at I plus 2. (In other words, take the decimal representation of VX, place the
				// hundreds digit in memory at location in I, the tens digit at location I+1, and the ones
				// digit at location I+2.)
				memory[regs.index] = byte(regs.v[b1&0xF] / 100)
				memory[regs.index+1] = byte(regs.v[b1&0xF]%100) / 10
				memory[regs.index+2] = regs.v[b1&0xF] % 10
				fmt.Printf("Store BCD of V%X (val:%d) at memory index:0x%X)\n", b1&0xF, regs.v[b1&0xF], regs.index)
			case 0x55:
				// FX55	MEM	reg_dump(Vx,&I)	Stores V0 to VX (including VX) in memory starting at address I.
				// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
				var valSlice []byte
				for i := uint16(0); i <= b1&0xF; i++ {
					memory[regs.index+i] = regs.v[i]
					valSlice = append(valSlice, regs.v[i])
				}
				fmt.Printf("Store V0 to V%X (valSlice:%d) in memory starting 0x%X\n", b1&0xF, valSlice, regs.index)
			case 0x65:
				// FX65	MEM	reg_load(Vx,&I)	Fills V0 to VX (including VX) with values from memory starting
				// at address I. The offset from I is increased by 1 for each value written, but I
				// itself is left unmodified.
				var valSlice []byte
				for i := uint16(0); i <= b1&0xF; i++ {
					regs.v[i] = memory[regs.index+i]
					valSlice = append(valSlice, regs.v[i])
				}
				fmt.Printf("Fill V0 to V%X with values (valSlice:%d) at memory\n", b1&0xF, valSlice)
			default:
				fmt.Printf("-------------> Unknown statement !!! Data:0x%X\n", val)
				return nil
			}

		default:
			fmt.Printf("-------------> Unknown statement !!! Data:0x%X\n", val)
			return nil
		}
		regs.progCounter += 2
	}
	return nil
}

func readKeyboard() (byte, error) {
	reader := bufio.NewReader(os.Stdin)
	for true {
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		chStr := unicode.ToUpper(char)
		switch chStr {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F':
			result, _ := strconv.ParseUint(string(chStr), 16, 64)
			return byte(result), nil
		}
	}

	return 0, nil
}

func getDigits(x byte) [8]bool {
	var val [8]bool
	index := 7
	for x > 0 {
		if x%2 == 1 {
			val[index] = true
		}
		x = x >> 1
		index--
	}
	return val
}

func setDelayTimer(variableIndex byte) {
	regs.delayTimer = regs.v[variableIndex]
}

func getDelay() byte {
	return regs.delayTimer
}

func setSoundTimer(variableIndex byte) {
	regs.soundTimer = regs.v[variableIndex]
}

func initSprites() {
	sprites := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, //0
		0x20, 0x60, 0x20, 0x20, 0x70, //1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
		0x90, 0x90, 0xF0, 0x10, 0x10, //4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
		0xF0, 0x10, 0x20, 0x40, 0x40, //7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
		0xF0, 0x90, 0xF0, 0x90, 0x90, //A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
		0xF0, 0x80, 0x80, 0x80, 0xF0, //C
		0xE0, 0x90, 0x90, 0x90, 0xE0, //D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
		0xF0, 0x80, 0xF0, 0x80, 0x80, //F
	}
	copy(memory[screenMemoryStart:], sprites)
}
