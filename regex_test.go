package main

import (
	"testing"
)

var config = `parent {
    type {
        subtype TestSubType {
            element TestElement {
                attributes {
                    name testname;
                    description testdescription;
                    count 1234;
                    value testvalue;
                }
                owner person1;
            }
        }
        subtype TestSubType1 {
            element TestElement1 {
                attributes {
                    name testname1;
                    description testdescription1;
                    count 1;
                    value 1
                }
                owner person2;
            }
        }
        subtype TestSubType2 {
            element TestElement2 {
                attributes {
                    name testname2;
                    description testdescription2;
                    count 2;
                    value 2
                }
                owner unknown;
            }
        }
        subtype TestSubType3 {
            element TestElement3 {
                attributes {
                    name testname3;
                    description testdescription3;
                    count 3;
                    value 3
                }
                owner unknown;
            }
        }
        subtype TestSubType3 {
            element TestElement4 {
                attributes {
                    name testname4;
                    description testdescription4;
                    count 4;
                    value 4
                }
                owner unknown;
            }
        }
        subtype TestSubType4 {
            element TestElement5 {
                attributes {
                    name testname5;
                    description testdescription5;
                    count 5;
                    value 5
                }
                owner unknown;
            }
        }
        subtype TestSubType5 {
            element TestElement6 {
                attributes {
                    name testname6;
                    description testdescription6;
                    count 6;
                    value 6
                }
                owner unknown;
            }
        }
        subtype TestSubType6 {
            element TestElement7 {
                attributes {
                    name testname7;
                    description testdescription7;
                    count 7;
                    value 7
                }
                owner unknown;
            }
        }
        subtype TestSubType7 {
            element TestElement8 {
                attributes {
                    name testname8;
                    description testdescription8;
                    count 8;
                    value 8
                }
                owner unknown;
            }
        }
        subtype TestSubType8 {
            element TestElement9 {
                attributes {
                    name testname9;
                    description testdescription9;
                    count 9;
                    value 9
                }
                owner unknown;
            }
        }
        subtype TestSubType9 {
            element TestElement10 {
                attributes {
                    name testname10;
                    description testdescription10;
                    count 10;
                    value 10
                }
                owner unknown;
            }
        }
        subtype TestSubType10 {
            element TestElement11 {
                attributes {
                    name testname11;
                    description testdescription11;
                    count 11;
                    value 11
                }
                owner unknown;
            }
        }
    }
}
`

func TestRegexPerformance(t *testing.T) {
	t.Skip()
	l := New("pattern")
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)
	if got := Match(nfa, config); got != true {
		t.Errorf("Match(%q) = %v, want %v", config, got, false)
	}
}

func TestRegexMatch(t *testing.T) {
	cases := map[string][]struct {
		input string
		match bool
	}{
		"aa.*": {
			{"aac", true},
			{"aab", true},
			{"aabbcc", true},
			{"aa", true},
			{"aaabbb", true},
			{"a", false},
			{"ab", false},
			{"abbbbb", false},
			{"xabc", false},
			{"abccd", false},
		},
		"ab*cd*": {
			{"ac", true},
			{"abc", true},
			{"abbbbbc", true},
			{"acd", true},
			{"abbbbbcdddd", true},
			{"a", false},
			{"ab", false},
			{"abbbbb", false},
			{"xabc", false},
			{"abccd", false},
		},
		"ab*": {
			{"a", true},
			{"ab", true},
			{"abb", true},
			{"abbb", true},
			{"b", false},
			{"ba", false},
			{"", false},
			{"abbbbcd", false},
			{"xab", false},
		},
		"ab": {
			{"ab", true},
			{"a", false},
			{"b", false},
			{"", false},
			{"abc", false},
			{"xab", false},
			{"abx", false},
		},
		"a.*b*": {
			{"a", true},
			{"ab", true},
			{"a123", true},
			{"axxxbbbb", true},
			{"aBBBB", true},
			{"a.b", true},
			{"abbbbbb", true},
			{"b", false},
			{"x", false},
			{"", false},
		},
		"a[bc]d": {
			{"abd", true},
			{"acd", true},
			{"aad", false},
			{"abcd", false},
			{"abbd", false},
			{"ad", false},
			{"abd ", false},
			{" abd", false},
			{"Abd", false},
			{"aCd", false},
			{"", false},
		},
		"pa\\sb": {
			{"pa b", true},       
			{"pa\tb", true},      
			{"pa\nb", true},      
			{"pa\r\nb", false},   
			{"pa\v b", false},    
			{"pa\vb", true},      
			{"pa\fb", true},      
			{"pab", false},       
			{"pa  b", false},     
			{" pa b", false},     
			{"pa b ", false},     
			{"p a b", false},     
			{"", false},          
		},
		"a[bc]*d": {
			{"ad", true},           
			{"abd", true},          
			{"acd", true},          
			{"abcd", true},         
			{"abcbcd", true},       
			{"abccbd", true},       
			{"a", false},           
			{"d", false},           
			{"abxd", false},        
			{"axcd", false},        
			{"ab cd", false},       
			{"abcbcbcbcd", true},   
			{"abcbdx", false},      
			{"abccd", true},        
			{"", false},            
		},
	}

	for key, val := range cases {
		l := New(key)
		parser := NewParser(l)
		ast := parser.Ast()
		nfa := Compile(ast)
		for _, c := range val {
			if got := Match(nfa, c.input); got != c.match {
				t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
			}
		}
	}
}
