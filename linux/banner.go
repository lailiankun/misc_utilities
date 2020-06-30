// ported from sysvbanner.c by Brian Wallis
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		for {
			var s string
			_, err := fmt.Scan(&s)
			if err != nil {
				break
			}
			printGlyphs(s)
		}
	} else {
		for _, s := range os.Args[1:] {
			printGlyphs(s)
		}
	}
}

func printGlyphs(s string) {
	var line [80]byte

	l := len(s)
	if l > 10 {
		l = 10
	}

	for a := 0; a < 7; a++ {
		for b := 0; b < l; b++ {
			ind := int(s[b] - ' ')
			if ind < 0 {
				ind = 0
			}
			for c := 0; c < 7; c++ {
				line[b*8+c] = glyphs[(ind/8*7)+a][(ind%8*7)+c]
			}
			line[b*8+7] = ' '
		}

		for b := l*8 - 1; b >= 0; b-- {
			if line[b] != ' ' {
				break
			}
			line[b] = 0
		}

		for i := range line {
			if line[i] == 0 {
				break
			}
			fmt.Printf("%c", line[i])
		}
		fmt.Println("")
	}
	fmt.Println("")

}

var glyphs = []string{
	"         ###  ### ###  # #   ##### ###   #  ##     ###  ",
	"         ###  ### ###  # #  #  #  ## #  #  #  #    ###   ",
	"         ###   #   # ########  #   ### #    ##      #   ",
	"          #            # #   #####    #    ###     #    ",
	"                     #######   #  #  # ####   # #       ",
	"         ###           # #  #  #  # #  # ##    #        ",
	"         ###           # #   ##### #   ### #### #       ",

	"   ##    ##                                            #",
	"  #        #   #   #    #                             # ",
	" #          #   # #     #                            #  ",
	" #          # ### ### #####   ###   #####           #   ",
	" #          #   # #     #     ###           ###    #    ",
	"  #        #   #   #    #      #            ###   #     ",
	"   ##    ##                   #             ###  #      ",

	"  ###     #    #####  ##### #      ####### ##### #######",
	" #   #   ##   #     ##     ##    # #      #     ##    # ",
	"# #   # # #         #      ##    # #      #          #  ",
	"#  #  #   #    #####  ##### ####### ##### ######    #   ",
	"#   # #   #   #            #     #       ##     #  #    ",
	" #   #    #   #      #     #     # #     ##     #  #    ",
	"  ###   ##### ####### #####      #  #####  #####   #    ",

	" #####  #####    #     ###      #           #     ##### ",
	"#     ##     #  # #    ###     #             #   #     #",
	"#     ##     #   #            #     #####     #        #",
	" #####  ######         ###   #                 #     ## ",
	"#     #      #   #     ###    #     #####     #     #   ",
	"#     ##     #  # #     #      #             #          ",
	" #####  #####    #     #        #           #       #   ",

	" #####    #   ######  ##### ###### ############## ##### ",
	"#     #  # #  #     ##     ##     ##      #      #     #",
	"# ### # #   # #     ##      #     ##      #      #      ",
	"# # # ##     ####### #      #     ######  #####  #  ####",
	"# #### ########     ##      #     ##      #      #     #",
	"#     ##     ##     ##     ##     ##      #      #     #",
	" ##### #     #######  ##### ###### ########       ##### ",

	"#     #  ###        ##    # #      #     ##     ########",
	"#     #   #         ##   #  #      ##   ####    ##     #",
	"#     #   #         ##  #   #      # # # ## #   ##     #",
	"#######   #         ####    #      #  #  ##  #  ##     #",
	"#     #   #   #     ##  #   #      #     ##   # ##     #",
	"#     #   #   #     ##   #  #      #     ##    ###     #",
	"#     #  ###   ##### #    # ########     ##     ########",

	"######  ##### ######  ##### ########     ##     ##     #",
	"#     ##     ##     ##     #   #   #     ##     ##  #  #",
	"#     ##     ##     ##         #   #     ##     ##  #  #",
	"###### #     #######  #####    #   #     ##     ##  #  #",
	"#      #   # ##   #        #   #   #     # #   # #  #  #",
	"#      #    # #    # #     #   #   #     #  # #  #  #  #",
	"#       #### ##     # #####    #    #####    #    ## ## ",

	"#     ##     ######## ##### #       #####    #          ",
	" #   #  #   #      #  #      #          #   # #         ",
	"  # #    # #      #   #       #         #  #   #        ",
	"   #      #      #    #        #        #               ",
	"  # #     #     #     #         #       #               ",
	" #   #    #    #      #          #      #               ",
	"#     #   #   ####### #####       # #####        #######",

	"  ###                                                   ",
	"  ###     ##   #####   ####  #####  ###### ######  #### ",
	"   #     #  #  #    # #    # #    # #      #      #    #",
	"    #   #    # #####  #      #    # #####  #####  #     ",
	"        ###### #    # #      #    # #      #      #  ###",
	"        #    # #    # #    # #    # #      #      #    #",
	"        #    # #####   ####  #####  ###### #       #### ",

	"                                                        ",
	" #    #    #        # #    # #      #    # #    #  #### ",
	" #    #    #        # #   #  #      ##  ## ##   # #    #",
	" ######    #        # ####   #      # ## # # #  # #    #",
	" #    #    #        # #  #   #      #    # #  # # #    #",
	" #    #    #   #    # #   #  #      #    # #   ## #    #",
	" #    #    #    ####  #    # ###### #    # #    #  #### ",

	"                                                        ",
	" #####   ####  #####   ####   ##### #    # #    # #    #",
	" #    # #    # #    # #         #   #    # #    # #    #",
	" #    # #    # #    #  ####     #   #    # #    # #    #",
	" #####  #  # # #####       #    #   #    # #    # # ## #",
	" #      #   #  #   #  #    #    #   #    #  #  #  ##  ##",
	" #       ### # #    #  ####     #    ####    ##   #    #",

	"                       ###     #     ###   ##    # # # #",
	" #    #  #   # ###### #        #        # #  #  # # # # ",
	"  #  #    # #      #  #        #        #     ## # # # #",
	"   ##      #      #  ##                 ##        # # # ",
	"   ##      #     #    #        #        #        # # # #",
	"  #  #     #    #     #        #        #         # # # ",
	" #    #    #   ######  ###     #     ###         # # # #",
}
