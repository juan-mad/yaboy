# yaboy
It is Yet Another Game Boy emulator

# Done so far
 - Compiler of instructions into the corresponding instruction bytes
 - Decompiler of ROMs or general binary files into an instruction file

# Next steps
 - Simulate RAM memory, registries, etc.
 - CLI tool that can:
    - Inspect a specific registry or address in memory
    - Execute an instruction
    - Have in mind that this could work as a debugger
 - Execute instruction by instruction a ROM, taking into account the cycles
    - Control via CLI each cycle

