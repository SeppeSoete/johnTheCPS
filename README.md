## John The Chess Puzzle Spoiler
### What does this do?
This program attempts to find a chess move described in [algebraic notation](https://en.wikipedia.org/wiki/Algebraic_notation_(chess)), when given the md5sum of that move (which could be generated by executing `echo "move" | md5sum` on a linux system)

### But why?
Because me and some buddies like to share puzzles and give the md5sum of the solution so that anyone can check their answer but they can't easily look it up.
Since I like being 'that guy' I decided to make a bruteforce tool.

### How do I use it then?
Build it using the go toolchain, if you don't know how to do so: there's google for that.
execute the binary with an optional parameter -n specifying the amount of [goroutines](https://en.wikipedia.org/wiki/Go_(programming_language)#Concurrency) you want to use for calculating the md5sums (default = 4), and with the md5sum of a chess move as argument

example: `./johnTheCPS -n 4 5b99cf37f882ea71d7e13db08c6503b7`
