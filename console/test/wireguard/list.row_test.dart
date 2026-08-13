import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/wireguard/list.row.dart';
import 'package:retrovibed/wireguard/edit.dart';
import 'package:retrovibed/wireguard/meta.wireguard.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Wireguard _wireguard({
  String id = 'wg-1',
  String description = 'test config',
  int port = 51820,
}) {
  return Wireguard(id: id, description: description, port: port);
}

void main() {
  group('ListRow Widget Tests', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(ListRow(_wireguard()));
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.byType(ListRow), findsOneWidget);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 300,
            height: 400,
            child: ListRow(_wireguard()),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('displays description text', (WidgetTester tester) async {
        await tester.pumpApp(ListRow(_wireguard(description: 'my vpn')));
        await tester.pumpAndSettle();

        expect(find.text('my vpn'), findsOneWidget);
      });

      testWidgets('shows checkmark when active', (WidgetTester tester) async {
        await tester.pumpApp(ListRow(_wireguard(), active: true));
        await tester.pumpAndSettle();

        final opacity = tester.widget<Opacity>(find.byType(Opacity));
        expect(opacity.opacity, 1.0);
      });

      testWidgets('hides checkmark when inactive', (WidgetTester tester) async {
        await tester.pumpApp(ListRow(_wireguard(), active: false));
        await tester.pumpAndSettle();

        final opacity = tester.widget<Opacity>(find.byType(Opacity));
        expect(opacity.opacity, 0.1);
      });

      testWidgets('renders leading widget between checkmark and text', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ListRow(
            _wireguard(),
            leading: Icon(Icons.vpn_key, key: Key('leading')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byKey(Key('leading')), findsOneWidget);
      });

      testWidgets('renders trailing widget after text', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ListRow(
            _wireguard(),
            trailing: Icon(Icons.info, key: Key('trailing')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byKey(Key('trailing')), findsOneWidget);
      });
    });

    group('expand/collapse', () {
      testWidgets('edit panel is hidden by default', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(ListRow(_wireguard()));
        await tester.pumpAndSettle();

        final visibility = tester.widget<Visibility>(find.byType(Visibility).last);
        expect(visibility.visible, isFalse);
      });

      testWidgets('tapping row toggles edit panel visibility', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(ListRow(_wireguard()));
        await tester.pumpAndSettle();

        // Tap the row to expand
        await tester.tap(find.text('test config').first);
        await tester.pumpAndSettle();

        final visibilityAfterTap = tester.widget<Visibility>(
          find.byType(Visibility).last,
        );
        expect(visibilityAfterTap.visible, isTrue);

        // Tap again to collapse
        await tester.tap(find.text('test config').first);
        await tester.pumpAndSettle();

        final visibilityAfterSecondTap = tester.widget<Visibility>(
          find.byType(Visibility).last,
        );
        expect(visibilityAfterSecondTap.visible, isFalse);
      });

      testWidgets('expanded panel contains Edit widget', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(ListRow(_wireguard()));
        await tester.pumpAndSettle();

        await tester.tap(find.text('test config').first);
        await tester.pump();

        expect(find.byType(Edit), findsOneWidget);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ListRow(_wireguard()),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);

      testWidgets('renders without overflow with long description', (
        WidgetTester tester,
      ) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ListRow(
            _wireguard(
              description: 'A very long wireguard configuration name that should overflow',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);

      testWidgets('renders expanded panel without overflow', (
        WidgetTester tester,
      ) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(child: ListRow(_wireguard())),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.text('test config').first);
        await tester.pump();

        expect(tester.takeException(), isNull);
        expect(find.byType(Edit), findsOneWidget);
      }, variant: _resolutions);
    });
  });
}
