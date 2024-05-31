package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	b := CompilerBF(">++++++[<+++++++++++++++++++>-]<")
	out, err := b.Exectue()
	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	fmt.Println(out)
}

const (
	TAPE_SIZE_DEFAULT uint32 = 500
	DEBUG             bool   = false
	data_size         int    = 65535
)

type BF struct {
	Code    string
	Stack   []uint16
	loops   []uint32
	windup  uint32
	Pointer int
	OutPut  string
}

func CompilerBF(code string) *BF {
	return &BF{
		Code:   code,
		Stack:  make([]uint16, len(code)),
		loops:  []uint32{},
		windup: 0,
	}
}

func (s *BF) Exectue() (string, error) {
	var err error
	for _, i := range s.Code {
		switch i {
		// +
		case 43:
			err = s.increment()

		// -
		case 45:
			err = s.decrement()

		// <
		case 60:
			err = s.shift(1)

		// >
		case 62:
			err = s.shift(0)

		// ,
		case 44:
			var reader *bufio.Reader = bufio.NewReader(os.Stdin)
			char, err := reader.ReadByte()
			if err != nil {
				return "", err
			}

			fmt.Println(char, uint16(char))

			s.Stack[s.Pointer] = uint16(char)

		// [
		case 91:
			// s.openLoop()

		// ]
		case 93:
			// s.closeLoop()

		default:
			err = errors.New("unknown operator")
		}
		if err != nil {
			return "", err
		}
	}

	fmt.Println(s.Stack[s.Pointer])
	return fmt.Sprintf("%c", s.Stack[s.Pointer]), nil
}

func (s *BF) increment() error {
	if s.Stack[s.Pointer] >= 255 {
		return errors.New("unable to increment stack")
	}

	s.Stack[s.Pointer]++
	return nil
}

func (s *BF) decrement() error {
	if s.Pointer <= 0  {
		return errors.New("unable to decrement stack")
	}

	s.Stack[s.Pointer]--
	return nil
}

func (s *BF) shift(mode int8) error {
	// 0 => Right == >
	// 1 ==> Left == <
	if mode == 0 {
		if uint32(s.Pointer) >= TAPE_SIZE_DEFAULT-1 {
			return errors.New("RangeError")
		} else {
			s.Pointer++
		}
	} else {
		if s.Pointer <= 0 {
			return errors.New("RangeError")
		} else {
			s.Pointer--
		}
	}

	return nil
}

func (s *BF) openLoop() {
	s.loops = append(s.loops, s.windup)
}

func (s *BF) closeLoop() {
	if s.Stack[s.Pointer] == 0 {
		s.popLoop()
		return
	}

	s.windup = s.peakLoop()
}

func (s *BF) popLoop() {
	length := len(s.loops)
	s.loops = s.loops[:length-1]
}

func (s *BF) peakLoop() uint32 {
	length := len(s.loops)
	return s.loops[length-1]
}

// func (s *BF) isValidToken(str string) bool {
// 	return strings.Contains("+-><,[].,", str)
// }
