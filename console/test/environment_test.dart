import 'dart:io' show Platform;
import 'package:flutter/foundation.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('ci environment validation', (WidgetTester tester) async {
    const cursor = SystemMouseCursors.forbidden;
    final platform =
        debugDefaultTargetPlatformOverride ?? defaultTargetPlatform;
    printOnFailure('Dart version: ${Platform.version}');
    printOnFailure('Target Platform: $platform');
    printOnFailure('Screen Size: ${tester.view.physicalSize}');
    printOnFailure('Device Pixel Ratio: ${tester.view.devicePixelRatio}');
    printOnFailure('Cursor kind: ${cursor.kind}');

    // validate initial expectations for the environment.
    expect(Platform.version, startsWith("3.12"));
    expect(tester.view.physicalSize.width, 2400.0);
    expect(tester.view.physicalSize.height, 1800.0);
    expect(tester.view.devicePixelRatio, 3.0);
    expect(
      cursor.kind,
      isNot('basic'),
      reason:
          "If this returns 'basic' in CI, the Docker image is missing cursor libraries/themes",
    );
  }, variant: TargetPlatformVariant.all());
}
