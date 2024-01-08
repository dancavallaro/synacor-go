# Synacor Challenge VM Implementation in Go

This is my work-in-progress implementation of a VM and text-based graphical 
debugger for the machine architecture described in the Synacor Challenge 
from OSCON 2012. The Synacor Challenge website is no longer alive, so I've 
been working from the archive at https://github.com/Aneurysm9/vm_challenge/tree/main. 

## TODOs

- Refactor to introduce "environment" concept with methods like
  Log/Halt/ReadInput, with different impls for VM vs. Debugger.
- Fix weird issue with reading from stdin, and make it work in the debugger.
- Implement scrolling in output view
- Finish TODOs for displaying memory contents
- Display disassembled program
- Keybindings to:
  - Reset start and restart execution
  - Set breakpoints/run until breakpoint
  - Edit registers/stack/memory
- Assembler, to help write new example programs