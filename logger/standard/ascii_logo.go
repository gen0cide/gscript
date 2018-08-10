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
		errorLevel.Sprint("***********************************************************"),
		infoMsg.Sprint("                             ____                         "),
		infoMsg.Sprint("                     __,-~~/~    `---.                    "),
		infoMsg.Sprint("                   _/_,---(      ,    )                   "),
		infoMsg.Sprint("               __ /        <    /   )  \\___               "),
		infoMsg.Sprint(" - ------===;;;'====-----------------===;;;===----- -  -  "),
		infoMsg.Sprint("                  \\/  ~\"~\"~\"~\"~\"~\\~\"~)~\"/                 "),
		infoMsg.Sprint("                  (_ (   \\  (     >    \\)                 "),
		infoMsg.Sprint("                   \\_( _ <         >_>'                   "),
		infoMsg.Sprint("                      ~ `-i' ::>|--\"                      "),
		infoMsg.Sprint("                          I;|.|.|                         "),
		infoMsg.Sprint("                         <|i::|i|`.                       "),
		fmt.Sprintf("            %s          %s          %s  ", color.HiGreenString("uL"), color.HiWhiteString("(` ^'\"`-' \")"), color.HiYellowString(")")),
		fmt.Sprintf("        %s          %s  ", color.HiGreenString(".ue888Nc.."), color.HiYellowString("(   (          ( /(")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("d88E`\"888E`"), color.HiYellowString("(    (  )(  )\\  `  )   )\\())")),
		fmt.Sprintf("       %s   %s  ", color.HiGreenString("888E  888E"), color.YellowString(")\\   )\\(()\\((_) /(/(  (_))/")),
		fmt.Sprintf("       %s  %s   ", color.HiGreenString("888E  888E"), fatalLevel.Sprint("((_) ((_)((_)(_)((_)_\\ | |_")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("888E  888E"), color.HiRedString("(_-</ _|| '_|| || '_ \\)|  _|")),
		fmt.Sprintf("       %s  %s %s ", color.HiGreenString("888& .888E"), fatalMsg.Sprint("/__/\\__||_|  |_|| .__/  \\__|"), defaultLevel.Sprintf("v%s", gscript.Version)),
		fmt.Sprintf("       %s                  %s           ", color.HiGreenString("*888\" 888&"), fatalMsg.Sprint("|_|")),
		fmt.Sprintf("        %s  %s        -- By --", color.HiGreenString("`\"   \"888E"), infoLevel.Sprint("G E N E S I S")),
		fmt.Sprintf("       %s   %s       %s", color.HiGreenString(".dWi   `88E"), infoLevel.Sprint("S C R I P T I N G"), debugLevel.Sprint("gen0cide")),
		fmt.Sprintf("       %s    %s            %s", color.HiGreenString("4888~  J8%%"), infoLevel.Sprint("E N G I N E"), debugLevel.Sprint("ahhh")),
		fmt.Sprintf("        %s                             %s", color.HiGreenString("^\"===*\"`"), debugLevel.Sprint("vyrus")),
		"                github.com/gen0cide/gscript",
		errorLevel.Sprint("***********************************************************"),
	}

	return strings.Join(lines, "\n")
}
