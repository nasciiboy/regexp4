//
// Recursive Regexp Raptor (go version)
// Available at http://github.com/nasciiboy/regexp4
//
// Copyright © 2017 nasciiboy <nasciiboy@gmail.com>.
// Distributed under the GNU GPL v3 License.
// See readme.org for details.
//

//
// Unit tests
//

package regexp4

import "testing"
import "fmt"
import "bytes"

func printASM( rexp *RE ){
  fmt.Printf( "                     init %q\n", rexp.re )

  for i, v := range rexp.asm {
    fmt.Printf( "[%3d][%3d]", i, v.close )

    switch v.inst {
    case  0: fmt.Printf( "[%-12s]", "asmPath"     )
    case  1: fmt.Printf( "[%-12s]", "asmPathEle"  )
    case  2: fmt.Printf( "[%-12s]", "asmPathEnd"  )
    case  3: fmt.Printf( "[%-12s]", "asmGroup"    )
    case  4: fmt.Printf( "[%-12s]", "asmGroupEnd" )
    case  5: fmt.Printf( "[%-12s]", "asmHook"     )
    case  6: fmt.Printf( "[%-12s]", "asmHookEnd"  )
    case  7: fmt.Printf( "[%-12s]", "asmSet"      )
    case  8: fmt.Printf( "[%-12s]", "asmSetEnd"   )
    case  9: fmt.Printf( "[%-12s]", "asmBackref"  )
    case 10: fmt.Printf( "[%-12s]", "asmMeta"     )
    case 11: fmt.Printf( "[%-12s]", "asmRangeab"  )
    case 12: fmt.Printf( "[%-12s]", "asmUTF8"     )
    case 13: fmt.Printf( "[%-12s]", "asmPoint"    )
    case 14: fmt.Printf( "[%-12s]", "asmSimple"   )
    case 15: fmt.Printf( "[%-12s]", "asmEnd"      )
    }

    fmt.Printf( " %-15q [%d-%d][%08b]\n", v.re.str, v.re.loopsMin, v.re.loopsMax, v.re.mods )
  }
}

func showCompile(t *testing.T) {
  re := new( RE )
  re.Compile( "<[:a]a>" )
  printASM( re )
  re.Compile( "#*cas[A-Z]" )
  printASM( re )
  re.Compile( "#^$<:b*:|(:|+#*:|)+>" )
  printASM( re )
}

func TestRegexp4(t *testing.T) {
  // showCompile( t )

  nTest( t )
  cTest( t )
  dTest( t )
  sTest( t )
  pTest( t )
  gTest( t )

  nTestUTF( t )
  cTestUTF( t )
  sTestUTF( t )
  pTestUTF( t )
  gTestUTF( t )
}

