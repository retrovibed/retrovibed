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
  final boolField = Boolean('hd', false, false, (_) {});
  final numField = Number('peers', 0, 0, (_) {});
  final fields = <Field<dynamic>>[boolField, numField];

  group('Input.consume', () {
    // Helpers to construct Input directly:
    // offset = position of '@', lastOffset = next char to process, partial = accumulated field name.

    test('single char appended to partial', () {
      // '@' at index 0, processing 'h' at index 1
      final state = sm.Input(sm.Context(fields, 0, '', _noop, lastOffset: 1));
      final (next, remaining) = state.consume(_ctrl('@h', cursor: 2));
      final i = next as sm.Input;
      expect(next, isA<sm.Input>());
      expect(i.ctx.partial, 'h');
      expect(i.ctx.lastOffset, 2);
      expect(remaining, 0);
    });

    test('advances lastOffset by 1 per char', () {
      final state = sm.Input(sm.Context(fields, 0, 'h', _noop, lastOffset: 2));
      final (next, remaining) = state.consume(_ctrl('@hd', cursor: 3));
      final i = next as sm.Input;
      expect(i.ctx.partial, 'hd');
      expect(i.ctx.lastOffset, 3);
      expect(remaining, 0);
    });

    test('remaining counts unprocessed chars', () {
      // '@' at 0, processing 'h' at lastOffset=1, cursor at 4
      final state = sm.Input(sm.Context(fields, 0, '', _noop, lastOffset: 1));
      final (_, remaining) = state.consume(_ctrl('@hd:', cursor: 4));
      // processes 'h', nextLastOffset=2, cursor=4 → remaining=2
      expect(remaining, 2);
    });

    test('offset is preserved as position of @', () {
      // '@' at index 6 in 'hello @hd'
      final state = sm.Input(sm.Context(fields, 6, '', _noop, lastOffset: 7));
      final (next, _) = state.consume(_ctrl('hello @hd', cursor: 9));
      expect((next as sm.Input).ctx.offset, 6);
    });

    test(': after field with autocomplete completes and returns to Query', () {
      ParserResult? captured;
      void capture(
        sm.Context ctx,
        TextRange range,
        String contents, {
        ParserResult? completed,
      }) {
        captured = completed;
      }

      final captureFields = <Field<dynamic>>[
        Boolean('hd', false, false, (_) {}),
      ];
      final state = sm.Input(
        sm.Context(captureFields, 0, 'hd', capture, lastOffset: 3),
      );
      final (next, remaining) = state.consume(_ctrl('@hd:', cursor: 4));
      expect(next, isA<sm.Query>());
      expect(remaining, 0);
      expect(captured, isA<ParserResultBool>());
    });

    test(': after field without autocomplete transitions to Value', () {
      // Number field has no autocomplete
      final numFields = <Field<dynamic>>[numField];
      final state = sm.Input(
        sm.Context(numFields, 0, 'peers', _noop, lastOffset: 6),
      );
      final (next, remaining) = state.consume(_ctrl('@peers:', cursor: 7));
      expect(next, isA<sm.Value>());
      expect(remaining, 0);
    });

    test(': after unknown field name transitions to UnknownFieldError', () {
      final state = sm.Input(
        sm.Context(fields, 0, 'foo', _noop, lastOffset: 4),
      );
      final (next, remaining) = state.consume(_ctrl('@foo:', cursor: 5));
      expect(next, isA<sm.UnknownFieldError>());
      expect(remaining, 0);
    });

    test('full "@peers:" typed at once via Parser reaches Value', () {
      final parser = sm.Parser([numField], _noop);
      final result = parser.consume(_ctrl('@peers:', cursor: 7));
      expect(result, isA<sm.Value>());
    });

    test(
      'full "@hd" typed at once via Parser reaches Input with full partial',
      () {
        final parser = sm.Parser(fields, _noop);
        final result = parser.consume(_ctrl('@hd', cursor: 3));
        expect(result, isA<sm.Input>());
        expect((result as sm.Input).ctx.partial, 'hd');
      },
    );

    test('replacing field char preserves Input state with updated partial', () {
      final parser = sm.Parser(fields, _noop);
      // Establish '@h' → Input with partial='h'
      parser.consume(_ctrl('@h', cursor: 2));
      // Replace 'h' with 'd' — cursor stays at 2
      final result = parser.consume(_ctrl('@d', cursor: 2));
      expect(result, isA<sm.Input>());
      expect((result as sm.Input).ctx.partial, 'd');
    });
  });
}
