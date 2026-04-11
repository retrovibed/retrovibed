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
  final numField = Number('peers', 0, 0, (_) {});
  final fields = <Field<dynamic>>[numField];

  group('Value.consume', () {
    test('accumulates chars into partial', () {
      final parser = sm.Parser(fields, _noop);
      final result = parser.consume(_ctrl('@peers:42', cursor: 9));
      expect(result, isA<sm.Value>());
      final v = result as sm.Value;
      final colon = v.ctx.partial.indexOf(':');
      final valuePartial =
          colon == -1 ? v.ctx.partial : v.ctx.partial.substring(colon + 1);
      expect(valuePartial, '42');
    });

    test('multi-char value accumulates fully', () {
      final parser = sm.Parser(fields, _noop);
      final result = parser.consume(_ctrl('@peers:100', cursor: 10));
      expect(result, isA<sm.Value>());
      final v = result as sm.Value;
      final colon = v.ctx.partial.indexOf(':');
      final valuePartial =
          colon == -1 ? v.ctx.partial : v.ctx.partial.substring(colon + 1);
      expect(valuePartial, '100');
    });

    test('cursor <= 0 transitions to Query', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:42', cursor: 9));
      final result = parser.consume(_ctrl('', cursor: 0));
      expect(result, isA<sm.Query>());
    });

    test('deleting colon transitions back to Input', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:', cursor: 7));
      final result = parser.consume(_ctrl('@peers', cursor: 6));
      expect(result, isA<sm.Input>());
    });

    test('Input partial is restored when colon is deleted', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:', cursor: 7));
      final result = parser.consume(_ctrl('@peers', cursor: 6));
      expect((result as sm.Input).ctx.partial, 'peers');
    });

    test('deleting colon with value typed transitions back to Input', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:42', cursor: 9));
      // Delete ':42', leaving '@peers'
      final result = parser.consume(_ctrl('@peers', cursor: 6));
      expect(result, isA<sm.Input>());
    });

    test('offset is preserved while accumulating value', () {
      final parser = sm.Parser(fields, _noop);
      final result = parser.consume(_ctrl('@peers:42', cursor: 9));
      expect((result as sm.Value).ctx.offset, 0);
    });

    test('space is accumulated into value (no whitespace termination)', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:42', cursor: 9));
      final result = parser.consume(_ctrl('@peers:42 ', cursor: 10));
      expect(result, isA<sm.Value>());
    });

    test('newline is accumulated into value (no whitespace termination)', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@peers:42', cursor: 9));
      final result = parser.consume(_ctrl('@peers:42\n', cursor: 10));
      expect(result, isA<sm.Value>());
    });
  });
}
