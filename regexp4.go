package regexp4

const inf = 1073741824 // 2^30

const (
  modAlpha      uint8 = 1
  modOmega      uint8 = 2
  modLonley     uint8 = 4
  modFwrByChar  uint8 = 8
  modCommunism  uint8 = 16
  modNegative   uint8 = 128
  modPositive   uint8 = ^modNegative
  modCapitalism uint8 = ^modCommunism
)

const (
  asmPath = iota; asmPathEle; asmPathEnd;
  asmGroup; asmGroupEnd; asmHook; asmHookEnd; asmSet; asmSetEnd;
  asmBackref; asmMeta; asmRangeab; asmUTF8; asmPoint; asmSimple; asmEnd
)

type reStruct struct {
  str                string
  reType             uint8
  mods               uint8
  loopsMin, loopsMax int
}

type catchInfo struct { init, end, id int }

type raptorASM struct {
  re    reStruct
  inst  uint8
  close int
}

type RE struct {
  txt, re      string
  compile      bool
  result       int

  txtInit      int
  txtLen       int
  txtPos       int

  catches      []catchInfo
  catchIndex   int
  catchIdIndex int

  asm          []raptorASM
  mods         uint8
}

func (r *RE) Compile( re string ){
  r.catchIndex = 1
  r.compile    = false
  if len(re) == 0 { return }

  rexp := reStruct{ str: re, reType: asmPath }
  r.re  = re
  r.asm = make( []raptorASM, 0, 32 )

  getMods( &rexp, &rexp )
  r.mods = rexp.mods

  r.genPaths( rexp )
  // if isPath( &rexp ) { r.genPaths( rexp )
  // } else { r.genTracks( &rexp ) }

  r.asm = append( r.asm, raptorASM{ inst: asmEnd, close: len(r.asm) } )
  r.compile = true
}

func isPath( rexp *reStruct ) bool {
  if len(rexp.str) == 0 { return false }

  for i, deep := 0, 0; walkMeta( rexp.str[i:], &i ) < len( rexp.str ); i++ {
    switch rexp.str[ i ] {
    case '(', '<': deep++
    case ')', '>': deep--
    case '[': i += walkSet( rexp.str[i:] )
    case '|': if deep == 0 { return true }
    }
  }

  return false
}

func (r *RE) genPaths( rexp reStruct ){
  var track reStruct
  pathIndex := len( r.asm )
  r.asm = append( r.asm, raptorASM{ inst: asmPath, re: rexp } )

  for cutByType( &rexp, &track, asmPath ) {
    trackIndex := len( r.asm )
    r.asm = append( r.asm, raptorASM{ inst: asmPathEle, re: track } )
    r.genTracks( &track )
    r.asm[trackIndex].close = len( r.asm )
  }

  r.asm[pathIndex].close = len( r.asm )
  r.asm = append( r.asm, raptorASM{ inst: asmPathEnd, close: len(r.asm) } )
}

func (r *RE) genTracks( rexp *reStruct ){
  var track reStruct
  for tracker( rexp, &track ) {
    trackIndex := len( r.asm )
    switch track.reType {
    case asmHook   :
      r.asm = append( r.asm, raptorASM{ inst: asmHook, re: track } )

      r.genPaths ( track )
      // if isPath( &track ) { r.genPaths ( track )
      // } else              { r.genTracks( &track ) }

      r.asm[trackIndex].close = len( r.asm )
      r.asm = append( r.asm, raptorASM{ inst: asmHookEnd, close: len(r.asm) } )
    case asmGroup  :
      r.asm = append( r.asm, raptorASM{ inst: asmGroup, re: track } )

      r.genPaths ( track )
      // if isPath( &track ) { r.genPaths ( track )
      // } else              { r.genTracks( &track ) }

      r.asm[trackIndex].close = len( r.asm )
      r.asm = append( r.asm, raptorASM{ inst: asmGroupEnd, close: len(r.asm) } )
    case asmPath   :
    case asmSet    : r.genSet( &track )
    case asmBackref: r.asm = append( r.asm, raptorASM{ inst: asmBackref, close: trackIndex, re: track } )
    case asmMeta   : r.asm = append( r.asm, raptorASM{ inst: asmMeta   , close: trackIndex, re: track } )
    case asmRangeab: r.asm = append( r.asm, raptorASM{ inst: asmRangeab, close: trackIndex, re: track } )
    case asmUTF8   : r.asm = append( r.asm, raptorASM{ inst: asmUTF8   , close: trackIndex, re: track } )
    case asmPoint  : r.asm = append( r.asm, raptorASM{ inst: asmPoint  , close: trackIndex, re: track } )
    default        : r.asm = append( r.asm, raptorASM{ inst: asmSimple , close: trackIndex, re: track } )
    }
  }
}

