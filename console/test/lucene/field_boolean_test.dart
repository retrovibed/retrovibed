import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;

void main() {
  group('Boolean', () {
    final field = lucene.Boolean.auto('hd', false, (_) {});

    group('parse', () {
      test('true', () => expect(field.parse('true'), isTrue));
      test('1', () => expect(field.parse('1'), isTrue));
      test('yes', () => expect(field.parse('yes'), isTrue));
      test('false', () => expect(field.parse('false'), isFalse));
      test('0', () => expect(field.parse('0'), isFalse));
      test('no', () => expect(field.parse('no'), isFalse));
      test('case insensitive YES', () => expect(field.parse('YES'), isTrue));
      test('case insensitive NO', () => expect(field.parse('NO'), isFalse));
      test('unknown returns defaultValue', () => expect(field.parse('maybe'), isFalse));
    });

    group('apply', () {
      test('true sets setter to true', () {
        bool result = false;
        final f = lucene.Boolean.auto('hd', false, (v) => result = v);
        f.apply(lucene.parse('hd:true').clauses.first.node);
        expect(result, isTrue);
      });

      test('false sets setter to false', () {
        bool result = true;
        final f = lucene.Boolean.auto('hd', true, (v) => result = v);
        f.apply(lucene.parse('hd:false').clauses.first.node);
        expect(result, isFalse);
      });

      test('unrelated field does not call setter', () {
        bool called = false;
        final f = lucene.Boolean.auto('hd', false, (_) => called = true);
        f.apply(lucene.parse('other:true').clauses.first.node);
        expect(called, isFalse);
      });
    });

    group('autocomplete', () {
      test('returns !defaultValue when defaultValue is false', () {
        final f = lucene.Boolean.auto('hd', false, (_) {});
        expect(f.autocomplete, isTrue);
      });

      test('returns !defaultValue when defaultValue is true', () {
        final f = lucene.Boolean.auto('hd', true, (_) {});
        expect(f.autocomplete, isFalse);
      });
    });

    group('suggestions', () {
      test('empty partial returns both on and off', () {
        final s = field.suggestions('');
        expect(s.map((s) => s.label), containsAll(['on', 'off']));
      });

      test('on prefix returns only on', () {
        final s = field.suggestions('on');
        expect(s.length, 1);
        expect(s.first.label, 'on');
      });

      test('of prefix returns only off', () {
        final s = field.suggestions('of');
        expect(s.length, 1);
        expect(s.first.label, 'off');
      });

      test('no match returns empty', () {
        expect(field.suggestions('xyz'), isEmpty);
      });
    });

    group('from', () {
      test('true produces Term hd:true', () {
        final node = field.from(true) as lucene.Term;
        expect(node.field, 'hd');
        expect(node.value, 'true');
      });

      test('false produces Term hd:false', () {
        final node = field.from(false) as lucene.Term;
        expect(node.field, 'hd');
        expect(node.value, 'false');
      });
    });
  });
}
