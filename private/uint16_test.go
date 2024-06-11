package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestUint16(t *testing.T) {
	stat := &svcStat{
		RWMutex: &sync.RWMutex{},
		value:   0,
	}
	printStat(stat.value)
	fmt.Println(stat.exists(Normal))
	stat.add(Normal)
	printStat(stat.value)
	fmt.Println(stat.exists(Normal))
	stat.add(ScoreLess)
	printStat(stat.value)
	stat.add(Exception)
	printStat(stat.value)
	stat.add(AppLock)
	printStat(stat.value)
	stat.add(OpenLock)
	printStat(stat.value)
	stat.add(OpenTimeout)
	printStat(stat.value)
	stat.add(MemoryAlarm)
	printStat(stat.value)
	stat.add(VirtualMachineError)
	printStat(stat.value)
	stat.add(Inactive)
	printStat(stat.value)
	stat.add(Broken)
	printStat(stat.value)

	stat.remove(Normal)
	stat.remove(ScoreLess)
	stat.remove(Exception)
	stat.remove(AppLock)
	stat.remove(OpenLock)
	stat.remove(OpenTimeout)
	stat.remove(MemoryAlarm)
	stat.remove(VirtualMachineError)
	stat.remove(Inactive)
	stat.remove(Broken)
	all := stat.getAll()
	fmt.Println(all)

}

func printStat(v uint16) {
	fmt.Printf("value:%016b\n", v)
}

const (
	Normal = 1 << (iota)
	Exception
	ScoreLess
	AppLock
	OpenLock
	OpenTimeout
	MemoryAlarm
	_
	_
	_
	_
	_
	_
	VirtualMachineError
	Inactive
	Broken
	_
)

type svcStat struct {
	*sync.RWMutex
	value uint16
}

func (s *svcStat) exists(stp uint16) bool {
	bits := s.bits(stp)
	s.RLock()
	defer s.RUnlock()
	return s.hasValueAtBit(bits)
}

func (s *svcStat) add(stp uint16) {
	if stp > 0 {
		s.Lock()
		defer s.Unlock()
		s.value |= uint16(1) << (s.bits(stp))
	}
}

func (s *svcStat) bits(stp uint16) uint8 {
	var bit uint8 = 0
	for result := uint16(1); result < stp; result <<= 1 {
		bit++
	}
	return bit
}

func (s *svcStat) remove(stp uint16) {
	s.Lock()
	defer s.Unlock()
	a := uint16(0xFFFF)
	b := uint16(1) << s.bits(stp)
	c := s.value & (a ^ b)
	// s.value &= uint16(0xFF) ^ (1 << s.bits(stp))
	s.value = c
}

func (s *svcStat) get() uint16 {
	s.RLock()
	defer s.RUnlock()
	if s.value != 0 {
		var bit int
		for bit = 16; bit >= 0; bit-- {
			if s.hasValueAtBit(uint8(bit)) {
				return 1 << bit
			}
		}
	}
	return 0
}

func (s *svcStat) getAll() []uint16 {
	s.RLock()
	defer s.RUnlock()
	if s.value != 0 {
		var tps []uint16
		var bit int
		for bit = 16; bit >= 0; bit-- {
			if s.hasValueAtBit(uint8(bit)) {
				tps = append(tps, 1<<bit)
			}
		}
		return tps
	}
	return nil
}

func (s *svcStat) hasValueAtBit(bitPos uint8) bool {
	return (s.value & (uint16(1) << bitPos)) != 0
}

func (s *svcStat) reset() {
	s.Lock()
	defer s.Unlock()
	s.value = 0
}
