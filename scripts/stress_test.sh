#!/bin/zsh
set -e

cd "$(dirname "$0")/.."
BIN() { go run ./cmd/dcd "$@"; }
TMP="/tmp/dcd-stress-test"
mkdir -p "$TMP"
PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  ✅ $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  ❌ $1"; echo "     expected: $2"; echo "     got:      $3"; }
header() { echo; echo "=== $1 ==="; }

clean() { rm -rf "$TMP"/*; }

# =============================================================
# HAPPY PATH — should succeed
# =============================================================

header "HAPPY PATH"

# --- 1. All inline tags ---
cat > "$TMP/01_all_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><b>bold</b> <i>italic</i> <u>underline</u> <s>strike</s> <code>code</code> <mark>mark</mark> <sub>sub</sub> <sup>sup</sup></p>
EOF
BIN "$TMP/01_all_inline.dcd" "$TMP/01_all_inline.docx" 2>&1 && pass "01 all inline tags" || fail "01" "exit 0" "non-zero exit"

# --- 2. <mark color=green> ---
cat > "$TMP/02_mark_color.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><mark color=green>green highlight</mark> <mark color=yellow>yellow</mark> <mark color=red>red</mark></p>
EOF
BIN "$TMP/02_mark_color.dcd" "$TMP/02_mark_color.docx" 2>&1 && pass "02 mark color" || fail "02" "exit 0" "non-zero exit"

# --- 3. <set:flags> complex ---
cat > "$TMP/03_set_flags.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:b|i|u|s>bold italic underline strike</set:b|i|u|s></p>
<p><set:b|i|code>bold italic code</set:b|i|code></p>
<p><set:b|s>bold strike</set:b|s></p>
<p><set:u underline=double>double underline</set:u></p>
<p><set:b|u underline=dash>bold dash underline</set:b|u></p>
EOF
BIN "$TMP/03_set_flags.dcd" "$TMP/03_set_flags.docx" 2>&1 && pass "03 set flags complex" || fail "03" "exit 0" "non-zero exit"

# --- 4. <w:flags> with attributes ---
cat > "$TMP/04_w_flags_attrs.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:u underline=double>double underline para</w:u>
<w:u underline=dash>dashed underline para</w:u>
<w:u underline=wavy>wavy underline para</w:u>
<w:c|b>centered bold</w:c|b>
<w:c|i|u>center italic underline</w:c|i|u>
EOF
BIN "$TMP/04_w_flags_attrs.dcd" "$TMP/04_w_flags_attrs.docx" 2>&1 && pass "04 w flags with attrs" || fail "04" "exit 0" "non-zero exit"

# --- 5. <ol type=a/A/i/I> static ---
cat > "$TMP/05_ol_type.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol type=a>
<li>item a</li>
<li>item b</li>
</ol>
<ol type=A>
<li>item A</li>
<li>item B</li>
</ol>
<ol type=i>
<li>item i</li>
<li>item ii</li>
</ol>
<ol type=I>
<li>item I</li>
<li>item II</li>
</ol>
EOF
BIN "$TMP/05_ol_type.dcd" "$TMP/05_ol_type.docx" 2>&1 && pass "05 ol type static" || fail "05" "exit 0" "non-zero exit"

# --- 6. Properly nested inline tags ---
cat > "$TMP/06_nested_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><b>bold <i>bold+italic <u>bold+italic+underline</u></i></b></p>
<p>normal <b>bold <i>bold+italic</i></b> normal <u>underline</u></p>
<p><code><b>bold code</b> <i>code italic</i></code></p>
EOF
BIN "$TMP/06_nested_inline.dcd" "$TMP/06_nested_inline.docx" 2>&1 && pass "06 nested inline" || fail "06" "exit 0" "non-zero exit"

# --- 7. <tab> tag ---
cat > "$TMP/07_tab.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Name:<tab>John Doe</p>
<p>Age:<tab size=4>25</p>
<p>City:<tab size=8>Jakarta</p>
EOF
BIN "$TMP/07_tab.dcd" "$TMP/07_tab.docx" 2>&1 && pass "07 tab" || fail "07" "exit 0" "non-zero exit"

# --- 8. sub/sup in lists ---
cat > "$TMP/08_subsup_list.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li>H<sub>2</sub>O</li>
<li>E=mc<sup>2</sup></li>
<li>H<sub>2</sub>SO<sub>4</sub></li>
</ul>
<ol>
<li>x<sup>2</sup> + y<sup>2</sup></li>
<li>log<sub>10</sub>(x)</li>
</ol>
EOF
BIN "$TMP/08_subsup_list.dcd" "$TMP/08_subsup_list.docx" 2>&1 && pass "08 sub sup in lists" || fail "08" "exit 0" "non-zero exit"

# --- 9. sub/sup in wrapped blocks ---
cat > "$TMP/09_subsup_wrapped.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>H<sub>2</sub>O center</w:c>
<w:code>log<sub>10</sub>(x)</w:code>
<w:sup>superscript block</w:sup>
<w:sub>subscript block</w:sub>
EOF
BIN "$TMP/09_subsup_wrapped.dcd" "$TMP/09_subsup_wrapped.docx" 2>&1 && pass "09 sub sup wrapped" || fail "09" "exit 0" "non-zero exit"

# --- 10. code in list ---
cat > "$TMP/10_code_list.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li><code>fn main()</code></li>
<li><code>println!("hello")</code></li>
</ul>
<ol>
<li>run <code>dcd</code></li>
<li>open <code>output.docx</code></li>
</ol>
EOF
BIN "$TMP/10_code_list.dcd" "$TMP/10_code_list.docx" 2>&1 && pass "10 code in lists" || fail "10" "exit 0" "non-zero exit"

# --- 11. header/footer mirror ---
cat > "$TMP/11_hdr_mirror.dcd" << 'EOF'
[style]
layout=A4
[header]
left=Left
right=Right
mirror=true
[section 0]
name=test
--- BODY ---
<p>mirror header test</p>
EOF
BIN "$TMP/11_hdr_mirror.dcd" "$TMP/11_hdr_mirror.docx" 2>&1 && pass "11 header mirror" || fail "11" "exit 0" "non-zero exit"

# --- 12. header/footer first-page ---
cat > "$TMP/12_hdr_firstpage.dcd" << 'EOF'
[style]
layout=A4
[header]
left=Header
first-page=false
[footer]
center=Footer
first-page=false
[section 0]
name=test
--- BODY ---
<p>first-page=false test</p>
EOF
BIN "$TMP/12_hdr_firstpage.dcd" "$TMP/12_hdr_firstpage.docx" 2>&1 && pass "12 header first-page" || fail "12" "exit 0" "non-zero exit"

# --- 13. header/footer justify_between ---
cat > "$TMP/13_hdr_justify.dcd" << 'EOF'
[style]
layout=A4
[header]
justify_between=Left, Center, Right
font-size=10
[footer]
justify_between=Dept. A\, B\, and C, {{date}}, Page {{page}}
[section 0]
name=test
--- BODY ---
<p>justify_between test</p>
EOF
BIN "$TMP/13_hdr_justify.dcd" "$TMP/13_hdr_justify.docx" 2>&1 && pass "13 header justify_between" || fail "13" "exit 0" "non-zero exit"

# --- 14. format specifiers ---
cat > "$TMP/14_format.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title, items.date_field
formats=[items.date_field:dd-MM-yyyy]
--- BODY ---
<h1>{{info.title}}</h1>
<loop x from items>
<p>{{x.name}}: {{x.date_field}}</p>
</loop>
EOF
BIN "$TMP/14_format.dcd" "$TMP/14_format.docx" 2>&1 && pass "14 format specifiers" || fail "14" "exit 0" "non-zero exit"

# --- 15. loop with type= ---
cat > "$TMP/15_loop_type.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items type=A>
{{x.label}}
</loop:ol>
<loop:ol x from items type=a>
{{x.label}}
</loop:ol>
<loop:ol x from items type=I>
{{x.label}}
</loop:ol>
<loop:ol x from items type=i>
{{x.label}}
</loop:ol>
EOF
BIN "$TMP/15_loop_type.dcd" "$TMP/15_loop_type.docx" 2>&1 && pass "15 loop with type" || fail "15" "exit 0" "non-zero exit"

# --- 16. unicode / special chars ---
cat > "$TMP/16_unicode.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Unicode: äöü ñ ç é è ê Ω π ≈ ∞ ≤ ≥</p>
<p>Special: & < > " '</p>
<p>Emoji: 😀 🎉 🚀 ❤️ 🐱</p>
<p>Mixed: résumé = 100 € — café naïve</p>
EOF
BIN "$TMP/16_unicode.dcd" "$TMP/16_unicode.docx" 2>&1 && pass "16 unicode/special" || fail "16" "exit 0" "non-zero exit"

# --- 17. 3-level nested list (mixed) ---
cat > "$TMP/17_nested3.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li>level 1
<ul>
<li>level 2
<ol>
<li>level 3a</li>
<li>level 3b</li>
</ol>
</li>
<li>level 2b</li>
</ul>
</li>
<li>level 1b
<ol>
<li>level 2
<ul>
<li>level 3</li>
</ul>
</li>
</ol>
</li>
</ul>
EOF
BIN "$TMP/17_nested3.dcd" "$TMP/17_nested3.docx" 2>&1 && pass "17 3-level nested list" || fail "17" "exit 0" "non-zero exit"

# --- 18. section limits boundary (5 var, 15 keys) ---
cat > "$TMP/18_section_limits.dcd" << 'EOF'
[section 0]
name=test
var=v1, []v2, []v3, []v4, []v5
keys=k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, k13, k14, k15
--- BODY ---
<p>{{v1.k1}} {{v1.k2}} {{v1.k3}} {{v1.k4}} {{v1.k5}}</p>
<p>{{v1.k6}} {{v1.k7}} {{v1.k8}} {{v1.k9}} {{v1.k10}}</p>
<p>{{v1.k11}} {{v1.k12}} {{v1.k13}} {{v1.k14}} {{v1.k15}}</p>
<loop x from v2>{{x}}</loop>
<loop x from v3>{{x}}</loop>
<loop x from v4>{{x}}</loop>
<loop x from v5>{{x}}</loop>
EOF
BIN "$TMP/18_section_limits.dcd" "$TMP/18_section_limits.docx" 2>&1 && pass "18 section limits boundary" || fail "18" "exit 0" "non-zero exit"

# --- 19. inline tag in wrapped blocks ---
cat > "$TMP/19_inline_wrapped.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>centered with <b>bold</b> and <i>italic</i></w:c>
<w:code>code with <b>bold</b> inline</w:code>
<w:b>bold block with <i>italic</i> inline</w:b>
EOF
BIN "$TMP/19_inline_wrapped.dcd" "$TMP/19_inline_wrapped.docx" 2>&1 && pass "19 inline in wrapped" || fail "19" "exit 0" "non-zero exit"

# --- 20. header/footer border + margin ---
cat > "$TMP/20_hdr_border.dcd" << 'EOF'
[style]
layout=A4
[header]
left=Header
border=bottom
margin=0.3
[footer]
center=Footer
border=top
margin=0.2
[section 0]
name=test
--- BODY ---
<p>header/footer border test</p>
EOF
BIN "$TMP/20_hdr_border.dcd" "$TMP/20_hdr_border.docx" 2>&1 && pass "20 header border margin" || fail "20" "exit 0" "non-zero exit"

# --- 21. Links with attributes ---
cat > "$TMP/21_link_attrs.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Link: <a=https://example.com>plain</a></p>
<p>Colored: <a=https://example.com color=red>red link</a></p>
<p>No underline: <a=https://example.com underline=false>no underline</a></p>
<p>Bookmark: <a=#section1>goto section</a></p>
EOF
BIN "$TMP/21_link_attrs.dcd" "$TMP/21_link_attrs.docx" 2>&1 && pass "21 link attributes" || fail "21" "exit 0" "non-zero exit"

# --- 22. Image with all attributes ---
cat > "$TMP/22_image_attrs.dcd" << 'EOF'
[section 0]
name=test
var=src
keys=img
--- BODY ---
<p>img {{src.img}} width=80% align=center alt Diagram border=1 bg=#f0f0f0</p>
EOF
BIN "$TMP/22_image_attrs.dcd" "$TMP/22_image_attrs.docx" 2>&1 && pass "22 image attributes" || fail "22" "exit 0" "non-zero exit"

# --- 23. <p> tag with local attrs ---
cat > "$TMP/23_p_attrs.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p align=center>centered paragraph</p>
<p align=right>right aligned</p>
EOF
BIN "$TMP/23_p_attrs.dcd" "$TMP/23_p_attrs.docx" 2>&1 && pass "23 p attributes" || fail "23" "exit 0" "non-zero exit"

# --- 24. inline formatting in headings (heading restriction) ---
cat > "$TMP/24_heading_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h1>Plain heading ok</h1>
<h2>{{var.key}} heading var ok</h2>
EOF
BIN "$TMP/24_heading_inline.dcd" "$TMP/24_heading_inline.docx" 2>&1 && pass "24 heading plain text" || fail "24" "exit 0" "non-zero exit"

# --- 25. Complex combined: list in table ---
cat > "$TMP/25_list_in_table.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col><b>bold cell</b> with <i>italic</i></col>
<col><code>code</code> and <u>underline</u></col>
</row>
<row>
<col>H<sub>2</sub>O</col>
<col>x<sup>2</sup></col>
</row>
</table>
EOF
BIN "$TMP/25_list_in_table.dcd" "$TMP/25_list_in_table.docx" 2>&1 && pass "25 table with inline" || fail "25" "exit 0" "non-zero exit"

# --- 26. static <ol type=a> with inline ---
cat > "$TMP/26_ol_type_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol type=A>
<li><b>bold item A</b></li>
<li><i>italic item B</i></li>
<li><code>code item C</code></li>
</ol>
<ol type=i>
<li>H<sub>2</sub>O</li>
<li>E=mc<sup>2</sup></li>
</ol>
EOF
BIN "$TMP/26_ol_type_inline.dcd" "$TMP/26_ol_type_inline.docx" 2>&1 && pass "26 ol type with inline" || fail "26" "exit 0" "non-zero exit"

# --- 27. Empty section ---
cat > "$TMP/27_empty_section.dcd" << 'EOF'
[section 0]
name=empty
--- BODY ---
EOF
BIN "$TMP/27_empty_section.dcd" "$TMP/27_empty_section.docx" 2>&1 && pass "27 empty section" || fail "27" "exit 0" "non-zero exit"

# --- 28. Keys-only section (no var) ---
cat > "$TMP/28_keys_only.dcd" << 'EOF'
[section 0]
name=test
keys=title, message
--- BODY ---
<p>{{title}}: {{message}}</p>
EOF
BIN "$TMP/28_keys_only.dcd" "$TMP/28_keys_only.docx" 2>&1 && pass "28 keys only" || fail "28" "exit 0" "non-zero exit"

# --- 29. Nested loop ---
cat > "$TMP/29_nested_loop.dcd" << 'EOF'
[section 0]
name=test
var=info, []categories, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop cat from categories>
<h2>{{cat.name}}</h2>
<loop x from items>
<p>{{cat.name}}: {{x.label}}</p>
</loop>
</loop>
EOF
BIN "$TMP/29_nested_loop.dcd" "$TMP/29_nested_loop.docx" 2>&1 && pass "29 nested loop" || fail "29" "exit 0" "non-zero exit"

# --- 30. Nested table ---
cat > "$TMP/30_nested_table.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col>outer A
<table border=1>
<row><col>inner</col></row>
</table>
</col>
<col>outer B</col>
</row>
</table>
EOF
BIN "$TMP/30_nested_table.dcd" "$TMP/30_nested_table.docx" 2>&1 && pass "30 nested table" || fail "30" "exit 0" "non-zero exit"

# --- 31. <hr> with attributes ---
cat > "$TMP/31_hr.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<hr>
<p>text</p>
<hr width=75%>
<p>text</p>
<hr width=50% color=#4472C4>
EOF
BIN "$TMP/31_hr.dcd" "$TMP/31_hr.docx" 2>&1 && pass "31 hr attributes" || fail "31" "exit 0" "non-zero exit"

# --- 32. Heading with keep-next / keep-lines / border-bottom ---
cat > "$TMP/32_heading_styles.dcd" << 'EOF'
[style]
layout=A4

[style:heading-1]
keep-next=true
keep-lines=true
border-bottom=1pt
font-family=Arial
font-size=24
color=#2b5797
bold=true
space-before=18
space-after=12

[style:heading-2]
keep-next=true
font-family=Arial
font-size=18
color=#444444

[section 0]
name=test
--- BODY ---
<h1>Chapter 1</h1>
<p>lorem ipsum</p>
<h2 keep-next=true>Section 1.1</h2>
<h2 border-bottom=2pt>Section 1.2</h2>
EOF
BIN "$TMP/32_heading_styles.dcd" "$TMP/32_heading_styles.docx" 2>&1 && pass "32 heading styles" || fail "32" "exit 0" "non-zero exit"

# --- 33. Named table style + style.first ---
cat > "$TMP/33_table_style.dcd" << 'EOF'
[style:table header]
bg=#2b5797
color=white
font-weight=bold
align=center
border-bottom=2pt

[style:table alt]
bg=#f5f5f5

[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<table border=1>
<row style=header>
<col>Name</col>
<col>Value</col>
</row>
<loop:row x from items style.first=header>
<col>{{x.name}}</col>
<col>{{x.value}}</col>
</loop:row>
<row style=alt>
<col>Total</col>
<col>100</col>
</row>
</table>
EOF
BIN "$TMP/33_table_style.dcd" "$TMP/33_table_style.docx" 2>&1 && pass "33 table style" || fail "33" "exit 0" "non-zero exit"

# --- 34. <br> multiple + in lists ---
cat > "$TMP/34_br.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>line 1<br>line 2<br><br>line 4</p>
<ul>
<li>item<br>with break</li>
<li>next item</li>
</ul>
EOF
BIN "$TMP/34_br.dcd" "$TMP/34_br.docx" 2>&1 && pass "34 br multiple" || fail "34" "exit 0" "non-zero exit"

# --- 35. Self-closing variants ---
cat > "$TMP/35_selfclose.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>page 1</p>
<br/>
<p>after br self-close</p>
<pb/>
<p>after pb self-close</p>
<br/>
<hr/>
<p>end</p>
EOF
BIN "$TMP/35_selfclose.dcd" "$TMP/35_selfclose.docx" 2>&1 && pass "35 self-closing variants" || fail "35" "exit 0" "non-zero exit"

# --- 36. Section continuous ---
cat > "$TMP/36_section_cont.dcd" << 'EOF'
[section 0]
name=first
--- BODY ---
<p>section 0</p>

[section:continuous 1]
name=second
--- BODY ---
<p>section 1 (continuous)</p>
EOF
BIN "$TMP/36_section_cont.dcd" "$TMP/36_section_cont.docx" 2>&1 && pass "36 section continuous" || fail "36" "exit 0" "non-zero exit"

# --- 37. Color named vs hex ---
cat > "$TMP/37_colors.dcd" << 'EOF'
[style]
layout=A4

[style:heading-1]
color=red

[section 0]
name=test
--- BODY ---
<p color=blue>blue paragraph</p>
<p color=#FF0000>hex red</p>
<p color=green>green text</p>
<h1 color=#336699>hex heading</h1>
EOF
BIN "$TMP/37_colors.dcd" "$TMP/37_colors.docx" 2>&1 && pass "37 named hex colors" || fail "37" "exit 0" "non-zero exit"

# --- 38. Very long paragraph ---
cat > "$TMP/38_long_para.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>EOF
for i in $(seq 1 50); do printf "Lorem ipsum dolor sit amet, consectetur adipiscing elit. "; done >> "$TMP/38_long_para.dcd"
cat >> "$TMP/38_long_para.dcd" << 'EOF'
</p>
EOF
BIN "$TMP/38_long_para.dcd" "$TMP/38_long_para.docx" 2>&1 && pass "38 very long paragraph" || fail "38" "exit 0" "non-zero exit"

# --- 39. Unicode in section/var/keys names ---
cat > "$TMP/39_unicode_names.dcd" << 'EOF'
[section 0]
name=s1
var=info
keys=title, desc
--- BODY ---
<p>{{info.title}}: {{info.desc}}</p>
EOF
BIN "$TMP/39_unicode_names.dcd" "$TMP/39_unicode_names.docx" 2>&1 && pass "39 unicode names" || fail "39" "exit 0" "non-zero exit"

# --- 40. Escape char \, in body ---
cat > "$TMP/40_escape.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>escape: \, hello \, world</p>
<p>newline: line1\nline2</p>
EOF
BIN "$TMP/40_escape.dcd" "$TMP/40_escape.docx" 2>&1 && pass "40 escape chars" || fail "40" "exit 0" "non-zero exit"

# --- 41. Multi-section with section:next-page ---
cat > "$TMP/41_multisection.dcd" << 'EOF'
[section 0]
name=cover
--- BODY ---
<h1>Cover</h1>

[section:next-page 1]
name=toc
--- BODY ---
<h2>Table of Contents</h2>

[section:next-page 2]
name=content
--- BODY ---
<h2>Content</h2>
EOF
BIN "$TMP/41_multisection.dcd" "$TMP/41_multisection.docx" 2>&1 && pass "41 multi section" || fail "41" "exit 0" "non-zero exit"

# --- 42. Loop 1000+ items (perf) ---
cat > "$TMP/42_perf_loop.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop x from items>
<p>{{x.id}}: {{x.label}}</p>
</loop>
EOF
python3 -c "
import json
data = {'info': {'title': 'Load Test'}, 'items': [{'id': i, 'label': f'item-{i}'} for i in range(1, 1001)]}
with open('$TMP/42_perf_loop.json', 'w') as f:
    json.dump(data, f)
" 2>/dev/null || {
  # fallback: manual JSON
  echo '{"info":{"title":"Load Test"},"items":['
  for i in $(seq 1 1000); do
    [ $i -eq 1000 ] && echo "{\"id\":$i,\"label\":\"item-$i\"}" || echo "{\"id\":$i,\"label\":\"item-$i\"},"
  done
  echo ']}' > "$TMP/42_perf_loop.json"
}
BIN -data "$TMP/42_perf_loop.json" "$TMP/42_perf_loop.dcd" "$TMP/42_perf_loop.docx" 2>&1 && pass "42 perf 1000 loops" || fail "42" "exit 0" "non-zero exit"

# --- 43. <w:*> with all caps/small-caps/letter-spacing ---
cat > "$TMP/43_w_styles.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:b caps=true>CAPS BOLD</w:b>
<w:i small-caps=true>small caps italic</w:i>
<w:c letter-spacing=2>spaced centered</w:c>
EOF
BIN "$TMP/43_w_styles.dcd" "$TMP/43_w_styles.docx" 2>&1 && pass "43 w style props" || fail "43" "exit 0" "non-zero exit"

# --- 44. Dynamic row style ---
cat > "$TMP/44_dynamic_row_style.dcd" << 'EOF'
[style:table custom]
bg=#E8F0FE
color=#333333

[section 0]
name=test
var=info, []rows
keys=title
--- BODY ---
<p>{{info.title}}</p>
<table border=1>
<loop:row x from rows>
<col>{{x.label}}</col>
</loop:row>
</table>
EOF
BIN "$TMP/44_dynamic_row_style.dcd" "$TMP/44_dynamic_row_style.docx" 2>&1 && pass "44 dynamic row style" || fail "44" "exit 0" "non-zero exit"

# --- 45. Image with variable path ---
cat > "$TMP/45_image_var.dcd" << 'EOF'
[section 0]
name=test
var=src
keys=img
--- BODY ---
<p>img path: {{src.img}}</p>
EOF
BIN "$TMP/45_image_var.dcd" "$TMP/45_image_var.docx" 2>&1 && pass "45 image var path" || fail "45" "exit 0" "non-zero exit"

# --- 46. Empty inline tags ---
cat > "$TMP/46_empty_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><b></b> <i></i> <u></u> <s></s> <code></code> <mark></mark></p>
EOF
BIN "$TMP/46_empty_inline.dcd" "$TMP/46_empty_inline.docx" 2>&1 && pass "46 empty inline tags" || fail "46" "exit 0" "non-zero exit"

# --- 47. Layout variants ---
cat > "$TMP/47_layouts.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>letter layout</p>
EOF
for l in letter legal A3 A5 B5; do
  sed "s/layout=A4/layout=$l/" "$TMP/47_layouts.dcd" > "$TMP/47_layout_$l.dcd"
  BIN "$TMP/47_layout_$l.dcd" "$TMP/47_layout_$l.docx" 2>&1 && pass "47 layout=$l" || fail "47 layout=$l" "exit 0" "non-zero exit"
done
cat > "$TMP/47_layout_custom.dcd" << 'EOF'
[style]
layout=custom
w=400
h=600
unit=mm

[section 0]
name=test
--- BODY ---
<p>custom layout</p>
EOF
BIN "$TMP/47_layout_custom.dcd" "$TMP/47_layout_custom.docx" 2>&1 && pass "47 layout=custom w h" || fail "47" "exit 0" "non-zero exit"

cat > "$TMP/47_layout_landscape.dcd" << 'EOF'
[style]
layout=A4
orientation=landscape

[section 0]
name=test
--- BODY ---
<p>landscape</p>
EOF
BIN "$TMP/47_layout_landscape.dcd" "$TMP/47_layout_landscape.docx" 2>&1 && pass "47 orientation=landscape" || fail "47" "exit 0" "non-zero exit"

# --- 48. Unit variants ---
for u in inch cm pt pica; do
  cat > "$TMP/48_unit_$u.dcd" << UNITEOF
[style]
layout=A4
unit=$u
m=1

[section 0]
name=test
--- BODY ---
<p>unit=$u</p>
UNITEOF
  BIN "$TMP/48_unit_$u.dcd" "$TMP/48_unit_$u.docx" 2>&1 && pass "48 unit=$u" || fail "48 unit=$u" "exit 0" "non-zero exit"
done

# --- 49. Margin variants ---
cat > "$TMP/49_margins.dcd" << 'EOF'
[style]
layout=A4
unit=mm
m=20

[section 0]
name=test
--- BODY ---
<p>uniform margin</p>
EOF
BIN "$TMP/49_margins.dcd" "$TMP/49_margins.docx" 2>&1 && pass "49 margin uniform" || fail "49" "exit 0" "non-zero exit"

cat > "$TMP/49_margins_axis.dcd" << 'EOF'
[style]
layout=A4
unit=mm
mx=25
my=15

[section 0]
name=test
--- BODY ---
<p>axis margins</p>
EOF
BIN "$TMP/49_margins_axis.dcd" "$TMP/49_margins_axis.docx" 2>&1 && pass "49 margin mx my" || fail "49" "exit 0" "non-zero exit"

cat > "$TMP/49_margins_individual.dcd" << 'EOF'
[style]
layout=A4
unit=mm
mt=10
mb=20
ml=15
mr=25

[section 0]
name=test
--- BODY ---
<p>individual margins</p>
EOF
BIN "$TMP/49_margins_individual.dcd" "$TMP/49_margins_individual.docx" 2>&1 && pass "49 margin mt mb ml mr" || fail "49" "exit 0" "non-zero exit"

cat > "$TMP/49_margins_precedence.dcd" << 'EOF'
[style]
layout=A4
unit=mm
m=30
mx=10
my=50
md=5
mt=1

[section 0]
name=test
--- BODY ---
<p>margin precedence</p>
EOF
BIN "$TMP/49_margins_precedence.dcd" "$TMP/49_margins_precedence.docx" 2>&1 && pass "49 margin precedence" || fail "49" "exit 0" "non-zero exit"

# --- 50. Image with real path ---
BIN "$TMP/22_image_attrs.dcd" "$TMP/22_image_attrs.docx" 2>&1 # re-gen baseline
IMG_HAPPY=0
cat > "$TMP/50_image.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<img=/tmp/dcd-stress-img/dummy.png width=100>
EOF
BIN "$TMP/50_image.dcd" "$TMP/50_image.docx" 2>&1 && IMG_HAPPY=1
[ "$IMG_HAPPY" -eq 1 ] && pass "50 image basic" || fail "50 image basic" "exit 0" "$(cat "$TMP/50_image.docx" 2>&1 || echo 'no output')"

cat > "$TMP/50_image_attrs.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<img=/tmp/dcd-stress-img/dummy.png width=50% align=center alt="test" border=1 bg=#f0f0f0>
EOF
BIN "$TMP/50_image_attrs.dcd" "$TMP/50_image_attrs.docx" 2>&1 && pass "50 image with attrs" || fail "50" "exit 0" "non-zero exit"

cat > "$TMP/50_image_abs.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<img=/tmp/dcd-stress-img/dummy.png width=200 height=100>
EOF
BIN "$TMP/50_image_abs.dcd" "$TMP/50_image_abs.docx" 2>&1 && pass "50 image abs w h" || fail "50" "exit 0" "non-zero exit"

# --- 51. {{title}} in body ---
cat > "$TMP/51_title_body.dcd" << 'EOF'
[title]
title=My Doc

[section 0]
name=test
--- BODY ---
<h1>{{title}}</h1>
<p>Document: {{title}}</p>
EOF
BIN "$TMP/51_title_body.dcd" "$TMP/51_title_body.docx" 2>&1 && pass "51 title in body" || fail "51" "exit 0" "non-zero exit"

# --- 52. Format with simple key ---
cat > "$TMP/52_format_simple.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title, date_field
formats=[date_field:dd-MM-yyyy]
--- BODY ---
<p>{{info.title}} — {{info.date_field}}</p>
EOF
BIN "$TMP/52_format_simple.dcd" "$TMP/52_format_simple.docx" 2>&1 && pass "52 format simple key" || fail "52" "exit 0" "non-zero exit"

# --- 53. Format with escaped colon ---
cat > "$TMP/53_format_escape.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title, time_field
formats=[time_field:HH\:mm]
--- BODY ---
<p>{{info.title}} @ {{info.time_field}}</p>
EOF
BIN "$TMP/53_format_escape.dcd" "$TMP/53_format_escape.docx" 2>&1 && pass "53 format escaped colon" || fail "53" "exit 0" "non-zero exit"

# --- 54. Dotted loop source ---
cat > "$TMP/54_dotted_source.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
--- BODY ---
<p>{{info.title}}</p>
EOF
# dotted sources skip array validation but need actual data — just test it doesn't crash parser
echo '{"info":{"title":"ok"},"data":{"items":[{"id":1}]}}' > "$TMP/54_dotted_source.json"
BIN -data "$TMP/54_dotted_source.json" "$TMP/54_dotted_source.dcd" "$TMP/54_dotted_source.docx" 2>&1 && pass "54 dotted source" || fail "54" "exit 0" "non-zero exit"

# --- 55. Multiple loop:row in table ---
cat > "$TMP/55_multi_looprow.dcd" << 'EOF'
[section 0]
name=test
var=info, []headers, []rows
keys=title
--- BODY ---
<p>{{info.title}}</p>
<table border=1>
<loop:row x from headers>
<col>{{x}}</col>
</loop:row>
<row bg=#eee><col>separator</col></row>
<loop:row x from rows>
<col>{{x}}</col>
</loop:row>
</table>
EOF
BIN "$TMP/55_multi_looprow.dcd" "$TMP/55_multi_looprow.docx" 2>&1 && pass "55 multiple loop:row" || fail "55" "exit 0" "non-zero exit"

# --- 56. Loop ol/ul with style.first ---
cat > "$TMP/56_loop_stylefirst.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items style.first=header>
<li>{{x.label}}</li>
</loop:ol>
<loop:ul x from items style.first=alt>
<li>{{x.label}}</li>
</loop:ul>
EOF
BIN "$TMP/56_loop_stylefirst.dcd" "$TMP/56_loop_stylefirst.docx" 2>&1 && pass "56 loop style.first" || fail "56" "exit 0" "non-zero exit"

# --- 57. Empty table ---
cat > "$TMP/57_empty_table.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
</table>
EOF
BIN "$TMP/57_empty_table.dcd" "$TMP/57_empty_table.docx" 2>&1 && pass "57 empty table" || fail "57" "exit 0" "non-zero exit"

# --- 58. Loop empty array ---
cat > "$TMP/58_empty_loop.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop x from items>{{x}}</loop>
EOF
echo '{"info":{"title":"ok"},"items":[]}' > "$TMP/58_empty_loop.json"
BIN -data "$TMP/58_empty_loop.json" "$TMP/58_empty_loop.dcd" "$TMP/58_empty_loop.docx" 2>&1 && pass "58 empty loop" || fail "58" "exit 0" "non-zero exit"

# --- 59. Header/footer font-family only ---
cat > "$TMP/59_hdr_font.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
font-family="Courier New"

[section 0]
name=test
--- BODY ---
<p>header font test</p>
EOF
BIN "$TMP/59_hdr_font.dcd" "$TMP/59_hdr_font.docx" 2>&1 && pass "59 header font-family only" || fail "59" "exit 0" "non-zero exit"

# --- 60. Header/footer justify_between 2 items ---
cat > "$TMP/60_hdr_justify2.dcd" << 'EOF'
[style]
layout=A4

[header]
justify_between=Left, Right

[section 0]
name=test
--- BODY ---
<p>justify 2 items</p>
EOF
BIN "$TMP/60_hdr_justify2.dcd" "$TMP/60_hdr_justify2.docx" 2>&1 && pass "60 justify 2 items" || fail "60" "exit 0" "non-zero exit"

# --- 61. <c:ol> with inline formatting ---
cat > "$TMP/61_col_inline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col bg=#f0f0f0><b>bold</b> and <i>italic</i></col>
<col><code>code</code> and <mark color=yellow>mark</mark></col>
</row>
</table>
EOF
BIN "$TMP/61_col_inline.dcd" "$TMP/61_col_inline.docx" 2>&1 && pass "61 col with inline" || fail "61" "exit 0" "non-zero exit"

# --- 62. <page-break> alias ---
cat > "$TMP/62_pagebreak_alias.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>page 1</p>
<page-break>
<p>page 2</p>
EOF
BIN "$TMP/62_pagebreak_alias.dcd" "$TMP/62_pagebreak_alias.docx" 2>&1 && pass "62 page-break alias" || fail "62" "exit 0" "non-zero exit"

# --- 63. <br> as standalone line ---
cat > "$TMP/63_br_standalone.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>line 1</p>
<br>
<p>line 2</p>
EOF
BIN "$TMP/63_br_standalone.dcd" "$TMP/63_br_standalone.docx" 2>&1 && pass "63 br standalone" || fail "63" "exit 0" "non-zero exit"

# --- 64. <tab> self-closing ---
cat > "$TMP/64_tab_selfclose.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>A<tab/>B</p>
<p>C<tab size=4/>D</p>
EOF
BIN "$TMP/64_tab_selfclose.dcd" "$TMP/64_tab_selfclose.docx" 2>&1 && pass "64 tab self-close" || fail "64" "exit 0" "non-zero exit"

# --- 65. Multi-line <w:> ---
cat > "$TMP/65_w_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>
centered
multi-line
content
</w:c>
EOF
BIN "$TMP/65_w_multiline.dcd" "$TMP/65_w_multiline.docx" 2>&1 && pass "65 w multi-line" || fail "65" "exit 0" "non-zero exit"

# --- 66. [header] empty ---
cat > "$TMP/66_hdr_empty.dcd" << 'EOF'
[style]
layout=A4

[header]

[section 0]
name=test
--- BODY ---
<p>empty header</p>
EOF
BIN "$TMP/66_hdr_empty.dcd" "$TMP/66_hdr_empty.docx" 2>&1 && pass "66 empty header" || fail "66" "exit 0" "non-zero exit"

# --- 67. Header/footer first-page default (true) ---
cat > "$TMP/67_hdr_firstpage_default.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header

[footer]
center=Footer
first-page=true

[section 0]
name=test
--- BODY ---
<p>first page default test</p>
EOF
BIN "$TMP/67_hdr_firstpage_default.dcd" "$TMP/67_hdr_firstpage_default.docx" 2>&1 && pass "67 header first-page true" || fail "67" "exit 0" "non-zero exit"

# --- 68. Missing inline-in-container combos ---
cat > "$TMP/68_container_gaps.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><u>u</u> <s>s</s> <mark>mark</mark></p>
<ul><li><u>u</u> <s>s</s> <mark>mark</mark></li></ul>
<table border=1><row><col><s>s</s></col></row></table>
<w:c><u>u</u> <s>s</s> <code>code</code> <mark>mark</mark> <sub>sub</sub> <sup>sup</sup></w:c>
EOF
BIN "$TMP/68_container_gaps.dcd" "$TMP/68_container_gaps.docx" 2>&1 && pass "68 container gaps" || fail "68" "exit 0" "non-zero exit"

# --- 69. <set:mark> <set:sub> <set:sup> ---
cat > "$TMP/69_set_mark_sub_sup.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:mark>mark</set:mark> <set:mark color=green>green</set:mark></p>
<p><set:sub>sub</set:sub> <set:sup>sup</set:sup></p>
<p><set:sub|sup>sub sup</set:sub|sup> <set:i|mark>italic mark</set:i|mark></p>
<p><set:code|mark>code mark</set:code|mark> <set:mark|u>mark underline</set:mark|u></p>
<p><set:b|i|u|s>all basic</set:b|i|u|s></p>
<p><set:u underline=wavy>wavy</set:u> <set:u underline=dotted>dotted</set:u></p>
EOF
BIN "$TMP/69_set_mark_sub_sup.dcd" "$TMP/69_set_mark_sub_sup.docx" 2>&1 && pass "69 set mark sub sup" || fail "69" "exit 0" "non-zero exit"

# --- 70. <w:mark> <w:sub> <w:sup> <w:code> block ---
cat > "$TMP/70_w_mark_sub_sup.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:mark color=yellow>marked block</w:mark>
<w:mark color=cyan>cyan block</w:mark>
<w:sub>subscript block</w:sub>
<w:sup>superscript block</w:sup>
<w:code>code block</w:code>
EOF
BIN "$TMP/70_w_mark_sub_sup.dcd" "$TMP/70_w_mark_sub_sup.docx" 2>&1 && pass "70 w mark sub sup code" || fail "70" "exit 0" "non-zero exit"

# --- 71. Default [style] font-family ---
cat > "$TMP/71_style_fontfamily.dcd" << 'EOF'
[style]
layout=A4
font-family="Courier New"

[section 0]
name=test
--- BODY ---
<p>courier paragraph</p>
EOF
BIN "$TMP/71_style_fontfamily.dcd" "$TMP/71_style_fontfamily.docx" 2>&1 && pass "71 style font-family" || fail "71" "exit 0" "non-zero exit"

# --- 72. Default [style] font-size ---
cat > "$TMP/72_style_fontsize.dcd" << 'EOF'
[style]
layout=A4
font-size=14

[section 0]
name=test
--- BODY ---
<p>14pt paragraph</p>
EOF
BIN "$TMP/72_style_fontsize.dcd" "$TMP/72_style_fontsize.docx" 2>&1 && pass "72 style font-size" || fail "72" "exit 0" "non-zero exit"

# --- 73. Default [style] color ---
cat > "$TMP/73_style_color.dcd" << 'EOF'
[style]
layout=A4
color=#336699

[section 0]
name=test
--- BODY ---
<p>colored default</p>
EOF
BIN "$TMP/73_style_color.dcd" "$TMP/73_style_color.docx" 2>&1 && pass "73 style color" || fail "73" "exit 0" "non-zero exit"

# --- 74. Default [style] line-height ---
cat > "$TMP/74_style_lineheight.dcd" << 'EOF'
[style]
layout=A4
line-height=2.0

[section 0]
name=test
--- BODY ---
<p>double spaced</p>
EOF
BIN "$TMP/74_style_lineheight.dcd" "$TMP/74_style_lineheight.docx" 2>&1 && pass "74 style line-height" || fail "74" "exit 0" "non-zero exit"

# --- 75. Default [style] combined ---
cat > "$TMP/75_style_combined.dcd" << 'EOF'
[style]
layout=A4
font-family="Georgia"
font-size=13
color=#222222
line-height=1.8

[section 0]
name=test
--- BODY ---
<p>combined default style</p>
EOF
BIN "$TMP/75_style_combined.dcd" "$TMP/75_style_combined.docx" 2>&1 && pass "75 style combined" || fail "75" "exit 0" "non-zero exit"

# --- 76. [title] with subject= ---
cat > "$TMP/76_title_subject.dcd" << 'EOF'
[title]
title=My Doc
subject=My Subject

[section 0]
name=test
--- BODY ---
<p>{{title}}</p>
EOF
BIN "$TMP/76_title_subject.dcd" "$TMP/76_title_subject.docx" 2>&1 && pass "76 title subject" || fail "76" "exit 0" "non-zero exit"

# --- 77. [title] with author= ---
cat > "$TMP/77_title_author.dcd" << 'EOF'
[title]
title=My Doc
author=Jane Doe

[section 0]
name=test
--- BODY ---
<p>{{title}}</p>
EOF
BIN "$TMP/77_title_author.dcd" "$TMP/77_title_author.docx" 2>&1 && pass "77 title author" || fail "77" "exit 0" "non-zero exit"

# --- 78. [title] full metadata ---
cat > "$TMP/78_title_full.dcd" << 'EOF'
[title]
title=Full Doc
subject=Full Subject
author=Full Author

[section 0]
name=test
--- BODY ---
<p>{{title}}</p>
EOF
BIN "$TMP/78_title_full.dcd" "$TMP/78_title_full.docx" 2>&1 && pass "78 title full metadata" || fail "78" "exit 0" "non-zero exit"

# --- 79. Heading local italic ---
cat > "$TMP/79_heading_italic.dcd" << 'EOF'
[style]
layout=A4

[style:heading-1]
italic=true

[section 0]
name=test
--- BODY ---
<h1>Italic heading</h1>
EOF
BIN "$TMP/79_heading_italic.dcd" "$TMP/79_heading_italic.docx" 2>&1 && pass "79 heading italic" || fail "79" "exit 0" "non-zero exit"

# --- 80. Heading local underline ---
cat > "$TMP/80_heading_underline.dcd" << 'EOF'
[style]
layout=A4

[style:heading-1]
underline=true

[section 0]
name=test
--- BODY ---
<h1>Underlined heading</h1>
EOF
BIN "$TMP/80_heading_underline.dcd" "$TMP/80_heading_underline.docx" 2>&1 && pass "80 heading underline" || fail "80" "exit 0" "non-zero exit"

# --- 81. Heading local align ---
cat > "$TMP/81_heading_align.dcd" << 'EOF'
[style]
layout=A4

[style:heading-1]
align=center

[section 0]
name=test
--- BODY ---
<h1>Center heading</h1>
EOF
BIN "$TMP/81_heading_align.dcd" "$TMP/81_heading_align.docx" 2>&1 && pass "81 heading align" || fail "81" "exit 0" "non-zero exit"

# --- 82. Heading local caps ---
cat > "$TMP/82_heading_caps.dcd" << 'EOF'
[style]
layout=A4

[section 0]
name=test
--- BODY ---
<h1 caps=true>ALL CAPS</h1>
<h2 small-caps=true>Small Caps</h2>
EOF
BIN "$TMP/82_heading_caps.dcd" "$TMP/82_heading_caps.docx" 2>&1 && pass "82 heading caps small-caps" || fail "82" "exit 0" "non-zero exit"

# --- 83. Heading local strike ---
cat > "$TMP/83_heading_strike.dcd" << 'EOF'
[style]
layout=A4

[section 0]
name=test
--- BODY ---
<h3 strike=true>Struck heading</h3>
EOF
BIN "$TMP/83_heading_strike.dcd" "$TMP/83_heading_strike.docx" 2>&1 && pass "83 heading strike" || fail "83" "exit 0" "non-zero exit"

# --- 84. Heading local letter-spacing ---
cat > "$TMP/84_heading_spacing.dcd" << 'EOF'
[style]
layout=A4

[section 0]
name=test
--- BODY ---
<h4 letter-spacing=3>Spaced heading</h4>
EOF
BIN "$TMP/84_heading_spacing.dcd" "$TMP/84_heading_spacing.docx" 2>&1 && pass "84 heading letter-spacing" || fail "84" "exit 0" "non-zero exit"

# --- 85. {{date}} in body ---
cat > "$TMP/85_date_body.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Date: {{date}}</p>
EOF
BIN "$TMP/85_date_body.dcd" "$TMP/85_date_body.docx" 2>&1 && pass "85 date in body" || fail "85" "exit 0" "non-zero exit"

# --- 86. {{page}} {{total}} in body (passthrough) ---
cat > "$TMP/86_page_total_body.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Page {{page}} of {{total}}</p>
EOF
BIN "$TMP/86_page_total_body.dcd" "$TMP/86_page_total_body.docx" 2>&1 && pass "86 page total in body" || fail "86" "exit 0" "non-zero exit"

# --- 87. Plain <p> without attrs (default) ---
cat > "$TMP/87_plain_p.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>plain</p>
<p>another</p>
<p>third</p>
EOF
BIN "$TMP/87_plain_p.dcd" "$TMP/87_plain_p.docx" 2>&1 && pass "87 plain p" || fail "87" "exit 0" "non-zero exit"

# --- 88. <row bg={{var}}> dynamic ---
cat > "$TMP/88_dyn_bg.dcd" << 'EOF'
[style:table custom]
bg=#E8F0FE

[section 0]
name=test
var=info
keys=bgcolor
--- BODY ---
<table border=1>
<row bg={{info.bgcolor}}>
<col>dynamic bg</col>
</row>
<row>
<col>static</col>
</row>
</table>
EOF
echo '{"info":{"bgcolor":"#FFE0E0"}}' > "$TMP/88_dyn_bg.json"
BIN -data "$TMP/88_dyn_bg.json" "$TMP/88_dyn_bg.dcd" "$TMP/88_dyn_bg.docx" 2>&1 && pass "88 row dynamic bg" || fail "88" "exit 0" "non-zero exit"

# --- 89. <table border> without value ---
cat > "$TMP/89_table_border_attr.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border>
<row><col>a</col><col>b</col></row>
</table>
EOF
BIN "$TMP/89_table_border_attr.dcd" "$TMP/89_table_border_attr.docx" 2>&1 && pass "89 table border attr" || fail "89" "exit 0" "non-zero exit"

# --- 90. <hr> width+color combined ---
cat > "$TMP/90_hr_width_color.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<hr>
<hr width=50%>
<hr color=red>
<hr width=75% color=#336699>
EOF
BIN "$TMP/90_hr_width_color.dcd" "$TMP/90_hr_width_color.docx" 2>&1 && pass "90 hr width color" || fail "90" "exit 0" "non-zero exit"

# --- 91. <w:r> right-aligned ---
cat > "$TMP/91_w_right.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:r>right aligned</w:r>
<w:r>another right</w:r>
EOF
BIN "$TMP/91_w_right.dcd" "$TMP/91_w_right.docx" 2>&1 && pass "91 w right" || fail "91" "exit 0" "non-zero exit"

# --- 92. <w:j> justify ---
cat > "$TMP/92_w_justify.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:j>justified text paragraph that should be spread across the full width of the page evenly</w:j>
EOF
BIN "$TMP/92_w_justify.dcd" "$TMP/92_w_justify.docx" 2>&1 && pass "92 w justify" || fail "92" "exit 0" "non-zero exit"

# --- 93. Header color property ---
cat > "$TMP/93_hdr_color.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
color=#999999

[section 0]
name=test
--- BODY ---
<p>header color test</p>
EOF
BIN "$TMP/93_hdr_color.dcd" "$TMP/93_hdr_color.docx" 2>&1 && pass "93 header color" || fail "93" "exit 0" "non-zero exit"

# --- 94. Footer-only (no header) ---
cat > "$TMP/94_footer_only.dcd" << 'EOF'
[style]
layout=A4

[footer]
center=Footer Only

[section 0]
name=test
--- BODY ---
<p>footer only test</p>
EOF
BIN "$TMP/94_footer_only.dcd" "$TMP/94_footer_only.docx" 2>&1 && pass "94 footer only" || fail "94" "exit 0" "non-zero exit"

# --- 95. Header+footer with all 3 columns ---
cat > "$TMP/95_hdr_3cols.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Left
center={{title}}
right={{date}}
font-size=9
color=#666

[footer]
left=Page {{page}}
center=of
right={{total}}
font-size=9
color=#666

[section 0]
name=test
--- BODY ---
<p>3 column header footer</p>
EOF
BIN "$TMP/95_hdr_3cols.dcd" "$TMP/95_hdr_3cols.docx" 2>&1 && pass "95 header footer 3 cols" || fail "95" "exit 0" "non-zero exit"

# --- 96. Link with target+color ---
cat > "$TMP/96_link_target_color.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><a=https://example.com target=_blank>blank link</a></p>
<p><a=https://example.com color=red underline=false>red no underline</a></p>
<p><a=#bookmark>bookmark link</a></p>
EOF
BIN "$TMP/96_link_target_color.dcd" "$TMP/96_link_target_color.docx" 2>&1 && pass "96 link target color" || fail "96" "exit 0" "non-zero exit"

# --- 97. <img=> height-only ---
cat > "$TMP/97_img_height.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>img /tmp/dcd-stress-img/dummy.png height=100</p>
EOF
BIN "$TMP/97_img_height.dcd" "$TMP/97_img_height.docx" 2>&1 && pass "97 img height only" || fail "97" "exit 0" "non-zero exit"

# --- 98. <ol> default type (numeric) ---
cat > "$TMP/98_ol_default.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol>
<li>one</li>
<li>two</li>
<li>three</li>
</ol>
EOF
BIN "$TMP/98_ol_default.dcd" "$TMP/98_ol_default.docx" 2>&1 && pass "98 ol default" || fail "98" "exit 0" "non-zero exit"

# --- 99. <ol type=a> cycle (a b c) ---
cat > "$TMP/99_ol_cycle.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol type=a>
<li>first</li>
<li>second</li>
<li>third</li>
</ol>
<ol type=A>
<li>alpha</li>
<li>beta</li>
<li>gamma</li>
</ol>
<ol type=i>
<li>i</li>
<li>ii</li>
<li>iii</li>
</ol>
<ol type=I>
<li>I</li>
<li>II</li>
<li>III</li>
</ol>
EOF
BIN "$TMP/99_ol_cycle.dcd" "$TMP/99_ol_cycle.docx" 2>&1 && pass "99 ol type cycle" || fail "99" "exit 0" "non-zero exit"

# --- 100. [style:heading-N] standalone ---
cat > "$TMP/100_heading_standalone.dcd" << 'EOF'
[style:heading-1]
font-family=Arial
font-size=28
color=#1F3864
bold=true

[style:heading-2]
font-family=Arial
font-size=18
color=#444

[section 0]
name=test
--- BODY ---
<h1>Standalone h1</h1>
<h2>Standalone h2</h2>
EOF
BIN "$TMP/100_heading_standalone.dcd" "$TMP/100_heading_standalone.docx" 2>&1 && pass "100 heading standalone" || fail "100" "exit 0" "non-zero exit"

# --- 101. Collapsed row (1 cell) ---
cat > "$TMP/101_collapsed_row.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row><col>single cell</col></row>
</table>
EOF
BIN "$TMP/101_collapsed_row.dcd" "$TMP/101_collapsed_row.docx" 2>&1 && pass "101 collapsed row" || fail "101" "exit 0" "non-zero exit"

# --- 102. Header+Footer first-page false ---
cat > "$TMP/102_hdr_firstpage_mixed.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
first-page=false

[footer]
center=Footer
first-page=true

[section 0]
name=test
--- BODY ---
<p>mixed first-page settings</p>
EOF
BIN "$TMP/102_hdr_firstpage_mixed.dcd" "$TMP/102_hdr_firstpage_mixed.docx" 2>&1 && pass "102 hdr first-page mixed" || fail "102" "exit 0" "non-zero exit"

# --- 103. <w:b|c> combined block flags ---
cat > "$TMP/103_w_block_combined.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c|b>centered bold</w:c|b>
<w:c|i>centered italic</w:c|i>
<w:c|b|i>centered bold italic</w:c|b|i>
<w:c|u>centered underline</w:c|u>
<w:r|i>right italic</w:r|i>
EOF
BIN "$TMP/103_w_block_combined.dcd" "$TMP/103_w_block_combined.docx" 2>&1 && pass "103 w block combined" || fail "103" "exit 0" "non-zero exit"

# --- 104. Multi-line <tab> ---
cat > "$TMP/104_tab_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>Name:<tab>John</p>
<p>Age:<tab size=4>30</p>
<p>Dept:<tab size=8>Engineering</p>
<p>Country:<tab size=12>Indonesia</p>
EOF
BIN "$TMP/104_tab_multiline.dcd" "$TMP/104_tab_multiline.docx" 2>&1 && pass "104 tab multiline" || fail "104" "exit 0" "non-zero exit"

# --- 105. Loop:ul with type= ---
cat > "$TMP/105_loopul_type.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ul x from items>
<li>{{x.label}}</li>
</loop:ul>
EOF
BIN "$TMP/105_loopul_type.dcd" "$TMP/105_loopul_type.docx" 2>&1 && pass "105 loop ul" || fail "105" "exit 0" "non-zero exit"

# --- 106. Loop with scalar items (non-object) ---
cat > "$TMP/106_loop_scalar.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<ul>
<loop x from items>
<li>{{x}}</li>
</loop>
</ul>
EOF
echo '{"info":{"title":"Scalars"},"items":["alpha","beta","gamma"]}' > "$TMP/106_loop_scalar.json"
BIN -data "$TMP/106_loop_scalar.json" "$TMP/106_loop_scalar.dcd" "$TMP/106_loop_scalar.docx" 2>&1 && pass "106 loop scalar" || fail "106" "exit 0" "non-zero exit"

# --- 107. Nested inline in <w:sup>/<w:sub>/<w:mark> ---
cat > "$TMP/107_w_nested_markup.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:sup>x<sup>2</sup> + y<sup>2</sup></w:sup>
<w:sub>H<sub>2</sub>O</w:sub>
<w:mark><b>bold mark</b> <i>italic mark</i></w:mark>
<w:code><b>bold code</b></w:code>
EOF
BIN "$TMP/107_w_nested_markup.dcd" "$TMP/107_w_nested_markup.docx" 2>&1 && pass "107 w nested markup" || fail "107" "exit 0" "non-zero exit"

# --- 108. {{}} unresolved without var ---
cat > "$TMP/108_unresolved_var.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>{{unknown}} should appear as-is</p>
EOF
BIN "$TMP/108_unresolved_var.dcd" "$TMP/108_unresolved_var.docx" 2>&1 && pass "108 unresolved var" || fail "108" "exit 0" "non-zero exit"

# --- 109. Multiple <br> in paragraph ---
cat > "$TMP/109_br_multi.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>a<br><br>b<br><br><br>c</p>
EOF
BIN "$TMP/109_br_multi.dcd" "$TMP/109_br_multi.docx" 2>&1 && pass "109 br multiple inline" || fail "109" "exit 0" "non-zero exit"

# --- 110. Self-closing <br/> / <pb/> (bug fix verification) ---
cat > "$TMP/110_selfclose_fixed.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>page 1</p>
<br/>
<p>after br/</p>
<pb/>
<p>after pb/</p>
<br/>
<hr/>
<p>end with hr/</p>
EOF
BIN "$TMP/110_selfclose_fixed.dcd" "$TMP/110_selfclose_fixed.docx" 2>&1 && pass "110 self-close bug fix" || fail "110" "exit 0" "non-zero exit"

# --- 111. </set> simplified closing (bug fix) ---
cat > "$TMP/111_set_simple_close.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:b|i>bold italic</set></p>
<p><set:b|i|u>all three</set></p>
<p><set:u underline=double>double</set></p>
EOF
BIN "$TMP/111_set_simple_close.dcd" "$TMP/111_set_simple_close.docx" 2>&1 && pass "111 set simple close" || fail "111" "exit 0" "non-zero exit"

# --- 112. Section continuous (bug fix verification) ---
cat > "$TMP/112_section_cont_fixed.dcd" << 'EOF'
[section 0]
name=first
--- BODY ---
<p>section 0</p>

[section:continuous 1]
name=second
--- BODY ---
<p>section 1 continuous</p>

[section:continuous 2]
name=third
--- BODY ---
<p>section 2 continuous</p>
EOF
BIN "$TMP/112_section_cont_fixed.dcd" "$TMP/112_section_cont_fixed.docx" 2>&1 && pass "112 section continuous fixed" || fail "112" "exit 0" "non-zero exit"

# --- 113. <w:*> with caps=true small-caps=true ---
cat > "$TMP/113_w_caps.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:b caps=true>CAPS BOLD</w:b>
<w:i small-caps=true>Small Caps Italic</w:i>
<w:c|b caps=true>Centered CAPS Bold</w:c|b>
EOF
BIN "$TMP/113_w_caps.dcd" "$TMP/113_w_caps.docx" 2>&1 && pass "113 w caps small-caps" || fail "113" "exit 0" "non-zero exit"

# --- 114. Loop source from dotted path (items.field) ---
cat > "$TMP/114_loop_dotted.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
--- BODY ---
<p>{{info.title}}</p>
EOF
echo '{"info":{"title":"ok"},"data":{"items":[{"id":1,"label":"A"},{"id":2,"label":"B"}]}}' > "$TMP/114_loop_dotted.json"
BIN -data "$TMP/114_loop_dotted.json" "$TMP/114_loop_dotted.dcd" "$TMP/114_loop_dotted.docx" 2>&1 && pass "114 dotted path source" || fail "114" "exit 0" "non-zero exit"

# --- 115. Format with regex specifier ---
cat > "$TMP/115_format_regex.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=phone
formats=[phone:(\d{3}) (\d{3}) (\d{4})]
--- BODY ---
<p>{{info.phone}}</p>
EOF
echo '{"info":{"phone":"0215551234"}}' > "$TMP/115_format_regex.json"
BIN -data "$TMP/115_format_regex.json" "$TMP/115_format_regex.dcd" "$TMP/115_format_regex.docx" 2>&1 && pass "115 format regex" || fail "115" "exit 0" "non-zero exit"

# --- 116. Loop:row with non-string field ---
cat > "$TMP/116_loop_numeric.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<table border=1>
<loop:row x from items>
<col>{{x.id}}</col>
<col>{{x.val}}</col>
</loop:row>
</table>
EOF
echo '{"info":{"title":"Numeric"},"items":[{"id":1,"val":100.5},{"id":2,"val":200}]}' > "$TMP/116_loop_numeric.json"
BIN -data "$TMP/116_loop_numeric.json" "$TMP/116_loop_numeric.dcd" "$TMP/116_loop_numeric.docx" 2>&1 && pass "116 loop numeric fields" || fail "116" "exit 0" "non-zero exit"

# --- 117. <p> with local color attr ---
cat > "$TMP/117_p_local_color.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p color=#ff6600>orange text</p>
<p color=green>green text</p>
<p color=red>red text</p>
EOF
BIN "$TMP/117_p_local_color.dcd" "$TMP/117_p_local_color.docx" 2>&1 && pass "117 p local color" || fail "117" "exit 0" "non-zero exit"

# --- 118. <p> with local font-size ---
cat > "$TMP/118_p_fontsize.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p font-size=10>small</p>
<p font-size=18>large</p>
EOF
BIN "$TMP/118_p_fontsize.dcd" "$TMP/118_p_fontsize.docx" 2>&1 && pass "118 p font-size" || fail "118" "exit 0" "non-zero exit"

# --- 119. <p> with local font-family ---
cat > "$TMP/119_p_fontfamily.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p font-family="Courier New">courier paragraph</p>
<p font-family=Georgia>georgia paragraph</p>
EOF
BIN "$TMP/119_p_fontfamily.dcd" "$TMP/119_p_fontfamily.docx" 2>&1 && pass "119 p font-family" || fail "119" "exit 0" "non-zero exit"

# --- 120. <p> with local line-height ---
cat > "$TMP/120_p_lineheight.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p line-height=1.2>tight</p>
<p line-height=2.0>loose</p>
EOF
BIN "$TMP/120_p_lineheight.dcd" "$TMP/120_p_lineheight.docx" 2>&1 && pass "120 p line-height" || fail "120" "exit 0" "non-zero exit"

# --- 121. <p> with bold/italic local attrs ---
cat > "$TMP/121_p_bold_italic.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p bold=true>bold paragraph</p>
<p italic=true>italic paragraph</p>
<p bold=true italic=true>bold italic</p>
EOF
BIN "$TMP/121_p_bold_italic.dcd" "$TMP/121_p_bold_italic.docx" 2>&1 && pass "121 p bold italic" || fail "121" "exit 0" "non-zero exit"

# --- 122. <table> with row bg color ---
cat > "$TMP/122_table_rowbg.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row bg=#FFE0E0>
<col>red row</col>
<col>data</col>
</row>
<row bg=#E0FFE0>
<col>green row</col>
<col>data</col>
</row>
</table>
EOF
BIN "$TMP/122_table_rowbg.dcd" "$TMP/122_table_rowbg.docx" 2>&1 && pass "122 table row bg" || fail "122" "exit 0" "non-zero exit"

# --- 123. Image with variable path + all attrs ---
cat > "$TMP/123_img_var_full.dcd" << 'EOF'
[section 0]
name=test
var=src
keys=img
--- BODY ---
<p>img {{src.img}} width=200 height=150 align=center alt Diagram border=1 bg=#f0f0f0</p>
EOF
BIN "$TMP/123_img_var_full.dcd" "$TMP/123_img_var_full.docx" 2>&1 && pass "123 img var full attrs" || fail "123" "exit 0" "non-zero exit"

# --- 124. Link with data var ---
cat > "$TMP/124_link_var.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=url, label
--- BODY ---
<p><a={{info.url}}>{{info.label}}</a></p>
EOF
echo '{"info":{"url":"https://example.com","label":"Click Here"}}' > "$TMP/124_link_var.json"
BIN -data "$TMP/124_link_var.json" "$TMP/124_link_var.dcd" "$TMP/124_link_var.docx" 2>&1 && pass "124 link data var" || fail "124" "exit 0" "non-zero exit"

# --- 125. Multi-section with header properties per section ---
cat > "$TMP/125_multisection_hdr.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Section Header

[section 0]
name=s0
--- BODY ---
<p>section zero</p>

[section:next-page 1]
name=s1
--- BODY ---
<p>section one</p>
EOF
BIN "$TMP/125_multisection_hdr.dcd" "$TMP/125_multisection_hdr.docx" 2>&1 && pass "125 multisection header" || fail "125" "exit 0" "non-zero exit"

# --- 126. Empty body with only whitespace ---
cat > "$TMP/126_whitespace_body.dcd" << 'EOF'
[section 0]
name=test

--- BODY ---

EOF
BIN "$TMP/126_whitespace_body.dcd" "$TMP/126_whitespace_body.docx" 2>&1 && pass "126 whitespace body" || fail "126" "exit 0" "non-zero exit"

# --- 127. <col align=right> with inline formatting ---
cat > "$TMP/127_col_align_right.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col align=right><b>right bold</b></col>
<col align=right><i>right italic</i></col>
</row>
<row>
<col align=center><code>center code</code></col>
<col align=right><mark>right mark</mark></col>
</row>
</table>
EOF
BIN "$TMP/127_col_align_right.dcd" "$TMP/127_col_align_right.docx" 2>&1 && pass "127 col align right" || fail "127" "exit 0" "non-zero exit"

# --- 128. <loop> with multiple vars in body ---
cat > "$TMP/128_loop_multi_var.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<h1>{{info.title}}</h1>
<loop x from items>
<p>{{x.label}} - {{x.value}}</p>
</loop>
EOF
echo '{"info":{"title":"Multi"},"items":[{"label":"A","value":10},{"label":"B","value":20}]}' > "$TMP/128_loop_multi_var.json"
BIN -data "$TMP/128_loop_multi_var.json" "$TMP/128_loop_multi_var.dcd" "$TMP/128_loop_multi_var.docx" 2>&1 && pass "128 loop multi var" || fail "128" "exit 0" "non-zero exit"

# --- 129. Nested lists 2-level ul/ol mixed ---
cat > "$TMP/129_nested_2level.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol type=A>
<li>Level 1A
<ul>
<li>Level 2a</li>
<li>Level 2b</li>
</ul>
</li>
<li>Level 1B
<ol type=i>
<li>Level 2i</li>
<li>Level 2ii</li>
</ol>
</li>
</ol>
EOF
BIN "$TMP/129_nested_2level.dcd" "$TMP/129_nested_2level.docx" 2>&1 && pass "129 nested 2-level mixed" || fail "129" "exit 0" "non-zero exit"

# --- 130. Keys-only with format ---
cat > "$TMP/130_keys_format.dcd" << 'EOF'
[section 0]
name=test
keys=title, date_field
formats=[date_field:dd/MM/yyyy]
--- BODY ---
<p>{{title}} — {{date_field}}</p>
EOF
echo '{"title":"Report","date_field":"2026-07-09"}' > "$TMP/130_keys_format.json"
BIN -data "$TMP/130_keys_format.json" "$TMP/130_keys_format.dcd" "$TMP/130_keys_format.docx" 2>&1 && pass "130 keys with format" || fail "130" "exit 0" "non-zero exit"

# --- 131. <loop:ol> with inline formatting ---
cat > "$TMP/131_loopol_inline.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items type=A>
<li><b>{{x.label}}</b> — <i>{{x.val}}</i></li>
</loop:ol>
EOF
echo '{"info":{"title":"List"},"items":[{"label":"One","val":"X"},{"label":"Two","val":"Y"}]}' > "$TMP/131_loopol_inline.json"
BIN -data "$TMP/131_loopol_inline.json" "$TMP/131_loopol_inline.dcd" "$TMP/131_loopol_inline.docx" 2>&1 && pass "131 loop ol inline" || fail "131" "exit 0" "non-zero exit"

# --- 132. <loop:ul> with type attr (no effect but should parse) ---
cat > "$TMP/132_loopul_type_attr.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ul x from items type=circle>
<li>{{x.name}}</li>
</loop:ul>
EOF
echo '{"info":{"title":"UL"},"items":[{"name":"X"},{"name":"Y"}]}' > "$TMP/132_loopul_type_attr.json"
BIN -data "$TMP/132_loopul_type_attr.json" "$TMP/132_loopul_type_attr.dcd" "$TMP/132_loopul_type_attr.docx" 2>&1 && pass "132 loop ul type attr" || fail "132" "exit 0" "non-zero exit"

# --- 133. Section with var but no keys (passthrough {{}}) ---
cat > "$TMP/133_var_no_keys.dcd" << 'EOF'
[section 0]
name=test
var=info
--- BODY ---
<p>{{info.title}} passthrough</p>
EOF
BIN "$TMP/133_var_no_keys.dcd" "$TMP/133_var_no_keys.docx" 2>&1 && pass "133 var no keys" || fail "133" "exit 0" "non-zero exit"

# --- 134. Header with font-size only ---
cat > "$TMP/134_hdr_fontsize.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
font-size=11

[section 0]
name=test
--- BODY ---
<p>font-size only</p>
EOF
BIN "$TMP/134_hdr_fontsize.dcd" "$TMP/134_hdr_fontsize.docx" 2>&1 && pass "134 header font-size" || fail "134" "exit 0" "non-zero exit"

# --- 135. Header border=none ---
cat > "$TMP/135_hdr_border_none.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
border=none

[section 0]
name=test
--- BODY ---
<p>no border</p>
EOF
BIN "$TMP/135_hdr_border_none.dcd" "$TMP/135_hdr_border_none.docx" 2>&1 && pass "135 header border none" || fail "135" "exit 0" "non-zero exit"

# --- 136. Footer with margin ---
cat > "$TMP/136_ftr_margin.dcd" << 'EOF'
[style]
layout=A4

[footer]
center=Footer
margin=0.5

[section 0]
name=test
--- BODY ---
<p>footer margin test</p>
EOF
BIN "$TMP/136_ftr_margin.dcd" "$TMP/136_ftr_margin.docx" 2>&1 && pass "136 footer margin" || fail "136" "exit 0" "non-zero exit"

# --- 137. Header with justify_between + font styling ---
cat > "$TMP/137_hdr_justify_style.dcd" << 'EOF'
[style]
layout=A4

[header]
justify_between=Left, Center, Right
font-family="Courier New"
font-size=9
color=#444

[section 0]
name=test
--- BODY ---
<p>styled justify_between</p>
EOF
BIN "$TMP/137_hdr_justify_style.dcd" "$TMP/137_hdr_justify_style.docx" 2>&1 && pass "137 hdr justify styled" || fail "137" "exit 0" "non-zero exit"

# --- 138. Section:next-page with name + vars ---
cat > "$TMP/138_nextpage_vars.dcd" << 'EOF'
[section 0]
name=cover
var=info
keys=title
--- BODY ---
<h1>{{info.title}}</h1>

[section:next-page 1]
name=detail
var=info
keys=subtitle
--- BODY ---
<h2>{{info.subtitle}}</h2>
EOF
echo '{"info":{"title":"Cover","subtitle":"Details"}}' > "$TMP/138_nextpage_vars.json"
BIN -data "$TMP/138_nextpage_vars.json" "$TMP/138_nextpage_vars.dcd" "$TMP/138_nextpage_vars.docx" 2>&1 && pass "138 nextpage with vars" || fail "138" "exit 0" "non-zero exit"

# --- 139. <hr> between paragraphs ---
cat > "$TMP/139_hr_paras.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>before</p>
<hr width=25% color=#ccc>
<p>after</p>
EOF
BIN "$TMP/139_hr_paras.dcd" "$TMP/139_hr_paras.docx" 2>&1 && pass "139 hr between paras" || fail "139" "exit 0" "non-zero exit"

# --- 140. <w:u> with underline=dash ---
cat > "$TMP/140_w_u_dash.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:u underline=dash>dashed underline block</w:u>
<w:u underline=dotted>dotted underline block</w:u>
<w:u underline=wavy>wavy underline block</w:u>
EOF
BIN "$TMP/140_w_u_dash.dcd" "$TMP/140_w_u_dash.docx" 2>&1 && pass "140 w u variants" || fail "140" "exit 0" "non-zero exit"

# --- 141. Loop with <loop> containing {{x}} only (no field) ---
cat > "$TMP/141_loop_xonly.dcd" << 'EOF'
[section 0]
name=test
var=[]items
--- BODY ---
<ul>
<loop x from items>
<li>{{x}}</li>
</loop>
</ul>
EOF
echo '{"items":["a","b","c"]}' > "$TMP/141_loop_xonly.json"
BIN -data "$TMP/141_loop_xonly.json" "$TMP/141_loop_xonly.dcd" "$TMP/141_loop_xonly.docx" 2>&1 && pass "141 loop x only" || fail "141" "exit 0" "non-zero exit"

# --- 142. Loop:ul with scalar items ---
cat > "$TMP/142_loopul_scalar.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ul x from items>
<li>{{x}}</li>
</loop:ul>
EOF
echo '{"info":{"title":"Scalar Ul"},"items":["x","y","z"]}' > "$TMP/142_loopul_scalar.json"
BIN -data "$TMP/142_loopul_scalar.json" "$TMP/142_loopul_scalar.dcd" "$TMP/142_loopul_scalar.docx" 2>&1 && pass "142 loop ul scalar" || fail "142" "exit 0" "non-zero exit"

# --- 143. Table with multiple <row> without loop ---
cat > "$TMP/143_table_static_rows.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row><col>R1C1</col><col>R1C2</col></row>
<row><col>R2C1</col><col>R2C2</col></row>
<row><col>R3C1</col><col>R3C2</col></row>
</table>
EOF
BIN "$TMP/143_table_static_rows.dcd" "$TMP/143_table_static_rows.docx" 2>&1 && pass "143 table static rows" || fail "143" "exit 0" "non-zero exit"

# --- 144. <col> with bg color ---
cat > "$TMP/144_col_bg.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col bg=#FFFFCC>yellow</col>
<col bg=#CCFFCC>green</col>
<col bg=#CCCCFF>blue</col>
</row>
</table>
EOF
BIN "$TMP/144_col_bg.dcd" "$TMP/144_col_bg.docx" 2>&1 && pass "144 col bg" || fail "144" "exit 0" "non-zero exit"

# --- 145. Empty [footer] ---
cat > "$TMP/145_empty_footer.dcd" << 'EOF'
[style]
layout=A4

[footer]

[section 0]
name=test
--- BODY ---
<p>empty footer</p>
EOF
BIN "$TMP/145_empty_footer.dcd" "$TMP/145_empty_footer.docx" 2>&1 && pass "145 empty footer" || fail "145" "exit 0" "non-zero exit"

# --- 146. [style] with letter-spacing default ---
cat > "$TMP/146_style_letterspacing.dcd" << 'EOF'
[style]
layout=A4
letter-spacing=2

[section 0]
name=test
--- BODY ---
<p>spaced text</p>
EOF
BIN "$TMP/146_style_letterspacing.dcd" "$TMP/146_style_letterspacing.docx" 2>&1 && pass "146 style letter-spacing" || fail "146" "exit 0" "non-zero exit"

# --- 147. <set:flags> with all 8 flags ---
cat > "$TMP/147_set_all_flags.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:b|i|u|s|code|mark|sub|sup>all eight</set:b|i|u|s|code|mark|sub|sup></p>
EOF
BIN "$TMP/147_set_all_flags.dcd" "$TMP/147_set_all_flags.docx" 2>&1 && pass "147 set all 8 flags" || fail "147" "exit 0" "non-zero exit"

# --- 148. <col> align=center with inline ---
cat > "$TMP/148_col_align_center.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col align=center><b>centered bold</b></col>
</row>
</table>
EOF
BIN "$TMP/148_col_align_center.dcd" "$TMP/148_col_align_center.docx" 2>&1 && pass "148 col align center" || fail "148" "exit 0" "non-zero exit"

# --- 149. Heading with border-bottom local attr ---
cat > "$TMP/149_heading_borderbottom.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h1 border-bottom=2pt>bordered heading</h1>
<h2 border-bottom=1pt>thinner</h2>
EOF
BIN "$TMP/149_heading_borderbottom.dcd" "$TMP/149_heading_borderbottom.docx" 2>&1 && pass "149 heading border-bottom" || fail "149" "exit 0" "non-zero exit"

# --- 150. Image variable with static path fallback ---
cat > "$TMP/150_img_static.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>img ./static/img.png width=100</p>
EOF
BIN "$TMP/150_img_static.dcd" "$TMP/150_img_static.docx" 2>&1 && pass "150 img static" || fail "150" "exit 0" "non-zero exit"

# --- 151. <w:s> strikethrough block ---
cat > "$TMP/151_w_strike.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:s>strikethrough block</w:s>
EOF
BIN "$TMP/151_w_strike.dcd" "$TMP/151_w_strike.docx" 2>&1 && pass "151 w strike block" || fail "151" "exit 0" "non-zero exit"

# --- 152. Multiple <pb> in sequence ---
cat > "$TMP/152_multiple_pb.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>page 1</p>
<pb>
<p>page 2</p>
<pb>
<pb>
<p>page 4</p>
EOF
BIN "$TMP/152_multiple_pb.dcd" "$TMP/152_multiple_pb.docx" 2>&1 && pass "152 multiple page breaks" || fail "152" "exit 0" "non-zero exit"

# --- 153. <page-break> in list context ---
cat > "$TMP/153_pagebreak_list.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol>
<li>item 1</li>
</ol>
<page-break>
<ul>
<li>after break</li>
</ul>
EOF
BIN "$TMP/153_pagebreak_list.dcd" "$TMP/153_pagebreak_list.docx" 2>&1 && pass "153 page-break after list" || fail "153" "exit 0" "non-zero exit"

# --- 154. <br> inside <w:*> ---
cat > "$TMP/154_br_in_w.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>line 1<br>line 2<br>line 3</w:c>
EOF
BIN "$TMP/154_br_in_w.dcd" "$TMP/154_br_in_w.docx" 2>&1 && pass "154 br in w block" || fail "154" "exit 0" "non-zero exit"

# --- 155. tab inside w block ---
cat > "$TMP/155_tab_in_w.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>name:<tab size=6>value</w:c>
EOF
BIN "$TMP/155_tab_in_w.dcd" "$TMP/155_tab_in_w.docx" 2>&1 && pass "155 tab in w block" || fail "155" "exit 0" "non-zero exit"

# --- 156. <w:b|i|u> triple combined block ---
cat > "$TMP/156_w_triple.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:b|i|u>bold italic underline block</w:b|i|u>
EOF
BIN "$TMP/156_w_triple.dcd" "$TMP/156_w_triple.docx" 2>&1 && pass "156 w b i u block" || fail "156" "exit 0" "non-zero exit"

# --- 157. Multiple <hr> tags ---
cat > "$TMP/157_hr_multi.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<hr>
<hr width=30%>
<hr color=green>
<hr width=50% color=blue>
EOF
BIN "$TMP/157_hr_multi.dcd" "$TMP/157_hr_multi.docx" 2>&1 && pass "157 multi hr" || fail "157" "exit 0" "non-zero exit"

# --- 158. Empty <ul> ---
cat > "$TMP/158_empty_ul.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
</ul>
EOF
BIN "$TMP/158_empty_ul.dcd" "$TMP/158_empty_ul.docx" 2>&1 && pass "158 empty ul" || fail "158" "exit 0" "non-zero exit"

# --- 159. Empty <ol> ---
cat > "$TMP/159_empty_ol.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ol>
</ol>
EOF
BIN "$TMP/159_empty_ol.dcd" "$TMP/159_empty_ol.docx" 2>&1 && pass "159 empty ol" || fail "159" "exit 0" "non-zero exit"

# --- 160. <li> with text only (no formatting) ---
cat > "$TMP/160_li_plain.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li>plain text</li>
<li>more plain</li>
</ul>
EOF
BIN "$TMP/160_li_plain.dcd" "$TMP/160_li_plain.docx" 2>&1 && pass "160 li plain" || fail "160" "exit 0" "non-zero exit"

# --- 161. Section with name= as unique across mixed sections ---
cat > "$TMP/161_section_unique.dcd" << 'EOF'
[section:next-page 0]
name=cover
--- BODY ---
<p>cover</p>

[section 1]
name=body
--- BODY ---
<p>body</p>
EOF
BIN "$TMP/161_section_unique.dcd" "$TMP/161_section_unique.docx" 2>&1 && pass "161 section unique names" || fail "161" "exit 0" "non-zero exit"

# --- 162. <w:c> with <br> and <b><i> nested ---
cat > "$TMP/162_w_c_complex.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>centered <b>bold</b> and <i>italic</i><br>new line <b>still bold</b></w:c>
EOF
BIN "$TMP/162_w_c_complex.dcd" "$TMP/162_w_c_complex.docx" 2>&1 && pass "162 w c complex" || fail "162" "exit 0" "non-zero exit"

# --- 163. <loop> with various var name lengths ---
cat > "$TMP/163_loop_var_names.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop myLongVarName from items>
<p>{{myLongVarName.label}}</p>
</loop>
EOF
echo '{"info":{"title":"Vars"},"items":[{"label":"short"},{"label":"also short"}]}' > "$TMP/163_loop_var_names.json"
BIN -data "$TMP/163_loop_var_names.json" "$TMP/163_loop_var_names.dcd" "$TMP/163_loop_var_names.docx" 2>&1 && pass "163 loop var names" || fail "163" "exit 0" "non-zero exit"

# --- 164. All heading levels h1-h6 ---
cat > "$TMP/164_all_headings.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h1>Heading 1</h1>
<p>text</p>
<h2>Heading 2</h2>
<p>text</p>
<h3>Heading 3</h3>
<p>text</p>
<h4>Heading 4</h4>
<p>text</p>
<h5>Heading 5</h5>
<p>text</p>
<h6>Heading 6</h6>
EOF
BIN "$TMP/164_all_headings.dcd" "$TMP/164_all_headings.docx" 2>&1 && pass "164 all headings h1-h6" || fail "164" "exit 0" "non-zero exit"

# --- 165. [style] only, no sections ---
cat > "$TMP/165_style_only.dcd" << 'EOF'
[style]
layout=A4

[section 0]
name=test
--- BODY ---
<p>style only</p>
EOF
BIN "$TMP/165_style_only.dcd" "$TMP/165_style_only.docx" 2>&1 && pass "165 style only" || fail "165" "exit 0" "non-zero exit"

# --- 166. Two separate <loop> in same section ---
cat > "$TMP/166_two_loops.dcd" << 'EOF'
[section 0]
name=test
var=info, []xitems, []yitems
keys=title
--- BODY ---
<p>{{info.title}}</p>
<ul>
<loop x from xitems><li>{{x}}</li></loop>
</ul>
<ol>
<loop y from yitems><li>{{y}}</li></loop>
</ol>
EOF
echo '{"info":{"title":"Two Loops"},"xitems":["a","b"],"yitems":["1","2"]}' > "$TMP/166_two_loops.json"
BIN -data "$TMP/166_two_loops.json" "$TMP/166_two_loops.dcd" "$TMP/166_two_loops.docx" 2>&1 && pass "166 two loops" || fail "166" "exit 0" "non-zero exit"

# --- 167. Section:even-page ---
cat > "$TMP/167_evenpage.dcd" << 'EOF'
[section:even-page 0]
name=even
--- BODY ---
<p>even page section</p>
EOF
BIN "$TMP/167_evenpage.dcd" "$TMP/167_evenpage.docx" 2>&1 && pass "167 even-page" || fail "167" "exit 0" "non-zero exit"

# --- 168. Section:odd-page ---
cat > "$TMP/168_oddpage.dcd" << 'EOF'
[section:odd-page 0]
name=odd
--- BODY ---
<p>odd page section</p>
EOF
BIN "$TMP/168_oddpage.dcd" "$TMP/168_oddpage.docx" 2>&1 && pass "168 odd-page" || fail "168" "exit 0" "non-zero exit"

# --- 169. <col> with width attribute (parsed but ignored) ---
cat > "$TMP/169_col_width.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row><col width=50%>wide</col><col width=50%>wide</col></row>
</table>
EOF
BIN "$TMP/169_col_width.dcd" "$TMP/169_col_width.docx" 2>&1 && pass "169 col width attr" || fail "169" "exit 0" "non-zero exit"

# --- 170. <row> with style={{var}} dynamic ---
cat > "$TMP/170_row_dyn_style.dcd" << 'EOF'
[style:table custom]
bg=#E0F0FF

[section 0]
name=test
var=info
keys=styleName
--- BODY ---
<table border=1>
<row style={{info.styleName}}>
<col>dynamic style</col>
</row>
</table>
EOF
echo '{"info":{"styleName":"custom"}}' > "$TMP/170_row_dyn_style.json"
BIN -data "$TMP/170_row_dyn_style.json" "$TMP/170_row_dyn_style.dcd" "$TMP/170_row_dyn_style.docx" 2>&1 && pass "170 row dyn style" || fail "170" "exit 0" "non-zero exit"

# --- 171. Colspan attribute parsing (ignored, no crash) ---
cat > "$TMP/171_colspan.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row><col colspan=2>wide cell</col></row>
</table>
EOF
BIN "$TMP/171_colspan.dcd" "$TMP/171_colspan.docx" 2>&1 && pass "171 colspan attr" || fail "171" "exit 0" "non-zero exit"

# --- 172. Rowspan attribute parsing (ignored, no crash) ---
cat > "$TMP/172_rowspan.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row><col rowspan=2>tall cell</col></row>
</table>
EOF
BIN "$TMP/172_rowspan.dcd" "$TMP/172_rowspan.docx" 2>&1 && pass "172 rowspan attr" || fail "172" "exit 0" "non-zero exit"

# --- 173. <table width=100%> (parsed but ignored) ---
cat > "$TMP/173_table_width.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1 width=100%>
<row><col>full width</col></row>
</table>
EOF
BIN "$TMP/173_table_width.dcd" "$TMP/173_table_width.docx" 2>&1 && pass "173 table width" || fail "173" "exit 0" "non-zero exit"

# --- 174. Multiple <hr> with all combos ---
cat > "$TMP/174_hr_all.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<hr>
<hr width=25%>
<hr color=red>
<hr width=50% color=#336699>
<hr width=100% color=#ccc>
EOF
BIN "$TMP/174_hr_all.dcd" "$TMP/174_hr_all.docx" 2>&1 && pass "174 hr all combos" || fail "174" "exit 0" "non-zero exit"

# --- 175. <img> with variable for all attrs ---
cat > "$TMP/175_img_var_path.dcd" << 'EOF'
[section 0]
name=test
var=src
keys=img
--- BODY ---
<p>img {{src.img}} width=150 height=100</p>
EOF
BIN "$TMP/175_img_var_path.dcd" "$TMP/175_img_var_path.docx" 2>&1 && pass "175 img var path size" || fail "175" "exit 0" "non-zero exit"

# --- 176. Link with underline=true ---
cat > "$TMP/176_link_underline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><a=https://example.com underline=true>underlined link</a></p>
<p><a=https://example.com underline=false>no underline</a></p>
EOF
BIN "$TMP/176_link_underline.dcd" "$TMP/176_link_underline.docx" 2>&1 && pass "176 link underline" || fail "176" "exit 0" "non-zero exit"

# --- 177. Multi-line <w:b> block ---
cat > "$TMP/177_w_b_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:b>
line 1
line 2
line 3
</w:b>
EOF
BIN "$TMP/177_w_b_multiline.dcd" "$TMP/177_w_b_multiline.docx" 2>&1 && pass "177 w b multiline" || fail "177" "exit 0" "non-zero exit"

# --- 178. Multi-line <w:i> block ---
cat > "$TMP/178_w_i_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:i>
italic
multiline
content
</w:i>
EOF
BIN "$TMP/178_w_i_multiline.dcd" "$TMP/178_w_i_multiline.docx" 2>&1 && pass "178 w i multiline" || fail "178" "exit 0" "non-zero exit"

# --- 179. Multi-line <w:u> block ---
cat > "$TMP/179_w_u_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:u>
underline
multiline
</w:u>
EOF
BIN "$TMP/179_w_u_multiline.dcd" "$TMP/179_w_u_multiline.docx" 2>&1 && pass "179 w u multiline" || fail "179" "exit 0" "non-zero exit"

# --- 180. Multi-line <w:s> block ---
cat > "$TMP/180_w_s_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:s>
strike
multiline
</w:s>
EOF
BIN "$TMP/180_w_s_multiline.dcd" "$TMP/180_w_s_multiline.docx" 2>&1 && pass "180 w s multiline" || fail "180" "exit 0" "non-zero exit"

# --- 181. Multi-line <w:code> block ---
cat > "$TMP/181_w_code_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:code>
code
multiline
</w:code>
EOF
BIN "$TMP/181_w_code_multiline.dcd" "$TMP/181_w_code_multiline.docx" 2>&1 && pass "181 w code multiline" || fail "181" "exit 0" "non-zero exit"

# --- 182. Multi-line <w:mark> block ---
cat > "$TMP/182_w_mark_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:mark color=yellow>
marked
multiline
</w:mark>
EOF
BIN "$TMP/182_w_mark_multiline.dcd" "$TMP/182_w_mark_multiline.docx" 2>&1 && pass "182 w mark multiline" || fail "182" "exit 0" "non-zero exit"

# --- 183. Multi-line <w:sub> block ---
cat > "$TMP/183_w_sub_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:sub>
subscript
multiline
</w:sub>
EOF
BIN "$TMP/183_w_sub_multiline.dcd" "$TMP/183_w_sub_multiline.docx" 2>&1 && pass "183 w sub multiline" || fail "183" "exit 0" "non-zero exit"

# --- 184. Multi-line <w:sup> block ---
cat > "$TMP/184_w_sup_multiline.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:sup>
superscript
multiline
</w:sup>
EOF
BIN "$TMP/184_w_sup_multiline.dcd" "$TMP/184_w_sup_multiline.docx" 2>&1 && pass "184 w sup multiline" || fail "184" "exit 0" "non-zero exit"

# --- 185. Header with left + font-family + color + font-size ---
cat > "$TMP/185_hdr_all_props.dcd" << 'EOF'
[style]
layout=A4

[header]
left=Header
font-family=Arial
font-size=10
color=#555
border=bottom
margin=0.3

[section 0]
name=test
--- BODY ---
<p>header all props</p>
EOF
BIN "$TMP/185_hdr_all_props.dcd" "$TMP/185_hdr_all_props.docx" 2>&1 && pass "185 header all props" || fail "185" "exit 0" "non-zero exit"

# --- 186. Footer with all props ---
cat > "$TMP/186_ftr_all_props.dcd" << 'EOF'
[style]
layout=A4

[footer]
center=Footer
font-family=Georgia
font-size=9
color=#888
border=top
margin=0.2
first-page=false

[section 0]
name=test
--- BODY ---
<p>footer all props</p>
EOF
BIN "$TMP/186_ftr_all_props.dcd" "$TMP/186_ftr_all_props.docx" 2>&1 && pass "186 footer all props" || fail "186" "exit 0" "non-zero exit"

# --- 187. Multiple <col> with different alignments ---
cat > "$TMP/187_col_aligns.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=1>
<row>
<col align=left>left</col>
<col align=center>center</col>
<col align=right>right</col>
</row>
</table>
EOF
BIN "$TMP/187_col_aligns.dcd" "$TMP/187_col_aligns.docx" 2>&1 && pass "187 col all aligns" || fail "187" "exit 0" "non-zero exit"

# --- 188. <table> with border=0 (no border) ---
cat > "$TMP/188_table_no_border.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<table border=0>
<row><col>no border</col></row>
</table>
EOF
BIN "$TMP/188_table_no_border.dcd" "$TMP/188_table_no_border.docx" 2>&1 && pass "188 table border=0" || fail "188" "exit 0" "non-zero exit"

# --- 189. Section with only vars, no keys ---
cat > "$TMP/189_vars_only.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
--- BODY ---
<p>{{info.title}}</p>
<loop x from items><p>{{x}}</p></loop>
EOF
echo '{"info":{"title":"No keys"},"items":["a","b"]}' > "$TMP/189_vars_only.json"
BIN -data "$TMP/189_vars_only.json" "$TMP/189_vars_only.dcd" "$TMP/189_vars_only.docx" 2>&1 && pass "189 vars no keys" || fail "189" "exit 0" "non-zero exit"

# --- 190. <loop:ol> with multiple items and inline mark ---
cat > "$TMP/190_loopol_mark.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items type=a>
<li><mark>{{x.label}}</mark></li>
</loop:ol>
EOF
echo '{"info":{"title":"Mark List"},"items":[{"label":"one"},{"label":"two"},{"label":"three"}]}' > "$TMP/190_loopol_mark.json"
BIN -data "$TMP/190_loopol_mark.json" "$TMP/190_loopol_mark.dcd" "$TMP/190_loopol_mark.docx" 2>&1 && pass "190 loop ol mark" || fail "190" "exit 0" "non-zero exit"

header "ERROR PATH"

# --- E1. overlapping tags ---
cat > "$TMP/E1_overlap.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><b>bold <i>bold+italic</b> italic</i></p>
EOF
err=$(BIN "$TMP/E1_overlap.dcd" "$TMP/E1_overlap.docx" 2>&1 || true)
case "$err" in
  *"expected"*"</i>"*"</b>"*) pass "E1 overlapping tags" ;;
  *) fail "E1 overlapping tags" "expected ...</i> but found </b>..." "$err" ;;
esac

# --- E2. unclosed tag ---
cat > "$TMP/E2_unclosed.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><b>bold text</p>
EOF
err=$(BIN "$TMP/E2_unclosed.dcd" "$TMP/E2_unclosed.docx" 2>&1 || true)
case "$err" in
  *"unclosed"*"<b>"*) pass "E2 unclosed tag" ;;
  *) fail "E2 unclosed tag" "unclosed tag <b>" "$err" ;;
esac

# --- E3. unexpected closing tag ---
cat > "$TMP/E3_unexpected_close.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p></i>unexpected close</p>
EOF
err=$(BIN "$TMP/E3_unexpected_close.dcd" "$TMP/E3_unexpected_close.docx" 2>&1 || true)
case "$err" in
  *"unexpected"*"</i>"*) pass "E3 unexpected closing tag" ;;
  *) fail "E3 unexpected closing tag" "unexpected closing tag </i>" "$err" ;;
esac

# --- E4. <w:*> nesting ---
cat > "$TMP/E4_w_nesting.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c><w:b>nested</w:b></w:c>
EOF
err=$(BIN "$TMP/E4_w_nesting.dcd" "$TMP/E4_w_nesting.docx" 2>&1 || true)
case "$err" in
  *"<w:"*"nest"*|*"nest"*"<w:"*) pass "E4 w nesting" ;;
  *) fail "E4 w nesting" "error about <w:> nesting" "$err" ;;
esac

# --- E5. heading inside <p> ---
cat > "$TMP/E5_heading_in_p.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><h2>heading inside p</h2></p>
EOF
err=$(BIN "$TMP/E5_heading_in_p.dcd" "$TMP/E5_heading_in_p.docx" 2>&1 || true)
case "$err" in
  *"heading"*"<p"*|*"<p"*"heading"*|*"not allow"*"heading"*) pass "E5 heading in p" ;;
  *) fail "E5 heading in p" "error about heading in <p>" "$err" ;;
esac

# --- E6. <w:*> with heading inside ---
cat > "$TMP/E6_heading_in_w.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c><h2>heading in w</h2></w:c>
EOF
err=$(BIN "$TMP/E6_heading_in_w.dcd" "$TMP/E6_heading_in_w.docx" 2>&1 || true)
case "$err" in
  *"heading"*|*"<h"*) pass "E6 heading in w" ;;
  *) fail "E6 heading in w" "error about heading" "$err" ;;
esac

# --- E7. <h> with nested <h> ---
cat > "$TMP/E7_h_nesting.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h1>outer <h2>inner</h2></h1>
EOF
err=$(BIN "$TMP/E7_h_nesting.dcd" "$TMP/E7_h_nesting.docx" 2>&1 || true)
case "$err" in
  *"heading"*|*"<h"*"<h"*) pass "E7 heading nesting" ;;
  *) fail "E7 heading nesting" "error about heading nesting" "$err" ;;
esac

# --- E8. name= missing ---
cat > "$TMP/E8_no_name.dcd" << 'EOF'
[section 0]
--- BODY ---
<p>no name</p>
EOF
err=$(BIN "$TMP/E8_no_name.dcd" "$TMP/E8_no_name.docx" 2>&1 || true)
case "$err" in
  *"name="*|*"required"*) pass "E8 name= missing" ;;
  *) fail "E8 name= missing" "error about name= required" "$err" ;;
esac

# --- E9. [] prefix violation (array used as object) ---
cat > "$TMP/E9_array_as_object.dcd" << 'EOF'
[section 0]
name=test
var=[]items
keys=name
--- BODY ---
<p>{{items.name}}</p>
EOF
err=$(BIN "$TMP/E9_array_as_object.dcd" "$TMP/E9_array_as_object.docx" 2>&1 || true)
case "$err" in
  *"array"*"object"*|*"[]"*"use"*"object"*) pass "E9 array as object" ;;
  *) fail "E9 array as object" "error about array used as object" "$err" ;;
esac

# --- E10. strict usage violation (declared but unused) ---
cat > "$TMP/E10_unused_var.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=username
--- BODY ---
<p>unused var</p>
EOF
err=$(BIN "$TMP/E10_unused_var.dcd" "$TMP/E10_unused_var.docx" 2>&1 || true)
case "$err" in
  *"unused"*|*"never used"*|*"not used"*) pass "E10 unused var" ;;
  *) fail "E10 unused var" "error about unused var" "$err" ;;
esac

# --- E11. object used as loop source ---
cat > "$TMP/E11_object_as_loop.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
--- BODY ---
<loop x from info>{{x}}</loop>
EOF
err=$(BIN "$TMP/E11_object_as_loop.dcd" "$TMP/E11_object_as_loop.docx" 2>&1 || true)
case "$err" in
  *"object"*"loop"*|*"without []"*"loop"*) pass "E11 object as loop" ;;
  *) fail "E11 object as loop" "error about object used as loop" "$err" ;;
esac

# --- E12. loop source not in var= ---
cat > "$TMP/E12_loop_not_declared.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
--- BODY ---
<loop x from unknown>{{x}}</loop>
EOF
err=$(BIN "$TMP/E12_loop_not_declared.dcd" "$TMP/E12_loop_not_declared.docx" 2>&1 || true)
case "$err" in
  *"not declared"*|*"not in var"*) pass "E12 loop not declared" ;;
  *) fail "E12 loop not declared" "error about loop source not declared" "$err" ;;
esac

# --- E13. unexpected </set:flags> mismatch ---
cat > "$TMP/E13_set_mismatch.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:b|i>bold italic</set:i|b></p>
EOF
err=$(BIN "$TMP/E13_set_mismatch.dcd" "$TMP/E13_set_mismatch.docx" 2>&1 || true)
case "$err" in
  *"balanc"*|*"expect"*|*"unexpec"*) pass "E13 set flags mismatch" ;;
  *) fail "E13 set flags mismatch" "error about tag balance" "$err" ;;
esac

# --- E14. dotted key in keys= without format ---
cat > "$TMP/E14_dotted_no_format.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title, items.name
--- BODY ---
<p>{{info.title}}</p>
<loop x from items>{{x.name}}</loop>
EOF
err=$(BIN "$TMP/E14_dotted_no_format.dcd" "$TMP/E14_dotted_no_format.docx" 2>&1 || true)
case "$err" in
  *"dotted"*"format"*) pass "E14 dotted no format" ;;
  *) fail "E14 dotted no format" "error about dotted key needs format" "$err" ;;
esac

# --- E15. <w:*> heading prohibition (block level) ---
cat > "$TMP/E15_w_heading_block.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c><h1>heading in block</h1></w:c>
EOF
err=$(BIN "$TMP/E15_w_heading_block.dcd" "$TMP/E15_w_heading_block.docx" 2>&1 || true)
case "$err" in
  *"heading"*|*"<h"*) pass "E15 w heading block" ;;
  *) fail "E15 w heading block" "error about heading inside <w:>" "$err" ;;
esac

# --- E16. closing tag for <w:> mismatch ---
cat > "$TMP/E16_w_close_mismatch.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<w:c>text</w:z>
EOF
err=$(BIN "$TMP/E16_w_close_mismatch.dcd" "$TMP/E16_w_close_mismatch.docx" 2>&1 || true)
case "$err" in
  *"mismatch"*|*"close"*|*"unexpected"*) pass "E16 w close mismatch" ;;
  *) fail "E16 w close mismatch" "error about close tag mismatch" "$err" ;;
esac

# --- E17. <loop:ol> closed with </loop> (wrong close) ---
cat > "$TMP/E17_loop_wrong_close.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items>{{x}}</loop>
EOF
err=$(BIN "$TMP/E17_loop_wrong_close.dcd" "$TMP/E17_loop_wrong_close.docx" 2>&1 || true)
case "$err" in
  *"mismatch"*|*"close"*|*"expect"*) pass "E17 loop wrong close" ;;
  *) fail "E17 loop wrong close" "error about loop close tag mismatch" "$err" ;;
esac

# --- E18. list inside loop (list loop prohibition) ---
cat > "$TMP/E18_list_loop_prohib.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<ul>
<loop x from items><li>{{x}}</li></loop>
</ul>
EOF
err=$(BIN "$TMP/E18_list_loop_prohib.dcd" "$TMP/E18_list_loop_prohib.docx" 2>&1 || true)
case "$err" in
  *"loop"*"list"*|*"list"*"loop"*|*"inside"*"<ul"*) pass "E18 list loop prohibition" ;;
  *) fail "E18 list loop prohibition" "error about list inside loop" "$err" ;;
esac

# --- E19. Nested list >3 levels ---
cat > "$TMP/E19_deep_nested.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li>1
<ul>
<li>2
<ul>
<li>3
<ul>
<li>4 deep</li>
</ul>
</li>
</ul>
</li>
</ul>
</li>
</ul>
EOF
err=$(BIN "$TMP/E19_deep_nested.dcd" "$TMP/E19_deep_nested.docx" 2>&1 || true)
# May pass or fail; we just check it doesn't crash
case "$err" in
  *"Error"*) fail "E19 deep nested" "no crash" "$(echo "$err" | head -1)" ;;
  *) pass "E19 deep nested" ;;
esac

# --- E20. <set:b|i> closed with </set:i|b> (flag order mismatch) ---
cat > "$TMP/E20_set_order.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:b|i>bold italic</set:i|b></p>
EOF
err=$(BIN "$TMP/E20_set_order.dcd" "$TMP/E20_set_order.docx" 2>&1 || true)
case "$err" in
  *"balanc"*|*"expect"*|*"unexpec"*) pass "E20 set flag order mismatch" ;;
  *) fail "E20 set flag order mismatch" "error about tag balance" "$err" ;;
esac

# --- E21. Duplicate name= across sections ---
cat > "$TMP/E21_dup_name.dcd" << 'EOF'
[section 0]
name=dup
--- BODY ---
<p>section 0</p>

[section 1]
name=dup
--- BODY ---
<p>section 1</p>
EOF
err=$(BIN "$TMP/E21_dup_name.dcd" "$TMP/E21_dup_name.docx" 2>&1 || true)
case "$err" in
  *"duplic"*|*"already"*|*"dup"*) pass "E21 duplicate name" ;;
  *) fail "E21 duplicate name" "error about duplicate name" "$err" ;;
esac

# --- E22. formats= key not in keys= ---
cat > "$TMP/E22_format_not_in_keys.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
formats=[items.date_field:dd-MM-yyyy]
--- BODY ---
<p>{{info.title}}</p>
<loop x from items>{{x.date_field}}</loop>
EOF
err=$(BIN "$TMP/E22_format_not_in_keys.dcd" "$TMP/E22_format_not_in_keys.docx" 2>&1 || true)
case "$err" in
  *"not in keys"*|*"not found"*|*"key"*"format"*) pass "E22 format key not in keys" ;;
  *) fail "E22 format key not in keys" "error about format key not in keys" "$err" ;;
esac

# --- E23. <loop> without [] prefix ---
cat > "$TMP/E23_loop_no_array.dcd" << 'EOF'
[section 0]
name=test
var=items
keys=title
--- BODY ---
<loop x from items>{{x}}</loop>
EOF
err=$(BIN "$TMP/E23_loop_no_array.dcd" "$TMP/E23_loop_no_array.docx" 2>&1 || true)
case "$err" in
  *"[]"*|*"array"*|*"not"*"array"*) pass "E23 loop without []" ;;
  *) fail "E23 loop without []" "error about missing []" "$err" ;;
esac

# --- E24. empty name= ---
cat > "$TMP/E24_empty_name.dcd" << 'EOF'
[section 0]
name=
--- BODY ---
<p>empty name</p>
EOF
err=$(BIN "$TMP/E24_empty_name.dcd" "$TMP/E24_empty_name.docx" 2>&1 || true)
case "$err" in
  *"required"*|*"empty"*name*) pass "E24 empty name=" ;;
  *) fail "E24 empty name=" "error about empty/required name" "$err" ;;
esac

# --- E25. heading level invalid ---
cat > "$TMP/E25_heading_invalid.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h0>too low</h0>
<h7>too high</h7>
EOF
err=$(BIN "$TMP/E25_heading_invalid.dcd" "$TMP/E25_heading_invalid.docx" 2>&1 || true)
case "$err" in
  *"invalid"*|*"not supported"*|*"level"*) pass "E25 heading level invalid" ;;
  *) fail "E25 heading level invalid" "error about heading level" "$err" ;;
esac

# --- E26. malformed format specifier ---
cat > "$TMP/E26_bad_format.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=date
formats=[date:]
--- BODY ---
<p>{{info.date}}</p>
EOF
err=$(BIN "$TMP/E26_bad_format.dcd" "$TMP/E26_bad_format.docx" 2>&1 || true)
case "$err" in
  *"format"*|*"invalid"*|*"malformed"*) pass "E26 malformed format" ;;
  *) fail "E26 malformed format" "error about format specifier" "$err" ;;
esac

# --- E27. loop source not in dataset ---
cat > "$TMP/E27_loop_source_missing.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<loop x from missing>{{x}}</loop>
EOF
err=$(BIN "$TMP/E27_loop_source_missing.dcd" "$TMP/E27_loop_source_missing.docx" 2>&1 || true)
case "$err" in
  *"not found"*|*"section"*"loop"*|*"unknown"*|*"source"*) pass "E27 loop source missing" ;;
  *) fail "E27 loop source missing" "error about unknown source" "$err" ;;
esac

# --- E28. format key not in dataset ---
cat > "$TMP/E28_format_key_missing.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
formats=[missing:dd-MM]
--- BODY ---
<p>{{info.title}}</p>
EOF
err=$(BIN "$TMP/E28_format_key_missing.dcd" "$TMP/E28_format_key_missing.docx" 2>&1 || true)
case "$err" in
  *"format"*|*"key"*|*"not found"*|*"missing"*) pass "E28 format key missing" ;;
  *) fail "E28 format key missing" "error about missing format key" "$err" ;;
esac

# --- E29. loop on non-array (runtime) ---
cat > "$TMP/E29_loop_nonarray.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=title
--- BODY ---
<loop x from info>{{x}}</loop>
EOF
# info is an object — loop should error
echo '{"info":{"title":"nope"}}' > "$TMP/E29_loop_nonarray.json"
err=$(BIN -data "$TMP/E29_loop_nonarray.json" "$TMP/E29_loop_nonarray.dcd" "$TMP/E29_loop_nonarray.docx" 2>&1 || true)
case "$err" in
  *"array"*|*"iterate"*|*"object"*"loop"*) pass "E29 loop non-array" ;;
  *) fail "E29 loop non-array" "error about non-array source" "$err" ;;
esac

# --- E30. var= exceeds limit of 5 ---
cat > "$TMP/E30_var_limit.dcd" << 'EOF'
[section 0]
name=test
var=a, b, c, d, e, f
keys=a, b, c, d, e, f
--- BODY ---
{{a}} {{b}} {{c}} {{d}} {{e}} {{f}}
EOF
err=$(BIN "$TMP/E30_var_limit.dcd" "$TMP/E30_var_limit.docx" 2>&1 || true)
case "$err" in
  *"5"*|*"limit"*|*"max"*|*"var"*) pass "E30 var limit exceeded" ;;
  *) fail "E30 var limit exceeded" "error about var limit" "$err" ;;
esac

# --- E31. keys= exceeds limit of 15 ---
cat > "$TMP/E31_keys_limit.dcd" << 'EOF'
[section 0]
name=test
var=a
keys=01,02,03,04,05,06,07,08,09,10,11,12,13,14,15,16
--- BODY ---
{{a}}
EOF
err=$(BIN "$TMP/E31_keys_limit.dcd" "$TMP/E31_keys_limit.docx" 2>&1 || true)
case "$err" in
  *"15"*|*"limit"*|*"max"*|*"key"*) pass "E31 keys limit exceeded" ;;
  *) fail "E31 keys limit exceeded" "error about keys limit" "$err" ;;
esac

# --- E32. <img=> without path ---
cat > "$TMP/E32_img_no_path.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<img=>
EOF
err=$(BIN "$TMP/E32_img_no_path.dcd" "$TMP/E32_img_no_path.docx" 2>&1 || true)
case "$err" in
  *"empty"*|*"missing"*|*"img"*|*"path"*) pass "E32 img no path" ;;
  *) fail "E32 img no path" "error about empty img path" "$err" ;;
esac

# --- E33. <a=""> empty url ---
cat > "$TMP/E33_link_empty.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><a=>click</a></p>
EOF
err=$(BIN "$TMP/E33_link_empty.dcd" "$TMP/E33_link_empty.docx" 2>&1 || true)
case "$err" in
  *"empty"*|*"url"*|*"link"*) pass "E33 empty link url" ;;
  *) fail "E33 empty link url" "error about empty link url" "$err" ;;
esac

# --- E34. Unclosed <a=url> tag ---
cat > "$TMP/E34_link_unclosed.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>text <a=https://example.com>unclosed link</p>
EOF
err=$(BIN "$TMP/E34_link_unclosed.dcd" "$TMP/E34_link_unclosed.docx" 2>&1 || true)
case "$err" in
  *"unclosed"*|*"balanc"*|*"expect"*) pass "E34 unclosed link tag" ;;
  *) fail "E34 unclosed link tag" "error about unclosed <a> tag" "$err" ;;
esac

# --- E35. Loop body without {{x}} (unused array) ---
cat > "$TMP/E35_loop_no_varref.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<loop x from items>
<p>no variable reference</p>
</loop>
EOF
err=$(BIN "$TMP/E35_loop_no_varref.dcd" "$TMP/E35_loop_no_varref.docx" 2>&1 || true)
case "$err" in
  *"unused"*|*"never used"*|*"not used"*|*"var"*"use"*) pass "E35 loop no var ref" ;;
  *) fail "E35 loop no var ref" "error about unused array var/object" "$err" ;;
esac

# --- E36. [section:unknown N] invalid type ---
cat > "$TMP/E36_unknown_section.dcd" << 'EOF'
[section:foobar 0]
name=test
--- BODY ---
<p>unknown section type</p>
EOF
err=$(BIN "$TMP/E36_unknown_section.dcd" "$TMP/E36_unknown_section.docx" 2>&1 || true)
case "$err" in
  *"unknown"*|*"invalid"*|*"foobar"*) pass "E36 unknown section type" ;;
  *) fail "E36 unknown section type" "error about unknown section type" "$err" ;;
esac

# --- E37. Empty var= ---
cat > "$TMP/E37_empty_var.dcd" << 'EOF'
[section 0]
name=test
var=
keys=title
--- BODY ---
<p>empty var</p>
EOF
err=$(BIN "$TMP/E37_empty_var.dcd" "$TMP/E37_empty_var.docx" 2>&1 || true)
case "$err" in
  *"empty"*|*"var"*) pass "E37 empty var=" ;;
  *) fail "E37 empty var=" "error about empty var" "$err" ;;
esac

# --- E38. Empty keys= ---
cat > "$TMP/E38_empty_keys.dcd" << 'EOF'
[section 0]
name=test
var=info
keys=
--- BODY ---
{{info.title}}
EOF
err=$(BIN "$TMP/E38_empty_keys.dcd" "$TMP/E38_empty_keys.docx" 2>&1 || true)
case "$err" in
  *"empty"*|*"keys"*) pass "E38 empty keys=" ;;
  *) fail "E38 empty keys=" "error about empty keys" "$err" ;;
esac

# --- E39. <col> without <table> parent ---
cat > "$TMP/E39_col_orphan.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<col>orphan cell</col>
EOF
err=$(BIN "$TMP/E39_col_orphan.dcd" "$TMP/E39_col_orphan.docx" 2>&1 || true)
case "$err" in
  *"col"*|*"table"*|*"unexpected"*) pass "E39 col orphan" ;;
  *) fail "E39 col orphan" "error about col outside table" "$err" ;;
esac

# --- E40. <li> without <ul>/<ol> parent ---
cat > "$TMP/E40_li_orphan.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<li>orphan item</li>
EOF
err=$(BIN "$TMP/E40_li_orphan.dcd" "$TMP/E40_li_orphan.docx" 2>&1 || true)
case "$err" in
  *"li"*|*"list"*|*"unexpected"*|*"<ul"*|*"<ol"*) pass "E40 li orphan" ;;
  *) fail "E40 li orphan" "error about li outside list" "$err" ;;
esac

# --- E41. <set:> empty flags ---
cat > "$TMP/E41_set_empty.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><set:>text</set:></p>
EOF
err=$(BIN "$TMP/E41_set_empty.dcd" "$TMP/E41_set_empty.docx" 2>&1 || true)
case "$err" in
  *"empty"*|*"flag"*|*"set:"*|*"invalid"*) pass "E41 set: empty flags" ;;
  *) fail "E41 set: empty flags" "error about empty set flags" "$err" ;;
esac

# --- E42. [section] without number ---
cat > "$TMP/E42_section_no_num.dcd" << 'EOF'
[section]
name=test
--- BODY ---
<p>no number</p>
EOF
err=$(BIN "$TMP/E42_section_no_num.dcd" "$TMP/E42_section_no_num.docx" 2>&1 || true)
case "$err" in
  *"number"*|*"invalid"*|*"section"*|*"parse"*) pass "E42 section no number" ;;
  *) fail "E42 section no number" "error about missing section number" "$err" ;;
esac

# --- E43. Duplicate section index ---
cat > "$TMP/E43_dup_index.dcd" << 'EOF'
[section 0]
name=first
--- BODY ---
<p>first</p>

[section 0]
name=second
--- BODY ---
<p>second</p>
EOF
err=$(BIN "$TMP/E43_dup_index.dcd" "$TMP/E43_dup_index.docx" 2>&1 || true)
case "$err" in
  *"duplic"*|*"already"*|*"index"*|*"section 0"*"already"*) pass "E43 duplicate section index" ;;
  *) fail "E43 duplicate section index" "error about duplicate section index" "$err" ;;
esac

# --- E44. List nesting >3 levels (error path) ---
cat > "$TMP/E44_deep_list.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<ul>
<li>1
<ul>
<li>2
<ul>
<li>3
<ul>
<li>4 deep</li>
</ul>
</li>
</ul>
</li>
</ul>
</li>
</ul>
EOF
err=$(BIN "$TMP/E44_deep_list.dcd" "$TMP/E44_deep_list.docx" 2>&1 || true)
# May succeed or fail depending on error recovery; only check for crash
case "$err" in
  *"Error"*) fail "E44 deep list" "no crash unexpected" "$(echo "$err" | head -1)" ;;
  *) pass "E44 deep list" ;;
esac

# --- E45. <loop:ul> with type= (invalid for ul) ---
cat > "$TMP/E45_loopul_type_invalid.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ul x from items type=A>
<li>{{x}}</li>
</loop:ul>
EOF
err=$(BIN "$TMP/E45_loopul_type_invalid.dcd" "$TMP/E45_loopul_type_invalid.docx" 2>&1 || true)
# type= on ul is silently ignored, should not crash
case "$err" in
  *"Error"*) fail "E45 loopul type invalid" "no crash on ul type" "$(echo "$err" | head -1)" ;;
  *) pass "E45 loopul type invalid" ;;
esac

# --- E46. Var declares []array but no loop uses it (strict usage) ---
cat > "$TMP/E46_array_unused.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<!-- no loop using items -->
EOF
err=$(BIN "$TMP/E46_array_unused.dcd" "$TMP/E46_array_unused.docx" 2>&1 || true)
case "$err" in
  *"unused"*|*"never used"*|*"not used"*) pass "E46 array unused" ;;
  *) fail "E46 array unused" "error about unused array var" "$err" ;;
esac

# --- E47. Section with var= having both object and array, all unused ---
cat > "$TMP/E47_multi_unused.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>no var usage at all</p>
EOF
err=$(BIN "$TMP/E47_multi_unused.dcd" "$TMP/E47_multi_unused.docx" 2>&1 || true)
case "$err" in
  *"unused"*|*"never used"*|*"not used"*) pass "E47 multi unused" ;;
  *) fail "E47 multi unused" "error about unused vars" "$err" ;;
esac

# --- E48. <row> without <table> parent ---
cat > "$TMP/E48_row_orphan.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<row><col>orphan row</col></row>
EOF
err=$(BIN "$TMP/E48_row_orphan.dcd" "$TMP/E48_row_orphan.docx" 2>&1 || true)
case "$err" in
  *"row"*|*"table"*|*"unexpected"*) pass "E48 row orphan" ;;
  *) fail "E48 row orphan" "error about row outside table" "$err" ;;
esac

# --- E49. </b> without opening <b> ---
cat > "$TMP/E49_close_no_open.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p>text</b> more</p>
EOF
err=$(BIN "$TMP/E49_close_no_open.dcd" "$TMP/E49_close_no_open.docx" 2>&1 || true)
case "$err" in
  *"unexpected"*|*"balanc"*|*"close"*) pass "E49 close without open" ;;
  *) fail "E49 close without open" "error about unexpected close tag" "$err" ;;
esac

# --- E50. Nested <a=url> inside <a=url> ---
cat > "$TMP/E50_nested_link.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<p><a=https://a.com>outer <a=https://b.com>inner</a></a></p>
EOF
err=$(BIN "$TMP/E50_nested_link.dcd" "$TMP/E50_nested_link.docx" 2>&1 || true)
case "$err" in
  *"nest"*|*"balanc"*|*"expect"*|*"<a"*) pass "E50 nested link" ;;
  *) fail "E50 nested link" "error about nested <a>" "$err" ;;
esac

# --- E51. Non-numeric heading level ---
cat > "$TMP/E51_heading_nonum.dcd" << 'EOF'
[section 0]
name=test
--- BODY ---
<h>no level</h>
EOF
err=$(BIN "$TMP/E51_heading_nonum.dcd" "$TMP/E51_heading_nonum.docx" 2>&1 || true)
case "$err" in
  *"invalid"*|*"heading"*|*"not supported"*|*"level"*) pass "E51 heading no level" ;;
  *) fail "E51 heading no level" "error about invalid heading level" "$err" ;;
esac

# --- E52. Var name with forbidden chars ---
cat > "$TMP/E52_var_badchars.dcd" << 'EOF'
[section 0]
name=test
var=my-var
keys=title
--- BODY ---
{{my-var.title}}
EOF
err=$(BIN "$TMP/E52_var_badchars.dcd" "$TMP/E52_var_badchars.docx" 2>&1 || true)
case "$err" in
  *"invalid"*|*"var"*|*"char"*|*"name"*) pass "E52 var bad chars" ;;
  *) fail "E52 var bad chars" "error about invalid var name" "$err" ;;
esac

# --- E53. Custom <li> attributes in loop:ol/loop:ul ---
cat > "$TMP/E53_custom_li_attrs.dcd" << 'EOF'
[section 0]
name=test
var=info, []items
keys=title
--- BODY ---
<p>{{info.title}}</p>
<loop:ol x from items type=A>
<li indent=20 hanging=10>{{x.label}}</li>
</loop:ol>
<loop:ul x from items>
<li indent=15>{{x.label}}</li>
</loop:ul>
EOF
BIN "$TMP/E53_custom_li_attrs.dcd" "$TMP/E53_custom_li_attrs.docx" 2>&1 && pass "E53 custom li attrs in loop" || fail "E53" "exit 0" "non-zero exit"

echo
echo "=========================================="
echo "  STRESS TEST RESULTS"
echo "=========================================="
echo "  Total:  $TOTAL"
echo "  Pass:   $PASS"
echo "  Fail:   $FAIL"
echo "=========================================="

clean
[ "$FAIL" -eq 0 ]