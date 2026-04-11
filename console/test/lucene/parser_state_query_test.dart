import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene/field.dart';
import 'package:retrovibed/lucene/parser.results.dart';
import 'package:retrovibed/lucene/parser.states.dart' as sm;

TextEditingController _ctrl(String text, {int? cursor}) {
  final c = TextEditingController(text: text);
  c.selection = TextSelection.collapsed(offset: cursor ?? text.length);
  return c;
}

void _noop(
  sm.Context ctx,
  TextRange range,
  String contents, {
  ParserResult? completed,
}) {}

void main() {
  final fields = <Field<dynamic>>[Boolean('hd', false, false, (_) {})];

  group('Query.consume', () {
    test('cursor == 0 returns Query with 0 remaining', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop));
      final (next, remaining) = state.consume(_ctrl('', cursor: 0));
      expect(next, isA<sm.Query>());
      expect(remaining, 0);
    });

    test('cursor < lastOffset resets lastOffset to ctx.offset', () {
      final state = sm.Query(
        sm.Context(fields, 2, 'hel', _noop, lastOffset: 5),
      );
      final (next, remaining) = state.consume(_ctrl('hello world', cursor: 3));
      final q = next as sm.Query;
      expect(next, isA<sm.Query>());
      expect(q.ctx.lastOffset, 2); // reset to ctx.offset
      expect(q.ctx.offset, 2); // offset preserved
      expect(q.ctx.partial, ''); // partial cleared
      expect(remaining, 1); // cursor(3) - resetLastOffset(2)
    });

    test('cursor == lastOffset returns Query with 0 remaining', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop, lastOffset: 5));
      final (next, remaining) = state.consume(_ctrl('hello', cursor: 5));
      expect(next, isA<sm.Query>());
      expect(remaining, 0);
    });

    test('non-@ char advances lastOffset by 1', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop));
      final (next, remaining) = state.consume(_ctrl('hello world'));
      expect(next, isA<sm.Query>());
      expect((next as sm.Query).ctx.lastOffset, 1);
      expect(remaining, 10);
    });

    test('@ at cursor transitions to Input', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop));
      final (next, remaining) = state.consume(_ctrl('@', cursor: 1));
      expect(next, isA<sm.Input>());
      expect(remaining, 0);
    });

    test('@ followed by more text returns Input with remaining chars', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop));
      final (next, remaining) = state.consume(_ctrl('@hd', cursor: 3));
      expect(next, isA<sm.Input>());
      expect(remaining, 2); // 'h' and 'd' after @
    });

    test('@ anchor offset is the position of @ in text', () {
      // set lastOffset to point directly at the @ character
      final state = sm.Query(sm.Context(fields, -1, '', _noop, lastOffset: 6));
      final (next, _) = state.consume(_ctrl('hello @hd'));
      expect(next, isA<sm.Input>());
      expect((next as sm.Input).ctx.offset, 6);
    });

    test('transitions on first @ when multiple @ present', () {
      final state = sm.Query(sm.Context(fields, -1, '', _noop));
      final (next, _) = state.consume(_ctrl('@foo @hd'));
      expect(next, isA<sm.Input>());
      expect((next as sm.Input).ctx.offset, 0);
    });

    test('remaining counts unprocessed chars after @ anchor', () {
      // set lastOffset to point directly at the @ character
      final state = sm.Query(sm.Context(fields, -1, '', _noop, lastOffset: 6));
      final (_, remaining) = state.consume(_ctrl('hello @hd'));
      // '@' at index 6, lastOffset becomes 7, cursor at 9 → remaining = 2
      expect(remaining, 2);
    });

    test('replacing a non-@ char with @ transitions to Input', () {
      final parser = sm.Parser(fields, _noop);
      // Establish partial='a' at lastOffset=1
      parser.consume(_ctrl('a', cursor: 1));
      // Replace 'a' with '@' — cursor stays at 1
      final result = parser.consume(_ctrl('@', cursor: 1));
      expect(result, isA<sm.Input>());
    });

    test('full rewrite from "abc a" to "@abc" transitions to Input', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('abc a', cursor: 5));
      final result = parser.consume(_ctrl('@abc', cursor: 4));
      expect(result, isA<sm.Input>());
      expect((result as sm.Input).ctx.partial, 'abc');
    });
  });
}