func (r *RE) genSet( rexp *reStruct ){
  if len(rexp.str) == 0 { return }

  if rexp.str[0] == '^' {
    rexp.str = rexp.str[1:]
    if rexp.mods & modNegative > 0 { rexp.mods &= modPositive
    } else                         { rexp.mods |= modNegative }
  }

  setIndex := len( r.asm )
  r.asm = append( r.asm, raptorASM{ inst: asmSet, re: *rexp } )

  var track reStruct
  for trackerSet( rexp, &track ) {
    switch track.reType {
    case asmMeta   : r.asm = append( r.asm, raptorASM{ inst: asmMeta   , close: len(r.asm), re: track } )
    case asmRangeab: r.asm = append( r.asm, raptorASM{ inst: asmRangeab, close: len(r.asm), re: track } )
    case asmUTF8   : r.asm = append( r.asm, raptorASM{ inst: asmUTF8   , close: len(r.asm), re: track } )
    default        : r.asm = append( r.asm, raptorASM{ inst: asmSimple , close: len(r.asm), re: track } )
    }
  }

  r.asm[ setIndex ].close = len( r.asm )
  r.asm = append( r.asm, raptorASM{ inst: asmSetEnd, close: len(r.asm) } )
}

func trackerSet( rexp, track *reStruct ) bool {
  if len( rexp.str ) == 0 { return false }

  if rexp.str[0] > 127 {
    cutByLen( rexp, track, utf8meter( rexp.str ), asmUTF8 )
  } else if rexp.str[0] == ':' {
    cutByLen ( rexp, track, 2, asmMeta  )
  } else {
    for i := 0; i < len( rexp.str ); i++ {
      if rexp.str[i] > 127 {
        cutByLen( rexp, track, i, asmSimple  ); goto setLM;
      } else {
        switch rexp.str[i] {
        case ':': cutByLen( rexp, track, i, asmSimple  ); goto setLM;
        case '-':
          if i == 1 { cutByLen( rexp, track,     3, asmRangeab )
          } else    { cutByLen( rexp, track, i - 1, asmSimple  ) }

          goto setLM;
        }
      }
    }

    cutByLen( rexp, track, len( rexp.str ), asmSimple  );
  }

 setLM:
  track.loopsMin, track.loopsMax = 1, 1
  track.mods &= modPositive
  return true
}

func tracker( rexp, track *reStruct ) bool {
  if len( rexp.str ) == 0 { return false }

  if rexp.str[0] > 127 {
    cutByLen( rexp, track, utf8meter( rexp.str ), asmUTF8 )
  } else {
    switch rexp.str[0] {
    case ':': cutByLen ( rexp, track, 2,     asmMeta    )
    case '.': cutByLen ( rexp, track, 1,     asmPoint   )
    case '@': cutByLen ( rexp, track, 1 +
            countCharDigits( rexp.str[1:] ), asmBackref )
    case '(': cutByType( rexp, track,        asmGroup   )
    case '<': cutByType( rexp, track,        asmHook    )
    case '[': cutByType( rexp, track,        asmSet     )
    default : cutSimple( rexp, track                    )
    }
  }

  getLoops( rexp, track );
  getMods ( rexp, track );
  return true
}

func cutSimple( rexp, track *reStruct ){
  for i, c := range rexp.str {
    if c > 127 {
      cutByLen( rexp, track, i, asmSimple  ); return
    } else {
      switch c {
      case '(', '<', '[', '@', ':', '.':
        cutByLen( rexp, track, i, asmSimple  ); return
      case '?', '+', '*', '{', '#':
        if i == 1 { cutByLen( rexp, track,     1, asmSimple  )
        } else    { cutByLen( rexp, track, i - 1, asmSimple  ) }
        return
      }
    }
  }

  cutByLen( rexp, track, len(rexp.str), asmSimple  );
}

func cutByLen( rexp, track *reStruct, length int, reType uint8 ){
  *track       = *rexp
  track.str    = rexp.str[:length]
  rexp.str     = rexp.str[length:]
  track.reType = reType;
}

