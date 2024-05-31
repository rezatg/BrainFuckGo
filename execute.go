package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var s string = ">++++++[<+++++++++++++++++++>-]<.>->++++++[<+++++++++++++++++>-]<.>>++[<+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++>-]<.>->+++++++[<++++++++++++++>-]<."
	b := CompilerBF(s)
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
	Stack   []uint32
	loops   []uint32
	windup  uint32
	Pointer int
	OutPut  string
}

func CompilerBF(code string) *BF {
	return &BF{
		Code:   strings.ReplaceAll(code, " ", ""),
		Stack:  make([]uint32, len(code)),
		loops:  []uint32{},
		windup: 0,
	}
}

func (s *BF) Exectue() (string, error) {
	var err error
	var result strings.Builder

	for i := 0; i < len(s.Code); i++ {
		switch s.Code[i] {
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

			s.Stack[s.Pointer] = uint32(char)

		// .
		case 46:
			result.WriteByte(byte(s.Stack[s.Pointer]))

		// [
		case 91:
			if s.Stack[s.Pointer] == 0 {
				var skips int
				for skips = 1; skips > 0; i++ {
					if s.Code[i] == '[' {
						skips++
					} else if s.Code[i] == ']' {
						skips--
					}
				}
				i--
			} else {
				s.loops = append(s.loops, uint32(i))
			}

		// ]
		case 93:
			if s.Stack[s.Pointer] != 0 {
				i = int(s.peekLoop()) - 1
			} else {
				s.popLoop()
			}

		default:
			err = errors.New("unknown operator")
		}
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func (s *BF) increment() error {
	s.Stack[s.Pointer]++
	if s.Stack[s.Pointer] == 0 {
		s.Stack[s.Pointer] = 1
	}

	return nil
}

func (s *BF) decrement() error {
	if s.Stack[s.Pointer] == 0 {
		s.Stack[s.Pointer] = 255
	} else {
		s.Stack[s.Pointer]--
	}

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

func (s *BF) peekLoop() uint32 {
	length := len(s.loops)
	return s.loops[length-1]
}

func (s *BF) popLoop() {
	s.loops = s.loops[:len(s.loops)-1]
}

// func (s *BF) isValidToken(str string) bool {
// 	return strings.Contains("+-><,[].,", str)
// }
