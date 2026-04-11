import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;

void main() {
  group('Number', () {
    final field = lucene.Number.auto('peers', 0, (_) {});

    test('parse plain number', () => expect(field.parse('5'), 5));
    test('parse with > operator stripped', () => expect(field.parse('>5'), 5));
    test(
      'parse with >= operator stripped',
      () => expect(field.parse('>=10'), 10),
    );
    test('parse float', () => expect(field.parse('3.14'), 3.14));
    test(
      'parse non-number returns defaultValue',
      () => expect(field.parse('abc'), 0),
    );

    test('apply from term sets setter', () {
      num? result;
      final f = lucene.Number.auto('peers', 0, (v) => result = v);
      f.apply(lucene.parse('peers:5').clauses.first.node);
      expect(result, 5);
    });
  });
}