func cutByType( rexp, track *reStruct, reType uint8 ) bool {
  if len(rexp.str) == 0 { return false }

  *track       = *rexp
  track.reType = reType
  for i , deep, cut := 0, 0, false; walkMeta( rexp.str[i:], &i ) < len( rexp.str ); i++ {
    switch rexp.str[ i ] {
    case '(', '<': deep++
    case ')', '>': deep--
    case '[': i += walkSet( rexp.str[i:] )
    }

    switch reType {
    case asmHook, asmGroup: cut = deep == 0
    case asmSet          : cut = rexp.str[ i ] == ']'
    case asmPath         : cut = rexp.str[ i ] == '|' && deep == 0
    }

    if cut {
      track.str  = rexp.str[:i]
      rexp.str   = rexp.str[i + 1:]
      if reType != asmPath { track.str = track.str[1:] }
      return true
    }
  }

  rexp.str = ""
  return true
}

func walkSet( str string ) int {
  for i := 0; walkMeta( str[i:], &i ) < len( str ); i++ {
    if str[i] == ']' { return i }
  }

  return len(str);
}

func walkMeta( str string, n *int ) int {
  for i := 0; i < len( str ); i += 2 {
    if str[i] != ':' { *n += i; return *n }
  }

  *n += len( str )
  return *n
}

func getMods( rexp, track *reStruct ){
  track.mods &= modPositive

  if len( rexp.str ) > 0 && rexp.str[ 0 ] == '#' {
    for i, c := range rexp.str[1:] {
      switch c {
      case '^': track.mods |= modAlpha
      case '$': track.mods |= modOmega
      case '?': track.mods |= modLonley
      case '~': track.mods |= modFwrByChar
      case '*': track.mods |= modCommunism
      case '/': track.mods &= modCapitalism
      case '!': track.mods |= modNegative
      default : rexp.str    = rexp.str[i+1:]; return
      }
    }

    rexp.str = ""
  }
}

func getLoops( rexp, track *reStruct ){
  pos := 0;
  track.loopsMin, track.loopsMax = 1, 1

  if len( rexp.str ) > 0 {
    switch rexp.str[0] {
    case '?' : pos = 1; track.loopsMin = 0; track.loopsMax =   1;
    case '+' : pos = 1; track.loopsMin = 1; track.loopsMax = inf;
    case '*' : pos = 1; track.loopsMin = 0; track.loopsMax = inf;
    case '{' : pos = 1
      track.loopsMin = aToi( rexp.str[pos:] )
      pos += countCharDigits( rexp.str[pos:] )

      if rexp.str[pos] == '}' {
        track.loopsMax = track.loopsMin;
        pos += 1
      } else if rexp.str[pos:pos+2] == ",}" {
        pos += 2
        track.loopsMax = inf
      } else if rexp.str[pos] == ',' {
        pos += 1
        track.loopsMax = aToi( rexp.str[pos:] )
        pos += countCharDigits( rexp.str[pos:] ) + 1
      }
    }

    rexp.str = rexp.str[pos:]
  }
}

//-! match

func (r *RE) Find( txt, re string ) bool {
  return r.Match( txt, re ) > 0
}

func (r *RE) Match( txt, re string ) int {
  r.Compile( re )
  return r.MatchString( txt )
}

func (r *RE) FindString( txt string ) bool {
  return r.MatchString( txt ) > 0
}

func (r *RE) MatchString( txt string ) int {
  loops       := len(txt)
  r.txt        = txt
  r.result     = 0
  r.catches    = make( []catchInfo, 32 )
  r.catchIndex = 1
  if len(txt) == 0  || !r.compile { return 0 }

  if (r.mods & modAlpha) > 0 { loops = 1 }

  for forward, i, ocindex := 0, 0, 0; i < loops; i += forward {
    forward, r.catchIdIndex       = utf8meter( txt[i:] ), 1
    r.txtPos, r.txtInit, r.txtLen = 0, i, len( txt[i:] )
    ocindex                       = r.catchIndex

    if r.trekking( 0 ) {
      if (r.mods & modOmega) > 0 {
        if r.txtPos == r.txtLen                              { r.result = 1; return 1
        } else { r.catchIndex = 1 }
      } else if (r.mods & modLonley   ) > 0                  { r.result = 1; return 1
      } else if (r.mods & modFwrByChar) > 0 || r.txtPos == 0 { r.result++
      } else {   forward = r.txtPos;                           r.result++; }
    } else { r.catchIndex = ocindex }
  }

  return r.result
}

