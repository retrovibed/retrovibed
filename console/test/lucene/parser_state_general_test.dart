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
  final numFields = <Field<dynamic>>[Number('peers', 0, 0, (_) {})];

  group('UnknownFieldError', () {
    test('stays in error when text is unchanged', () {
      final parser = sm.Parser(fields, _noop);
      final result = parser.consume(_ctrl('@foo:', cursor: 5));
      expect(result, isA<sm.UnknownFieldError>());
      expect(
        parser.consume(_ctrl('@foo:', cursor: 5)),
        isA<sm.UnknownFieldError>(),
      );
    });

    test('partial updates when text changes but field is still unknown', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('@bar:', cursor: 5));
      expect(result, isA<sm.UnknownFieldError>());
      expect((result as sm.UnknownFieldError).ctx.partial, 'bar');
    });

    test('transitions to Query when @ anchor is gone', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('foo:', cursor: 4));
      expect(result, isA<sm.Query>());
    });

    test('transitions to Input when colon is removed', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('@foo', cursor: 4));
      expect(result, isA<sm.Input>());
    });

    test('Input partial is restored when colon is removed', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('@foo', cursor: 4));
      expect((result as sm.Input).ctx.partial, 'foo');
    });

    test('transitions to Value when partial corrected to a known field', () {
      final parser = sm.Parser(numFields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('@peers:', cursor: 7));
      expect(result, isA<sm.Value>());
    });

    test('transitions to Query when cursor is at zero', () {
      final parser = sm.Parser(fields, _noop);
      parser.consume(_ctrl('@foo:', cursor: 5));
      final result = parser.consume(_ctrl('', cursor: 0));
      expect(result, isA<sm.Query>());
    });

    test('Input returns UnknownFieldError for unrecognised field', () {
      final parser = sm.Parser(fields, _noop);
      final result = parser.consume(_ctrl('@foo:', cursor: 5));
      expect(result, isA<sm.UnknownFieldError>());
    });

    test(
      'stays in error with original partial when typing continues after colon',
      () {
        final parser = sm.Parser(fields, _noop);
        parser.consume(_ctrl('@foo:', cursor: 5));
        final inputs = [
          '@foo:@',
          '@foo:@d',
          '@foo:@de',
          '@foo:@der',
          '@foo:@derp',
        ];
        for (final text in inputs) {
          final result = parser.consume(_ctrl(text));
          expect(
            result,
            isA<sm.UnknownFieldError>(),
            reason: 'expected error for "$text"',
          );
          expect(
            (result as sm.UnknownFieldError).ctx.partial,
            'foo',
            reason: 'partial should remain "foo" for "$text"',
          );
        }
      },
    );
  });
}