func nTest( t *testing.T ){
  numTest := []struct {
    txt, re string
    n int
  }{
    { "a", "a", 1 },
    { "aa", "aa", 1 },
    { "raptor", "raptor", 1 },
    { "a", "(a)", 1 },
    { "a", "<a>", 1 },
    { "a", "((a))", 1 },
    { "a", "<<a>>", 1 },
    { "a", "((((((a))))))", 1 },
    { "a", "<<<<<<a>>>>>>", 1 },
    { "a", "b|a", 1 },
    { "a", "c|b|a", 1 },
    { "a", "(b|a)", 1 },
    { "a", "(c|b|a)", 1 },
    { "a", "(c|b)|a", 1 },
    { "a", "((<c>|<b>)|a)", 1 },
    { "raptor", "b|raptor", 1 },
    { "raptor", "c|b|raptor", 1 },
    { "raptor", "(b|raptor)", 1 },
    { "raptor", "(c|raptor)|a", 1 },
    { "raptor", "((<c>|<raptor>)|a)", 1 },
    { "ab", "a(b|c)|A(B|C)", 1 },
    { "ac", "a(b|c)|A(B|C)", 1 },
    { "AB", "a(b|c)|A(B|C)", 1 },
    { "AC", "a(b|c)|A(B|C)", 1 },
    { "ab", "a<b|c>|A<B|C>", 1 },
    { "ac", "a<b|c>|A<B|C>", 1 },
    { "AB", "a<b|c>|A<B|C>", 1 },
    { "AC", "a<b|c>|A<B|C>", 1 },
    { "ab"    , "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "ac"    , "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "AB"    , "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "AC"    , "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "ab"    , "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "ac"    , "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "AB"    , "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "AC"    , "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "1234ea", "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "1234eb", "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "1234ec", "(a(b|c)|A(B|C))|1234(ea|eb|ec)", 1 },
    { "1234ea", "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "1234eb", "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "1234ec", "<a<b|c>|A<B|C>>|1234<ea|eb|ec>", 1 },
    { "abd", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "abe", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "acd", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "ace", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "ABD", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "ABE", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "ACD", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "ACE", "a(b|c)(d|e)|A(B|C)(D|E)", 1 },
    { "raptor", "(c|r)(e|a)(p|q)(t|u)(0|o)(t|r)", 1 },

    { "",  "",  0 },
    { "", "a",  0 },
    { "a", "",  0 },
    { "a", "o", 0 },
    { "a", "a", 1 },
    { "aaa", "a", 3 },
    { "a", "aaa", 0 },
    { "a aaa aaa", "aaa", 2 },
    { "Raptor Test", "a", 1 },
    { "Raptor Test", "t", 2 },
    { "aeiou", "a|e|i|o|u", 5 },
    { "aeiou", "(a|e|i|o|u)", 5 },
    { "aeiou", "(a|e)|i|(o|u)", 5 },
    { "aeiou", "(a|(e))|(i|(o|u))", 5 },
    { "aa ae ai ao au", "a(a|e|i|o|u)", 5 },
    { "aa ae ai ao au", "a(0|1|2|3|4)", 0 },
    { "a1 a2 a3 ao au", "a(1|2|3|4|5)", 3 },
    { "a1 a2 a3 a4 a5", "a(1|2|3|4|5)", 5 },
    { "aa ae ai ao au", "a(a|e|i|o|u) ", 4 },
    { "aa ae Ai ao au", "A(a|e|i|o|u)", 1 },
    { "aa ae Ai ao au", "(A|a)(a|e|i|o|u)", 5 },
    { "aae aei Aio aoa auu", "(A|a)(a|e|i|o|u)(a|e|i|o|u)", 5 },

    { "aa aaaa aaaa", "a", 10 },
    { "aa aaaa aaaa", "aa", 5 },
    { "aa aaaa aaaa", "aaa", 2 },
    { "aa aaaa aaaa", "aaaa", 2 },
    { "aaaaaaaaaaaaaaaaaaaa", "a", 20 },
    { "abababababababababababababababababababab", "a"  , 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "(a)", 20 },
    { "abababababababababababababababababababab", "(a)", 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "<a>", 20 },
    { "abababababababababababababababababababab", "<a>", 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a+",   1 },
    { "abababababababababababababababababababab", "a+" , 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a?", 20 },
    { "abababababababababababababababababababab", "a?" , 40 },
    { "aaaaaaaaaaaaaaaaaaaa", "a*", 1 },
    { "abababababababababababababababababababab", "a*" , 40 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{1}", 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{1}", 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{5}", 4 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{1,5}", 4 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{5,5}", 4 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{10}", 2 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{1,100}", 1 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{001,00100}", 1 },
    { "abababababababababababababababababababab", "a{1}" , 20 },
    { "abababababababababababababababababababab", "a{001}" , 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{1,1}", 20 },
    { "abababababababababababababababababababab", "a{1,1}" , 20 },
    { "abababababababababababababababababababab", "a{001,000001}" , 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "a{20}", 1 },
    { "abababababababababababababababababababab", "(a|b){1,1}" , 40 },

    { "aaaaaaaaaaaaaaaaaaaa", "a{1}b{0}", 20 },
    { "aaaaaaaaaaaaaaaaaaaa", "b{0}a{1}", 20 },

    { "aaaaaaaaaaaaaaaaaaaa", "b{0}", 20 },
    { "bbbbbbbbbbbbbbbbbbbb", "b{0}", 20 },
    { "bbbbbbbbbbbbbbbbbbbb", "b{1}", 20 },
    { "bbbbbbbbbbbbbbbbbbbb", "b{2}", 10 },

    { "abc", "<b>", 1 },
    { "abc", "a<b>", 1 },
    { "abc", "<b>c", 1 },
    { "abc", "a<b>c", 1 },
    { "abc", "<a|b>", 2 },
    { "abc", "<a|b|c>", 3 },
    { "abc", "<(a|b)|c>", 3 },
    { "aa aaaa aaaa", "<aa>", 5 },
    { "abc", "a<x>", 0 },
    { "abc", "<a>x", 0 },
    { "abc", "<a|b>x", 0 },
    { "abc", "<<a|b>x|abc>", 1 },
    { "abc", "<x<a|b>|abc>", 1 },
    { "abc abc abc", "<a|b|c>", 9 },
    { "abc abc abc", "<(a|b|c)(a|b|c)(a|b|c)>", 3 },
    { "abc abc abc", "<(a|b|c)(a|b|c)(a|b|c)> ", 2 },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 3 },

    { "a", "a?", 1 },
    { "a", "b?", 1 },
    { "a", "a+", 1 },
    { "a", "a*", 1 },
    { "a", "b*", 1 },
    { "a", "aa?", 1 },
    { "a", "ab?", 1 },
    { "a", "aa+", 0 },
    { "a", "aa*", 1 },
    { "a", "ab*", 1 },
    { "a", "a{1,2}", 1 },
    { "aaa", "a+", 1 },
    { "aaa", "a*", 1 },
    { "aaa", "a+", 1 },
    { "aaa", "a?", 3 },
    { "aaab", "a+", 1 },
    { "aaab", "a*", 2 },
    { "aaab", "a?", 4 },
    { "aaab", "a+b", 1 },
    { "aaab", "a*b", 1 },
    { "aaab", "a?b", 1 },
    { "aaab", "a+b?", 1 },
    { "aaab", "a*b?", 1 },
    { "aaab", "a?b?", 3 },
    { "aaab", "a+b+", 1 },
    { "aaab", "a*b+", 1 },
    { "aaab", "a?b+", 1 },
    { "aaab", "a+b*", 1 },
    { "aaab", "a*b*", 1 },
    { "aaab", "a?b*", 3 },
    { "aaabaaa", "a+", 2 },
    { "aaabaaa", "a*", 3 },
    { "aaabaaa", "a*", 3 },
    { "aaabaaa", "a*", 3 },
    { "a", "(a)?", 1 },
    { "a", "(b)?", 1 },
    { "a", "(a)+", 1 },
    { "a", "(a)*", 1 },
    { "a", "(b)*", 1 },
    { "aaa", "(a)+", 1 },
    { "aaa", "(a)*", 1 },

    { "Raptor Test",     "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1 },
    { "Raaaaptor TFest", "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1 },
    { "CaptorTest",      "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1 },
    { "Cap CaptorTest",  "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1 },
    { "Rap Captor Fest", "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1 },

    { "a", ":a", 1 },
    { "a", ":A", 0 },
    { "a", ":d", 0 },
    { "a", ":D", 1 },
    { "a", ":w", 1 },
    { "a", ":W", 0 },
    { "a", ":s", 0 },
    { "a", ":S", 1 },
    { "A", ":a", 1 },
    { "A", ":A", 0 },
    { "A", ":d", 0 },
    { "A", ":D", 1 },
    { "A", ":w", 1 },
    { "A", ":W", 0 },
    { "A", ":s", 0 },
    { "A", ":S", 1 },
    { "4", ":a", 0 },
    { "4", ":A", 1 },
    { "4", ":d", 1 },
    { "4", ":D", 0 },
    { "4", ":w", 1 },
    { "4", ":W", 0 },
    { "4", ":s", 0 },
    { "4", ":S", 1 },
    { " ", ":a", 0 },
    { " ", ":A", 1 },
    { " ", ":d", 0 },
    { " ", ":D", 1 },
    { " ", ":w", 0 },
    { " ", ":W", 1 },
    { " ", ":s", 1 },
    { " ", ":S", 0 },
    { "\t", ":a", 0 },
    { "\t", ":A", 1 },
    { "\t", ":d", 0 },
    { "\t", ":D", 1 },
    { "\t", ":w", 0 },
    { "\t", ":W", 1 },
    { "\t", ":s", 1 },
    { "\t", ":S", 0 },

    { "abc", ":a", 3 },
    { "abc", ":A", 0 },
    { "abc", ":d", 0 },
    { "abc", ":D", 3 },
    { "abc", ":w", 3 },
    { "abc", ":W", 0 },
    { "abc", ":s", 0 },
    { "abc", ":S", 3 },
    { "ABC", ":a", 3 },
    { "ABC", ":A", 0 },
    { "ABC", ":d", 0 },
    { "ABC", ":D", 3 },
    { "ABC", ":w", 3 },
    { "ABC", ":W", 0 },
    { "ABC", ":s", 0 },
    { "ABC", ":S", 3 },
    { "123", ":a", 0 },
    { "123", ":A", 3 },
    { "123", ":d", 3 },
    { "123", ":D", 0 },
    { "123", ":w", 3 },
    { "123", ":W", 0 },
    { "123", ":s", 0 },
    { "123", ":S", 3 },
    { " \n\t", ":a", 0 },
    { " \n\t", ":A", 3 },
    { " \n\t", ":d", 0 },
    { " \n\t", ":D", 3 },
    { " \n\t", ":w", 0 },
    { " \n\t", ":W", 3 },
    { " \n\t", ":s", 3 },
    { " \n\t", ":S", 0 },
    { " \n\t", ":a", 0 },
    { " \n\t", ":A", 3 },
    { " \n\t", ":d", 0 },
    { " \n\t", ":D", 3 },
    { " \n\t", ":w", 0 },
    { " \n\t", ":W", 3 },
    { " \n\t", ":s", 3 },
    { " \n\t", ":S", 0 },

    { "abc", ":a+", 1 },
    { "abc", ":A+", 0 },
    { "abc", ":d+", 0 },
    { "abc", ":D+", 1 },
    { "abc", ":w+", 1 },
    { "abc", ":W+", 0 },
    { "abc", ":s+", 0 },
    { "abc", ":S+", 1 },
    { "ABC", ":a+", 1 },
    { "ABC", ":A+", 0 },
    { "ABC", ":d+", 0 },
    { "ABC", ":D+", 1 },
    { "ABC", ":w+", 1 },
    { "ABC", ":W+", 0 },
    { "ABC", ":s+", 0 },
    { "ABC", ":S+", 1 },
    { "123", ":a+", 0 },
    { "123", ":A+", 1 },
    { "123", ":d+", 1 },
    { "123", ":D+", 0 },
    { "123", ":w+", 1 },
    { "123", ":W+", 0 },
    { "123", ":s+", 0 },
    { "123", ":S+", 1 },
    { " \n\t", ":a+", 0 },
    { " \n\t", ":A+", 1 },
    { " \n\t", ":d+", 0 },
    { " \n\t", ":D+", 1 },
    { " \n\t", ":w+", 0 },
    { " \n\t", ":W+", 1 },
    { " \n\t", ":s+", 1 },
    { " \n\t", ":S+", 0 },
    { " \n\t", ":a+", 0 },
    { " \n\t", ":A+", 1 },
    { " \n\t", ":d+", 0 },
    { " \n\t", ":D+", 1 },
    { " \n\t", ":w+", 0 },
    { " \n\t", ":W+", 1 },
    { " \n\t", ":s+", 1 },
    { " \n\t", ":S+", 0 },

    { "aeiou", ":a", 5 },

    { "(((", ":(", 3 },
    { ")))", ":)", 3 },
    { "<<<", ":<", 3 },
    { ">>>", ":>", 3 },
    { ":::", "::", 3 },
    { "|||", ":|", 3 },
    { "###", ":#", 3 },
    { ":#()|<>", ":::#:(:):|:<:>", 1 },
    { ":#()|<>", "(::|:#|:(|:)|:||:<|:>)", 7 },
    { "()<>[]{}*?+", "[:(:):<:>:[:]:{:}:*:?:+]", 11 },
    { "()<>[]|{}*#@?+", "[()<>:[:]|{}*?+#@]", 14 },
    { "12)", "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { "12.", "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { "a)" , "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { "a." , "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { "-"  , "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { "+"  , "#^:b*(-|:+|(:d+|:a+)[.)])", 1 },
    { ")>}", "[)>}]", 3 },
    { "(test1)(test2)", ":(test:d:)", 2 },

    { "",  ".",  0 },
    { "a", ".",  1 },
    { "aaa", ".", 3 },
    { "a", "...", 0 },
    { "a aaa aaa", ".", 9 },
    { "a aaa aaa", "...", 3 },
    { "a aaa aaa", ".aa", 2 },
    { "a aaa aaa", "aa.", 2 },
    { "Raptor Test", ".a", 1 },
    { "Raptor Test", ".t", 2 },
    { "Raptor Test", ".z", 0 },
    { "Raptor Test", "a.", 1 },
    { "Raptor Test", " .", 1 },
    { "Raptor Test", "z.", 0 },
    { "a", ".?", 1 },
    { "a", ".+", 1 },
    { "a", ".*", 1 },
    { "a", ".{1}", 1 },
    { "a aaa aaa", ".?", 9 },
    { "a aaa aaa", ".+", 1 },
    { "a aaa aaa", ".*", 1 },
    { "a aaa aaa", ".{1}", 9 },
    { "a", "a.?", 1 },
    { "a", "a.+", 0 },
    { "a", "a.*", 1 },
    { "a", "a.{1}", 0 },
    { "aeiou", "a|.", 5 },
    { "aeiou", "a|.?", 5 },
    { "aeiou", "a|.+", 2 },
    { "aeiou", "a|.*", 2 },
    { "aeiou", ".|a", 5 },
    { "aeiou", ".?|a", 5 },
    { "aeiou", ".+|a", 1 },
    { "aeiou", ".*|a", 1},
    { "aeiou", "(a|.)", 5 },
    { "aeiou", "(a|.?)", 5 },
    { "aeiou", "(a|.+)", 2 },
    { "aeiou", "(a|.*)", 2 },
    { "aeiou", "(.|a)", 5 },
    { "aeiou", "(.?|a)", 5 },
    { "aeiou", "(.+|a)", 1 },
    { "aeiou", "(.*|a)", 1},
    { "aeiou", "a|(.)", 5 },
    { "aeiou", "a|(.?)", 5 },
    { "aeiou", "a|(.+)", 2 },
    { "aeiou", "a|(.*)", 2 },
    { "aeiou", "(.)|a", 5 },
    { "aeiou", "(.?)|a", 5 },
    { "aeiou", "(.+)|a", 1 },
    { "aeiou", "(.*)|a", 1},
    { "aeiou", "a|(.)", 5 },
    { "aeiou", "a|(.)?", 5 },
    { "aeiou", "a|(.)+", 2 },
    { "aeiou", "a|(.)*", 2 },
    { "aeiou", "(.)|a", 5 },
    { "aeiou", "(.)?|a", 5 },
    { "aeiou", "(.)+|a", 1 },
    { "aeiou", "(.)*|a", 1},
    { "abababababababababababababababababababab", "." , 40 },
    { "abababababababababababababababababababab", "(a.)" , 20 },
    { "abababababababababababababababababababab", "(.a)" , 19 },
    { "abababababababababababababababababababab", "(:a.)" , 20 },
    { "abababababababababababababababababababab", "(.:a)" , 20 },
    { "abababababababababababababababababababab", "(.{5}:a{5})" , 4 },

    { "",  "a-z",  0 },
    { "a", "a-z",  0 },
    { "-", "-",  1 },
    { "-", "-a",  0 },
    { "-a", "-a",  1 },
    { "a-z", "a-z", 1 },
    { "A-Z", "A-Z", 1 },
    { "a-c", "a-z", 0 },
    { "A-c", "A-Z", 0 },
    { "a-zA-Z", "a-zA-Z", 1 },
    { "a-zB-Z", "a-zA-Z", 0 },
    { "a-", "a-z?", 1 },
    { "a-z", "a-z+", 1 },
    { "a-", "a-z*", 1 },
    { "a-z", "a-z?", 1 },
    { "a-zzzz", "a-z+", 1 },
    { "a-zz", "a-z*", 1 },
    { "a-b", "a-z?", 1 },
    { "a-bzzzz", "a-z+", 0 },
    { "a-bzz", "a-z*", 1 },

    { "",  "[a]",  0 },
    { "a", "[a]",  1 },
    { "a", "[.]",  0 },
    { ".", "[.]",  1 },
    { "a", "[A]",  0 },
    { "A", "[A]",  1 },
    { "1", "[A]",  0 },
    { "1", "[1]",  1 },
    { "a", "[:a]", 1 },
    { "A", "[:D]", 1 },
    { "aaa", "[a-z]", 3 },
    { "a", "[a-z][a-z][a-z]", 0 },
    { "a aaa aaa", "[a-z]", 7 },
    { "a aaa aaa", "[ a-z]", 9 },
    { "a aaa aaa", "[a-z][a-z][a-z]", 2 },
    { "a aaa aaa", "[a-z]aa", 2 },
    { "a aaa aaa", "aa[a-z]", 2 },
    { "Raptor Test", "[:w]a", 1 },
    { "Raptor Test", "[:w]t", 2 },
    { "Raptor Test", "[a-z]z", 0 },
    { "Raptor Test", "a[a-z]", 1 },
    { "Raptor Test", " [A-Z]", 1 },
    { "Raptor Test", "z[a-z]", 0 },
    { "a", "[a]?", 1 },
    { "a", "[a]+", 1 },
    { "a", "[a]*", 1 },
    { "a", "[a]{1}", 1 },
    { "a aaa aaa", "[a-z]?", 9 },
    { "a aaa aaa", "[a-z]+", 3 },
    { "a aaa aaa", "[a-z]*", 5 },
    { "a aaa aaa", "[a-z]{1}", 7 },
    { "a", "a[a-z]?", 1 },
    { "a", "a[a-z]+", 0 },
    { "a", "a[a-z]*", 1 },
    { "a", "a[a-z]{1}", 0 },
    { "aeiou", "a|[aeiou]", 5 },
    { "aeiou", "a|[aeiou]?", 5 },
    { "aeiou", "a|[aeiou]+", 2 },
    { "aeiou", "a|[aeiou]*", 2 },
    { "aeiou", "[aeiou]|a", 5 },
    { "aeiou", "[aeiou]?|a", 5 },
    { "aeiou", "[aeiou]+|a", 1 },
    { "aeiou", "[aeiou]*|a", 1},
    { "aeiou", "(a|[aeiou])", 5 },
    { "aeiou", "(a|[aeiou]?)", 5 },
    { "aeiou", "(a|[aeiou]+)", 2 },
    { "aeiou", "(a|[aeiou]*)", 2 },
    { "aeiou", "([aeiou]|a)", 5 },
    { "aeiou", "([aeiou]?|a)", 5 },
    { "aeiou", "([aeiou]+|a)", 1 },
    { "aeiou", "([aeiou]*|a)", 1},
    { "aeiou", "a|([aeiou])", 5 },
    { "aeiou", "a|([aeiou]?)", 5 },
    { "aeiou", "a|([aeiou]+)", 2 },
    { "aeiou", "a|([aeiou]*)", 2 },
    { "aeiou", "([aeiou])|a", 5 },
    { "aeiou", "([aeiou]?)|a", 5 },
    { "aeiou", "([aeiou]+)|a", 1 },
    { "aeiou", "([aeiou]*)|a", 1},
    { "aeiou", "a|([aeiou])", 5 },
    { "aeiou", "a|([aeiou])?", 5 },
    { "aeiou", "a|([aeiou])+", 2 },
    { "aeiou", "a|([aeiou])*", 2 },
    { "aeiou", "([aeiou])|a", 5 },
    { "aeiou", "([aeiou])?|a", 5 },
    { "aeiou", "([aeiou])+|a", 1 },
    { "aeiou", "([aeiou])*|a", 1},
    { "1a2a3a4a5a6a", "[1-6]a", 6 },
    { "1a2a3a4a5a6a", "[1-3]a", 3 },
    { "1a2b3c4d5e6f", "[123456][abcdef]", 6 },
    { "1a2b3c4d5e6f", "[123][abcdef]", 3 },
    { ".b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b", "[:.]",  20 },
    { ".b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b", "[:.b]",  40 },
    { ".b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b", "[.]",  20 },
    { ".b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b.b", "[.b]",  40 },
    { "abababababababababababababababababababab", "(a[ab])" , 20 },
    { "abababababababababababababababababababab", "([ab]a)" , 19 },
    { "abababababababababababababababababababab", "(:a[ab])" , 20 },
    { "abababababababababababababababababababab", "([ab]:a)" , 20 },
    { "abababababababababababababababababababab", "([ab]{5}:a{5})" , 4 },

    { "",  "[^a]",  0 },
    { "a", "[^1]",  1 },
    { "a", "[^a]",  0 },
    { "A", "[^a]",  1 },
    { "1", "[^1]",  0 },
    { "1", "[^A]",  1 },
    { "a", "[^:a]", 0 },
    { "A", "[^:A]", 1 },
    { "aaa", "[^z]", 3 },
    { "a", "[^z][^z][^z]", 0 },
    { "a aaa aaa", "[^ ]", 7 },
    { "a aaa aaa", "[^ a]", 0 },
    { "a aaa aaa", "[^:d]", 9 },
    { "a aaa aaa", "[^:d:s]", 7 },
    { "a aaa aaa", "[^:d:s][^:d:s][^:d:s]", 2 },
    { "a aaa aaa", "[^:d:s]aa", 2 },
    { "a aaa aaa", "aa[^:d:s]", 2 },
    { "Raptor Test", "[^:d:s]a", 1 },
    { "Raptor Test", "[^A-Z]t", 2 },
    { "Raptor Test", "[^:s]z", 0 },
    { "Raptor Test", "a[^ ]", 1 },
    { "Raptor Test", " [^t]", 1 },
    { "Raptor Test", "z[^a]", 0 },
    { "a", "[^z]?", 1 },
    { "a", "[^z]+", 1 },
    { "a", "[^z]*", 1 },
    { "a", "[^z]{1}", 1 },
    { "a aaa aaa", "[^ ]?", 9 },
    { "a aaa aaa", "[^ ]+", 3 },
    { "a aaa aaa", "[^ ]*", 5 },
    { "a aaa aaa", "[^ ]{1}", 7 },
    { "a", "a[^ ]?", 1 },
    { "a", "a[^ ]+", 0 },
    { "a", "a[^ ]*", 1 },
    { "a", "a[^ ]{1}", 0 },
    { "aeiou", "a|[^ ]", 5 },
    { "aeiou", "a|[^ ]?", 5 },
    { "aeiou", "a|[^ ]+", 2 },
    { "aeiou", "a|[^ ]*", 2 },
    { "aeiou", "[^ ]|a", 5 },
    { "aeiou", "[^ ]?|a", 5 },
    { "aeiou", "[^ ]+|a", 1 },
    { "aeiou", "[^ ]*|a", 1},
    { "aeiou", "(a|[^ ])", 5 },
    { "aeiou", "(a|[^ ]?)", 5 },
    { "aeiou", "(a|[^ ]+)", 2 },
    { "aeiou", "(a|[^ ]*)", 2 },
    { "aeiou", "([^ ]|a)", 5 },
    { "aeiou", "([^ ]?|a)", 5 },
    { "aeiou", "([^ ]+|a)", 1 },
    { "aeiou", "([^ ]*|a)", 1},
    { "aeiou", "a|([^ ])", 5 },
    { "aeiou", "a|([^ ]?)", 5 },
    { "aeiou", "a|([^ ]+)", 2 },
    { "aeiou", "a|([^ ]*)", 2 },
    { "aeiou", "([^ ])|a", 5 },
    { "aeiou", "([^ ]?)|a", 5 },
    { "aeiou", "([^ ]+)|a", 1 },
    { "aeiou", "([^ ]*)|a", 1},
    { "aeiou", "a|([^ ])", 5 },
    { "aeiou", "a|([^ ])?", 5 },
    { "aeiou", "a|([^ ])+", 2 },
    { "aeiou", "a|([^ ])*", 2 },
    { "aeiou", "([^ ])|a", 5 },
    { "aeiou", "([^ ])?|a", 5 },
    { "aeiou", "([^ ])+|a", 1 },
    { "aeiou", "([^ ])*|a", 1},
    { "1a2a3a4a5a6a", "[^:a]a", 6 },
    { "1a2a3a4a5a6a", "[^4-6]a", 3 },
    { "1a2b3c4d5e6f", "[^:a][^:d]", 6 },
    { "1a2b3c4d5e6f", "[^4-6][^:d]", 3 },
    { "abababababababababababababababababababab", "(a[^a])" , 20 },
    { "abababababababababababababababababababab", "([^a]a)" , 19 },
    { "abababababababababababababababababababab", "(:a[^a])" , 20 },
    { "abababababababababababababababababababab", "([^b]:a)" , 20 },
    { "abababababababababababababababababababab", "([^x]{5}:a{5})" , 4 },
    { "()<>[]{}*?+", "[^:w]", 11 },

    { "ABC", "#^A", 1 },
    { "ABC", "#^AB", 1 },
    { "ABC", "#^ABC", 1 },
    { "ABC", "#^(b|A)", 1 },
    { "ABC", "#^A(B|C)(B|C)", 1 },
    { "ABC", "#^(A(B|C))(B|C)", 1 },
    { "ABC", "#$C", 1 },
    { "ABC", "#$BC", 1 },
    { "ABC", "#$ABC", 1 },
    { "ABC", "#$(b|C)", 1 },
    { "ABC", "#$A(B|C)(B|C)", 1 },
    { "ABC", "#$(A(B|C))(B|C)", 1 },
    { "ABC", "#^$ABC", 1 },
    { "ABC", "#^$A(c|B)(b|C)", 1 },
    { "ABC", "#^$A(B|C)(B|C)", 1 },
    { "ABC", "#^$(A(B|C))(B|C)", 1 },
    { "ABC", "#^$AB([^C]+)", 0 },
    { "ABC", "#^$AB(A)+", 0 },

    { "ABC", "#^E", 0 },
    { "ABC", "#^EB", 0 },
    { "ABC", "#^EBC", 0 },
    { "ABC", "#^(b|E)", 0 },
    { "ABC", "#^A(B|C)(B|E)", 0 },
    { "ABC", "#^(A(B|C))(B|E)", 0 },
    { "ABC", "#$E", 0 },
    { "ABC", "#$BE", 0 },
    { "ABC", "#$ABE", 0 },
    { "ABC", "#$(b|E)", 0 },
    { "ABC", "#$A(B|C)(B|E)", 0 },
    { "ABC", "#$(A(B|C))(B|E)", 0 },
    { "ABC", "#^$ABE", 0 },
    { "ABC", "#^$A(c|B)(b|E)", 0 },
    { "ABC", "#^$A(B|C)(B|E)", 0 },
    { "ABC", "#^$(A(B|C))(B|E)", 0 },

    { "A", "a#*", 1 },
    { "A", "a?#*", 1 },
    { "A", "b?#*", 1 },
    { "A", "a+#*", 1 },
    { "A", "a*#*", 1 },
    { "A", "b*#*", 1 },
    { "aAa", "a#*", 3 },
    { "aAa", "a+#*", 1 },
    { "aAa", "a*#*", 1 },
    { "aAa", "a+#*", 1 },
    { "aAa", "a?#*", 3 },
    { "aAab", "a+#*", 1 },
    { "aAab", "a*#*", 2 },
    { "aAab", "a?#*", 4 },
    { "aAab", "a+#*?^$~b", 1 },
    { "aAab", "a*#*?^$~b", 1 },
    { "aAab", "a?#*?^$~b", 1 },
    { "aAab", "a+#*?^$~b?", 1 },
    { "aAab", "a*#*?^$~b?", 1 },
    { "aAab", "a?#*?^$~b?", 3 },
    { "aAab", "a+#*?^$~b+", 1 },
    { "aAab", "a*#*?^$~b+", 1 },
    { "aAab", "a?#*?^$~b+", 1 },
    { "aAab", "a+#*?^$~b*", 1 },
    { "aAab", "a*#*?^$~b*", 1 },
    { "aAab", "a?#*?^$~b*", 3 },

    { "a", "a#*", 1 },
    { "a", "A#*", 1 },
    { "a", "#*A", 1 },
    { "a", "#*a", 1 },
    { "a", "#*(A)", 1 },
    { "a", "#*(a)", 1 },
    { "a", "#*[A]", 1 },
    { "a", "#*[a]", 1 },
    { "a-Z", "#*A-Z", 1 },
    { "a-z", "a-Z#*", 1 },
    { "a-Z", "a-Z#*", 1 },
    { "a", "(a)#*", 1 },
    { "a", "(A)#*", 1 },
    { "a", "[a]#*", 1 },
    { "a", "[A]#*", 1 },
    { "a", "#*[A-Z]", 1 },
    { "a", "[A-Z]#*", 1 },

    { "aAaA", "a#*", 4 },
    { "aAaA", "A#*", 4 },
    { "aAaA", "#*A", 4 },
    { "aAaA", "#*a", 4 },
    { "aAaA", "#*(A)", 4 },
    { "aAaA", "#*(a)", 4 },
    { "aAaA", "#*[A]", 4 },
    { "aAaA", "#*[a]", 4 },
    { "aAaA", "#*Aa", 2 },
    { "aAaa", "aA#*", 2 },
    { "aAaA", "(a)#*", 4 },
    { "aAaA", "(A)#*", 4 },
    { "aAaA", "[a]#*", 4 },
    { "aAaA", "[A]#*", 4 },
    { "aAaA", "#*[A-Z]", 4 },
    { "aAaA", "[A-Z]#*", 4 },
    { "aAaA", "(a#*)", 4 },
    { "aAaA", "(A#*)", 4 },
    { "aAaA", "(a)#*", 4 },
    { "aAaA", "(A)#*", 4 },
    { "aAbB", "#*a|b", 4 },
    { "aAbB", "#*A|B", 4 },
    { "aAbB", "#*(a|b)", 4 },
    { "aAbB", "#*(A|B)", 4 },
    { "aAbB", "(a#*|b#*)", 4 },
    { "aAbB", "(A#*|B#*)", 4 },
    { "aAbB", "(a|b)#*", 4 },
    { "aAbB", "(A|B)#*", 4 },
    { "TesT", "test", 0 },
    { "TesT", "test#*", 0 },
    { "TesT", "t#*est#*", 1 },
    { "TesT", "#*test", 1 },
    { "TesT", "#*tESt", 1 },
    { "TesT", "#*(tESt)", 1 },
    { "TesT", "(tESt)#*", 1 },

    { "a aaa aaa", "#^aaa", 0 },
    { "a aaa aaa", "#$aaa", 1 },
    { "a aaa aaa", "#?aaa", 1 },
    { "a aaa aaa", "#~aaa", 2 },
    { "a aaa aaa", "#^?aaa", 0 },
    { "a aaa aaa", "#?^aaa", 0 },
    { "a aaa aaa", "#?$aaa", 1 },
    { "a aaa aaa", "#^?$aaa", 0 },
    { "a aaa aaa", "#?$^aaa", 0 },
    { "a aaa aaa", "#^?$a aaa aaa", 1 },
    { "aa aaaa aaaa", "#~a", 10 },
    { "aa aaaa aaaa", "#~aa", 7 },
    { "aa aaaa aaaa", "#~aaa", 4 },
    { "aaaaaaaaaaaaaaaaaaaa", "#?a+", 1 },
    { "abababababababababababababababababababab", "#?a+" , 1 },
    { "aaaaaaaaaaaaaaaaaaaa", "#~a+", 20 },
    { "abababababababababababababababababababab", "#~a+" , 20 },

    { "Raptor TesT Fest", "RapTor (tESt)#* fEST", 0 },
    { "Raptor TesT Fest", "#*rapTor (tESt) fEST", 1 },
    { "Raptor TesT Fest", "(RapTor)#* (tESt)#* (fEST)#*", 1 },
    { "Raptor TesT Fest", "((Rap#*Tor)#* (t#*ESt)#* (fEST)#*)#*", 1 },
    { "Raptor TesT Fest", "#*[a-z]#*apTor (tESt) [A-Z]#*EST", 1 },

    { "a", "a#/", 1 },
    { "a", "A#/", 0 },
    { "a", "#/A", 0 },
    { "a", "#/a", 1 },
    { "a", "#/(A)", 0 },
    { "a", "#/(a)", 1 },
    { "a", "#/[A]", 0 },
    { "a", "#/[a]", 1 },
    { "a", "#/A-Z", 0 },
    { "a", "A-Z#/", 0 },
    { "a", "(a)#/", 1 },
    { "a", "(A)#/", 0 },
    { "a", "[a]#/", 1 },
    { "a", "[A]#/", 0 },
    { "a", "#/[A-Z]", 0 },
    { "a", "[A-Z]#/", 0 },

    { "aAaA", "a#/", 2 },
    { "aAaA", "A#/", 2 },
    { "aAaA", "#/A", 2 },
    { "aAaA", "#/a", 2 },
    { "aAaA", "#/(A)", 2 },
    { "aAaA", "#/(a)", 2 },
    { "aAaA", "#/[A]", 2 },
    { "aAaA", "#/[a]", 2 },
    { "aAaA", "#/Aa", 1 },
    { "aAaA", "#/aA", 2 },
    { "aAaA", "Aa#/", 1 },
    { "aAaA", "aA#/", 2 },
    { "aAaA", "(a)#/", 2 },
    { "aAaA", "(A)#/", 2 },
    { "aAaA", "[a]#/", 2 },
    { "aAaA", "[A]#/", 2 },
    { "aAaA", "#/[A-Z]", 2 },
    { "aAaA", "[A-Z]#/", 2 },
    { "aAaA", "(a#/)", 2 },
    { "aAaA", "(A#/)", 2 },
    { "aAaA", "(a)#/", 2 },
    { "aAaA", "(A)#/", 2 },
    { "aAbB", "#/a|b", 2 },
    { "aAbB", "#/A|B", 2 },
    { "aAbB", "#/(a|b)", 2 },
    { "aAbB", "#/(A|B)", 2 },
    { "aAbB", "(a#/|b#/)", 2 },
    { "aAbB", "(A#/|B#/)", 2 },
    { "aAbB", "(a|b)#/", 2 },
    { "aAbB", "(A|B)#/", 2 },

    { "Raptor TesT Fest", "#*rapTor (tESt)#/ fEST", 0 },
    { "Raptor tESt Fest", "#*rapTor (tESt)#/ fEST", 1 },
    { "Raptor TesT Fest", "#*rapTor (tE#/S#/t)#* fEST", 0 },
    { "Raptor tESt Fest", "#*rapTor (tE#/S#/t)#* fEST", 1 },

    { "a aaa aaa", "[^ a]", 0 },
    { "a aaa aaa", "[^:d:s]", 7 },
    { "a aaa aaa", "[^:d:s][^:d:s][^:d:s]", 2 },
    { "a aaa aaa", "[^:d:s]aa", 2 },
    { "a aaa aaa", "aa[^:d:s]", 2 },
    { "Raptor Test", "[^:d:s]a", 1 },
    { "Raptor Test", "[^A-Z]t", 2 },
    { "a", "a[^ ]?", 1 },
    { "a", "a[^ ]+", 0 },
    { "a", "a[^ ]*", 1 },
    { "a", "a[^ ]{1}", 0 },
    { "aeiou", "a|[^ ]", 5 },
    { "aeiou", "a|[^ ]?", 5 },
    { "aeiou", "a|[^ ]+", 2 },
    { "aeiou", "a|[^ ]*", 2 },
    { "aeiou", "[^ ]|a", 5 },
    { "aeiou", "[^ ]?|a", 5 },
    { "aeiou", "[^ ]+|a", 1 },
    { "aeiou", "[^ ]*|a", 1},
    { "aeiou", "(a|[^ ])", 5 },
    { "aeiou", "(a|[^ ]?)", 5 },
    { "aeiou", "(a|[^ ]+)", 2 },
    { "aeiou", "(a|[^ ]*)", 2 },
    { "aeiou", "([^ ]|a)", 5 },
    { "aeiou", "([^ ]?|a)", 5 },
    { "aeiou", "([^ ]+|a)", 1 },
    { "aeiou", "([^ ]*|a)", 1},
    { "aeiou", "a|([^ ])", 5 },
    { "aeiou", "a|([^ ]?)", 5 },
    { "aeiou", "a|([^ ]+)", 2 },
    { "aeiou", "a|([^ ]*)", 2 },
    { "aeiou", "([^ ])|a", 5 },
    { "aeiou", "([^ ]?)|a", 5 },
    { "aeiou", "([^ ]+)|a", 1 },
    { "aeiou", "([^ ]*)|a", 1},
    { "aeiou", "a|([^ ])", 5 },
    { "aeiou", "a|([^ ])?", 5 },
    { "aeiou", "a|([^ ])+", 2 },
    { "aeiou", "a|([^ ])*", 2 },
    { "aeiou", "([^ ])|a", 5 },
    { "aeiou", "([^ ])?|a", 5 },
    { "aeiou", "([^ ])+|a", 1 },
    { "aeiou", "([^ ])*|a", 1},

    { "31/13-1331", "<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>", 0 },
    { "71-17/1177", "<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>", 0 },

    { "",  "@1",  0 },
    { "a", "@1",  0 },
    { "a", "@a",  0 },
    { "A", "@100", 0 },
    { "1", "@1",  0 },
    { "",  "[@1]",  0 },
    { "a", "[@1]",  0 },
    { "a", "[@a]",  1 },
    { "A", "[@100]", 0 },
    { "1", "[@1]",  1 },
    { "@", "[@1]",  1 },
    { "1@@a", "[a@1]",  4 },
    { "",  "[^@1]",  0 },
    { "a", "[^@1]",  1 },
    { "a", "[^@a]",  0 },
    { "A", "[^@100]", 1 },
    { "1", "[^@1]",  0 },
    { "@", "[^2@]",  0 },

    { "",  "(@1)",  0 },
    { "a", "(@1)",  0 },
    { "a", "(@a)",  0 },
    { "A", "(@100)", 0 },
    { "1", "(@1)",  0 },
    { "",  "([@1])",  0 },
    { "a", "([@1])",  0 },
    { "a", "([@a])",  1 },
    { "A", "([@100])", 0 },
    { "1", "([@1])",  1 },

    { "a", "<a>@1",  0 },
    { "a", "<a>@1?", 1 },
    { "a", "<a>@1*", 1 },
    { "a", "<a>@1+", 0 },
    { "a", "<a>@1{1}", 0 },
    { "aa", "<a>@1",  1 },
    { "aa", "<a>@1?", 1 },
    { "aa", "<a>@1+", 1 },
    { "aa", "<a>@1*", 1 },
    { "aa", "<a>@1{1}", 1 },
    { "aaaaa", "<a>@1",  2 },
    { "aaaaa", "<a>@1?", 3 },
    { "aaaaa", "<a>@1+", 1 },
    { "aaaaa", "<a>@1*", 1 },
    { "aaaaa", "<a>@1{1}", 2 },

    { "a-a", "<a|:d|o_O!>:-@1",  1 },
    { "1-1", "<a|:d|o_O!>:-@1", 1 },
    { "o_O!-o_O!", "<a|:d|o_O!>:-@1", 1 },

    { "ae_ea", "<a><e>_@2@1", 1 },
    { "ae_ea", "<<a><e>>_@2@1", 0 },
    { "ae_aae", "<<a><e>>_@2@1", 1 },
    { "ae_eaae_ea", "<a><e>_@2@1", 2 },
    { "ae_eaae_ea", "<<a><e>>_@2@1", 0 },
    { "ae_aaeae_aae", "<<a><e>>_@2@1", 2 },
    { "ae_aaeae_aa1", "<<a><e>>_@2@1", 1 },
    { "aaaaa", "@1<a>", 0 },

    { "012345678910012345678910", "<0><1><2><3><4><5><6><7><8><9><10>@1@2@3@4@5@6@7@8@9@10@11", 1 },
  }

  done := make(chan struct{})
  for _, c := range numTest {
    go func( txt, re string, n int ){
      r := Compile( re )
      x := r.MatchString( txt )

      if x != n  {
        t.Errorf( "Regexp4( %q, %q ) == %d, expected %d", txt, re, x, n )
      }
      done <- struct{}{}
    }( c.txt, c.re, c.n )
  }

  for range numTest { <-done }
}

func cTest( t *testing.T ){
  catchTest := []struct {
    txt, re string
    n int
    catch string
  }{
    { "a", "<a>", 1, "a" },
    { "a", "<a>", 1, "a" },
    { "aa", "<aa>", 1, "aa" },
    { "a a a", "<a>", 2, "a" },
    { "abcd", "<a|b|c|d>", 1, "a" },
    { "abcd", "<a|b|c|d>", 2, "b" },
    { "abcd", "<a|b|c|d>", 3, "c" },
    { "abcd", "<a|b|c|d>", 4, "d" },
    { "abcd", "<a|b|c|d>", 5, "" },
    { "abc", "a<x>", 1, "" },
    { "abc", "<a>x", 1, "" },
    { "abc", "<a|b>x", 1, "" },
    { "abc", "<<a|b>x|abc>", 1, "abc" },
    { "abc", "<<a|b>x|abc>", 2, "" },
    { "abc", "<x<a|b>|abc>", 1, "abc" },
    { "abc", "<x<a|b>|abc>", 2, "" },
    { "abc abc abc", "<a|b|c>", 9, "c" },
    { "abc abc abc", "<(a|b|c)(a|b|c)(a|b|c)>", 3, "abc" },
    { "abc abc abc", "<(a|b|c)(a|b|c)(a|b|c)> ", 2, "abc" },
    { "abc abc abc", "#?<(a|b|c)(a|b|c)(a|b|c)>", 1, "abc" },
    { "abc abc abc", "#?<(a|b|c)(a|b|c)((a|b)|x)>", 1, "" },
    { "abc abc abx", "#?<(a|b|c)(a|b|c)((a|b)|x)>", 1, "abx" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 1, "abc" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 2, "iec" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 3, "i" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 4, "c" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 5, "oeb" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 6, "o" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 7, "b" },
    { "abc iecc oeb", "<<(a|e)|(i|o)>e<b|c>|abc>", 8, "" },

    { "A", "#$<.{5}>", 1, "" },
    { "AB", "#$<.{5}>", 1, "" },
    { "ABC", "#$<.{5}>", 1, "" },
    { "ABCD", "#$<.{5}>", 1, "" },
    { "ABCDE", "#$<.{5}>", 1, "ABCDE" },
    { "ABCDEF", "#$<.{5}>", 1, "BCDEF" },
    { "ABCDEFG", "#$<.{5}>", 1, "CDEFG" },
    { "ABCDEFGH", "#$<.{5}>", 1, "DEFGH" },
    { "ABCDEFGHI", "#$<.{5}>", 1, "EFGHI" },
    { "ABCDEFGHIJ", "#$<.{5}>", 1, "FGHIJ" },
    { "ABCDEFGHIJK", "#$<.{5}>", 1, "GHIJK" },
    { "ABCDEFGHIJKL", "#$<.{5}>", 1, "HIJKL" },
    { "ABCDEFGHIJKLM", "#$<.{5}>", 1, "IJKLM" },
    { "ABCDEFGHIJKLMN", "#$<.{5}>", 1, "JKLMN" },
    { "ABCDEFGHIJKLMNO", "#$<.{5}>", 1, "KLMNO" },
    { "ABCDEFGHIJKLMNOP", "#$<.{5}>", 1, "LMNOP" },
    { "ABCDEFGHIJKLMNOPQ", "#$<.{5}>", 1, "MNOPQ" },
    { "ABCDEFGHIJKLMNOPQR", "#$<.{5}>", 1, "NOPQR" },
    { "ABCDEFGHIJKLMNOPQRS", "#$<.{5}>", 1, "OPQRS" },
    { "ABCDEFGHIJKLMNOPQRST", "#$<.{5}>", 1, "PQRST" },
    { "ABCDEFGHIJKLMNOPQRSTU", "#$<.{5}>", 1, "QRSTU" },
    { "ABCDEFGHIJKLMNOPQRSTUV", "#$<.{5}>", 1, "RSTUV" },
    { "ABCDEFGHIJKLMNOPQRSTUVW", "#$<.{5}>", 1, "STUVW" },
    { "ABCDEFGHIJKLMNOPQRSTUVWX", "#$<.{5}>", 1, "TUVWX" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXY", "#$<.{5}>", 1, "UVWXY" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "#$<.{5}>", 1, "VWXYZ" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ[", "#$<.{5}>", 1, "WXYZ[" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ][", "#$<.{5}>", 1, "XYZ][" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ^][", "#$<.{5}>", 1, "YZ^][" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_^][", "#$<.{5}>", 1, "Z_^][" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][", "#$<.{5}>", 1, "_`^][" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][a", "#$<.{5}>", 1, "`^][a" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][ab", "#$<.{5}>", 1, "^][ab" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abc", "#$<.{5}>", 1, "][abc" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcd", "#$<.{5}>", 1, "[abcd" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcde", "#$<.{5}>", 1, "abcde" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdef", "#$<.{5}>", 1, "bcdef" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefg", "#$<.{5}>", 1, "cdefg" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefgh", "#$<.{5}>", 1, "defgh" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghi", "#$<.{5}>", 1, "efghi" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghij", "#$<.{5}>", 1, "fghij" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijk", "#$<.{5}>", 1, "ghijk" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijkl", "#$<.{5}>", 1, "hijkl" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklm", "#$<.{5}>", 1, "ijklm" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmn", "#$<.{5}>", 1, "jklmn" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmno", "#$<.{5}>", 1, "klmno" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnop", "#$<.{5}>", 1, "lmnop" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopq", "#$<.{5}>", 1, "mnopq" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqr", "#$<.{5}>", 1, "nopqr" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrs", "#$<.{5}>", 1, "opqrs" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrst", "#$<.{5}>", 1, "pqrst" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstu", "#$<.{5}>", 1, "qrstu" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuv", "#$<.{5}>", 1, "rstuv" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvw", "#$<.{5}>", 1, "stuvw" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwx", "#$<.{5}>", 1, "tuvwx" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxy", "#$<.{5}>", 1, "uvwxy" },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxyz", "#$<.{5}>", 1, "vwxyz" },

    { "a", "<a>?", 1, "a" },
    { "a", "<b>?", 1, "" },
    { "a", "<a>+", 1, "a" },
    { "a", "<a>*", 1, "a" },
    { "a", "<b>*", 1, "" },
    { "aaa", "<a>+", 1, "aaa" },
    { "aaa", "<a>*", 1, "aaa" },
    { "aaa", "#~<a>+", 1, "aaa" },
    { "aaa", "#~<a>*", 1, "aaa" },
    { "aaab", "#~<a+>", 1, "aaa" },
    { "aaab", "#~<a*>", 1, "aaa" },
    { "aaab", "#~<a?>", 4, "" },
    { "aaab", "#~<a+b>", 1, "aaab" },
    { "aaab", "#~<a*b>", 1, "aaab" },
    { "aaab", "#~<a?b>", 1, "ab" },
    { "aaab", "#~<a+b?>", 1, "aaab" },
    { "aaab", "#~<a*b?>", 1, "aaab" },
    { "aaab", "#~<a?b?>", 3, "ab" },
    { "aaab", "#~<a+b+>", 1, "aaab" },
    { "aaab", "#~<a*b+>", 1, "aaab" },
    { "aaab", "#~<a?b+>", 1, "ab" },
    { "aaab", "#~<a+b*>", 1, "aaab" },
    { "aaab", "#~<a*b*>", 1, "aaab" },
    { "aaab", "#~<a?b*>", 3, "ab" },
    { "aaabaaa", "#~<a+>", 4, "aaa" },
    { "aaabaaa", "#~<a*>", 5, "aaa" },

    { "Raptor Test",     "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1, "Raptor Test" },
    { "Raptor Test",     "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 2, "T" },
    { "Raaaaptor TFest", "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1, "Raaaaptor TFest" },
    { "Raaaaptor TFest", "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 2, "TF" },
    { "CaptorTest",      "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1, "CaptorTest" },
    { "Cap CaptorTest",  "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1, "Cap CaptorTest" },
    { "Cap CaptorTest",  "#~<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 3, "CaptorTest" },
    { "Rap Captor Fest", "<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 1, "Rap Captor Fest" },
    { "Rap Captor Fest", "#~<((C|R)ap C|C|R)(a+p{1}tor) ?((<T|F>+e)(st))>", 3, "Captor Fest" },
    { "012345678910109876501234", "<0><1><2><3><4><5><6><7><8><9><10><@11@10@9@8@7@6@1@2@3@4@5>", 12, "109876501234" },


    { "| | text", "#^$<:s*>", 1, "" },
  }

  done := make(chan struct{})
  for _, c := range catchTest {
    go func( txt, re string, nCatch int, eCatch string ){
      var r RE
      r.Compile( re )
      r.FindString( txt )
      catch := r.GetCatch( nCatch )
      if catch != eCatch  {
        t.Errorf( "Regexp4( %q, %q )\nGetCatch( %d ) == %q, expected %q",
                  txt, re, nCatch, catch, eCatch )
      }
      done <- struct{}{}
    }( c.txt, c.re, c.n, c.catch )
  }

  for range catchTest { <-done }
}

func dTest( t *testing.T ){
  catchTest := []struct {
    txt, re string
    n int
  }{
    { "A", "<.>", 1 },
    { "AB", "<.>", 2 },
    { "ABC", "<.>", 3 },
    { "ABCD", "<.>", 4 },
    { "ABCDE", "<.>", 5 },
    { "ABCDEF", "<.>", 6 },
    { "ABCDEFG", "<.>", 7 },
    { "ABCDEFGH", "<.>", 8 },
    { "ABCDEFGHI", "<.>", 9 },
    { "ABCDEFGHIJ", "<.>", 10 },
    { "ABCDEFGHIJK", "<.>", 11 },
    { "ABCDEFGHIJKL", "<.>", 12 },
    { "ABCDEFGHIJKLM", "<.>", 13 },
    { "ABCDEFGHIJKLMN", "<.>", 14 },
    { "ABCDEFGHIJKLMNO", "<.>", 15 },
    { "ABCDEFGHIJKLMNOP", "<.>", 16 },
    { "ABCDEFGHIJKLMNOPQ", "<.>", 17 },
    { "ABCDEFGHIJKLMNOPQR", "<.>", 18 },
    { "ABCDEFGHIJKLMNOPQRS", "<.>", 19 },
    { "ABCDEFGHIJKLMNOPQRST", "<.>", 20 },
    { "ABCDEFGHIJKLMNOPQRSTU", "<.>", 21 },
    { "ABCDEFGHIJKLMNOPQRSTUV", "<.>", 22 },
    { "ABCDEFGHIJKLMNOPQRSTUVW", "<.>", 23 },
    { "ABCDEFGHIJKLMNOPQRSTUVWX", "<.>", 24 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXY", "<.>", 25 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "<.>", 26 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ[", "<.>", 27 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ][", "<.>", 28 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ^][", "<.>", 29 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_^][", "<.>", 30 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][", "<.>", 31 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][a", "<.>", 32 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][ab", "<.>", 33 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abc", "<.>", 34 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcd", "<.>", 35 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcde", "<.>", 36 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdef", "<.>", 37 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefg", "<.>", 38 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefgh", "<.>", 39 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghi", "<.>", 40 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghij", "<.>", 41 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijk", "<.>", 42 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijkl", "<.>", 43 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklm", "<.>", 44 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmn", "<.>", 45 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmno", "<.>", 46 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnop", "<.>", 47 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopq", "<.>", 48 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqr", "<.>", 49 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrs", "<.>", 50 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrst", "<.>", 51 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstu", "<.>", 52 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuv", "<.>", 53 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvw", "<.>", 54 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwx", "<.>", 55 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxy", "<.>", 56 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxyz", "<.>", 57 },

    { "A", "<:S>", 1 },
    { "AB", "<:S>", 2 },
    { "ABC", "<:S>", 3 },
    { "ABCD", "<:S>", 4 },
    { "ABCDE", "<:S>", 5 },
    { "ABCDEF", "<:S>", 6 },
    { "ABCDEFG", "<:S>", 7 },
    { "ABCDEFGH", "<:S>", 8 },
    { "ABCDEFGHI", "<:S>", 9 },
    { "ABCDEFGHIJ", "<:S>", 10 },
    { "ABCDEFGHIJK", "<:S>", 11 },
    { "ABCDEFGHIJKL", "<:S>", 12 },
    { "ABCDEFGHIJKLM", "<:S>", 13 },
    { "ABCDEFGHIJKLMN", "<:S>", 14 },
    { "ABCDEFGHIJKLMNO", "<:S>", 15 },
    { "ABCDEFGHIJKLMNOP", "<:S>", 16 },
    { "ABCDEFGHIJKLMNOPQ", "<:S>", 17 },
    { "ABCDEFGHIJKLMNOPQR", "<:S>", 18 },
    { "ABCDEFGHIJKLMNOPQRS", "<:S>", 19 },
    { "ABCDEFGHIJKLMNOPQRST", "<:S>", 20 },
    { "ABCDEFGHIJKLMNOPQRSTU", "<:S>", 21 },
    { "ABCDEFGHIJKLMNOPQRSTUV", "<:S>", 22 },
    { "ABCDEFGHIJKLMNOPQRSTUVW", "<:S>", 23 },
    { "ABCDEFGHIJKLMNOPQRSTUVWX", "<:S>", 24 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXY", "<:S>", 25 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "<:S>", 26 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ[", "<:S>", 27 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ][", "<:S>", 28 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ^][", "<:S>", 29 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_^][", "<:S>", 30 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][", "<:S>", 31 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][a", "<:S>", 32 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][ab", "<:S>", 33 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abc", "<:S>", 34 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcd", "<:S>", 35 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcde", "<:S>", 36 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdef", "<:S>", 37 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefg", "<:S>", 38 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefgh", "<:S>", 39 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghi", "<:S>", 40 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghij", "<:S>", 41 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijk", "<:S>", 42 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijkl", "<:S>", 43 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklm", "<:S>", 44 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmn", "<:S>", 45 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmno", "<:S>", 46 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnop", "<:S>", 47 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopq", "<:S>", 48 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqr", "<:S>", 49 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrs", "<:S>", 50 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrst", "<:S>", 51 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstu", "<:S>", 52 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuv", "<:S>", 53 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvw", "<:S>", 54 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwx", "<:S>", 55 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxy", "<:S>", 56 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ_`^][abcdefghijklmnopqrstuvwxyz", "<:S>", 57 },

    { "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@", "<:S>", 1024 },

    { "A", "#$<.{5}>", 0 },
    { "AB", "#$<.{5}>", 0 },
    { "ABC", "#$<.{5}>", 0 },
    { "ABCD", "#$<.{5}>", 0 },
    { "ABCDE", "#$<.{5}>", 1 },
    { "ABCDEF", "#$<.{5}>", 1 },
    { "ABCDEFG", "#$<.{5}>", 1 },
    { "ABCDEFGH", "#$<.{5}>", 1 },
    { "ABCDEFGHI", "#$<.{5}>", 1 },
  }

  done := make(chan struct{})
  for _, c := range catchTest {
    go func( txt, re string, n int ){
      var r RE
      r.Match( txt, re )
      x := r.TotCatch()
      if x != n  {
        t.Errorf( "Regexp4( %q, %q )\nTotCatch() == %d, expected %d",
                  txt, re, x, n )
      }
      done <- struct{}{}
    }( c.txt, c.re, c.n )
  }

  for range catchTest { <-done }
}

func sTest( t *testing.T ){
  swapTest := []struct {
    txt, re string
    n int
    swap, expected string
  }{
    { "aaab", "<a>"   , 1, "e", "eeeb" },
    { "a", "<a>"   , 1, "", "" },
    { "a", "<a>?"  , 1, "", "" },
    { "a", "<a>+"  , 1, "", "" },
    { "a", "<a>*"  , 1, "", "" },
    { "a", "<a>{1}", 1, "", "" },
    { "a", "<a?>"  , 1, "", "" },
    { "a", "<a+>"  , 1, "", "" },
    { "a", "<a*>"  , 1, "", "" },
    { "a", "<a{1}>", 1, "", "" },

    { "a", "<a>"   , 1, "e", "e" },
    { "a", "<a>?"  , 1, "e", "e" },
    { "a", "<a>+"  , 1, "e", "e" },
    { "a", "<a>*"  , 1, "e", "e" },
    { "a", "<a>{1}", 1, "e", "e" },
    { "a", "<a?>"  , 1, "e", "e" },
    { "a", "<a+>"  , 1, "e", "e" },
    { "a", "<a*>"  , 1, "e", "e" },
    { "a", "<a{1}>", 1, "e", "e" },

    { "a", "<x>"   , 1, "z", "a" },
    { "a", "<x>?"  , 1, "z", "za" },
    { "a", "<x>+"  , 1, "z", "a" },
    { "a", "<x>*"  , 1, "z", "za" },
    { "a", "<x>{1}", 1, "z", "a" },
    { "a", "<x?>"  , 1, "z", "za" },
    { "a", "<x+>"  , 1, "z", "a" },
    { "a", "<x*>"  , 1, "z", "za" },
    { "a", "<x{1}>", 1, "z", "a" },

    { "aaa", "<a>"   , 1, "", "" },
    { "aaa", "<a>?"  , 1, "", "" },
    { "aaa", "<a>+"  , 1, "", "" },
    { "aaa", "<a>*"  , 1, "", "" },
    { "aaa", "<a>{1}", 1, "", "" },
    { "aaa", "<a?>"  , 1, "", "" },
    { "aaa", "<a+>"  , 1, "", "" },
    { "aaa", "<a*>"  , 1, "", "" },
    { "aaa", "<a{1}>", 1, "", "" },

    { "aaa", "<a>"   , 1, "e", "eee" },
    { "aaa", "<a>?"  , 1, "e", "eee" },
    { "aaa", "<a>+"  , 1, "e", "e" },
    { "aaa", "<a>*"  , 1, "e", "e" },
    { "aaa", "<a>{1}", 1, "e", "eee" },
    { "aaa", "<a?>"  , 1, "e", "eee" },
    { "aaa", "<a+>"  , 1, "e", "e" },
    { "aaa", "<a*>"  , 1, "e", "e" },
    { "aaa", "<a{1}>", 1, "e", "eee" },

    { "aaa", "<x>"   , 1, "z", "aaa" },

    { "aaa", "<x>?"  , 1, "z", "zazaza" },
    { "aaa", "<x>+"  , 1, "z", "aaa" },
    { "aaa", "<x>*"  , 1, "z", "zazaza" },
    { "aaa", "<x>{1}", 1, "z", "aaa" },
    { "aaa", "<x?>"  , 1, "z", "zazaza" },
    { "aaa", "<x+>"  , 1, "z", "aaa" },
    { "aaa", "<x*>"  , 1, "z", "zazaza" },
    { "aaa", "<x{1}>", 1, "z", "aaa" },

    { "aaab", "<a>"   , 1, "e", "eeeb" },
    { "aaab", "<a>?"  , 1, "e", "eeeeb" },
    { "aaab", "<a>+"  , 1, "e", "eb" },
    { "aaab", "<a>*"  , 1, "e", "eeb" },
    { "aaab", "<a>{1}", 1, "e", "eeeb" },
    { "aaab", "<a?>"  , 1, "e", "eeeeb" },
    { "aaab", "<a+>"  , 1, "e", "eb" },
    { "aaab", "<a*>"  , 1, "e", "eeb" },
    { "aaab", "<a{1}>", 1, "e", "eeeb" },

    { "aaab", "<x>"   , 1, "z", "aaab" },
    { "aaab", "<x>?"  , 1, "z", "zazazazb" },
    { "aaab", "<x>+"  , 1, "z", "aaab" },
    { "aaab", "<x>*"  , 1, "z", "zazazazb" },
    { "aaab", "<x>{1}", 1, "z", "aaab" },
    { "aaab", "<x?>"  , 1, "z", "zazazazb" },
    { "aaab", "<x+>"  , 1, "z", "aaab" },
    { "aaab", "<x*>"  , 1, "z", "zazazazb" },
    { "aaab", "<x{1}>", 1, "z", "aaab" },

    { "aaabaaa", "<a>"   , 1, "e", "eeebeee" },
    { "aaabaaa", "<a>?"  , 1, "e", "eeeebeee" },
    { "aaabaaa", "<a>+"  , 1, "e", "ebe" },
    { "aaabaaa", "<a>*"  , 1, "e", "eebe" },
    { "aaabaaa", "<a>{1}", 1, "e", "eeebeee" },
    { "aaabaaa", "<a?>"  , 1, "e", "eeeebeee" },
    { "aaabaaa", "<a+>"  , 1, "e", "ebe" },
    { "aaabaaa", "<a*>"  , 1, "e", "eebe" },
    { "aaabaaa", "<a{1}>", 1, "e", "eeebeee" },

    { "aaabaaa", "<x>"   , 1, "z", "aaabaaa" },
    { "aaabaaa", "<x>?"  , 1, "z", "zazazazbzazaza" },
    { "aaabaaa", "<x>+"  , 1, "z", "aaabaaa" },
    { "aaabaaa", "<x>*"  , 1, "z", "zazazazbzazaza" },
    { "aaabaaa", "<x>{1}", 1, "z", "aaabaaa" },
    { "aaabaaa", "<x?>"  , 1, "z", "zazazazbzazaza" },
    { "aaabaaa", "<x+>"  , 1, "z", "aaabaaa" },
    { "aaabaaa", "<x*>"  , 1, "z", "zazazazbzazaza" },
    { "aaabaaa", "<x{1}>", 1, "z", "aaabaaa" },

    { "Raptor Test", "<Raptor>", 1, "Captor", "Captor Test"   },
    { "Raptor Test", "<Raptor>", 0, "Captor", "Raptor Test"   },
    { "Raptor Test", "<Raptor|Test>", 0, "Captor", "Raptor Test"   },
    { "Raptor Test", "<Raptor|Test>", 1, "Captor", "Captor Captor"   },
    { "Raptor Test", "<Raptor|Test>", 2, "Captor", "Raptor Test"   },
    { "Raptor Test", "<Raptor|<Test>>", 2, "Fest", "Raptor Fest"   },
    { "Raptor Raptors Raptoring", "<Raptor:w*>", 1, "Test", "Test Test Test" },
    { "Raptor Raptors Raptoring", "<Raptor>:w*", 1, "Test", "Test Tests Testing" },
    { "Raptor Raptors Raptoring", "<<<R>a>ptor>:w*", 3, "C", "Captor Captors Captoring" },
    { "Raptor Raptors Raptoring", "<<<R>a>ptor>:w*", 2, "4", "4ptor 4ptors 4ptoring" },
  }

  var re RE
  for _, c := range swapTest {
    re.Match( c.txt, c.re )
    swap := re.RplCatch( c.swap, c.n )
    if swap != c.expected {
      t.Errorf( "Regexp4( %q, %q )\nRplCatch( %q, %d ) == %q\n             expected %q",
                c.txt, c.re, c.swap, c.n, swap, c.swap )
    }
  }
}

func pTest( t *testing.T ){
  putTest := []struct {
    txt, re string
    put, expected string
  }{
    { "a", "<a>", "#1", "a" },
    { "a", "<a>", "#x", "x" },
    { "a", "<a>", "#xx", "xx" },
    { "a", "<a>", "###1##", "#a#" },
    { "a", "<a>", "[#0][#1][#2#3#1000000]", "[][a][]" },
    { "aa", "<aa>", "#1", "aa" },
    { "a a a", "<a>", "#1#2#3", "aaa" },
    { "abcd", "<a|b|c|d>", "#4 #3 #2 #1", "d c b a" },
    { "1 2 3 4 5 6 7 8 9", "<1|2|3|4|5|6|7|8|9>", "#5 #6 #7 #8 #9 #1 #2 #3 #4", "5 6 7 8 9 1 2 3 4" },
    { "Raptor Test", "<aptor|est>", "C#1 F#2", "Captor Fest" },
    { "Raptor Test", "<aptor|est>", "C#5 F#2", "C Fest" },
    { "Raptor Test", "<aptor|est>", "C#a F#2", "Ca Fest" },
    { "Raptor Test", "<aptor|est>", "C#0 F#2", "C Fest" },
    { "Raptor Test", "<aptor|est>", "C#43 F#43", "C F" },
    { "Raptor Test", "<aptor|est>", "C##43 ##F#43##", "C#43 #F#" },
    { "Raptor Test", "<aptor|est>", "C##43 ##1##2", "C#43 #1#2" },
    { "Raptor Test", "<aptor|est>", "##Raptor ##Test", "#Raptor #Test" },
    { "Raptor Test Fest", "<Raptor> <Test>", "#1_#2", "Raptor_Test" },

    { "nasciiboy@gmail.com", "<[_A-Za-z0-9:-]+(:.[_A-Za-z0-9:-]+)*>:@<[A-Za-z0-9]+>:.<[A-Za-z0-9]+><:.[A-Za-z0-9]{2}>*", "[#1][#2][#3]", "[nasciiboy][gmail][com]" },
    { "<mail>nasciiboy@gmail.com</mail>", "<[_A-Za-z0-9:-]+(:.[_A-Za-z0-9:-]+)*>:@<[A-Za-z0-9]+>:.<[A-Za-z0-9]+><:.[A-Za-z0-9]{2}>*", "[#1][#2][#3]", "[nasciiboy][gmail][com]" },
    { "u.s.r_43@ru.com.jp", "<[_A-Za-z0-9:-]+(:.[_A-Za-z0-9:-]+)*>:@<[A-Za-z0-9]+>:.<[A-Za-z0-9]+><:.[A-Za-z0-9]{2}>*", "[#1][#2][#3]", "[u.s.r_43][ru][com]" },
    { "07-07-1777", "<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>", "d:#1 m:#3 y:#4", "d:07 m:07 y:1777" },
    { "fecha: 07-07-1777", "<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>", "d:#1 m:#3 y:#4", "d:07 m:07 y:1777" },

  }

  var re RE
  for _, c := range putTest {
    re.Find( c.txt, c.re )
    put := re.PutCatch( c.put )
    if put != c.expected {
      t.Errorf( "Regexp4( %q, %q )\nPutCatch( %q ) == %q, expected %q",
                c.txt, c.re, c.put, put, c.expected )
    }
  }
}

func gTest( t *testing.T ){
  putTest := []struct {
    txt, re string
    catch int
    pos int
  }{
    { "", "<0>", 1, 0 },
    { "0", "<0>", 1, 0 },
    { "0123456", "<6>", 1, 6 },
    { "0123456789", "<9>", 1, 9 },
    { "0123456789", "<:d>", 1, 0 },
    { "0123456789", "<:d>", 2, 1 },
    { "0123456789", "<:d>", 3, 2 },
    { "0123456789", "<:d>", 4, 3 },
    { "0123456789", "<:d>", 5, 4 },
    { "0123456789", "<:d>", 6, 5 },
    { "0123456789", "<:d>", 7, 6 },
    { "0123456789", "<9><.?>", 1, 9 },
    { "0123456789", "<9><.?>", 2, 10 },
    { "0123456789", "<9><.>?", 2, 10 },
    { "0123456789", "<9><.>", 2, 0 },
    { "0123456789", "<9><:.>*<:a>*", 2, 10 },
    { "0123456789", "<9><:.>*<:a>*", 3, 10 },
    { "0123456789.", "<9><:.>*<:a>*", 2, 10 },
    { "0123456789.", "<9><:.>*<:a>*", 3, 11 },
    { "0123456789.e", "<9><:.>*<:a>*", 3, 11 },
    { "0123456789...e", "<9><:.>*<:a>*", 3, 13 },
    { "0123456789^-^!e", "<9><:.>*<:a>*", 3, 10 },
  }

  var re RE
  for _, c := range putTest {
    re.Match( c.txt, c.re )
    pos := re.GpsCatch( c.catch )
    if pos != c.pos {
      t.Errorf( "Regexp4( %q, %q )\nGpsCatch( %d ) == %d, expected %d",
                c.txt, c.re, c.catch, pos, c.pos )
    }
  }
}

func nTestUTF( t *testing.T ){
  numTest := []struct {
    txt, re string
    n int
  }{
    { "a", "▲", 0 },
    { "a", "▲?", 1 },
    { "a", "△?", 1 },
    { "a", "▲+", 0 },
    { "a", "▲*", 1 },
    { "a", "△*", 1 },
    { "▲", "▲", 1 },
    { "▲", "▲?", 1 },
    { "▲", "△?", 1 },
    { "▲", "▲+", 1 },
    { "▲", "▲*", 1 },
    { "▲", "△*", 1 },
    { "▲▲▲", "▲+", 1 },
    { "▲▲▲", "▲*", 1 },
    { "▲▲▲", "▲+", 1 },
    { "▲▲▲", "▲?", 3 },
    { "▲▲▲△", "▲+", 1 },
    { "▲▲▲△", "▲*", 2 },
    { "▲▲▲△", "▲?", 4 },
    { "▲▲▲△", "▲+△", 1 },
    { "▲▲▲△", "▲*△", 1 },
    { "▲▲▲△", "▲?△", 1 },
    { "▲▲▲△", "▲+△?", 1 },
    { "▲▲▲△", "▲*△?", 1 },
    { "▲▲▲△", "▲?△?", 3 },
    { "▲▲▲△", "▲+△+", 1 },
    { "▲▲▲△", "▲*△+", 1 },
    { "▲▲▲△", "▲?△+", 1 },
    { "▲▲▲△", "▲+△*", 1 },
    { "▲▲▲△", "▲*△*", 1 },
    { "▲▲▲△", "▲?△*", 3 },
    { "▲▲▲△▲▲▲", "▲+", 2 },
    { "▲▲▲△▲▲▲", "▲*", 3 },
    { "▲▲▲△▲▲▲", "▲*", 3 },
    { "▲▲▲△▲▲▲", "▲*", 3 },
    { "▲", "(▲)?", 1 },
    { "▲", "(△)?", 1 },
    { "▲", "(▲)+", 1 },
    { "▲", "(▲)*", 1 },
    { "▲", "(△)*", 1 },
    { "▲▲▲", "(▲)+", 1 },
    { "▲▲▲", "(▲)*", 1 },
    { "▲▲▲", "#~(▲)+", 3 },
    { "▲▲▲", "#~(▲)*", 3 },
    { "▲▲▲△", "#~(▲+)", 3 },
    { "▲▲▲△", "#~(▲*)", 4 },
    { "▲▲▲△", "#~(▲?)", 4 },
    { "▲▲▲△", "#~(▲+△)", 3 },
    { "▲▲▲△", "#~(▲*△)", 4 },
    { "▲▲▲△", "#~(▲?△)", 2 },
    { "▲▲▲△", "#~(▲+△?)", 3 },
    { "▲▲▲△", "#~(▲*△?)", 4 },
    { "▲▲▲△", "#~(▲?△?)", 4 },
    { "▲▲▲△", "#~(▲+△+)", 3 },
    { "▲▲▲△", "#~(▲*△+)", 4 },
    { "▲▲▲△", "#~(▲?△+)", 2 },
    { "▲▲▲△", "#~(▲+△*)", 3 },
    { "▲▲▲△", "#~(▲*△*)", 4 },
    { "▲▲▲△", "#~(▲?△*)", 4 },
    { "▲▲▲△▲▲▲", "#~(▲+)", 6 },
    { "▲▲▲△▲▲▲", "#~(▲*)", 7 },
    { "▲", "[▲]?", 1 },
    { "▲", "[△]?", 1 },
    { "▲", "[▲]+", 1 },
    { "▲", "[▲]*", 1 },
    { "▲", "[△]*", 1 },
    { "▲▲▲", "[▲]?", 3 },
    { "▲▲▲", "[▲]+", 1 },
    { "▲▲▲", "[▲]*", 1 },
    { "▲▲▲", "#~[▲]?", 3 },
    { "▲▲▲", "#~[▲]+", 3 },
    { "▲▲▲", "#~[▲]*", 3 },
    { "▲▲▲△", "#~[▲△]", 4 },
    { "▲▲▲△", "#~[▲△]?", 4 },
    { "▲▲▲△", "#~[▲△]+", 4 },
    { "▲▲▲△", "#~[▲△]*", 4 },
    { "▲▲▲△▲▲▲", "#~[▲△]", 7 },
    { "▲", ":&", 1 },
    { "▲", ":&?", 1 },
    { "▲", ":&+", 1 },
    { "▲", ":&*", 1 },
    { "▲▲▲", ":&?", 3 },
    { "▲▲▲", ":&+", 1 },
    { "▲▲▲", ":&*", 1 },
    { "▲▲▲", "#~:&?", 3 },
    { "▲▲▲", "#~:&+", 3 },
    { "▲▲▲", "#~:&*", 3 },
    { "▲▲▲△", "#~:&", 4 },
    { "▲▲▲△", "#~:&?", 4 },
    { "▲▲▲△", "#~:&+", 4 },
    { "▲▲▲△", "#~:&*", 4 },
    { "▲▲▲△▲▲▲", "#~:&", 7 },
    { "▲", ":w", 0 },
    { "▲", ":w?", 1 },
    { "▲", ":w+", 0 },
    { "▲", ":w*", 1 },
    { "▲▲▲", ":w?", 3 },
    { "▲▲▲", ":w+", 0 },
    { "▲▲▲", ":w*", 3 },
    { "▲▲▲", "#~:w?", 3 },
    { "▲▲▲", "#~:w+", 0 },
    { "▲▲▲", "#~:w*", 3 },
    { "▲▲▲△", "#~:w", 0 },
    { "▲▲▲△", "#~:w?", 4 },
    { "▲▲▲△", "#~:w+", 0 },
    { "▲▲▲△", "#~:w*", 4 },
    { "▲▲▲△▲▲▲", "#~:w", 0 },
    { "▲", ":W", 1 },
    { "▲", ":W?", 1 },
    { "▲", ":W+", 1 },
    { "▲", ":W*", 1 },
    { "▲▲▲", ":W?", 3 },
    { "▲▲▲", ":W+", 1 },
    { "▲▲▲", ":W*", 1 },
    { "▲▲▲", "#~:W?", 3 },
    { "▲▲▲", "#~:W+", 3 },
    { "▲▲▲", "#~:W*", 3 },
    { "▲▲▲△", "#~:W", 4 },
    { "▲▲▲△", "#~:W?", 4 },
    { "▲▲▲△", "#~:W+", 4 },
    { "▲▲▲△", "#~:W*", 4 },
    { "▲▲▲△▲▲▲", "#~:W", 7 },

    { "△▲3△567△9", ".", 9 },
    { "△▲3△567△9", "(.)", 9 },
    { "△▲3△567△9", "[.]", 0 },
    { "△▲3△567△9", "(.+)", 1 },
    { "△▲3△567△9", ":&", 4 },
    { "△▲3△567△9", ":w", 5 },
    { "△▲3△567△9", ":W", 4 },
    { "△▲3△567△9", ":d", 5 },
    { "△▲3△567△9", ":a", 0 },
    { "△▲3△567△9", "[△5]", 4 },
    { "△▲3△567△9", "[▲1]", 1 },
    { "△▲3△567△9", "[3-9]", 5 },
    { "△▲3△567△9", "[▲1-7]", 5 },
    { "△▲3△567△9", "[^3-9]", 4 },
    { "△▲3△567△9", "[^a-z]", 9 },
    { "△▲3△567△9", "[^▲1-7]", 4 },
    { "△▲3△567△9", "[^:d]", 4 },
    { "△▲3△567△9", "[^:D]", 5 },
    { "△▲3△567△9", "[^:w]", 4 },
    { "△▲3△567△9", "[^:W]", 5 },
    { "△▲3△567△9", "[^:&]", 5 },
    { "△▲3△567△9", "[^:a]", 9 },
    { "△▲3△567△9", "[^:A]", 0 },

    { "Rááptor Test", "R.áptor", 1 },
    { "Rááptor Test", "Rá{2}ptor", 1 },
    { "Rááptor Test", "R(á){2}ptor", 1 },
    { "R△△△ptor Test", "R△{3}ptor", 1 },
    { "R△△△ptor Test", "R[^a]{3}ptor", 1 },
    { "R▲△ptor Test", "R[△▲]{2}ptor", 1 },
    { "R▲△ptor Test", "R[^ae]{2}ptor", 1 },
    { "R▲△ptor Test", "R.{2}ptor", 1 },
    { "R▲△ptor Test", "R[:W]{2}ptor", 1 },
    { "R▲△ptor Test", "R[^:w]{2}ptor", 1 },

    { "Σὲ γνωρίζω ἀπὸ τὴν κόψη", ".", 23 },
    { "Σὲ γνωρίζω ἀπὸ τὴν κόψη", ":&", 19 },
    { "Σὲ γνωρίζω ἀπὸ τὴν κόψη", "[:&]", 19 },
    { "Σὲ γνωρίζω ἀπὸ τὴν κόψη", "[^:&]", 4 },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", ".", 36 },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", ":&", 7 },
    { "ต้องรบราฆ่าฟันจนบรรลัย", ".", 22 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ /0123456789", ".", 38 },
    { "abcdefghijklmnopqrstuvwxyz £©µÀÆÖÞßéöÿ", ".", 38 },
    { "–—‘“”„†•…‰™œŠŸž€ ΑΒΓΔΩαβγδω АБВГДабвгд", ".", 38 },
    { "∀∂∈ℝ∧∪≡∞ ↑↗↨↻⇣ ┐┼╔╘░►☺♀ ﬁ�⑀₂ἠḂӥẄɐː⍎אԱა", ".", 38 },
    { "⡍⠜⠇⠑⠹ ⠺⠁⠎ ⠙⠑⠁⠙⠒ ⠞⠕ ⠃⠑⠛⠔ ⠺⠊⠹⠲ ⡹⠻⠑ ⠊⠎ ⠝⠕ ⠙⠳⠃⠞", ".", 43 },
  }

  done := make(chan struct{})
  for _, c := range numTest {
    go func( txt, re string, n int ){
      x := new( RE ).Match( txt, re )

      if x != n  {
        t.Errorf( "Regexp4( %q, %q ) == %d, expected %d", txt, re, x, n )
      }
      done <- struct{}{}
    }( c.txt, c.re, c.n )
  }

  for range numTest { <-done }
}

func cTestUTF( t *testing.T ){
  catchTest := []struct {
    txt, re string
    n int
    catch string
  }{
    { "▲", "<▲>", 1, "▲" },
    { "▲▲", "<▲▲>", 1, "▲▲" },
    { "▲ ▲ ▲", "<▲>", 3, "▲" },
    { "▲bcd", "<▲|b|c|d>", 1, "▲" },
    { "▲bcd", "<▲|b|c|d>", 2, "b" },
    { "▲bcd", "<▲|b|c|d>", 3, "c" },
    { "▲bcd", "<▲|b|c|d>", 4, "d" },
    { "▲bcd", "<▲|b|c|d>", 5, "" },
    { "▲bc", "▲<x>", 1, "" },
    { "▲bc", "<▲>x", 1, "" },
    { "▲bc", "<▲|b>x", 1, "" },
    { "▲bc", "<<▲|b>x|▲bc>", 1, "▲bc" },
    { "▲bc", "<<▲|b>x|▲bc>", 2, "" },
    { "▲bc", "<x<▲|b>|▲bc>", 1, "▲bc" },
    { "▲bc", "<x<▲|b>|▲bc>", 2, "" },
    { "▲bc ▲bc ▲bc", "<▲|b|c>", 9, "c" },
    { "▲bc ▲bc ▲bc", "<(▲|b|c)(▲|b|c)(▲|b|c)>", 3, "▲bc" },
    { "▲bc ▲bc ▲bc", "<(▲|b|c)(▲|b|c)(▲|b|c)> ", 2, "▲bc" },
    { "▲bc ▲bc ▲bc", "#?<(▲|b|c)(▲|b|c)(▲|b|c)>", 1, "▲bc" },
    { "▲bc ▲bc ▲bc", "#?<(▲|b|c)(▲|b|c)((▲|b)|x)>", 1, "" },
    { "▲bc ▲bc ▲bx", "#?<(▲|b|c)(▲|b|c)((▲|b)|x)>", 1, "▲bx" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 1, "▲bc" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 2, "iec" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 3, "i" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 4, "c" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 5, "oeb" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 6, "o" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 7, "b" },
    { "▲bc iecc oeb", "<<(▲|e)|(i|o)>e<b|c>|▲bc>", 8, "" },

    { "▲", "<▲>?", 1, "▲" },
    { "▲", "<b>?", 1, "" },
    { "▲", "<▲>+", 1, "▲" },
    { "▲", "<▲>*", 1, "▲" },
    { "▲", "<b>*", 1, "" },
    { "▲▲▲", "<▲>+", 1, "▲▲▲" },
    { "▲▲▲", "<▲>*", 1, "▲▲▲" },
    { "▲▲▲", "#~<▲>+", 1, "▲▲▲" },
    { "▲▲▲", "#~<▲>*", 1, "▲▲▲" },
    { "▲▲▲b", "#~<▲+>", 1, "▲▲▲" },
    { "▲▲▲b", "#~<▲*>", 1, "▲▲▲" },
    { "▲▲▲b", "#~<▲?>", 4, "" },
    { "▲▲▲b", "#~<▲+b>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲*b>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲?b>", 1, "▲b" },
    { "▲▲▲b", "#~<▲+b?>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲*b?>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲?b?>", 3, "▲b" },
    { "▲▲▲b", "#~<▲+b+>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲*b+>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲?b+>", 1, "▲b" },
    { "▲▲▲b", "#~<▲+b*>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲*b*>", 1, "▲▲▲b" },
    { "▲▲▲b", "#~<▲?b*>", 3, "▲b" },
    { "▲▲▲b▲▲▲", "#~<▲+>", 4, "▲▲▲" },
    { "▲▲▲b▲▲▲", "#~<▲*>", 5, "▲▲▲" },

    { "R▲ptor Test",     "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 1, "R▲ptor Test" },
    { "R▲ptor Test",     "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 2, "T" },
    { "R▲▲▲▲ptor TFest", "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 1, "R▲▲▲▲ptor TFest" },
    { "R▲▲▲▲ptor TFest", "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 2, "TF" },
    { "C▲ptorTest",      "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 1, "C▲ptorTest" },
    { "C▲p C▲ptorTest",  "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 1, "C▲p C▲ptorTest" },
    { "C▲p C▲ptorTest",  "#~<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 3, "C▲ptorTest" },
    { "R▲p C▲ptor Fest", "<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 1, "R▲p C▲ptor Fest" },
    { "R▲p C▲ptor Fest", "#~<((C|R)▲p C|C|R)(▲+p{1}tor) ?((<T|F>+e)(st))>", 3, "C▲ptor Fest" },

    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 1, "Λ" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 2, "̊" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 3, "̇" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 4, "̈" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 5, "⃑" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 6, "⊥" },
    { "STARGΛ̊TE SG-1, a = v̇ = r̈, a⃑ ⊥ b⃑", "<:&>", 7, "⃑" },
  }

  done := make(chan struct{})
  for _, c := range catchTest {
    go func( txt, re string, nCatch int, eCatch string ){
      r := new( RE )
      r.Match( txt, re )
      catch := r.GetCatch( nCatch )
      if catch != eCatch  {
        t.Errorf( "Regexp4( %q, %q )\nGetCatch( %d ) == %q, expected %q",
                  txt, re, nCatch, catch, eCatch )
      }
      done <- struct{}{}
    }( c.txt, c.re, c.n, c.catch )
  }

  for range catchTest { <-done }
}

func sTestUTF( t *testing.T ){
  swapTest := []struct {
    txt, re string
    n int
    swap, expected string
  }{
    { "▲", "<▲>"   , 1, "", "" },
    { "▲", "<▲>?"  , 1, "", "" },
    { "▲", "<▲>+"  , 1, "", "" },
    { "▲", "<▲>*"  , 1, "", "" },
    { "▲", "<▲>{1}", 1, "", "" },
    { "▲", "<▲?>"  , 1, "", "" },
    { "▲", "<▲+>"  , 1, "", "" },
    { "▲", "<▲*>"  , 1, "", "" },
    { "▲", "<▲{1}>", 1, "", "" },

    { "▲", "<▲>"   , 1, "e", "e" },
    { "▲", "<▲>?"  , 1, "e", "e" },
    { "▲", "<▲>+"  , 1, "e", "e" },
    { "▲", "<▲>*"  , 1, "e", "e" },
    { "▲", "<▲>{1}", 1, "e", "e" },
    { "▲", "<▲?>"  , 1, "e", "e" },
    { "▲", "<▲+>"  , 1, "e", "e" },
    { "▲", "<▲*>"  , 1, "e", "e" },
    { "▲", "<▲{1}>", 1, "e", "e" },

    { "▲", "<x>"   , 1, "z", "▲" },
    { "▲", "<x>?"  , 1, "z", "z▲" },
    { "▲", "<x>+"  , 1, "z", "▲" },
    { "▲", "<x>*"  , 1, "z", "z▲" },
    { "▲", "<x>{1}", 1, "z", "▲" },
    { "▲", "<x?>"  , 1, "z", "z▲" },
    { "▲", "<x+>"  , 1, "z", "▲" },
    { "▲", "<x*>"  , 1, "z", "z▲" },
    { "▲", "<x{1}>", 1, "z", "▲" },

    { "▲▲▲", "<▲>"   , 1, "", "" },
    { "▲▲▲", "<▲>?"  , 1, "", "" },
    { "▲▲▲", "<▲>+"  , 1, "", "" },
    { "▲▲▲", "<▲>*"  , 1, "", "" },
    { "▲▲▲", "<▲>{1}", 1, "", "" },
    { "▲▲▲", "<▲?>"  , 1, "", "" },
    { "▲▲▲", "<▲+>"  , 1, "", "" },
    { "▲▲▲", "<▲*>"  , 1, "", "" },
    { "▲▲▲", "<▲{1}>", 1, "", "" },

    { "▲▲▲", "<▲>"   , 1, "e", "eee" },
    { "▲▲▲", "<▲>?"  , 1, "e", "eee" },
    { "▲▲▲", "<▲>+"  , 1, "e", "e" },
    { "▲▲▲", "<▲>*"  , 1, "e", "e" },
    { "▲▲▲", "<▲>{1}", 1, "e", "eee" },
    { "▲▲▲", "<▲?>"  , 1, "e", "eee" },
    { "▲▲▲", "<▲+>"  , 1, "e", "e" },
    { "▲▲▲", "<▲*>"  , 1, "e", "e" },
    { "▲▲▲", "<▲{1}>", 1, "e", "eee" },

    { "▲▲▲", "<x>"   , 1, "z", "▲▲▲" },
    { "▲▲▲", "<x>?"  , 1, "z", "z▲z▲z▲" },
    { "▲▲▲", "<x>+"  , 1, "z", "▲▲▲" },
    { "▲▲▲", "<x>*"  , 1, "z", "z▲z▲z▲" },
    { "▲▲▲", "<x>{1}", 1, "z", "▲▲▲" },
    { "▲▲▲", "<x?>"  , 1, "z", "z▲z▲z▲" },
    { "▲▲▲", "<x+>"  , 1, "z", "▲▲▲" },
    { "▲▲▲", "<x*>"  , 1, "z", "z▲z▲z▲" },
    { "▲▲▲", "<x{1}>", 1, "z", "▲▲▲" },

    { "▲▲▲b", "<▲>"   , 1, "e", "eeeb" },
    { "▲▲▲b", "<▲>?"  , 1, "e", "eeeeb" },
    { "▲▲▲b", "<▲>+"  , 1, "e", "eb" },
    { "▲▲▲b", "<▲>*"  , 1, "e", "eeb" },
    { "▲▲▲b", "<▲>{1}", 1, "e", "eeeb" },
    { "▲▲▲b", "<▲?>"  , 1, "e", "eeeeb" },
    { "▲▲▲b", "<▲+>"  , 1, "e", "eb" },
    { "▲▲▲b", "<▲*>"  , 1, "e", "eeb" },
    { "▲▲▲b", "<▲{1}>", 1, "e", "eeeb" },

    { "▲▲▲b", "<x>"   , 1, "z", "▲▲▲b" },
    { "▲▲▲b", "<x>?"  , 1, "z", "z▲z▲z▲zb" },
    { "▲▲▲b", "<x>+"  , 1, "z", "▲▲▲b" },
    { "▲▲▲b", "<x>*"  , 1, "z", "z▲z▲z▲zb" },
    { "▲▲▲b", "<x>{1}", 1, "z", "▲▲▲b" },
    { "▲▲▲b", "<x?>"  , 1, "z", "z▲z▲z▲zb" },
    { "▲▲▲b", "<x+>"  , 1, "z", "▲▲▲b" },
    { "▲▲▲b", "<x*>"  , 1, "z", "z▲z▲z▲zb" },
    { "▲▲▲b", "<x{1}>", 1, "z", "▲▲▲b" },

    { "▲▲▲b▲▲▲", "<▲>"   , 1, "e", "eeebeee" },
    { "▲▲▲b▲▲▲", "<▲>?"  , 1, "e", "eeeebeee" },
    { "▲▲▲b▲▲▲", "<▲>+"  , 1, "e", "ebe" },
    { "▲▲▲b▲▲▲", "<▲>*"  , 1, "e", "eebe" },
    { "▲▲▲b▲▲▲", "<▲>{1}", 1, "e", "eeebeee" },
    { "▲▲▲b▲▲▲", "<▲?>"  , 1, "e", "eeeebeee" },
    { "▲▲▲b▲▲▲", "<▲+>"  , 1, "e", "ebe" },
    { "▲▲▲b▲▲▲", "<▲*>"  , 1, "e", "eebe" },
    { "▲▲▲b▲▲▲", "<▲{1}>", 1, "e", "eeebeee" },

    { "▲▲▲b▲▲▲", "<x>"   , 1, "z", "▲▲▲b▲▲▲" },
    { "▲▲▲b▲▲▲", "<x>?"  , 1, "z", "z▲z▲z▲zbz▲z▲z▲" },
    { "▲▲▲b▲▲▲", "<x>+"  , 1, "z", "▲▲▲b▲▲▲" },
    { "▲▲▲b▲▲▲", "<x>*"  , 1, "z", "z▲z▲z▲zbz▲z▲z▲" },
    { "▲▲▲b▲▲▲", "<x>{1}", 1, "z", "▲▲▲b▲▲▲" },
    { "▲▲▲b▲▲▲", "<x?>"  , 1, "z", "z▲z▲z▲zbz▲z▲z▲" },
    { "▲▲▲b▲▲▲", "<x+>"  , 1, "z", "▲▲▲b▲▲▲" },
    { "▲▲▲b▲▲▲", "<x*>"  , 1, "z", "z▲z▲z▲zbz▲z▲z▲" },
    { "▲▲▲b▲▲▲", "<x{1}>", 1, "z", "▲▲▲b▲▲▲" },

    { "R▲ptor Test", "<R▲ptor>", 1, "C▲ptor", "C▲ptor Test"   },
    { "R▲ptor Test", "<R▲ptor>", 0, "C▲ptor", "R▲ptor Test"   },
    { "R▲ptor Test", "<R▲ptor|Test>", 0, "C▲ptor", "R▲ptor Test"   },
    { "R▲ptor Test", "<R▲ptor|Test>", 1, "C▲ptor", "C▲ptor C▲ptor"   },
    { "R▲ptor Test", "<R▲ptor|Test>", 2, "C▲ptor", "R▲ptor Test"   },
    { "R▲ptor Test", "<R▲ptor|<Test>>", 2, "Fest", "R▲ptor Fest"   },
    { "R▲ptor R▲ptors R▲ptoring", "<R▲ptor:w*>", 1, "Test", "Test Test Test" },
    { "R▲ptor R▲ptors R▲ptoring", "<R▲ptor>:w*", 1, "Test", "Test Tests Testing" },
    { "R▲ptor R▲ptors R▲ptoring", "<<<R>▲>ptor>:w*", 3, "C", "C▲ptor C▲ptors C▲ptoring" },
    { "R▲ptor R▲ptors R▲ptoring", "<<<R>▲>ptor>:w*", 2, "4", "4ptor 4ptors 4ptoring" },
  }

  var re RE
  for _, c := range swapTest {
    re.Match( c.txt, c.re )
    swap := re.RplCatch( c.swap, c.n )
    if swap != c.expected {
      t.Errorf( "Regexp4( %q, %q )\nRplCatch( %q, %d ) == %q\n             expected %q",
                c.txt, c.re, c.swap, c.n, swap, c.swap )
    }
  }
}

func pTestUTF( t *testing.T ){
  putTest := []struct {
    txt, re string
    put, expected string
  }{
    { "▲", "<▲>", "#1", "▲" },
    { "▲", "<▲>", "#x", "x" },
    { "▲", "<▲>", "#xx", "xx" },
    { "▲", "<▲>", "###1##", "#▲#" },
    { "▲", "<▲>", "[#0][#1][#2#3#1000000]", "[][▲][]" },
    { "▲▲", "<▲▲>", "#1", "▲▲" },
    { "▲ ▲ ▲", "<▲>", "#1#2#3", "▲▲▲" },
    { "▲bcd", "<▲|b|c|d>", "#4 #3 #2 #1", "d c b ▲" },
    { "1 2 3 4 5 6 7 8 9", "<1|2|3|4|5|6|7|8|9>", "#5 #6 #7 #8 #9 #1 #2 #3 #4", "5 6 7 8 9 1 2 3 4" },
    { "R▲ptor Test", "<▲ptor|est>", "C#1 F#2", "C▲ptor Fest" },
    { "R▲ptor Test", "<▲ptor|est>", "C#5 F#2", "C Fest" },
    { "R▲ptor Test", "<▲ptor|est>", "C#▲ F#2", "C▲ Fest" },
    { "R▲ptor Test", "<▲ptor|est>", "C#0 F#2", "C Fest" },
    { "R▲ptor Test", "<▲ptor|est>", "C#43 F#43", "C F" },
    { "R▲ptor Test", "<▲ptor|est>", "C##43 ##F#43##", "C#43 #F#" },
    { "R▲ptor Test", "<▲ptor|est>", "C##43 ##1##2", "C#43 #1#2" },
    { "R▲ptor Test", "<▲ptor|est>", "##R▲ptor ##Test", "#R▲ptor #Test" },
    { "R▲ptor Test Fest", "<R▲ptor> <Test>", "#1_#2", "R▲ptor_Test" },
  }

  var re RE
  for _, c := range putTest {
    re.Match( c.txt, c.re )
    put := re.PutCatch( c.put )
    if put != c.expected {
      t.Errorf( "Regexp4( %q, %q )\nPutCatch( %q ) == %q, expected %q",
                c.txt, c.re, c.put, put, c.expected )
    }
  }
}

func gTestUTF( t *testing.T ){
  putTest := []struct {
    txt, re string
    catch int
    pos, len int
  }{
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ /0123456789", "<8>", 1, 36, 1 },
    { "abcdefghijklmnopqrstuvwxyz £©µÀÆÖÞßéöÿ", "<ö>", 1, 45, 2 },
    { "–—‘“”„†•…‰™œŠŸž€ ΑΒΓΔΩαβγδω АБВГДабвгд", "<г>", 1, 82, 2 },
    { "∀∂∈ℝ∧∪≡∞ ↑↗↨↻⇣ ┐┼╔╘░►☺♀ ﬁ�⑀₂ἠḂӥẄɐː⍎אԱა", "<Ա>", 1, 98, 2 },
    { "⡍⠜⠇⠑⠹ ⠺⠁⠎ ⠙⠑⠁⠙⠒ ⠞⠕ ⠃⠑⠛⠔ ⠺⠊⠹⠲ ⡹⠻⠑ ⠊⠎ ⠝⠕ ⠙⠳⠃⠞", "<⡹>", 1, 75, 3 },
    { "ABCDEFGHIJKLMNOPQRSTUVWXYZ /0123456789", "<89>", 1, 36, 2 },
    { "abcdefghijklmnopqrstuvwxyz £©µÀÆÖÞßéöÿ", "<öÿ>", 1, 45, 4 },
    { "–—‘“”„†•…‰™œŠŸž€ ΑΒΓΔΩαβγδω АБВГДабвгд", "<гд>", 1, 82, 4 },
    { "∀∂∈ℝ∧∪≡∞ ↑↗↨↻⇣ ┐┼╔╘░►☺♀ ﬁ�⑀₂ἠḂӥẄɐː⍎אԱა", "<Աა>", 1, 98, 5 },
    { "⡍⠜⠇⠑⠹ ⠺⠁⠎ ⠙⠑⠁⠙⠒ ⠞⠕ ⠃⠑⠛⠔ ⠺⠊⠹⠲ ⡹⠻⠑ ⠊⠎ ⠝⠕ ⠙⠳⠃⠞", "<⡹⠻>", 1, 75, 6 },
  }

  var re RE
  for _, c := range putTest {
    re.Find( c.txt, c.re )
    pos := re.GpsCatch( c.catch )
    len := re.LenCatch( c.catch )
    if pos != c.pos || len != c.len {
      t.Errorf( "Regexp4( %q, %q )\nGpsCatch( %d ) == %d, expected %d\nLenCatch( %d ) == %d, expected %d",
                c.txt, c.re, c.catch, pos, c.pos, c.catch, len, c.len )
    }
  }
}

////////////// INTERNAL-COMPARATIVE-BENCHMARKS
/// Find vs [Compile() + Copy().FindStirng()]

const rebe = "<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>"
const reco = "07-07-1777"

func BenchmarkFind(b *testing.B) {
  var re RE
  for i := 0; i < b.N; i++ {
    if !re.Find( reco, rebe ) {
      b.Errorf( "BenchmarkFind: re.Find(): no-match" )
    }
  }
}

var reFi = Compile( rebe )

func BenchmarkFindCopy(b *testing.B) {
  for i := 0; i < b.N; i++ {
    if !reFi.Copy().FindString( reco ) {
      b.Errorf( "BenchmarkFindCopy: re.Find(): no-match" )
    }
  }
}

const rebe2 = "#^<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>"
const rebe3 = "#*$<0?[1-9]|[12][0-9]|3[01]><[/:-\\]><0?[1-9]|1[012]>@2<[12][0-9]{3}>"

func BenchmarkFind3X(b *testing.B) {
  var re RE
  for i := 0; i < b.N; i++ {
    if !re.Find( reco, rebe ) {
      b.Errorf( "BenchmarkFind: re.Find(): no-match" )
    }
    if !re.Find( reco, rebe2 ) {
      b.Errorf( "BenchmarkFind: re.Find(): no-match" )
    }
    if !re.Find( reco, rebe3 ) {
      b.Errorf( "BenchmarkFind: re.Find(): no-match" )
    }
  }
}

var reFi2 = Compile( rebe2 )
var reFi3 = Compile( rebe3 )

func BenchmarkFindCopy3X(b *testing.B) {
  for i := 0; i < b.N; i++ {
    if !reFi.Copy().FindString( reco ) {
      b.Errorf( "BenchmarkFindCopy: re.Find(): no-match" )
    }
    if !reFi2.Copy().FindString( reco ) {
      b.Errorf( "BenchmarkFindCopy: re.Find(): no-match" )
    }
    if !reFi3.Copy().FindString( reco ) {
      b.Errorf( "BenchmarkFindCopy: re.Find(): no-match" )
    }
  }
}

const srebe = "#^text"
const sreco = "text"

func BenchmarkFindSimple(b *testing.B) {
  var re RE
  for i := 0; i < b.N; i++ {
    if !re.Find( sreco, srebe ) {
      b.Errorf( "BenchmarkFind: re.Find(): no-match" )
    }
  }
}

var reSi = Compile( srebe )

func BenchmarkFindCopySimple(b *testing.B) {
  for i := 0; i < b.N; i++ {
    if !reSi.Copy().FindString( sreco ) {
      b.Errorf( "BenchmarkFindCopy: re.Find(): no-match" )
    }
  }
}

/// RplCatch (string vs []byte vs bytes.Buffer)

func (r *RE) OldRplCatch( rplStr string, id int ) (result string) {
  last := 0

  for index := 1; index < r.catchIndex; index++ {
    if r.catches[index].id == id {
      if last > r.catches[index].init { last = r.catches[index].init }

      result += r.txt[last:r.catches[index].init]
      result += rplStr
      last    = r.catches[index].end
    }
  }

  if last < len(r.txt) { result += r.txt[last:] }

  return
}

func (r *RE) BufferRplCatch( rplStr string, id int ) string {
  last := 0
  var b bytes.Buffer

  for index := 1; index < r.catchIndex; index++ {
    if r.catches[index].id == id {
      if last > r.catches[index].init { last = r.catches[index].init }

      b.WriteString( r.txt[last:r.catches[index].init] )
      b.WriteString( rplStr )
      last    = r.catches[index].end
    }
  }

  if last < len(r.txt) { b.WriteString( r.txt[last:] ) }

  return b.String()
}

var   rerpl = Compile( "<:s>+" )
const ssIn  = "  \nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\n"
const ssSwp = "––"
const ssOut = "––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––line-a––line-b––line-c––line-d––line-en––"


func BenchmarkRplCatchOld( b *testing.B ){
  r := rerpl.Copy()
  r.FindString( ssIn )

  for i := 0; i < b.N; i++ {
    if( r.OldRplCatch( ssSwp, 1 ) != ssOut ){
      b.Fatalf( "BenchmarkRplCatchOld(): no match" )
    }
  }
}

func BenchmarkRplCatchBuffer( b *testing.B ){
  r := rerpl.Copy()
  r.FindString( ssIn )

  for i := 0; i < b.N; i++ {
    if( r.BufferRplCatch( ssSwp, 1 ) != ssOut ){
      b.Fatalf( "BenchmarkRplCatchBuffer(): no match" )
    }
  }
}

func BenchmarkRplCatch( b *testing.B ){
  r := rerpl.Copy()
  r.FindString( ssIn )

  for i := 0; i < b.N; i++ {
    if( r.RplCatch( ssSwp, 1 ) != ssOut ){
      b.Fatalf( "BenchmarkRplCatch(): no match" )
    }
  }
}