func (r *RE) trekking( index int ) (result bool) {
  for ; r.asm[ index ].inst != asmEnd; index = r.asm[ index ].close + 1 {
    switch r.asm[ index ].inst {
    case asmPathEnd, asmPathEle, asmGroupEnd, asmHookEnd, asmSetEnd: return true
    case asmHook :
      iCatch := r.openCatch();
      result  = r.loopGroup( index )
      if result { r.closeCatch( iCatch ) }
    case asmGroup: result = r.loopGroup( index )
    case asmPath : result = r.walker   ( index )
    default      : result = r.looper   ( index )
    }

    if !result { return false }
  }

  return true
}

func (r *RE) walker( index int ) bool {
  index++
  for oTextPos, oCatchIndex, oCatchIdIndex := r.txtPos, r.catchIndex, r.catchIdIndex;
      r.asm[ index ].inst == asmPathEle
      index, r.txtPos, r.catchIndex, r.catchIdIndex = r.asm[ index ].close, oTextPos, oCatchIndex, oCatchIdIndex {
    if r.trekking( index + 1 ) { return true }
  }

  return false
}

func (r *RE) looper( index int ) bool {
  loops := 0

  if (r.asm[ index ].re.mods & modNegative) > 0 {
    for forward := 0; loops < r.asm[ index ].re.loopsMax && r.txtPos < r.txtLen && !r.match( index, r.txt[r.txtInit + r.txtPos:], &forward ); {
      r.txtPos += utf8meter( r.txt[r.txtInit + r.txtPos:] )
      loops++;
    }
  } else {
    for forward := 0; loops < r.asm[ index ].re.loopsMax && r.txtPos < r.txtLen &&  r.match( index, r.txt[r.txtInit + r.txtPos:], &forward ); {
      r.txtPos += forward
      loops++;
    }
  }

  if loops < r.asm[ index ].re.loopsMin { return false }
  return true
}

func (r *RE) loopGroup( index int ) bool {
  loops, textxtPos := 0, r.txtPos;

  if (r.asm[ index ].re.mods & modNegative) > 0 {
    for loops < r.asm[ index ].re.loopsMax && !r.trekking( index + 1 ) {
      textxtPos++;
      r.txtPos = textxtPos;
      loops++;
    }

    r.txtPos = textxtPos;
  } else {
    for loops < r.asm[ index ].re.loopsMax && r.trekking( index + 1 ) {
      loops++;
    }
  }

  if loops < r.asm[ index ].re.loopsMin { return false  }
  return true
}

func (r *RE) match( index int, txt string, forward *int ) bool {
  switch r.asm[ index ].inst {
  case asmPoint  : *forward = utf8meter( txt );  return true
  case asmSet    : return r.matchSet    ( index, txt, forward )
  case asmBackref: return r.matchBackRef( &r.asm[ index ].re, txt, forward )
  case asmRangeab: return matchRange    ( &r.asm[ index ].re, txt, forward )
  case asmMeta   : return matchMeta     ( &r.asm[ index ].re, txt, forward )
  default        : return matchText     ( &r.asm[ index ].re, txt, forward )
  }
}

func matchText( rexp *reStruct, txt string, forward *int ) bool {
  *forward = len(rexp.str)

  if len(txt) < *forward { return false }

  if (rexp.mods & modCommunism) > 0 {
    return strnEqlCommunist( txt, rexp.str, *forward )
  }

  return txt[:*forward] == rexp.str
}

func matchRange( rexp *reStruct, txt string, forward *int ) bool {
  *forward = 1
  if (rexp.mods & modCommunism) > 0 {
    chr := toLower( rune(txt[0]) )
    return chr >= toLower( rune(rexp.str[ 0 ]) ) && chr <= toLower( rune(rexp.str[ 2 ]) )
  }

  return txt[0] >= rexp.str[ 0 ] && txt[0] <= rexp.str[ 2 ];
}

func matchMeta( rexp *reStruct, txt string, forward *int ) bool {
  var f func( r rune ) bool
  *forward = 1

  switch rexp.str[1] {
  case 'a' : return isAlpha( rune(txt[0]) )
  case 'A' : f = isAlpha
  case 'd' : return isDigit( rune(txt[0]) )
  case 'D' : f = isDigit
  case 'w' : return isAlnum( rune(txt[0]) )
  case 'W' : f = isAlnum
  case 's' : return isSpace( rune(txt[0]) )
  case 'S' : f = isSpace
  case 'b' : return isBlank( rune(txt[0]) )
  case 'B' : f = isBlank
  case '&' : if txt[0] < 128 { return false }
    *forward = utf8meter( txt )
    return true
  default  : return txt[0] == rexp.str[1]
  }

  if f( rune(txt[0]) ) { return false }
  *forward = utf8meter( txt )
  return true
}

