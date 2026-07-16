import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/noop.dart';

class _Item {
  final String id;
  final String label;
  const _Item(this.id, this.label);
}

void main() {
  group('fnNoop', () {
    test('does nothing and returns void', () {
      expect(() => fnNoop<int>(1), returnsNormally);
    });
  });

  group('fnOnChange', () {
    final items = [
      _Item('a', 'A'),
      _Item('b', 'B'),
      _Item('c', 'C'),
    ];

    test('replaces the matching item when v is non-null', () {
      final updated = fnOnChange(items, _Item('b', 'B-updated'), (o) => o.id == 'b');

      expect(updated.map((o) => o.label).toList(), ['A', 'B-updated', 'C']);
      expect(updated, isNot(same(items)));
    });

    test('leaves the list unchanged (by value) when no item matches and v is non-null', () {
      final updated = fnOnChange(items, _Item('z', 'Z'), (o) => o.id == 'z');

      expect(updated.map((o) => o.id).toList(), ['a', 'b', 'c']);
    });

    test('removes the matching item when v is null', () {
      final updated = fnOnChange<_Item>(items, null, (o) => o.id == 'b');

      expect(updated.map((o) => o.id).toList(), ['a', 'c']);
    });

    test('leaves the list unchanged when v is null and no item matches', () {
      final updated = fnOnChange<_Item>(items, null, (o) => o.id == 'z');

      expect(updated.map((o) => o.id).toList(), ['a', 'b', 'c']);
    });

    test('replaces every matching item when the predicate matches more than one', () {
      final dupes = [_Item('a', 'A1'), _Item('a', 'A2'), _Item('b', 'B')];

      final updated = fnOnChange(dupes, _Item('a', 'A-merged'), (o) => o.id == 'a');

      expect(updated.map((o) => o.label).toList(), ['A-merged', 'A-merged', 'B']);
    });

    test('removes every matching item when v is null and the predicate matches more than one', () {
      final dupes = [_Item('a', 'A1'), _Item('a', 'A2'), _Item('b', 'B')];

      final updated = fnOnChange<_Item>(dupes, null, (o) => o.id == 'a');

      expect(updated.map((o) => o.id).toList(), ['b']);
    });

    test('handles an empty source list for both branches', () {
      expect(fnOnChange<_Item>([], _Item('a', 'A'), (o) => o.id == 'a'), isEmpty);
      expect(fnOnChange<_Item>([], null, (o) => o.id == 'a'), isEmpty);
    });
  });
}
