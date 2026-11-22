// ignore: deprecated_member_use
import 'package:flutter/material.dart' as sdk show Theme, IconTheme;
import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/design.kit/material.legacy.dart';

void main() {
  testWidgets('sdk theme mirrors material_ui dark theme', (WidgetTester tester) async {
    late Brightness brightness;
    late Color? icon;
    await tester.pumpWidget(
      MaterialApp(
        darkTheme: ThemeData(brightness: Brightness.dark),
        themeMode: ThemeMode.dark,
        home: MaterialLegacy(
          Builder(
            builder: (outer) => sdk.Theme(
              // mirrors how media_kit_video wraps its controls.
              data: sdk.Theme.of(outer).copyWith(),
              child: Builder(
                builder: (context) {
                  brightness = sdk.Theme.of(context).brightness;
                  icon = sdk.IconTheme.of(context).color;
                  return const SizedBox();
                },
              ),
            ),
          ),
        ),
      ),
    );

    expect(brightness, Brightness.dark);
    expect(icon, isNotNull);
    expect(icon!.computeLuminance(), greaterThan(0.5));
  });
}
