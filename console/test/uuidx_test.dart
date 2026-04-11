import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

void main() {
  group('fromString', () {
    test('empty string returns nil UUID', () {
      expect(uuidx.isMin(uuidx.fromString('')), isTrue);
    });

    test('nil UUID string returns nil UUID', () {
      expect(uuidx.isMin(uuidx.fromString('00000000-0000-0000-0000-000000000000')),
          isTrue);
    });

    test('max UUID string returns max UUID', () {
      expect(uuidx.isMax(uuidx.fromString('ffffffff-ffff-ffff-ffff-ffffffffffff')),
          isTrue);
    });

    test('valid UUID string is neither min nor max', () {
      final v = uuidx.fromString('550e8400-e29b-41d4-a716-446655440000');
      expect(uuidx.isMinMax(v), isFalse);
    });
  });

  group('isMinMax', () {
    test('empty string is treated as min', () {
      expect(uuidx.isMinMax(uuidx.fromString('')), isTrue);
    });

    test('nil UUID is min', () {
      expect(uuidx.isMinMax(uuidx.fromString(uuidx.min())), isTrue);
    });

    test('max UUID is max', () {
      expect(uuidx.isMinMax(uuidx.fromString(uuidx.max())), isTrue);
    });

    test('regular UUID is not min or max', () {
      expect(uuidx.isMinMax(uuidx.fromString('550e8400-e29b-41d4-a716-446655440000')),
          isFalse);
    });
  });
}
