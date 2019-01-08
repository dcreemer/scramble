# Scamble Square Solver

I recently came across a [Scramble Square](http://www.b-dazzle.com/puzzdetail.asp?a=1&PuzzID=70&CatID=9) puzzle of Penguins.

![image of puzzle](http://www.b-dazzle.com/store/pc/catalog/10044penguins_1067_detail.jpg)

The puzzle consists of nine cards, each with either the lower or upper half of some
penguins. There are four different penguin designs (Adelie, Gentoo, Emperor, and
Chinstrap), and the top of bottom half of one of the designs is found on one of the
four edges of each square card.

The object of the puzzle is to arrange the nine cards in a 3x3 grid so that all of
the connecting top and bottom halves match up. I think this means that there are 4^9
or 232,144 combinations to try. I spent a couple hours playing with the puzzle one
evening, and then decided to write a program to solve it :-).
