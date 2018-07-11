package standard

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
)

// PrintLogo prints the genesis logo to the current standard out
func PrintLogo() {
	fmt.Fprintf(color.Output, "%s\n", ASCIILogo())
}

// ASCIILogo returns a string containing a color formatted logo block
func ASCIILogo() string {
	lines := []string{
		color.YellowString("***********************************************************"),
		color.HiWhiteString("                             ____                         "),
		color.HiWhiteString("                     __,-~~/~    `---.                    "),
		color.HiWhiteString("                   _/_,---(      ,    )                   "),
		color.HiWhiteString("               __ /        <    /   )  \\___               "),
		color.HiWhiteString("- ------===;;;'====------------------===;;;===----- -  -  "),
		color.HiWhiteString("                  \\/  ~\"~\"~\"~\"~\"~\\~\"~)~\"/                 "),
		color.HiWhiteString("                  (_ (   \\  (     >    \\)                 "),
		color.HiWhiteString("                   \\_( _ <         >_>'                   "),
		color.HiWhiteString("                      ~ `-i' ::>|--\"                      "),
		color.HiWhiteString("                          I;|.|.|                         "),
		color.HiWhiteString("                         <|i::|i|`.                       "),
		fmt.Sprintf("            %s          %s          %s  ", color.HiGreenString("uL"), color.HiWhiteString("(` ^'\"`-' \")"), color.HiYellowString(")")),
		fmt.Sprintf("        %s          %s  ", color.HiGreenString(".ue888Nc.."), color.HiYellowString("(   (          ( /(")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("d88E`\"888E`"), color.HiYellowString("(    (  )(  )\\  `  )   )\\())")),
		fmt.Sprintf("       %s   %s  ", color.HiGreenString("888E  888E"), color.YellowString(")\\   )\\(()\\((_) /(/(  (_))/")),
		fmt.Sprintf("       %s  %s   ", color.HiGreenString("888E  888E"), color.HiRedString("((_) ((_)((_)(_)((_)_\\ | |_")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("888E  888E"), color.RedString("(_-</ _|| '_|| || '_ \\)|  _|")),
		fmt.Sprintf("       %s  %s %s ", color.HiGreenString("888& .888E"), color.RedString("/__/\\__||_|  |_|| .__/  \\__|"), color.WhiteString(gscript.Version)),
		fmt.Sprintf("       %s                  %s           ", color.HiGreenString("*888\" 888&"), color.RedString("|_|")),
		fmt.Sprintf("        %s  %s        -- By --", color.HiGreenString("`\"   \"888E"), color.HiWhiteString("G E N I S I S")),
		fmt.Sprintf("       %s   %s       %s", color.HiGreenString(".dWi   `88E"), color.HiWhiteString("S C R I P T I N G"), color.CyanString("gen0cide")),
		fmt.Sprintf("       %s    %s            %s", color.HiGreenString("4888~  J8%%"), color.HiWhiteString("E N G I N E"), color.CyanString("ahhh")),
		fmt.Sprintf("        %s                             %s", color.HiGreenString("^\"===*\"`"), color.CyanString("vyrus")),
		"                github.com/gen0cide/gscript",
		color.YellowString("***********************************************************"),
	}

	return strings.Join(lines, "\n")
}
