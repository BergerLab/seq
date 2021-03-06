package seq

import (
	"fmt"
	"strings"
	"testing"
)

func TestNeedlemanWunsch(t *testing.T) {
	type test struct {
		seq1, seq2 string
		out1, out2 string
	}

	tests := []test{
		{
			"ABCD",
			"ABCD",
			"ABCD",
			"ABCD",
		},
		{
			"GHIKLMNPQR",
			"GAAAHIKLMN",
			"---GHIKLMNPQR",
			"GAAAHIKLMN---",
		},
		{
			"GHIKLMNPQRSTVW",
			"GAAAHIKLMNPQRSTVW",
			"---GHIKLMNPQRSTVW",
			"GAAAHIKLMNPQRSTVW",
		},
		{
			"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		},
		{
			"NNNNNNNN",
			"NNNNNNNN",
			"NNNNNNNN",
			"NNNNNNNN",
		},
		{
			"NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
			"NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
			"NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
			"NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
		},
		{
			"ABCDEFGWXYZ",
			"ABCDEFMNPQRSTZABEGWXYZ",
			"ABCDEF-----------GWXYZ",
			"ABCDEFMNPQRSTZABEGWXYZ",
		},
		{
			"ASAECVSNENVEIEAPKTNIWTSLAKEEVQEVLDLLHSTYNITEVTKADFFSNYVLWIETLKPN" +
				"KTEALTYLDEDGDLPPRNARTVVYFGEGEEGYFEELKVGPLPVSDETTIEPLSFYNTNGK" +
				"SKLPFEVGHLDRIKSAAKSSFLNKNLNTTIMRDVLEGLIGVPYEDMGCHSAAPQLHDPAT" +
				"GATVDYGTCNINTENDAENLVPTGFFFKFDMTGRDVSQWKMLEYIYNNKVYTSAEELYEA" +
				"MQKDDFVTLPKIDVDNLDWTVIQRNDSAPVRHLDDRKSPRLVEPEGRRWAYDGDEEYFSW" +
				"MDWGFYTSWSRDTGISFYDITFKGERIVYELSLQELIAEYGSDDPFNQHTFYSDISYGVG" +
				"NRFSLVPGYDCPSTAGYFTTDTFEYDEFYNRTLSYCVFENQEDYSLLRHTGASYSAITQN" +
				"PTLNVRFISTIGNYDYNFLYKFFLDGTLEVSVRAAGYIQAGYWNPETSAPYGLKIHDVLS" +
				"GSFHDHVLNYKVDLDVGGTKNRASQYVMKDVDVEYPWAPGTVYNTKQIAREVFENEDFNG" +
				"INWPENGQGILLIESAEETNSFGNPRAYNIMPGGGGVHRIVKNSRSGPETQNWARSNLFL" +
				"TKHKDTELRSSTALNTNALYDPPVNFNAFLDDESLDGEDIVAWVNLGLHHLPNSNDLPNT" +
				"IFSTAHASFMLTPFNYFDSENSRDTTQQVFYTYDDETEESNWEFYGNDWSSCGVEVAEPN" +
				"FEDYTYGRGTRINKKMTNSDEVY",
			"AECVSNENVEIEAPKTNIWTSLAKEEVQEVLDLLHSTYNITEVTKADFFSNYVLWIETLKPNKT" +
				"EALTYLDEDGDLPPRNARTVVYFGEGEEGYFEELKVGPLPVSDETTIEPLSFYNTNGKSK" +
				"LPFEVGHLDRIKSAAKSSFLNKNLNTTIMRDVLEGLIGVPYEDMGCHSAAPQLHDPATGA" +
				"TVDYGTCNINTENDAENLVPTGFFFKFDMTGRDVSQWKMLEYIYNNKVYTSAEELYEAMQ" +
				"KDDFVTLPKIDVDNLDWTVIQRNDSAPVRHLDDRKSPRLVEPEGRRWAYDGDEEYFSWMD" +
				"WGFYTSWSRDTGISFYDITFKGERIVYELSLQELIAEYGSDDPFNQHTFYSDISYGVGNR" +
				"FSLVPGYDCPSTAGYFTTDTFEYDEFYNRTLSYCVFENQEDYSLLRHTGASYSAITQNPT" +
				"LNVRFISTIGNDYNFLYKFFLDGTLEVSVRAAGYIQAGYWNPETSAPYGLKIHDVLSGSF" +
				"HDHVLNYKVDLDVGGTKNRASQYVMKDVDVEYPWAPGTVYNTKQIAREVFENEDFNGINW" +
				"PENGQGILLIESAEETNSFGNPRAYNIMPGGGGVHRIVKNSRSGPETQNWARSNLFLTKH" +
				"KDTELRSSTALNTNALYDPPVNFNAFLDDESLDGEDIVAWVNLGLHHLPNSNDLPNTIFS" +
				"TAHASFMLTPFNYFDSENSRDTTQQVFYTYDDETEESNWEFYGNDWSSCGVEVAEPNFED" +
				"YTYGRGTRINKK",
			"ASAECVSNENVEIEAPKTNIWTSLAKEEVQEVLDLLHSTYNITEVTKADFFSNYVLWIETLKPN" +
				"KTEALTYLDEDGDLPPRNARTVVYFGEGEEGYFEELKVGPLPVSDETTIEPLSFYNTNGK" +
				"SKLPFEVGHLDRIKSAAKSSFLNKNLNTTIMRDVLEGLIGVPYEDMGCHSAAPQLHDPAT" +
				"GATVDYGTCNINTENDAENLVPTGFFFKFDMTGRDVSQWKMLEYIYNNKVYTSAEELYEA" +
				"MQKDDFVTLPKIDVDNLDWTVIQRNDSAPVRHLDDRKSPRLVEPEGRRWAYDGDEEYFSW" +
				"MDWGFYTSWSRDTGISFYDITFKGERIVYELSLQELIAEYGSDDPFNQHTFYSDISYGVG" +
				"NRFSLVPGYDCPSTAGYFTTDTFEYDEFYNRTLSYCVFENQEDYSLLRHTGASYSAITQN" +
				"PTLNVRFISTIGNYDYNFLYKFFLDGTLEVSVRAAGYIQAGYWNPETSAPYGLKIHDVLS" +
				"GSFHDHVLNYKVDLDVGGTKNRASQYVMKDVDVEYPWAPGTVYNTKQIAREVFENEDFNG" +
				"INWPENGQGILLIESAEETNSFGNPRAYNIMPGGGGVHRIVKNSRSGPETQNWARSNLFL" +
				"TKHKDTELRSSTALNTNALYDPPVNFNAFLDDESLDGEDIVAWVNLGLHHLPNSNDLPNT" +
				"IFSTAHASFMLTPFNYFDSENSRDTTQQVFYTYDDETEESNWEFYGNDWSSCGVEVAEPN" +
				"FEDYTYGRGTRINKKMTNSDEVY",
			"--AECVSNENVEIEAPKTNIWTSLAKEEVQEVLDLLHSTYNITEVTKADFFSNYVLWIETLKPN" +
				"KTEALTYLDEDGDLPPRNARTVVYFGEGEEGYFEELKVGPLPVSDETTIEPLSFYNTNGK" +
				"SKLPFEVGHLDRIKSAAKSSFLNKNLNTTIMRDVLEGLIGVPYEDMGCHSAAPQLHDPAT" +
				"GATVDYGTCNINTENDAENLVPTGFFFKFDMTGRDVSQWKMLEYIYNNKVYTSAEELYEA" +
				"MQKDDFVTLPKIDVDNLDWTVIQRNDSAPVRHLDDRKSPRLVEPEGRRWAYDGDEEYFSW" +
				"MDWGFYTSWSRDTGISFYDITFKGERIVYELSLQELIAEYGSDDPFNQHTFYSDISYGVG" +
				"NRFSLVPGYDCPSTAGYFTTDTFEYDEFYNRTLSYCVFENQEDYSLLRHTGASYSAITQN" +
				"PTLNVRFISTIGN-DYNFLYKFFLDGTLEVSVRAAGYIQAGYWNPETSAPYGLKIHDVLS" +
				"GSFHDHVLNYKVDLDVGGTKNRASQYVMKDVDVEYPWAPGTVYNTKQIAREVFENEDFNG" +
				"INWPENGQGILLIESAEETNSFGNPRAYNIMPGGGGVHRIVKNSRSGPETQNWARSNLFL" +
				"TKHKDTELRSSTALNTNALYDPPVNFNAFLDDESLDGEDIVAWVNLGLHHLPNSNDLPNT" +
				"IFSTAHASFMLTPFNYFDSENSRDTTQQVFYTYDDETEESNWEFYGNDWSSCGVEVAEPN" +
				"FEDYTYGRGTRINKK--------",
		},
	}
	sep := strings.Repeat("-", 45)
	for _, test := range tests {
		s1, s2 := stringToSeq(test.seq1), stringToSeq(test.seq2)
		aligned := NeedlemanWunsch(s1, s2)
		sout1 := fmt.Sprintf("%s", aligned.A)
		sout2 := fmt.Sprintf("%s", aligned.B)

		if sout1 != test.out1 || sout2 != test.out2 {
			t.Fatalf(
				`Alignment for:
%s
%s
%s
%s
resulted in
%s
%s
%s
%s
but should have been
%s
%s
%s
%s`,
				sep, test.seq1, test.seq2, sep,
				sep, sout1, sout2, sep,
				sep, test.out1, test.out2, sep)
		}
	}
}

func stringToSeq(s string) []Residue {
	residues := make([]Residue, len(s))
	for i := range s {
		residues[i] = Residue(s[i])
	}
	return residues
}
