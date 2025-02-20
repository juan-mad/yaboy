# yaboy
It is Yet Another Game Boy emulator

# Brainstorming
What follows are some brainstorming ideas of pieces of the emulator that I could tackle somewhat
independently. For now I am restricting myself to DMG.

 - **RAM memory**: this should be easy, creating one or several arrays of appropriate size.
 Different banks and their handling adds some complexity to this.
 - **CPU**: this seems to entail creating the registers that are in the CPU, as well as imeplementing
 the different instructions of the original CPU. This requires the RAM memory, as we can manipulate it
 using CPU instructions.
 - **Assembly code parser**: perhaps useful for testing, a parser of `.asm` source files into the CPU
 instructions.
 
 The above already looks like it would take a while. Other parts to be implemented would be:
  - Everything **graphics related**
  - Everything **audio related**
  - Decoder of `.gb` files

