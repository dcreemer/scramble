# Penguin Puzzle

ref: http://www.b-dazzle.com/scramble.asp

There are 9 cards, each with 1/2 of a penguin on each of the four edges. There are
four penguin designs, A, B, C, and D. Cards must be places in a 3x3 square, so that
the facing edges have a corresponding match, i.e. A1+A2 or B2+B2, but not A1+B1.

There are 4^9 or 262144 combinations - an easy problem for a computer.

Cards have a number 1-9, and a N, S, E, and W edge value. Each value is one of a
top or bottom penguin, 1a, 1b, 2a, 2b etc.

card map {:n :1a, :s, 1b...}
rotation 0, 90, 180, 270
