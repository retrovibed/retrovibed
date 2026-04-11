import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/wireguard/edit.dart';
import 'package:retrovibed/wireguard/meta.wireguard.pb.dart';

Wireguard _wireguard({
  String id = 'wg-1',
  String description = 'test config',
  int port = 51820,
}) {
  return Wireguard(id: id, description: description, port: port);
}

void main() {
  group('Edit Widget Tests', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          Scaffold(body: Edit(_wireguard())),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.byType(Edit), findsOneWidget);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: BoxConstraints(maxWidth: 300, maxHeight: 400),
            child: SingleChildScrollView(child: Edit(_wireguard())),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('displays description and port fields', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Scaffold(body: Edit(_wireguard())),
        );
        await tester.pumpAndSettle();

        expect(find.text('test config'), findsOneWidget);
        expect(find.text('51820'), findsOneWidget);
        expect(find.text('description'), findsOneWidget);
        expect(find.text('port'), findsOneWidget);
      });

      testWidgets('displays delete button', (WidgetTester tester) async {
        await tester.pumpApp(
          Scaffold(body: Edit(_wireguard())),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.delete), findsOneWidget);
      });
    });

    group('callbacks', () {
      testWidgets('calls onChange when description is modified', (
        WidgetTester tester,
      ) async {
        Wireguard? changed;

        await tester.pumpApp(
          Scaffold(
            body: Edit(
              _wireguard(),
              onChange: (wg) => changed = wg,
            ),
          ),
        );
        await tester.pumpAndSettle();

        final descField = find.widgetWithText(TextFormField, 'test config');
        await tester.enterText(descField, 'updated config');

        expect(changed, isNotNull);
        expect(changed?.description, 'updated config');
      });

      testWidgets('calls onChange when port is modified', (
        WidgetTester tester,
      ) async {
        Wireguard? changed;

        await tester.pumpApp(
          Scaffold(
            body: Edit(
              _wireguard(),
              onChange: (wg) => changed = wg,
            ),
          ),
        );
        await tester.pumpAndSettle();

        final portField = find.widgetWithText(TextFormField, '51820');
        await tester.enterText(portField, '12345');

        expect(changed, isNotNull);
        expect(changed?.port, 12345);
      });

      testWidgets('does not call onChange for non-numeric port input', (
        WidgetTester tester,
      ) async {
        Wireguard? changed;

        await tester.pumpApp(
          Scaffold(
            body: Edit(
              _wireguard(),
              onChange: (wg) => changed = wg,
            ),
          ),
        );
        await tester.pumpAndSettle();

        final portField = find.widgetWithText(TextFormField, '51820');
        await tester.enterText(portField, 'abc');

        expect(changed, isNull);
      });

      testWidgets('handles onChange being null', (WidgetTester tester) async {
        await tester.pumpApp(
          Scaffold(body: Edit(_wireguard())),
        );
        await tester.pumpAndSettle();

        final descField = find.widgetWithText(TextFormField, 'test config');
        await tester.enterText(descField, 'updated');

        expect(tester.takeException(), isNull);
      });
    });
  });
}
