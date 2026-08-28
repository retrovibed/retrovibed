import 'dart:io' as io;
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/retrovibed.dart' as retro;

class _Exited implements Exception {
  final int code;
  const _Exited(this.code);
  @override
  String toString() => 'Exited($code)';
}

Future<int> _captureFault(Future<void> Function() fn) async {
  int? code;
  try {
    await io.IOOverrides.runZoned(fn, exit: (c) {
      code = c;
      throw _Exited(c);
    });
  } on _Exited {
    // expected — exit was intercepted
  }
  return code!;
}

void _mockWindowManager() {
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    (call) async {
      switch (call.method) {
        // waitUntilReadyToShow reads these as bool and casts the channel
        // result directly - a null response (the default below) would throw.
        case 'isFullScreen':
        case 'isMaximized':
        case 'isMinimized':
          return false;
        default:
          return null;
      }
    },
  );
}

void _clearWindowManagerMock() {
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    null,
  );
}

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  group('fault', () {
    test('fault(0) exits with code 0', () async {
      expect(await _captureFault(() async => retro.fault(0)), equals(0));
    });

    test('fault(1) exits with code 1', () async {
      expect(await _captureFault(() async => retro.fault(1)), equals(1));
    });

    test('fault(2) exits with code 2', () async {
      expect(await _captureFault(() async => retro.fault(2)), equals(2));
    });

    test('fault(127) exits with code 127', () async {
      expect(await _captureFault(() async => retro.fault(127)), equals(127));
    });
  });

  group('run — fn uses fault', () {
    setUp(_mockWindowManager);
    tearDown(_clearWindowManagerMock);

    test('run calls fn and fault(0) exits cleanly', () async {
      expect(
        await _captureFault(() => retro.run(() => retro.fault(0))),
        equals(0),
      );
    });

    test('run calls fn and fault(1) exits with error', () async {
      expect(
        await _captureFault(() => retro.run(() => retro.fault(1))),
        equals(1),
      );
    });

    test('run calls fn and fault(2) exits with code 2', () async {
      expect(
        await _captureFault(() => retro.run(() => retro.fault(2))),
        equals(2),
      );
    });
  });
}
