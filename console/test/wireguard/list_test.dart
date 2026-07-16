import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/wireguard/list.dart';
import 'package:retrovibed/wireguard/list.row.dart';
import 'package:retrovibed/wireguard/api.dart' as api;
import 'package:retrovibed/wireguard/meta.wireguard.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Wireguard _wireguard({required String id, required String description}) {
  return Wireguard(id: id, description: description, port: 51820);
}

Future<api.WireguardSearchResponse> _mockSearch(api.WireguardSearchRequest req) async {
  return api.WireguardSearchResponse(
    items: [
      _wireguard(id: 'wg-1', description: 'Config One'),
      _wireguard(id: 'wg-2', description: 'Config Two'),
    ],
    next: req,
  );
}

void main() {
  group('wireguard.ListDisplay onChange/onDelete wiring', () {
    testWidgets('onChange replaces the matching peer in place', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(ListDisplay(search: _mockSearch));
      await tester.pumpAndSettle();

      final rows = tester.widgetList<ListRow>(find.byType(ListRow)).toList();
      expect(rows.length, 2);
      final target = rows.firstWhere((r) => r.current.id == 'wg-1');

      target.onChange(_wireguard(id: 'wg-1', description: 'Config One Updated'));
      await tester.pump();

      expect(find.text('Config One Updated'), findsOneWidget);
      expect(find.text('Config One'), findsNothing);
      expect(find.text('Config Two'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('onDelete removes the matching peer', (WidgetTester tester) async {
      await tester.pumpApp(ListDisplay(search: _mockSearch));
      await tester.pumpAndSettle();

      final rows = tester.widgetList<ListRow>(find.byType(ListRow)).toList();
      expect(rows.length, 2);
      final target = rows.firstWhere((r) => r.current.id == 'wg-1');

      target.onDelete(target.current);
      await tester.pump();

      final remaining = tester.widgetList<ListRow>(find.byType(ListRow)).toList();
      expect(remaining.length, 1);
      expect(remaining.single.current.id, 'wg-2');
      expect(find.text('Config One'), findsNothing);
      expect(find.text('Config Two'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
