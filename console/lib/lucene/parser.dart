import './ast.dart';

// Parse a Lucene query string into a Query AST.
Query parse(String input) => _Parser(_Lexer(input).tokenize()).parse();

// ─── Lexer ────────────────────────────────────────────────────────────────────

enum _TK {
  word,
  phrase,
  colon,
  lbracket,
  rbracket,
  lbrace,
  rbrace,
  lparen,
  rparen,
  gt,
  gte,
  lt,
  lte,
  plus,
  minus,
  star,
  eof,
}

class _Token {
  final _TK kind;
  final String value;
  const _Token(this.kind, this.value);
}

class _Lexer {
  final String _src;
  int _i = 0;

  _Lexer(this._src);

  bool get _done => _i >= _src.length;
  String get _c => _done ? '' : _src[_i];

  void _ws() {
    while (!_done && _src[_i].trim().isEmpty) _i++;
  }

  List<_Token> tokenize() {
    final out = <_Token>[];
    while (true) {
      _ws();
      if (_done) {
        out.add(const _Token(_TK.eof, ''));
        break;
      }
      switch (_c) {
        case ':':
          out.add(const _Token(_TK.colon, ':'));
          _i++;
        case '[':
          out.add(const _Token(_TK.lbracket, '['));
          _i++;
        case ']':
          out.add(const _Token(_TK.rbracket, ']'));
          _i++;
        case '{':
          out.add(const _Token(_TK.lbrace, '{'));
          _i++;
        case '}':
          out.add(const _Token(_TK.rbrace, '}'));
          _i++;
        case '(':
          out.add(const _Token(_TK.lparen, '('));
          _i++;
        case ')':
          out.add(const _Token(_TK.rparen, ')'));
          _i++;
        case '+':
          out.add(const _Token(_TK.plus, '+'));
          _i++;
        case '-':
          out.add(const _Token(_TK.minus, '-'));
          _i++;
        case '*':
          out.add(const _Token(_TK.star, '*'));
          _i++;
        case '>':
          if (_i + 1 < _src.length && _src[_i + 1] == '=') {
            out.add(const _Token(_TK.gte, '>='));
            _i += 2;
          } else {
            out.add(const _Token(_TK.gt, '>'));
            _i++;
          }
        case '<':
          if (_i + 1 < _src.length && _src[_i + 1] == '=') {
            out.add(const _Token(_TK.lte, '<='));
            _i += 2;
          } else {
            out.add(const _Token(_TK.lt, '<'));
            _i++;
          }
        case '"':
          out.add(_phrase());
        default:
          out.add(_word());
      }
    }
    return out;
  }

  _Token _phrase() {
    _i++; // skip opening "
    final buf = StringBuffer();
    while (!_done && _c != '"') {
      if (_c == '\\' && _i + 1 < _src.length) _i++;
      buf.write(_c);
      _i++;
    }
    if (!_done) _i++; // skip closing "
    return _Token(_TK.phrase, buf.toString());
  }

  _Token _word() {
    const stops = {
      ' ',
      '\t',
      '\n',
      '\r',
      ':',
      '[',
      ']',
      '{',
      '}',
      '(',
      ')',
      '"',
      '>',
      '<',
    };
    final buf = StringBuffer();
    while (!_done && !stops.contains(_c)) {
      buf.write(_c);
      _i++;
    }
    return _Token(_TK.word, buf.toString());
  }
}

// ─── Parser ───────────────────────────────────────────────────────────────────

const _keywords = {'AND', 'OR', 'NOT', 'TO'};

class _Parser {
  final List<_Token> _tokens;
  int _i = 0;

  _Parser(this._tokens);

  _Token get _cur => _tokens[_i.clamp(0, _tokens.length - 1)];
  bool _is(_TK k) => _cur.kind == k;
  bool _kw(String v) => _cur.kind == _TK.word && _cur.value.toUpperCase() == v;

  _Token _eat() {
    final t = _cur;
    if (_i < _tokens.length - 1) _i++;
    return t;
  }

  Query parse() {
    final clauses = <Clause>[];
    while (!_is(_TK.eof) && !_is(_TK.rparen)) {
      final c = _clause();
      if (c != null) clauses.add(c);
    }
    return Query(clauses);
  }

  Clause? _clause() {
    // skip connectors — sequence implies OR by default
    if (_kw('AND') || _kw('OR')) {
      _eat();
      return null;
    }

    var occur = Occur.should;
    if (_is(_TK.plus)) {
      _eat();
      occur = Occur.must;
    } else if (_is(_TK.minus)) {
      _eat();
      occur = Occur.mustNot;
    } else if (_kw('NOT')) {
      _eat();
      occur = Occur.mustNot;
    }

    final node = _node();
    return node == null ? null : Clause(occur, node);
  }

  Node? _node() {
    // group
    if (_is(_TK.lparen)) {
      _eat(); // (
      final q = parse();
      if (_is(_TK.rparen)) _eat(); // )
      return Group(q);
    }

    // field query: word/phrase followed by colon
    if ((_is(_TK.word) || _is(_TK.phrase)) && _i + 1 < _tokens.length && _tokens[_i + 1].kind == _TK.colon) {
      return _field();
    }

    // bare word
    if (_is(_TK.word)) {
      final v = _eat().value;
      return _keywords.contains(v.toUpperCase()) ? null : Term('', v);
    }

    // bare phrase
    if (_is(_TK.phrase)) return Term('', _eat().value);

    _eat(); // skip unexpected token
    return null;
  }

  Node _field() {
    final name = _eat().value; // field name
    _eat(); // colon

    if (_is(_TK.lbracket) || _is(_TK.lbrace)) return _range(name);
    if (_is(_TK.gt) || _is(_TK.gte) || _is(_TK.lt) || _is(_TK.lte)) {
      return _comparison(name);
    }
    if (_is(_TK.word) || _is(_TK.phrase) || _is(_TK.star)) {
      return Term(name, _eat().value);
    }
    return Term(name, '');
  }

  Range _range(String field) {
    final open = _eat(); // [ or {
    final minInclusive = open.kind == _TK.lbracket;

    final min = _rangeValue();

    if (_kw('TO')) _eat();

    final max = _rangeValue();

    final maxInclusive = _is(_TK.rbracket);
    if (_is(_TK.rbracket) || _is(_TK.rbrace)) _eat();

    return Range(
      field,
      min: min,
      max: max,
      minInclusive: minInclusive,
      maxInclusive: maxInclusive,
    );
  }

  // Read a range bound: '*' → null, or a sequence of word+colon tokens joined
  // together. Joining handles ISO-8601 timestamps like 2020-01-01T00:00:00.000Z
  // whose colons the lexer tokenises as separate colon tokens.
  String? _rangeValue() {
    if (_is(_TK.star)) {
      _eat();
      return null;
    }
    final buf = StringBuffer();
    while ((_is(_TK.word) || _is(_TK.colon)) && !_kw('TO')) {
      buf.write(_eat().value);
    }
    final s = buf.toString();
    return s.isEmpty ? null : s;
  }

  Range _comparison(String field) {
    final op = _eat();
    final v = (_is(_TK.word) || _is(_TK.star)) ? _eat().value : '';
    return switch (op.kind) {
      _TK.gt => Range(field, min: v, minInclusive: false),
      _TK.gte => Range(field, min: v, minInclusive: true),
      _TK.lt => Range(field, max: v, maxInclusive: false),
      _TK.lte => Range(field, max: v, maxInclusive: true),
      _ => Range(field, min: v),
    };
  }
}