func (r *RE) matchSet( index int, txt string, forward *int ) (result bool) {
  *forward = 1

  for index++; !result && r.asm[ index ].inst != asmSetEnd; index++ {
    switch r.asm[ index ].inst {
    case asmRangeab, asmUTF8, asmMeta:
      result = r.match( index, txt, forward )
    default:
      if (r.asm[ index ].re.mods & modCommunism)  > 0 {
        result = findRuneCommunist( r.asm[ index ].re.str, rune( txt[ 0 ] ) )
      } else {
        result = strnchr( r.asm[ index ].re.str, rune( txt[ 0 ] ) )
      }
    }

    if result { return true }
  }

  return false
}

func (r *RE) matchBackRef( rexp *reStruct, txt string, forward *int ) bool {
  backRefId    := aToi( rexp.str[1:] )
  backRefIndex := r.lastIdCatch( backRefId )
  strCatch     := r.GetCatch( backRefIndex )
  *forward      = len(strCatch)

  if strCatch == "" || len( txt ) < *forward || strCatch != txt[:*forward] { return false }

  return true
}

func (r *RE) lastIdCatch( id int ) int {
  for index := r.catchIndex - 1; index > 0; index-- {
    if r.catches[ index ].id == id { return index }
  }

  return len(r.catches);
}

func (r *RE) openCatch() (index int) {
  index = r.catchIndex

  if r.catchIndex < len(r.catches) {
    r.catches[index] = catchInfo{ r.txtInit + r.txtPos, r.txtInit + r.txtPos, r.catchIdIndex }
  } else {
    r.catches = append( r.catches, catchInfo{ r.txtInit + r.txtPos, r.txtInit + r.txtPos, r.catchIdIndex } )
  }

  r.catchIndex++
  r.catchIdIndex++
  return
}

func (r *RE) closeCatch( index int ){
  if index < r.catchIndex {
    r.catches[index].end = r.txtInit + r.txtPos
  }
}

func (r *RE) Result  () int { return r.result }

func (r *RE) TotCatch() int { return r.catchIndex - 1 }

func (r *RE) GetCatch( index int ) string {
  if index < 1 || index >= r.catchIndex { return "" }
  return r.txt[ r.catches[index].init : r.catches[index].end ]
}

func (r *RE) GpsCatch( index int ) int {
  if index < 1 || index >= r.catchIndex { return 0 }
  return r.catches[index].init
}

func (r *RE) LenCatch( index int ) int {
  if index < 1 || index >= r.catchIndex { return 0 }
  return r.catches[index].end - r.catches[index].init
}

func (r *RE) RplCatch( rplStr string, id int ) string {
  last, rpls, catchLens := 0, 0, 0
  for index := 1; index < r.catchIndex; index++ {
    if r.catches[index].id == id {
      rpls++
      catchLens += r.catches[index].end - r.catches[index].init
    }
  }

  if rpls == 0 { return r.txt }
  if (r.mods & modFwrByChar) > 0 { catchLens = 0 }

  result, gps := make( []byte, len( r.txt ) - catchLens + rpls * len( rplStr ) ), 0

  for index := 1; index < r.catchIndex; index++ {
    if r.catches[index].id == id {
      if last > r.catches[index].init { last = r.catches[index].init } // modFwrByChar

      gps += copy( result[gps:], r.txt[last:r.catches[index].init] )
      gps += copy( result[gps:], rplStr )
      last = r.catches[index].end
    }
  }

  if last < len(r.txt) { gps += copy( result[gps:], r.txt[last:] ) }

  return string( result[:gps] )
}

func (r *RE) PutCatch( pStr string ) (result string) {
  for i := 0; i < len(pStr); {
    if pStr[i] == '#' {
      i++
      if len(pStr[i:]) > 0 && pStr[i] == '#' {
        i++
        result += "#"
      } else {
        result += r.GetCatch( aToi( pStr[i:] ) )
        i      += countCharDigits ( pStr[i:] )
      }
    } else { result += pStr[i:i+1]; i++ }
  }

  return
}

func (r *RE) Copy() *RE {
  nre := RE{ txt: r.txt, re: r.re, compile: r.compile, result: r.result, catchIndex: r.catchIndex, mods: r.mods }
  nre.catches = make( []catchInfo, r.catchIndex )
  copy( nre.catches, r.catches )
  nre.asm     = make( []raptorASM, len( r.asm ) )
  copy( nre.asm, r.asm )

  return &nre
}

func Compile( re string ) *RE {
  nre :=  new( RE )
  nre.Compile( re )
  return nre
}
